package frontend

/**
 * A Data object that holds data of primitive type float.
 *
 * @see Data
 */
type FloatData struct {
	values            []float32
	sampleRate        int
	firstSampleNumber int64
	collectTime       int64
}

/**
 * Constructs a Data object with the given values, sample rate, collect time, and first sample number.
 *
 * @param values            the data values
 * @param sampleRate        the sample rate of the data
 * @param firstSampleNumber the position of the first sample in the original data
 */
func NewFloatData(values []float32, sampleRate int, firstSampleNumber int64) *FloatData {
	return NewFloatDataWithCollectTime(values, sampleRate, firstSampleNumber*1000/int64(sampleRate), firstSampleNumber)
}

/**
 * Constructs a Data object with the given values, sample rate, collect time, and first sample number.
 *
 * @param values            the data values
 * @param sampleRate        the sample rate of the data
 * @param collectTime       the time at which this data is collected
 * @param firstSampleNumber the position of the first sample in the original data
 */
func NewFloatDataWithCollectTime(values []float32, sampleRate int,
	collectTime, firstSampleNumber int64) *FloatData {
	this := new(FloatData)
	this.values = values
	this.sampleRate = sampleRate
	this.collectTime = collectTime
	this.firstSampleNumber = firstSampleNumber
	return this
}

/**
 * @return the values of this data.
 */
func (d *FloatData) Values() []float32 {
	return d.values
}
