package search

type AlternateHypothesisManager struct {
	viterbiLoserMap map[*Token][]*Token
	maxEdges        int
}
