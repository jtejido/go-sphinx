package search

import (
	"github.com/jtejido/go-sphinx/utils"
)

type SimpleActiveListFactory struct {
	BaseActiveListFactory
}

func NewDefaultSimpleActiveListFactory() *SimpleActiveListFactory {
	salf := new(SimpleActiveListFactory)
	salf.absoluteBeamWidth = -1
	salf.logRelativeBeamWidth = utils.LogMath.LinearToLog(1E-80)
	return salf
}

func NewSimpleActiveListFactory(absoluteBeamWidth int, relativeBeamWidth float64) *SimpleActiveListFactory {
	salf := new(SimpleActiveListFactory)
	salf.absoluteBeamWidth = absoluteBeamWidth
	salf.logRelativeBeamWidth = utils.LogMath.LinearToLog(relativeBeamWidth)
	return salf
}

func (salf *SimpleActiveListFactory) NewInstance() *SimpleActiveList {
	return NewSimpleActiveList(salf.absoluteBeamWidth, salf.logRelativeBeamWidth)
}
