package result

import (
	"github.com/jtejido/go-sphinx/decoder/search"
	"github.com/jtejido/go-sphinx/util"
)

/**
 * Provides recognition results. Results can be partial or final. A result
 * should not be modified before it is a final result. Note that a result may
 * not contain all possible information.
 * <p>
 * The following methods are not yet defined but should be:
 *
 * <pre>
 * public Result getDAG(int compressionLevel);
 * </pre>
 */
type Result struct {
	activeList                 search.ActiveList
	resultList                 []*search.Token
	alternateHypothesisManager *search.AlternateHypothesisManager
	isFinal                    bool
	wordTokenFirst             bool
	currentCollectTime         int64
	reference                  string
	logMath                    *util.LogMath
	toCreateLattice            bool
}
