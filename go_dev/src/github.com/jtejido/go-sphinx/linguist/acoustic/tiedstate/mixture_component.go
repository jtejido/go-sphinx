package tiedstate

import (
	"fmt"
	"math"

	"github.com/jtejido/go-sphinx/frontend"
	"github.com/jtejido/go-sphinx/util"
)

const (
	DEFAULT_VAR_FLOOR  float32 = 0.0001 // this also seems to be the default of SphinxTrain
	DEFAULT_DIST_FLOOR float32 = 0.0
)

/**
 * Defines the set of shared elements for a GaussianMixture. Since these elements are potentially
 * shared by a number of {@link GaussianMixture GaussianMixtures}, these elements should not be
 * written to. The GaussianMixture defines a single probability density function along with a set of
 * adaptation parameters.
 * <p>
 * Note that all scores and weights are in LogMath log base
 */
// TODO: Since many of the subcomponents of a MixtureComponent are shared, are
// there some potential opportunities to reduce the number of computations in scoring
// senones by sharing intermediate results for these subcomponents?
type MixtureComponent struct {
	mean []float32
	/** Mean after transformed by the adaptation parameters. */
	meanTransformed          []float32
	meanTransformationMatrix [][]float32
	meanTransformationVector []float32
	variance                 []float32
	/** Precision is the inverse of the variance. This includes adaptation. */
	precisionTransformed         []float32
	varianceTransformationMatrix [][]float32
	varianceTransformationVector []float32

	distFloor     float32
	varianceFloor float32

	logPreComputedGaussianFactor float32
}

/**
 * Create a MixtureComponent with the given sub components.
 *
 * @param mean     the mean vector for this PDF
 * @param variance the variance for this PDF
 */
func NewMixtureComponentFromMeanVar(mean, variance []float32) *MixtureComponent {
	return NewMixtureComponent(mean, nil, nil, variance, nil, nil, DEFAULT_DIST_FLOOR, DEFAULT_VAR_FLOOR)
}

/**
 * Create a MixtureComponent with the given sub components.
 *
 * @param mean                         the mean vector for this PDF
 * @param meanTransformationMatrix     transformation matrix for this pdf
 * @param meanTransformationVector     transform vector for this PDF
 * @param variance                     the variance for this PDF
 * @param varianceTransformationMatrix var. transform matrix for this PDF
 * @param varianceTransformationVector var. transform vector for this PDF
 */
func NewDefaultMixtureComponent(
	mean []float32,
	meanTransformationMatrix [][]float32,
	meanTransformationVector []float32,
	variance []float32,
	varianceTransformationMatrix [][]float32,
	varianceTransformationVector []float32) *MixtureComponent {
	return NewMixtureComponent(mean, meanTransformationMatrix, meanTransformationVector, variance,
		varianceTransformationMatrix, varianceTransformationVector, DEFAULT_DIST_FLOOR, DEFAULT_VAR_FLOOR)
}

/**
 * Create a MixtureComponent with the given sub components.
 *
 * @param mean                         the mean vector for this PDF
 * @param meanTransformationMatrix     transformation matrix for this pdf
 * @param meanTransformationVector     transform vector for this PDF
 * @param variance                     the variance for this PDF
 * @param varianceTransformationMatrix var. transform matrix for this PDF
 * @param varianceTransformationVector var. transform vector for this PDF
 * @param distFloor                    the lowest score value (in linear domain)
 * @param varianceFloor                the lowest value for the variance
 */
func NewMixtureComponent(
	mean []float32,
	meanTransformationMatrix [][]float32,
	meanTransformationVector []float32,
	variance []float32,
	varianceTransformationMatrix [][]float32,
	varianceTransformationVector []float32,
	distFloor float32,
	varianceFloor float32) *MixtureComponent {

	assert(len(variance) == len(mean))
	this := new(MixtureComponent)
	this.mean = mean
	this.meanTransformationMatrix = meanTransformationMatrix
	this.meanTransformationVector = meanTransformationVector
	this.variance = variance
	this.varianceTransformationMatrix = varianceTransformationMatrix
	this.varianceTransformationVector = varianceTransformationVector

	assert2(distFloor >= 0.0, "distFloot seems to be already in log-domain")
	this.distFloor = util.GetLogMath().LinearToLog(float64(distFloor))
	this.varianceFloor = varianceFloor

	this.TransformStats()

	this.logPreComputedGaussianFactor = this.PrecomputeDistance()
	return this
}

/**
 * Returns the mean for this component.
 *
 * @return the mean
 */
func (m *MixtureComponent) Mean() []float32 {
	return m.mean
}

/**
 * Returns the variance for this component.
 *
 * @return the variance
 */
func (m *MixtureComponent) Variance() []float32 {
	return m.variance
}

/**
 * Calculate the score for this mixture against the given feature.
 * <p>
 * Note: The support of <code>DoubleData</code>-features would require an array conversion to
 * float[]. Because getScore might be invoked with very high frequency, features are restricted
 * to be <code>FloatData</code>s.
 *
 * @param feature the feature to score
 * @return the score, in log, for the given feature
 */
func (m *MixtureComponent) Score(feature *frontend.FloatData) float32 {
	return m.ScoreFromValues(feature.Values())
}

/**
 * Calculate the score for this mixture against the given feature. We model the output
 * distributions using a mixture of Gaussians, therefore the current implementation is simply
 * the computation of a multi-dimensional Gaussian. <p> <b>Normal(x) = exp{-0.5 * (x-m)' *
 * inv(Var) * (x-m)} / {sqrt((2 * PI) ^ N) * det(Var))}</b></p>
 * <p>
 * where <b>x</b> and <b>m</b> are the incoming cepstra and mean vector respectively,
 * <b>Var</b> is the Covariance matrix, <b>det()</b> is the determinant of a matrix,
 * <b>inv()</b> is its inverse, <b>exp</b> is the exponential operator, <b>x'</b> is the
 * transposed vector of <b>x</b> and <b>N</b> is the dimension of the vectors <b>x</b> and
 * <b>m</b>.
 *
 * @param feature the feature to score
 * @return the score, in log, for the given feature
 */
func (m *MixtureComponent) ScoreFromValues(feature []float32) float32 {
	logDval := m.logPreComputedGaussianFactor

	// First, compute the argument of the exponential function in
	// the definition of the Gaussian, then convert it to the
	// appropriate base. If the log base is <code>Math.E</code>,
	// then no operation is necessary.

	for i := 0; i < len(feature); i++ {
		logDiff := feature[i] - m.meanTransformed[i]
		logDval += logDiff * logDiff * m.precisionTransformed[i]
	}
	// logDval = -logVal / 2;

	// At this point, we have the ln() of what we need, that is,
	// the argument of the exponential in the javadoc comment.

	// Convert to the appropriate base.
	logDval = util.GetLogMath().LnToLog(logDval)

	// System.out.println("MC: getscore " + logDval);

	// TODO: Need to use mean and variance transforms here

	if math.IsNaN(float64(logDval)) {
		fmt.Println("gs is Nan, converting to 0")
		logDval = float32(math.Inf(-1))
	}

	if logDval < m.distFloor {
		logDval = m.distFloor
	}

	return logDval
}

/**
 * Pre-compute factors for the Mahalanobis distance. Some of the Mahalanobis distance
 * computation can be carried out in advance. Specifically, the factor containing only variance
 * in the Gaussian can be computed in advance, keeping in mind that the the determinant of the
 * covariance matrix, for the degenerate case of a mixture with independent components - only
 * the diagonal elements are non-zero - is simply the product of the diagonal elements. <p>
 * We're computing the expression:
 * <pre>{sqrt((2 * PI) ^ N) * det(Var))}</pre>
 *
 * @return the precomputed distance
 */
func (m *MixtureComponent) PrecomputeDistance() float32 {
	var logPreComputedGaussianFactor float64 // = log(1.0)
	// Compute the product of the elements in the Covariance
	// matrix's main diagonal. Covariance matrix is assumed
	// diagonal - independent dimensions. In log, the product
	// becomes a summation.
	for i := 0; i < len(m.variance); i++ {
		logPreComputedGaussianFactor += math.Log(float64(m.precisionTransformed[i]) * -2)
		//	     variance[i] = 1.0f / (variance[i] * 2.0f);
	}

	// We need the minus sign since we computed
	// logPreComputedGaussianFactor based on precision, which is
	// the inverse of the variance. Therefore, in the log domain,
	// the two quantities have opposite signs.

	// The covariance matrix's dimension becomes a multiplicative
	// factor in log scale.
	logPreComputedGaussianFactor = math.Log(2.0*math.Pi)*float64(len(m.variance)) - logPreComputedGaussianFactor

	// The sqrt above is a 0.5 multiplicative factor in log scale.
	return float32(-logPreComputedGaussianFactor) * 0.5
}

/** Applies transformations to means and variances. */
func (m *MixtureComponent) TransformStats() {
	featDim := len(m.mean)
	/*
	 * The transformed mean vector is given by:
	 *
	 * <p><b>M = A * m + B</b></p>
	 *
	 * where <b>M</b> and <b>m</b> are the mean vector after and
	 * before transformation, respectively, and <b>A</b> and
	 * <b>B</b> are the transformation matrix and vector,
	 * respectively.
	 *
	 * if A or B are <code>null</code> the according substeps are skipped
	 */
	if m.meanTransformationMatrix != nil {
		m.meanTransformed = make([]float32, featDim)
		for i := 0; i < featDim; i++ {
			for j := 0; j < featDim; j++ {
				m.meanTransformed[i] += m.mean[j] * m.meanTransformationMatrix[i][j]
			}
		}
	} else {
		m.meanTransformed = m.mean
	}

	if m.meanTransformationVector != nil {
		for k := 0; k < featDim; k++ {
			m.meanTransformed[k] += m.meanTransformationVector[k]
		}
	}
	/**
	 * We do analogously with the variance. In this case, we also
	 * invert the variance, and work with precision instead of
	 * variance.
	 */
	if m.varianceTransformationMatrix != nil {
		m.precisionTransformed = make([]float32, len(m.variance))
		for i := 0; i < featDim; i++ {
			for j := 0; j < featDim; j++ {
				m.precisionTransformed[i] += m.variance[j] * m.varianceTransformationMatrix[i][j]
			}
		}
	} else {
		m.precisionTransformed = append([]float32(nil), m.variance...)
	}

	if m.varianceTransformationVector != nil {
		for k := 0; k < featDim; k++ {
			m.precisionTransformed[k] += m.varianceTransformationVector[k]
		}
	}
	for k := 0; k < featDim; k++ {
		var flooredPrecision float32
		if m.precisionTransformed[k] < m.varianceFloor {
			flooredPrecision = m.varianceFloor
		} else {
			flooredPrecision = m.precisionTransformed[k]
		}
		m.precisionTransformed[k] = 1.0 / (-2.0 * flooredPrecision)
	}
}

func (m *MixtureComponent) String() string {
	return fmt.Sprintf("mu=%v cov=%v", m.mean, m.variance)
}
