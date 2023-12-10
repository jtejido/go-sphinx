package decoder

import (
	"github.com/jtejido/go-sphinx/result"
	"github.com/jtejido/go-sphinx/util/props"
)

/** The listener interface for being informed when new results are generated. */
type ResultListener interface {
	props.Configurable
	/**
	 * Method called when a new result is generated
	 *
	 * @param result the new result
	 */
	NewResult(*result.Result)
}
