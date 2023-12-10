package props

import (
	"github.com/jtejido/go-sphinx/util"
)

/**
 * Wraps annotations
 *
 */
type S4PropWrapper struct {
	annotation util.Annotation
}

func NewS4PropWrapper(annotation util.Annotation) *S4PropWrapper {
	return &S4PropWrapper{
		annotation: annotation,
	}
}

func (w *S4PropWrapper) Annotation() util.Annotation {
	return w.annotation
}
