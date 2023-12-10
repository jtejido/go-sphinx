package util

import (
	"math"
	"sync"
)

// LogMath represents a log math utility
type LogMath struct {
	naturalLogBase        float32
	inverseNaturalLogBase float32
	theAddTable           []float32
}

var (
	logBase  = 1.0001
	useTable = true
	instance *LogMath
	once     sync.Once
)

// LogMath instance initialization
func initLogMath() {
	instance = &LogMath{
		naturalLogBase:        float32(math.Log(logBase)),
		inverseNaturalLogBase: 1.0 / float32(math.Log(logBase)),
	}
	if useTable {
		// Now create the addTable table.
		// summation needed in the loop
		var innerSummation float64
		// First decide number of elements.
		const (
			veryLargeNumberOfEntries = 150000
			verySmallNumberOfEntries = 0
		)
		// To decide the size of the table, take into account that a base
		// of 1.0001 or 1.0003 converts probabilities, which are
		// numbers less than 1, into integers. Therefore, a good
		// approximation for the smallest number in the table,
		// therefore the value with the highest index, is an
		// index that maps into 0.5: indices higher than that, if
		// they were present, would map to fewer values less than
		// 0.5, therefore they would be mapped to 0 as
		// integers. Since the table implements the expression:
		//
		// log(1.0 + base^(-index)))
		//
		// then the highest index would be:
		//
		// topIndex = - log(logBase^(0.5) - 1)
		//
		// where log is the log in the appropriate base.
		//
		// Added -math.Floor(...) to round to the nearest
		// integer. Added the negation to match the preceding
		// documentation
		entriesInTheAddTable := int(-math.Floor(float64(instance.LinearToLog(instance.LogToLinear(0.5) - 1))))
		// We reach this max if the log base is 1.00007. The
		// closer you get to 1, the higher the number of entries
		// in the table.
		if entriesInTheAddTable > veryLargeNumberOfEntries {
			entriesInTheAddTable = veryLargeNumberOfEntries
		}
		if entriesInTheAddTable <= verySmallNumberOfEntries {
			panic("The log base is too close to 1.0, resulting in a very small addTable.")
		}
		// PBL added this just to see how many entries really are
		// in the table
		instance.theAddTable = make([]float32, entriesInTheAddTable)
		for index := 0; index < entriesInTheAddTable; index++ {
			// This loop implements the expression:
			//
			// log(1.0 + power(base, index))
			//
			// needed to add two numbers in the log domain.
			innerSummation = instance.LogToLinear(-float32(index))
			innerSummation += 1.0
			instance.theAddTable[index] = instance.LinearToLog(innerSummation)
		}
	}
}

// GetLogMath returns the singleton instance of LogMath
func GetLogMath() *LogMath {
	once.Do(initLogMath)
	return instance
}

// SetLogBase sets the log base
func SetLogBase(logBaseValue float64) {
	logBase = logBaseValue
}

/**
 * Converts the source, which is a number in base Math.E, to a log value which base is the LogBase of this LogMath.
 *
 * @return converted value
 * @param logSource the number in base Math.E to convert
 */
func (lm *LogMath) LnToLog(logSource float32) float32 {
	return (logSource * lm.inverseNaturalLogBase)
}

// linearToLog converts linear scale to log scale
func (lm *LogMath) LinearToLog(linear float64) float32 {
	return float32(math.Log(linear)) * lm.inverseNaturalLogBase
}

/** Converts a vector from linear domain to log domain using a given <code>LogMath</code>-instance for conversion.
 * @param vector to convert in-place
 */
func (lm *LogMath) LinearToLogFromFloats(vector []float32) {
	for i := 0; i < len(vector); i++ {
		vector[i] = lm.LinearToLog(float64(vector[i]))
	}
}

// logToLinear converts log scale to linear scale
func (lm *LogMath) LogToLinear(log float32) float64 {
	return math.Exp(float64(log * lm.naturalLogBase))
}
