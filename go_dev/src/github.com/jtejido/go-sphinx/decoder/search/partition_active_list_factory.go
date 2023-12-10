package search

import (
	"github.com/jtejido/go-sphinx/utils"
)

type PartitionActiveListFactory struct {
	BaseActiveListFactory
}

func NewDefaultPartitionActiveListFactory() *PartitionActiveListFactory {
	palf := new(PartitionActiveListFactory)
	palf.absoluteBeamWidth = 20000
	palf.logRelativeBeamWidth = utils.LogMath.LinearToLog(1e-60)
	return palf
}

func NewPartitionActiveListFactory(absoluteBeamWidth int, relativeBeamWidth float64) *PartitionActiveListFactory {
	palf := new(PartitionActiveListFactory)
	palf.absoluteBeamWidth = absoluteBeamWidth
	palf.logRelativeBeamWidth = utils.LogMath.LinearToLog(relativeBeamWidth)
	return palf
}

func (palf *PartitionActiveListFactory) NewInstance() *PartitionActiveList {
	return NewPartitionActiveList(palf.absoluteBeamWidth, palf.logRelativeBeamWidth)
}
