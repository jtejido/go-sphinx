package lextree

import (
	"github.com/jtejido/go-sphinx/linguist"
)

type LexTreeSearchGraph struct {
	initialState linguist.SearchState
}

func NewLexTreeSearchGraph(initialState linguist.SearchState) *LexTreeSearchGraph {
	ltsg := new(LexTreeSearchGraph)
	ltsg.initialState = initialState

	return ltsg
}

func (ltsg LexTreeSearchGraph) GetInitialState() linguist.SearchState {
	return initialState
}

func (ltsg LexTreeSearchGraph) GetNumStateOrder() int {
	return 6
}

func (ltsg LexTreeSearchGraph) GetWordTokenFirst() bool {
	return false
}
