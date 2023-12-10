package pruner

import (
	"github.com/jtejido/go-sphinx/decoder/search"
)

// Performs the default pruning behavior which is to invoke the purge on the active list
type SimplePruner struct {
	name string
}

func NewDefaultSimplePruner() *SimplePruner {
	return new(SimplePruner)
}

func (sp *SimplePruner) StartRecognition() {}

func (sp *SimplePruner) Prune(activeList search.ActiveList) search.ActiveList {
	return activeList.Purge()
}

func (sp *SimplePruner) StopRecognition() {}

func (sp *SimplePruner) Allocate() {}

func (sp *SimplePruner) Deallocate() {}
