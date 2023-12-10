package decoder

import (
	"github.com/jtejido/go-sphinx/decoder/search"
	"github.com/jtejido/go-sphinx/result"
	"github.com/jtejido/go-sphinx/util/props"
)

type BaseDecoder struct {
	searchManager       search.SearchManager
	resultListeners     []ResultListener
	fireNonFinalResults bool
	name                string
}

func (bd *BaseDecoder) Allocate() {
	bd.searchManager.Allocate()
}

func (bd *BaseDecoder) Deallocate() {
	bd.searchManager.Deallocate()
}

func (bd *BaseDecoder) fireResultListeners(result *result.Result) {
	if bd.fireNonFinalResults || result.IsFinal() {
		for _, resultListener := range bd.resultListeners {
			resultListener.NewResult(result)
		}
	}

	// else {
	// 	logger.finer("skipping non-final result " + result)
	// }
}

func (bd *BaseDecoder) NewProperties(ps *props.PropertySheet) error {
	init( ps.getInstanceName(), ps.getLogger(), (SearchManager) ps.getComponent(PROP_SEARCH_MANAGER), ps.getBoolean(FIRE_NON_FINAL_RESULTS), ps.getBoolean(AUTO_ALLOCATE), ps.getComponentList(PROP_RESULT_LISTENERS, ResultListener.class));
}
