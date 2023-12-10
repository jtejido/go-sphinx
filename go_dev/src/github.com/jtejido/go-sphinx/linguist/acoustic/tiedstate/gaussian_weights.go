package tiedstate

import (
	"github.com/jtejido/go-sphinx/util"
)

/** Structure to store weights for all gaussians in AM.
 * Supposed to provide faster access in case of large models */
type GaussianWeights struct {
	weights                            [][]float32
	numStates, gauPerState, numStreams int
	name                               string
}

func NewGaussianWeights(name string, numStates, gauPerState, numStreams int) *GaussianWeights {
	ans := new(GaussianWeights)
	ans.numStates = numStates
	ans.gauPerState = gauPerState
	ans.numStreams = numStreams
	ans.name = name
	ans.weights = make([][]float32, gauPerState)
	for i := 0; i < len(ans.weights); i++ {
		ans.weights[i] = make([]float32, numStates*numStreams)
	}

	return ans
}

func (g *GaussianWeights) Put(stateID, streamID int, gauWeights []float32) {
	assert(len(gauWeights) == g.gauPerState)

	for i := 0; i < g.gauPerState; i++ {
		g.weights[i][stateID*g.numStreams+streamID] = gauWeights[i]
	}
}

func (g *GaussianWeights) Get(stateID, streamID, gaussianID int) float32 {
	return g.weights[gaussianID][stateID*g.numStreams+streamID]
}

func (g *GaussianWeights) StatesNum() int {
	return g.numStates
}

func (g *GaussianWeights) GauPerState() int {
	return g.gauPerState
}

func (g *GaussianWeights) StreamsNum() int {
	return g.numStreams
}

func (g *GaussianWeights) Name() string {
	return g.name
}

func (g *GaussianWeights) LogInfo(logger util.Logger) {
	logger.Infof("Gaussian weights: %s. Entries: %d", g.name, g.numStates*g.numStreams)
}

func (g *GaussianWeights) ConvertToPool() *Pool[float32] {
	return nil
}
