package tiedstate

import (
	"github.com/jtejido/go-sphinx/linguist/acoustic"
	"github.com/jtejido/go-sphinx/util"
)

/**
 * Manages HMMs. This HMMManager groups {@link edu.cmu.sphinx.linguist.acoustic.HMM HMMs} together by their {@link
 * edu.cmu.sphinx.linguist.acoustic.HMMPosition position} with the word.
 */
type HMMManager struct {
	allHMMs         []acoustic.HMM
	hmmsPerPosition map[acoustic.HMMPosition]map[*acoustic.Unit]acoustic.HMM
}

func NewHMMManager() *HMMManager {
	ans := new(HMMManager)
	ans.allHMMs = make([]acoustic.HMM, 0)
	ans.hmmsPerPosition = make(map[acoustic.HMMPosition]map[*acoustic.Unit]acoustic.HMM)
	for _, pos := range acoustic.Values() {
		ans.hmmsPerPosition[pos] = make(map[*acoustic.Unit]acoustic.HMM)
	}

	return ans
}

/**
 * Put an HMM into this manager
 *
 * @param hmm the hmm to manage
 */
func (m *HMMManager) Put(hmm acoustic.HMM) {
	m.hmmsPerPosition[hmm.Position()][hmm.Unit()] = hmm
	m.allHMMs = append(m.allHMMs, hmm)
}

/**
 * Retrieves an HMM by position and unit
 *
 * @param position the position of the HMM
 * @param unit     the unit that this HMM represents
 * @return the HMM for the unit at the given position or null if no HMM at the position could be found
 */
func (m *HMMManager) Get(position acoustic.HMMPosition, unit *acoustic.Unit) acoustic.HMM {
	return m.hmmsPerPosition[position][unit]
}

/**
 * Gets an iterator that iterates through all HMMs
 *
 * @return an iterator that iterates through all HMMs
 */
func (m *HMMManager) Iterator() util.Iterator[acoustic.HMM] {
	return util.NewIterator[acoustic.HMM](m.allHMMs)
}

/**
 * Returns the number of HMMS in this manager
 *
 * @return the number of HMMs
 */
func (m *HMMManager) NumHMMs() int {
	var count int

	for _, v := range m.hmmsPerPosition {
		if v != nil {
			count += len(v)
		}
	}
	return count
}

/**
 * Log information about this manager
 *
 * @param logger logger to use for this logInfo
 */
func (m *HMMManager) LogInfo(logger util.Logger) {
	logger.Infof("HMM Manager: %d hmms", m.NumHMMs())
}
