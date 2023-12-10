package search

const (
	// The property that specifies whether to build a word lattice.
	DEFAULT_BUILD_WORD_LATTICE = true
	// The property that controls whether or not we keep all tokens. If this is
	// set to false, only word tokens are retained, otherwise all tokens are
	// retained.
	DEFAULT_KEEP_ALL_TOKENS = false
)

type TokenSearchManager struct {
	buildWordLattice, keepAllTokens bool
}

// Find the token to use as a predecessor in resultList given a candidate
// predecessor. There are three cases here:
//
// 1. We want to store everything in resultList. In that case
// keepAllTokens is set to true and we just store everything that
// was built before.
//
// 2. We are only interested in sequence of words. In this case we just
// keep word tokens and ignore everything else. In this case timing and
// scoring information is lost since we keep scores in emitting tokens.
//
// 3. We want to keep words but we want to keep scores to build a lattice
// from the result list later and buildWordLattice is set to true.
// In this case we want to insert intermediate token to store the score and
// this token will be used during lattice path collapse to get score on
// edge. See Lattice for details of resultList
// compression.
func (tsm *TokenSearchManager) getResultListPredecessor(token *Token) *Token {

	if tsm.keepAllTokens {
		return token
	}

	if !tsm.buildWordLattice {
		if token.IsWord() {
			return token
		} else {
			return token.GetPredecessor()
		}
	}

	logAcousticScore := 0.0
	logLanguageScore := 0.0
	logInsertionScore := 0.0

	for token != nil && !token.IsWord() {
		logAcousticScore += token.GetAcousticScore()
		logLanguageScore += token.GetLanguageScore()
		logInsertionScore += token.GetInsertionScore()
		token = token.GetPredecessor()
	}

	return NewToken(token, token.GetScore(), logInsertionScore, logAcousticScore, logLanguageScore)
}
