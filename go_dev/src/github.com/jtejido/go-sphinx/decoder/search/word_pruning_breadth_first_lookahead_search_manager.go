package search

import (
	"fmt"
	"github.com/jtejido/go-sphinx/decoder/pruner"

	"github.com/jtejido/go-sphinx/decoder/scorer"
	"github.com/jtejido/go-sphinx/linguist"
	"github.com/jtejido/go-sphinx/linguist/acoustic/tiedstate"
	"github.com/jtejido/go-sphinx/linguist/allphone"
	"github.com/jtejido/go-sphinx/linguist/lextree"
	"github.com/jtejido/go-sphinx/utils"
)

const (
	// The property that controls size of lookahead window. Acceptable values
	// are in range [1..10].
	DEFAULT_LOOKAHEAD_WINDOW         = 5
	DEFAULT_LOOKAHEAD_PENALTY_WEIGHT = 1.0
)

type FrameCiScores struct {
	scores   []float64
	maxScore float64
}

func NewFrameCiScores(scores []float64, maxScore float64) *FrameCiScores {
	fcs := new(FrameCiScores)
	fcs.scores = scores
	fcs.maxScore = maxScore
}

// Provides the breadth first search with fast match heuristic included to
// reduce amount of tokens created.
//
// All scores and probabilities are maintained in the log math log domain.
type WordPruningBreadthFirstLookaheadSearchManager struct {
	WordPruningBreadthFirstSearchManager
	fastmatchLinguist           linguist.Liguist
	loader                      tiedstate.Loader
	fastmatchActiveListFactory  *ActiveListFactory
	lookaheadWindow             int
	lookaheadWeight             float64
	penalties                   map[int]float64
	ciScores                    []*FrameCiScores
	currentFastMatchFrameNumber int
	fastmatchActiveList         ActiveList
	fastMatchBestTokenMap       map[linguist.SearchState]*Token
	fastmatchStreamEnd          bool
}

func NewDefaultWordPruningBreadthFirstLookaheadSearchManager() *WordPruningBreadthFirstLookaheadSearchManager {
	wpbflsm := new(WordPruningBreadthFirstLookaheadSearchManager)
	wpbflsm.logMath = utils.LogMath{}
	wpbflsm.linguist = lextree.NewDefaultLexTreeLinguist()
	wpbflsm.pruner = pruner.NewDefaultSimplePruner()
	wpbflsm.scorer = scorer.NewDefaultSimpleAcousticScorer()
	wpbflsm.activeListManager = search.NewDefaultSimpleActiveListManager()
	wpbflsm.showTokenCount = DEFAULT_SHOW_TOKEN_COUNT
	wpbflsm.growSkipInterval = DEFAULT_GROW_SKIP_INTERVAL
	wpbflsm.checkStateOrder = DEFAULT_CHECK_STATE_ORDER
	wpbflsm.buildWordLattice = DEFAULT_BUILD_WORD_LATTICE
	wpbflsm.maxLatticeEdges = DEFAULT_MAX_LATTICE_EDGES
	wpbflsm.acousticLookaheadFrames = 1.7
	wpbflsm.keepAllTokens = DEFAULT_KEEP_ALL_TOKENS
	wpbflsm.relativeBeamWidth = wpbflsm.logMath.LinearToLog(1e-60)

	wpbflsm.loader = tiedstate.NewDefaultSphinx3Loader()
	wpbflsm.fastmatchLinguist = allphone.NewDefaultAllPhoneLinguist()
	wpbflsm.fastmatchActiveListFactory = NewDefaultPartitionActiveListFactory()
	wpbflsm.lookaheadWindow = DEFAULT_LOOKAHEAD_WINDOW
	wpbflsm.lookaheadWeight = 6

	if lookaheadWindow < 1 || lookaheadWindow > 10 {
		panic(fmt.Sprintf("Unsupported lookahead window size: %d. Value in range [1..10] is expected", lookaheadWindow))
	}

	wpbflsm.ciScores = make([]*FrameCiScores, 0)
	wpbflsm.penalties = make(map[int]float64)
	l, ok := wpbflsm.loader.(tiedstate.Sphinx3Loader)
	if ok && l.HasTiedMixtures() {
		l.SetGauScoresQueueLength(wpbflsm.lookaheadWindow + 2)
	}

	return wpbflsm
}

func NewWordPruningBreadthFirstLookaheadSearchManager(linguist, fastmatchLinguist linguist.Liguist, loader tiedstate.Loader,
	pruner pruner.Pruner, scorer scorer.AcousticScorer, activeListManager ActiveListManager,
	fastmatchActiveListFactory *ActiveListFactory, showTokenCount bool, relativeWordBeamWidth float64,
	growSkipInterval int, checkStateOrder bool, buildWordLattice bool, lookaheadWindow int, lookaheadWeight float64,
	maxLatticeEdges int, acousticLookaheadFrames float64, keepAllTokens bool) *WordPruningBreadthFirstLookaheadSearchManager {
	wpbflsm := new(WordPruningBreadthFirstLookaheadSearchManager)

	wpbflsm.logMath = utils.LogMath{}
	wpbflsm.linguist = linguist
	wpbflsm.pruner = pruner
	wpbflsm.scorer = scorer
	wpbflsm.activeListManager = activeListManager
	wpbflsm.showTokenCount = showTokenCount
	wpbflsm.growSkipInterval = growSkipInterval
	wpbflsm.checkStateOrder = checkStateOrder
	wpbflsm.buildWordLattice = buildWordLattice
	wpbflsm.maxLatticeEdges = maxLatticeEdges
	wpbflsm.acousticLookaheadFrames = acousticLookaheadFrames
	wpbflsm.keepAllTokens = keepAllTokens
	wpbflsm.relativeBeamWidth = wpbflsm.logMath.LinearToLog(relativeWordBeamWidth)

	wpbflsm.loader = loader
	wpbflsm.fastmatchLinguist = fastmatchLinguist
	wpbflsm.fastmatchActiveListFactory = fastmatchActiveListFactory
	wpbflsm.lookaheadWindow = lookaheadWindow
	wpbflsm.lookaheadWeight = lookaheadWeight

	if lookaheadWindow < 1 || lookaheadWindow > 10 {
		panic(fmt.Sprintf("Unsupported lookahead window size: %d. Value in range [1..10] is expected", lookaheadWindow))
	}

	wpbflsm.ciScores = make([]*FrameCiScores, 0)
	wpbflsm.penalties = make(map[int]float64)
	l, ok := wpbflsm.loader.(tiedstate.Sphinx3Loader)
	if ok && l.HasTiedMixtures() {
		l.SetGauScoresQueueLength(wpbflsm.lookaheadWindow + 2)
	}
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) Recognize(nFrames int) *result.Result {
	done := false
	result := nil
	wpbflsm.streamEnd = false

	for i := 0; i < nFrames && !done; i++ {
		if !wpbflsm.fastmatchStreamEnd {
			wpbflsm.fastMatchRecognize()
		}
		wpbflsm.penalties = make(map[int]float64)

		// remove head
		copy(wpbflsm.ciScores[i:], wpbflsm.ciScores[i+1:])
		wpbflsm.ciScores[len(wpbflsm.ciScores)-1] = nil
		wpbflsm.ciScores = wpbflsm.ciScores[:len(wpbflsm.ciScores)-1]

		done = wpbflsm.recognize()
	}

	if !streamEnd {
		result = result.NewResult(wpbflsm.loserManager, wpbflsm.activeList, wpbflsm.resultList, wpbflsm.currentCollectTime, done, wpbflsm.linguist.GetSearchGraph().GetWordTokenFirst(), true)
	}

	if wpbflsm.showTokenCount {
		wpbflsm.showTokenCount()
	}

	return result
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) fastMatchRecognize() {
	more := wpbflsm.scoreFastMatchTokens()

	if more {
		wpbflsm.pruneFastMatchBranches()
		wpbflsm.currentFastMatchFrameNumber++
		wpbflsm.createFastMatchBestTokenMap()
		wpbflsm.growFastmatchBranches()
	}
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) createFastMatchBestTokenMap() {
	mapSize := wpbflsm.fastmatchActiveList.Size() * 10
	if mapSize == 0 {
		mapSize = 1
	}
	wpbflsm.fastMatchBestTokenMap = make(map[linguist.SearchState]*Token, mapSize)
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) localStart() {
	wpbflsm.currentFastMatchFrameNumber = 0
	l, ok := wpbflsm.loader.(tiedstate.Sphinx3Loader)
	if ok && l.HasTiedMixtures() {
		l.ClearGauScores()
	}
	// prepare fast match active list
	wpbflsm.fastmatchActiveList = wpbflsm.fastmatchActiveListFactory.newInstance() //CHECK THIS
	fmInitState := wpbflsm.fastmatchLinguist.GetSearchGraph().GetInitialState()
	wpbflsm.fastmatchActiveList.Add(NewToken(fmInitState, wpbflsm.currentFastMatchFrameNumber))
	wpbflsm.createFastMatchBestTokenMap()
	wpbflsm.growFastmatchBranches()
	wpbflsm.fastmatchStreamEnd = false
	for i := 0; (i < wpbflsm.lookaheadWindow-1) && !wpbflsm.fastmatchStreamEnd; i++ {
		wpbflsm.fastMatchRecognize()
	}

	// CHECK FROM EMBEDDED
	// super.localStart()

	searchGraph := wpbflsm.linguist.GetSearchGraph()
	wpbflsm.currentFrameNumber = 0
	wpbflsm.curTokensScored.Value = 0
	wpbflsm.numStateOrder = searchGraph.GetNumStateOrder()
	wpbflsm.activeListManager.SetNumStateOrder(wpbflsm.numStateOrder)
	if buildWordLattice {
		wpbflsm.loserManager = NewAlternateHypothesisManager(wpbflsm.maxLatticeEdges)
	}

	state := searchGraph.GetInitialState()

	wpbflsm.activeList = wpbflsm.activeListManager.GetEmittingList()
	wpbflsm.activeList.Add(NewToken(state, -1))

	wpbflsm.clearCollectors()

	wpbflsm.growBranches()
	wpbflsm.growNonEmittingBranches()
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) growFastmatchBranches() {
	wpbflsm.growTimer.Start()
	oldActiveList := wpbflsm.fastmatchActiveList
	wpbflsm.fastmatchActiveList = wpbflsm.fastmatchActiveListFactory.newInstance()
	fastmathThreshold := oldActiveList.GetBeamThreshold()
	// TODO more precise range of baseIds, remove magic number
	frameCiScores := make([]float64, 1024)

	frameMaxCiScore := -math.MaxFloat64
	for _, token := range oldActiveList {
		tokenScore := token.GetScore()
		if tokenScore < fastmathThreshold {
			continue
		}
		// filling max ci scores array that will be used in general search
		// token score composing
		t, ok := token.GetSearchState().(allphone.PhoneHmmSearchState)
		if ok {
			baseId := t.GetBaseId()
			if frameCiScores[baseId] < tokenScore {
				frameCiScores[baseId] = tokenScore
			}
			if frameMaxCiScore < tokenScore {
				frameMaxCiScore = tokenScore
			}
		}
		wpbflsm.collectFastMatchSuccessorTokens(token)
	}
	wpbflsm.ciScores = append(wpbflsm.ciScores, NewFrameCiScores(frameCiScores, frameMaxCiScore))
	wpbflsm.growTimer.Stop()
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) scoreFastMatchTokens() bool {
	var moreTokens bool
	wpbflsm.scoreTimer.Start()
	data := wpbflsm.scorer.CalculateScoresAndStoreData(wpbflsm.fastmatchActiveList.GetTokens())
	wpbflsm.scoreTimer.Stop()

	bestToken := nil
	d, ok := data.(*Token)
	if ok {
		bestToken = d
	} else {
		wpbflsm.fastmatchStreamEnd = true
	}

	if bestToken != nil {
		moreTokens = true
	}

	wpbflsm.fastmatchActiveList.SetBestToken(bestToken)

	// monitorWords(activeList);
	wpbflsm.monitorStates(wpbflsm.fastmatchActiveList)

	// System.out.println("BEST " + bestToken);

	wpbflsm.curTokensScored.value += wpbflsm.fastmatchActiveList.Size()
	wpbflsm.totalTokensScored.value += wpbflsm.fastmatchActiveList.Size()

	return moreTokens
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) pruneFastMatchBranches() {
	wpbflsm.pruneTimer.Start()
	wpbflsm.fastmatchActiveList = wpbflsm.pruner.Prune(wpbflsm.fastmatchActiveList)
	wpbflsm.pruneTimer.Stop()
}

// CLEAN THIS USE OF MAP SHIT, CAN'T MAKE INTERFACE AS KEY STUPID
func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) getFastMatchBestToken(state linguist.SearchState) *Token {
	return wpbflsm.fastMatchBestTokenMap[state]
}

// CLEAN THIS USE OF MAP SHIT, CAN'T MAKE INTERFACE AS KEY STUPID
func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) setFastMatchBestToken(token *Token, state linguist.SearchState) {
	wpbflsm.fastMatchBestTokenMap[state] = token
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) collectFastMatchSuccessorTokens(token *Token) {
	state := token.GetSearchState()
	arcs := state.GetSuccessors()
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
		// We're actually multiplying the variables, but since
		// these come in log(), multiply gets converted to add
		logEntryScore := token.GetScore() + arc.GetProbability()
		predecessor := wpbflsm.getResultListPredecessor(token)

		// if not emitting, check to see if we've already visited
		// this state during this frame. Expand the token only if we
		// haven't visited it already. This prevents the search
		// from getting stuck in a loop of states with no
		// intervening emitting nodes. This can happen with nasty
		// jsgf grammars such as ((foo*)*)*
		if !nextState.IsEmitting() {
			newTok := newToken(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbflsm.currentFastMatchFrameNumber)
			wpbflsm.tokensCreated.value++
			if !wpbflsm.isVisited(newTok) {
				wpbflsm.collectFastMatchSuccessorTokens(newTok)
			}
			continue
		}

		bestToken := wpbflsm.getFastMatchBestToken(nextState)
		if bestToken == nil {
			newTok := newToken(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbflsm.currentFastMatchFrameNumber)
			wpbflsm.tokensCreated.value++
			wpbflsm.setFastMatchBestToken(newTok, nextState)
			wpbflsm.fastmatchActiveList.Add(newTok)
		} else {
			if bestToken.GetScore() <= logEntryScore {
				bestToken.Update(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbflsm.currentFastMatchFrameNumber)
			}
		}
	}
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) collectSuccessorTokens(token *Token) {
	// If this is a final state, add it to the final list

	if token.IsFinal() {
		wpbflsm.resultList.Add(wpbflsm.getResultListPredecessor(token))
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

	if !token.IsEmitting() && (wpbflsm.keepAllTokens && wpbflsm.isVisited(token)) {
		return
	}

	state := token.GetSearchState()
	arcs := state.GetSuccessors()
	predecessor := wpbflsm.getResultListPredecessor(token)

	// For each successor
	// calculate the entry score for the token based upon the
	// predecessor token score and the transition probabilities
	// if the score is better than the best score encountered for
	// the SearchState and frame then create a new token, add
	// it to the lattice and the SearchState.
	// If the token is an emitting token add it to the list,
	// otherwise recursively collect the new tokens successors.

	tokenScore := token.GetScore()
	beamThreshold := wpbflsm.activeList.GetBeamThreshold()
	stateProducesPhoneHmms := false

	switch v := state.(type) {
	case lextree.LexTreeNonEmittingHMMState, lextree.LexTreeWordState, lextree.LexTreeEndUnitState:
		stateProducesPhoneHmms = true
	default:
		stateProducesPhoneHmms = false
	}
	for _, arc := range arcs {
		nextState := arc.GetState()

		// prune states using lookahead heuristics
		if stateProducesPhoneHmms {
			lt, ok := nextState.LexTreeHMMState
			if ok {
				penalty := 0.
				baseId := lt.GetHMMState().GetHMM().GetBaseUnit().GetBaseID()
				if wpbflsm.penalties.Get(baseId) == nil {
					penalty = wpbflsm.updateLookaheadPenalty(baseId)
				}
				if (tokenScore + wpbflsm.lookaheadWeight*penalty) < beamThreshold {
					continue
				}
			}
		}

		if checkStateOrder {
			wpbflsm.checkStateOrder(state, nextState)
		}

		// We're actually multiplying the variables, but since
		// these come in log(), multiply gets converted to add
		logEntryScore := tokenScore + arc.GetProbability()

		bestToken := wpbflsm.getBestToken(nextState)

		_, wss_ok := nextState.(WordSearchState)

		if bestToken == nil {
			newBestToken = newToken(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbflsm.currentCollectTime)
			wpbflsm.tokensCreated.value++
			wpbflsm.setBestToken(newBestToken, nextState)
			wpbflsm.activeListAdd(newBestToken)
		} else if bestToken.GetScore() < logEntryScore {
			// System.out.println("Updating " + bestToken + " with " +
			// newBestToken);
			oldPredecessor := bestToken.GetPredecessor()
			bestToken.Update(predecessor, nextState, logEntryScore, arc.GetInsertionProbability(), arc.GetLanguageProbability(), wpbflsm.currentCollectTime)

			if buildWordLattice && wss_ok {
				wpbfsm.loserManager.AddAlternatePredecessor(bestToken, oldPredecessor)
			}
		} else if buildWordLattice && wss_ok {
			if predecessor != nil {
				wpbfsm.loserManager.AddAlternatePredecessor(bestToken, predecessor)
			}
		}
	}
}

func (wpbflsm *WordPruningBreadthFirstLookaheadSearchManager) updateLookaheadPenalty(baseId int) float64 {
	if len(wpbflsm.ciScores) <= 0 {
		return 0.0
	}

	penalty := -math.MaxFloat64
	for _, frameCiScores := range wpbflsm.ciScores {
		diff := frameCiScores.scores[baseId] - frameCiScores.maxScore
		if diff > penalty {
			penalty = diff
		}
	}

	wpbflsm.penalties[baseId] = penalty
	return penalty
}
