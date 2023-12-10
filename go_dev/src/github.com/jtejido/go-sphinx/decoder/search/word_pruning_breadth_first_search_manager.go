package search

import (
	"fmt"
	"github.com/jtejido/go-sphinx/decoder/pruner"
	"github.com/jtejido/go-sphinx/decoder/scorer"
	"github.com/jtejido/go-sphinx/linguist"
	"github.com/jtejido/go-sphinx/result"
	"github.com/jtejido/go-sphinx/utils"
	"math"
)

const (
	// The property that specifies the relative beam width
	DEFAULT_RELATIVE_BEAM_WIDTH = 0.0

	// The property that controls the amount of simple acoustic lookahead
	// performed. Setting the property to zero (the default) disables simple
	// acoustic lookahead. The lookahead need not be an integer.
	DEFAULT_ACOUSTIC_LOOKAHEAD_FRAMES = 0

	// The property that specifies the maximum lattice edges
	DEFAULT_MAX_LATTICE_EDGES = 100

	// The property that specifies the maximum lattice edges
	DEFAULT_CHECK_STATE_ORDER = false

	// The property that controls the number of frames processed for every time
	// the decode growth step is skipped. Setting this property to zero disables
	// grow skipping. Setting this number to a small integer will increase the
	// speed of the decoder but will also decrease its accuracy. The higher the
	// number, the less often the grow code is skipped. Values like 6-8 is known
	// to be the good enough for large vocabulary tasks. That means that one of
	// 6 frames will be skipped.
	DEFAULT_GROW_SKIP_INTERVAL = 0

	// The property than, when set to <code>true</code> will cause the
	// recognizer to count up all the tokens in the active list after every
	// frame.
	DEFAULT_SHOW_TOKEN_COUNT = false
)

// Provides the breadth first search. To perform recognition an application
// should call initialize before recognition begins, and repeatedly call
// recognize until Result.isFinal() returns true. Once a final
// result has been obtained, stopRecognition should be called.
//
// All scores and probabilities are maintained in the log math log domain.
type WordPruningBreadthFirstSearchManager struct {
	TokenSearchManager
	// -----------------------------------
	// Configured Subcomponents
	// -----------------------------------

	// a linguist for search space
	linguist linguist.Linguist

	// pruner to drop tokens
	pruner pruner.Pruner

	// scorer to estimate token probability
	scorer scorer.AcousticScorer

	// active list manager to store tokens
	activeListManager ActiveListManager
	logMath           utils.LogMath

	// -----------------------------------
	// Configuration data
	// -----------------------------------

	// show count during decoding
	showTokenCount bool

	// check order of states during growth
	checkStateOrder bool

	// skip interval for grown
	growSkipInterval int

	// relative beam for lookahead pruning
	relativeBeamWidth float64

	// frames to do lookahead
	acousticLookaheadFrames float64

	// max edges to keep in lattice
	maxLatticeEdges int

	// -----------------------------------
	// Instrumentation
	// -----------------------------------
	scoreTimer, pruneTimer, growTimer                 *utils.Timer
	totalTokensScored, curTokensScored, tokensCreated *utils.StatisticsVariable
	tokenSum                                          int64
	tokenCount                                        int
	// -----------------------------------
	// Working data
	// -----------------------------------

	// the current frame number
	currentFrameNumber int

	// the current frame number
	currentCollectTime int64

	// the list of active tokens
	activeList ActiveList

	// the current set of results
	resultList    []*Token
	bestTokenMap  map[linguist.SearchState]*Token
	loserManager  *AlternateHypothesisManager
	numStateOrder int
	streamEnd     bool
}

func NewWordPruningBreadthFirstSearchManager(linguist linguist.Linguist, pruner pruner.Pruner, scorer scorer.AcousticScorer,
	activeListManager ActiveListManager, showTokenCount bool, relativeWordBeamWidth float64, growSkipInterval int,
	checkStateOrder bool, buildWordLattice bool, maxLatticeEdges int, acousticLookaheadFrames float64,
	keepAllTokens bool) *WordPruningBreadthFirstSearchManager {
	wpbfsm := new(WordPruningBreadthFirstSearchManager)
	wpbfsm.logMath = utils.LogMath{}
	wpbfsm.linguist = linguist
	wpbfsm.pruner = pruner
	wpbfsm.scorer = scorer
	wpbfsm.activeListManager = activeListManager
	wpbfsm.showTokenCount = showTokenCount
	wpbfsm.growSkipInterval = growSkipInterval
	wpbfsm.checkStateOrder = checkStateOrder
	wpbfsm.buildWordLattice = buildWordLattice
	wpbfsm.maxLatticeEdges = maxLatticeEdges
	wpbfsm.acousticLookaheadFrames = acousticLookaheadFrames
	wpbfsm.keepAllTokens = keepAllTokens

	wpbfsm.relativeBeamWidth = wpbfsm.logMath.LinearToLog(relativeWordBeamWidth)
	return wpbfsm
}

func (wpbfsm *WordPruningBreadthFirstSearchManager) Allocate() {

	wpbfsm.scoreTimer = utils.NewTimer("Score")
	wpbfsm.pruneTimer = utils.NewTimer("Prune")
	wpbfsm.growTimer = utils.NewTimer("Grow")

	wpbfsm.totalTokensScored = &utils.StatisticsVariable{Name: "totalTokensScored"}
	wpbfsm.curTokensScored = &utils.StatisticsVariable{Name: "curTokensScored"}
	wpbfsm.tokensCreated = &utils.StatisticsVariable{Name: "tokensCreated"}

	wpbfsm.linguist.Allocate()
	wpbfsm.pruner.Allocate()
	wpbfsm.scorer.Allocate()
}

func (wpbfsm *WordPruningBreadthFirstSearchManager) Deallocate() {
	wpbfsm.scorer.Deallocate()
	wpbfsm.pruner.Deallocate()
	wpbfsm.linguist.Deallocate()
}

// Called at the start of recognition. Gets the search manager ready to recognize
func (wpbfsm *WordPruningBreadthFirstSearchManager) StartRecognition() {
	wpbfsm.linguist.StartRecognition()
	wpbfsm.pruner.StartRecognition()
	wpbfsm.scorer.StartRecognition()
	wpbfsm.localStart()
}

// Performs the recognition for the given number of frames.
func (wpbfsm *WordPruningBreadthFirstSearchManager) Recognize(nFrames int) (res *result.Result) {
	done := false
	wpbfsm.streamEnd = false

	for i := 0; i < nFrames && !done; i++ {
		done = wpbfsm.recognize()
	}

	if !wpbfsm.streamEnd {
		res = result.NewResult(wpbfsm.loserManager, wpbfsm.activeList, wpbfsm.resultList, wpbfsm.currentCollectTime, done, wpbfsm.linguist.GetSearchGraph().GetWordTokenFirst(), true)
	}

	if wpbfsm.showTokenCount {
		wpbfsm.showTokenCount()
	}
	return res
}

func (wpbfsm *WordPruningBreadthFirstSearchManager) recognize() bool {

	wpbfsm.activeList = wpbfsm.activeListManager.GetEmittingList()
	more := wpbfsm.scoreTokens()

	if more {
		wpbfsm.pruneBranches()
		wpbfsm.currentFrameNumber++
		if wpbfsm.growSkipInterval == 0 || (wpbfsm.currentFrameNumber%wpbfsm.growSkipInterval) != 0 {
			wpbfsm.clearCollectors()
			wpbfsm.growEmittingBranches()
			wpbfsm.growNonEmittingBranches()
		}
	}
	return !more
}

// Clears lists and maps before next expansion stage
func (wpbfsm *WordPruningBreadthFirstSearchManager) clearCollectors() {
	wpbfsm.resultList = make([]*Token, 0)
	wpbfsm.createBestTokenMap()
	wpbfsm.activeListManager.ClearEmittingList()
}

// creates a new best token map with the best size
func (wpbfsm *WordPruningBreadthFirstSearchManager) createBestTokenMap() {
	mapSize := wpbfsm.activeList.Size() * 10
	if mapSize == 0 {
		mapSize = 1
	}
	wpbfsm.bestTokenMap = make(map[linguist.SearchState]*Token, mapSize)
}

// Terminates a recognition
func (wpbfsm *WordPruningBreadthFirstSearchManager) StopRecognition() {
	//wpbfsm.localStop() this doesn't have any
	wpbfsm.scorer.StopRecognition()
	wpbfsm.pruner.StopRecognition()
	wpbfsm.linguist.StopRecognition()
}

// Gets the initial grammar node from the linguist and creates a GrammarNodeToken
func (wpbfsm *WordPruningBreadthFirstSearchManager) localStart() {
	searchGraph := wpbfsm.linguist.GetSearchGraph()
	wpbfsm.currentFrameNumber = 0
	wpbfsm.curTokensScored.Value = 0
	wpbfsm.numStateOrder = searchGraph.GetNumStateOrder()
	wpbfsm.activeListManager.SetNumStateOrder(wpbfsm.numStateOrder)
	if buildWordLattice {
		wpbfsm.loserManager = NewAlternateHypothesisManager(wpbfsm.maxLatticeEdges)
	}

	state := searchGraph.GetInitialState()

	wpbfsm.activeList = wpbfsm.activeListManager.GetEmittingList()
	wpbfsm.activeList.Add(NewToken(state, -1))

	wpbfsm.clearCollectors()

	wpbfsm.growBranches()
	wpbfsm.growNonEmittingBranches()

}

// Goes through the active list of tokens and expands each token, finding
// the set of successor tokens until all the successor tokens are emitting
// tokens.
func (wpbfsm *WordPruningBreadthFirstSearchManager) growBranches() {
	wpbfsm.growTimer.Start()
	relativeBeamThreshold := wpbfsm.activeList.GetBeamThreshold()
	// if (logger.isLoggable(Level.FINE)) {
	//     logger.fine("Frame: " + currentFrameNumber + " thresh : " + relativeBeamThreshold + " bs "
	//             + activeList.getBestScore() + " tok " + activeList.getBestToken());
	// }
	aList := wpbfsm.activeList.GetList()
	for _, token := range aList {
		if token.GetScore() >= relativeBeamThreshold && wpbfsm.allowExpansion(token) {
			wpbfsm.collectSuccessorTokens(token)
		}
	}
	wpbfsm.growTimer.Stop()
}

// Grows the emitting branches. This version applies a simple acoustic
// lookahead based upon the rate of change in the current acoustic score.
func (wpbfsm *WordPruningBreadthFirstSearchManager) growEmittingBranches() {
	if wpbfsm.acousticLookaheadFrames <= 0.0 {
		wpbfsm.growBranches()
		return
	}
	wpbfsm.growTimer.Start()
	bestScore := -math.MaxFloat64
	toks := wpbfsm.activeList.GetTokens()
	for _, t := range toks {
		score := t.GetScore() + t.GetAcousticScore()*wpbfsm.acousticLookaheadFrames
		if score > bestScore {
			bestScore = score
		}
	}
	relativeBeamThreshold := bestScore + wpbfsm.relativeBeamWidth
	for _, t := range toks {
		if t.GetScore()+t.GetAcousticScore()*wpbfsm.acousticLookaheadFrames > relativeBeamThreshold {
			wpbfsm.collectSuccessorTokens(t)
		}
	}
	wpbfsm.growTimer.Stop()
}

// Grow the non-emitting branches, until the tokens reach an emitting state.
func (wpbfsm *WordPruningBreadthFirstSearchManager) growNonEmittingBranches() {
	i := wpbfsm.activeListManager.GetNonEmittingListIterator()

	for i.HasNext() {
		wpbfsm.activeList = i.Next()
		if wpbfsm.activeList != nil {
			i.Remove()
			wpbfsm.pruneBranches()
			wpbfsm.growBranches()
		}

	}
}

// Calculate the acoustic scores for the active list. The active list should
// contain only emitting tokens.
func (wpbfsm *WordPruningBreadthFirstSearchManager) scoreTokens() bool {
	var moreTokens bool
	wpbfsm.scoreTimer.Start()
	data := wpbfsm.scorer.CalculateScores(wpbfsm.activeList.GetTokens())
	wpbfsm.scoreTimer.Stop()

	var bestToken *Token

	if data == nil {
		streamEnd = true
	} else {

		_, ok := data.(Token)
		if ok {
			bestToken = data
		}
	}

	if bestToken != nil {
		wpbfsm.currentCollectTime = bestToken.GetCollectTime()
	}

	moreTokens = (bestToken != nil)
	wpbfsm.activeList.SetBestToken(bestToken)

	wpbfsm.monitorStates(activeList)

	wpbfsm.curTokensScored.Value += wpbfsm.activeList.Size()
	wpbfsm.totalTokensScored.Value += wpbfsm.activeList.Size()

	return moreTokens
}

// Keeps track of and reports statistics about the number of active states
func (wpbfsm *WordPruningBreadthFirstSearchManager) monitorStates(activeList ActiveList) {

	wpbfsm.tokenSum += wpbfsm.activeList.Size()
	wpbfsm.tokenCount++

	// if (wpbfsm.tokenCount % 1000) == 0 {
	// 	logger.info("Average Tokens/State: " + (wpbfsm.tokenSum / wpbfsm.tokenCount))
	// }
}

// Removes unpromising branches from the active list
func (wpbfsm *WordPruningBreadthFirstSearchManager) pruneBranches() {
	wpbfsm.pruneTimer.Start()
	wpbfsm.activeList = wpbfsm.pruner.Prune(wpbfsm.activeList)
	wpbfsm.pruneTimer.Stop()
}

// Gets the best token for this state
func (wpbfsm *WordPruningBreadthFirstSearchManager) getBestToken(state) *Token {
	return wpbfsm.bestTokenMap.Get(state)
}

// Sets the best token for a given state
func (wpbfsm *WordPruningBreadthFirstSearchManager) setBestToken(token *Token, state linguist.SearchState) {
	wpbfsm.bestTokenMap.Put(state, token)
}

// Checks that the given two states are in legitimate order.
func (wpbfsm *WordPruningBreadthFirstSearchManager) checkStateOrder(fromState, toState linguist.SearchState) {
	if fromState.GetOrder() == numStateOrder-1 {
		return
	}

	if fromState.GetOrder() > toState.GetOrder() {
		panic(fmt.Sprintf("IllegalState order: from %s %s order: %d to %s %s order: %d", fromState, fromState.ToPrettyString(), fromState.GetOrder(), toState, toState.ToPrettyString(), toState.GetOrder()))
	}
}

// Collects the next set of emitting tokens from a token and accumulates
// them in the active or result lists
func (wpbfsm *WordPruningBreadthFirstSearchManager) collectSuccessorTokens(token *Token) {

	// If this is a final state, add it to the final list
	if token.IsFinal() {
		wpbfsm.resultList.Add(wpbfsm.getResultListPredecessor(token))
		return
	}

	// if this is a non-emitting token and we've already
	// visited the same state during this frame, then we
	// are in a grammar loop, so we don't continue to expand.
	// This check only works properly if we have kept all of the
	// tokens (instead of skipping the non-word tokens).
	// Note that certain linguists will never generate grammar loops
	// (lextree linguist for example). For these cases, it is perfectly
	// fine to disable this check by setting keepAllTokens to false

	if !token.IsEmitting() && (wpbfsm.keepAllTokens && wpbfsm.isVisited(token)) {
		return
	}

	state := token.GetSearchState()
	arcs := state.GetSuccessors()
	predecessor := wpbfsm.getResultListPredecessor(token)

	// For each successor
	// calculate the entry score for the token based upon the
	// predecessor token score and the transition probabilities
	// if the score is better than the best score encountered for
	// the SearchState and frame then create a new token, add
	// it to the lattice and the SearchState.
	// If the token is an emitting token add it to the list,
	// otherwise recursively collect the new tokens successors.

	for _, arc := range arcs {
		nextState := arc.GetState()

		if wpbfsm.checkStateOrder {
			wpbfsm.checkStateOrder(state, nextState)
		}

		// We're actually multiplying the variables, but since
		// these come in log(), multiply gets converted to add
		logEntryScore := token.GetScore() + arc.GetProbability()

		bestToken := wpbfsm.getBestToken(nextState)
		_, ok := nextState.(WordSearchState)
		if bestToken == nil {
			newBestToken := newToken(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbfsm.currentCollectTime)
			wpbfsm.tokensCreated.value++
			wpbfsm.setBestToken(newBestToken, nextState)
			wpbfsm.activeListAdd(newBestToken)
		} else if bestToken.GetScore() < logEntryScore {
			oldPredecessor := bestToken.GetPredecessor()
			bestToken.update(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbfsm.currentCollectTime)

			if wpbfsm.buildWordLattice && ok {
				wpbfsm.loserManager.AddAlternatePredecessor(bestToken, oldPredecessor)
			}
		} else if wpbfsm.buildWordLattice && ok {
			if predecessor != nil {
				wpbfsm.loserManager.AddAlternatePredecessor(bestToken, predecessor)
			}
		}
	}
}

// Determines whether or not we've visited the state associated with this
// token since the previous frame.
func (wpbfsm *WordPruningBreadthFirstSearchManager) isVisited(t *Token) bool {
	curState := t.GetSearchState()

	t = t.GetPredecessor()

	for t != nil && !t.IsEmitting() {
		if curState.Equals(t.GetSearchState()) {
			// System.out.println("CS " + curState + " match " + t.getSearchState());
			return true
		}
		t = t.GetPredecessor()
	}

	return false
}

func (wpbfsm *WordPruningBreadthFirstSearchManager) activeListAdd(token *Token) {
	wpbfsm.activeListManager.Add(token)
}

// Counts all the tokens in the active list (and displays them). This is an
// expensive operation.
func (wpbfsm *WordPruningBreadthFirstSearchManager) showTokenCount() {
	tokenSet := make([]*Token, 0)

	for _, token := range wpbfsm.activeList {
		for token != nil {
			tokenSet = append(tokenSet, token)
			token = token.GetPredecessor()
		}
	}

	fmt.Println("Token Lattice size: %d", len(tokenSet))

	tokenSet = make([]*Token, 0)

	for _, token := range wpbfsm.resultList {
		for token != nil {
			tokenSet = append(tokenSet, token)
			token = token.GetPredecessor()
		}
	}

	fmt.Println("Result Lattice size: %d", len(tokenSet))
}
