package adaptation

import "github.com/jtejido/go-sphinx/util"

const (
	MIN_FRAMES = 300 // Minimum number of frames to perform estimation
)

/**
 * This class is used for estimating a MLLR transform for each cluster of data.
 * The clustering must be previously performed using
 * ClusteredDensityFileData.java
 */
type Stats struct {
	means     *ClusteredDensityFileData
	regLs     [][][][][]float64
	regRs     [][][][]float64
	nClusters int
	loader    *tiedstate.Sphinx3Loader
	varFlor   float32
	logMath   *util.LogMath
	nFrames   int
}

func NewStats(loader tiedstate.Loader, means *ClusteredDensityFileData) *Stats {
	this := new(Stats)
	this.loader = loader.(*tiedstate.Sphinx3Loader)
	this.nClusters = means.NumberOfClusters()
	this.means = means
	this.varFlor = 1e-5
	this.invertVariances()
	this.init()
	this.nFrames = 0
	return this
}

func (s *Stats) ClusteredData() *ClusteredDensityFileData {
	return s.means
}

func (s *Stats) RegLs() [][][][][]float64 {
	return s.regLs
}

func (s *Stats) RegRs() [][][][]float64 {
	return s.regRs
}

func (s *Stats) init() {
	len := s.loader.VectorLength()[0]
	s.regLs = make([][][][][]float64, s.nClusters)
	s.regRs = make([][][][]float64, s.nClusters)

	for i := 0; i < s.nClusters; i++ {
		s.regLs[i] = make([][][][]float64, s.loader.NumStreams())
		s.regRs[i] = make([][][]float64, s.loader.NumStreams())

		for j := 0; j < s.loader.NumStreams(); j++ {
			len = s.loader.VectorLength()[j]
			s.regLs[i][j] = make([][][]float64, len)
			s.regRs[i][j] = make([][]float64, len)
			for k := 0; k < len; k++ {
				s.regLs[i][j][k] = make([][]float64, len+1)
				s.regRs[i][j][k] = make([]float64, len+1)
				for l := 0; l < len+1; l++ {
					s.regLs[i][j][k][l] = make([]float64, len+1)
				}
			}
		}
	}
}

/**
 * Used for inverting variances.
 */
func (s *Stats) invertVariances() {
	for i := 0; i < s.loader.NumStates(); i++ {
		for k := 0; k < s.loader.NumGaussiansPerState(); k++ {
			for l := 0; l < s.loader.VectorLength()[0]; l++ {
				if s.loader.VariancePool().Get(i*s.loader.NumGaussiansPerState() + k)[l] <= 0. {
					s.loader.VariancePool().Get(i*s.loader.NumGaussiansPerState() + k)[l] = 0.5
				} else if s.loader.VariancePool().Get(i*s.loader.NumGaussiansPerState() + k)[l] < s.varFlor {
					s.loader.VariancePool().Get(i*s.loader.NumGaussiansPerState() + k)[l] = (1. / s.varFlor)
				} else {
					s.loader.VariancePool().Get(i*s.loader.NumGaussiansPerState() + k)[l] = (1. / s.loader.VariancePool().Get(i*s.loader.NumGaussiansPerState() + k)[l])
				}
			}
		}
	}
}

/**
 * Fill lower part of Legetter's set of G matrices.
 */
func (s *Stats) FillRegLowerPart() {
	for i := 0; i < s.nClusters; i++ {
		for j := 0; j < s.loader.NumStreams(); j++ {
			for l := 0; l < s.loader.VectorLength()[j]; l++ {
				for p := 0; p <= s.loader.VectorLength()[j]; p++ {
					for q := p + 1; q <= s.loader.VectorLength()[j]; q++ {
						s.regLs[i][j][l][q][p] = s.regLs[i][j][l][p][q]
					}
				}
			}
		}
	}
}
