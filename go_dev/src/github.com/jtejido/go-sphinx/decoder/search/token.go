package search

import (
	"fmt"
	"github.com/jtejido/go-sphinx/decoder/scorer"
	"github.com/jtejido/go-sphinx/frontend"
	"github.com/jtejido/go-sphinx/linguist"
	"github.com/jtejido/go-sphinx/linguist/dictionary"
)

// Represents a single state in the recognition trellis. Subclasses of a token are used to represent the various
// emitting state.
//
// All scores are maintained in LogMath log base
type Token struct {
	curCount, lastCount                                                  int
	predecessor                                                          *Token
	logLanguageScore, logTotalScore, logInsertionScore, logAcousticScore float64
	searchState                                                          linguist.SearchState
	collectTime                                                          int64
	data                                                                 frontend.Data
}

// Internal constructor for a token. Used by classes Token, CombineToken, ParallelToken
func newToken(predecessor *Token, state linguist.SearchState, logTotalScore, logInsertionScore, logLanguageScore float64, collectTime int64) *Token {
	token := new(Token)
	token.predecessor = predecessor
	token.searchState = state
	token.logTotalScore = logTotalScore
	token.logInsertionScore = logInsertionScore
	token.logLanguageScore = logLanguageScore
	token.collectTime = collectTime
	token.curCount++

	return token
}

// Creates the initial token with the given word history depth
func NewToken(state linguist.SearchState, collectTime int64) *Token {
	return newToken(nil, state, 0.0, 0.0, 0.0, collectTime)
}

// Creates a Token with the given acoustic and language scores and predecessor
func NewTokenWithScores(predecessor *Token, logTotalScore, logAcousticScore, logInsertionScore, logLanguageScore float64) *Token {
	token := newToken(predecessor, nil, logTotalScore, logInsertionScore, logLanguageScore, 0)
	token.logAcousticScore = logAcousticScore
	return token
}

// Returns the predecessor for this token, or null if this token has no predecessors
func (tok *Token) GetPredecessor() *Token {
	return tok.predecessor
}

// Collect time is different from frame number because some frames might be skipped in silence detector
func (tok *Token) GetCollectTime() int64 {
	return tok.collectTime
}

// Sets the feature for this Token.
func (tok *Token) SetData(data frontend.Data) {
	tok.data = data
	fd, ok := data.(frontend.FloatData)

	if ok {
		tok.collectTime = fd.GetCollectTime()
	}
}

// Returns the feature for this Token.
func (tok *Token) GetData() frontend.Data {
	return tok.data
}

// Returns the score for the token. The score is a combination of language and acoustic scores
func (tok *Token) GetScore() float64 {
	return tok.logTotalScore
}

// Calculates a score against the given feature. The score can be retrieved
// with get score. The token will keep a reference to the scored feature-vector.
func (tok *Token) CalculateScore(feature frontend.Data) float64 {

	tok.logAcousticScore = tok.searchState.(scorer.ScoreProvider).GetScore(feature)

	tok.logTotalScore += tok.logAcousticScore

	tok.SetData(feature)

	return tok.logTotalScore
}

func (tok *Token) CalculateComponentScore(feature frontend.Data) []float64 {
	return searchState.(scorer.ScoreProvider).GetComponentScore(feature)
}

// Normalizes a previously calculated score
func (tok *Token) NormalizeScore(maxLogScore float64) float64 {
	tok.logTotalScore -= maxLogScore
	tok.logAcousticScore -= maxLogScore
	return tok.logTotalScore
}

// Sets the score for this token
func (tok *Token) SetScore(logScore float64) {
	tok.logTotalScore = logScore
}

// Returns the language score associated with this token
func (tok *Token) GetLanguageScore() float64 {
	return tok.logLanguageScore
}

// Returns the insertion score associated with this token.
// Insertion score is the score of the transition between
// states. It might be transition score from the acoustic model,
// phone insertion score or word insertion probability from
// the linguist.
func (tok *Token) GetInsertionScore() float64 {
	return tok.logInsertionScore
}

// Returns the acoustic score for this token (in logMath log base).
// Acoustic score is a sum of frame GMM.
func (tok *Token) GetAcousticScore() float64 {
	return tok.logAcousticScore
}

// Returns the SearchState associated with this token
func (tok *Token) GetSearchState() linguist.SearchState {
	return tok.searchState
}

// Determines if this token is associated with an emitting state. An emitting state is a state that can be scored
// acoustically.
func (tok *Token) IsEmitting() bool {
	return tok.searchState.IsEmitting()
}

// Determines if this token is associated with a final SentenceHMM state.
func (tok *Token) IsFinal() bool {
	return tok.searchState.IsFinal()
}

// Determines if this token marks the end of a word
func (tok *Token) IsWord() bool {
	_, ok := tok.searchState.(linguist.WordSearchState)
	return ok
}

// Retrieves the string representation of this object
func (tok *Token) String() string {
	return fmt.Sprintf("%d %.7f %.7f %.7f %s", tok.GetCollectTime(), tok.GetScore(), tok.GetAcousticScore(), tok.GetLanguageScore(), tok.GetSearchState())

}

// dumps a branch of tokens
// true include all sentence hmm states
func (tok *Token) DumpTokenPath(includeHMMStates bool) {
	token := tok
	list := make([]*Token, 0)

	for token != nil {
		list = append(list, token)
		token = token.GetPredecessor()
	}

	for i := len(list) - 1; i >= 0; i-- {
		token = list[i]
		_, ok := token.GetSearchState().(linguist.HMMSearchState)
		if includeHMMStates || (!ok) {
			fmt.Println("   %s", token)
		}
	}
	fmt.Println()
}

// Returns the string of words leading up to this token.
func (tok *Token) GetWordPath(wantFiller, wantPronunciations bool) string {
	var sb string
	token := tok

	for token != nil {
		if token.IsWord() {
			wordState := token.GetSearchState().(linguist.WordSearchState)
			pron := wordState.GetPronunciation()
			word := wordState.GetPronunciation().GetWord()

			if wantFiller || !word.IsFiller() {
				if wantPronunciations {
					sb += "]"
					u = pron.GetUnits()
					for i := len(u) - 1; i >= 0; i-- {
						if i < len(u)-1 {
							sb += ","
						}
						sb += u[i].GetName()
					}
					sb += "["
				}
				sb += word.GetSpelling()
				sb += " "
			}
		}
		token = token.GetPredecessor()
	}

	return sb
}

// Returns the string of words for this token, with no embedded filler words
func (tok *Token) GetWordPathNoFiller() string {
	return tok.GetWordPath(false, false)
}

// Returns the string of words for this token, with embedded silences
func (tok *Token) GetWordPath() string {
	return tok.GetWordPath(true, false)
}

// Returns the string of words and units for this token, with embedded silences.
func (tok *Token) GetWordUnitPath() string {
	var sb string
	token := tok

	for token != nil {
		searchState := token.GetSearchState()
		wss, ok_wss := searchState.(linguist.WordSearchState)
		uss, ok_uss := searchState.(linguist.UnitSearchState)
		if ok_wss {
			word := wss.GetPronunciation().GetWord()
			sb += " " + word.GetSpelling()
		} else if ok_uss {
			unit := uss.GetUnit()
			sb += " " + unit.GetName()
		}
		token = token.GetPredecessor()
	}
	return sb
}

// Returns the word of this Token, the search state is a WordSearchState. If the search state is not a
// WordSearchState, return null.
func (tok *Token) GetWord() *dictionary.Word {
	if tok.IsWord() {
		wordState := searchState.(linguist.WordSearchState)
		return wordState.GetPronunciation().GetWord()
	}

	return nil
}

// Shows the token count
func (tok *Token) ShowCount() {
	fmt.Println("Cur count: %d new %d", tok.curCount, (tok.curCount - tok.lastCount))
	tok.lastCount = tok.curCount
}

// Determines if this branch is valid
func (tok *Token) Validate() bool {
	return true
}

func (tok *Token) Update(predecessor *Token, nextState liguist.SearchState, logEntryScore, insertionProbability, languageProbability float64, collectTime int64) {
	tok.predecessor = predecessor
	tok.searchState = nextState
	tok.logTotalScore = logEntryScore
	tok.logInsertionScore = insertionProbability
	tok.logLanguageScore = languageProbability
	tok.collectTime = collectTime
}
