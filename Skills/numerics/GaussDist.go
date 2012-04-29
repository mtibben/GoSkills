package numerics

import (
	"fmt"
	"math"
)

const logSqrt2Pi = 0.91893853320467274178032973640562

func GaussCumulativeTo(x float64) float64 {
	return math.Erf(x/math.Sqrt2)/2 + 0.5
}

func GaussAt(x float64) float64 {
	return math.Exp(-x*x/2) / (math.Sqrt2 * math.SqrtPi)
}

type GaussDist struct {
	Mean          float64
	Stddev        float64
	Precision     float64
	PrecisionMean float64
	Variance      float64
}

func NewGaussDist(mean, stddev float64) *GaussDist {
	variance := stddev * stddev
	precision := 1 / variance
	return &GaussDist{
		Mean:          mean,
		Stddev:        stddev,
		Variance:      variance,
		Precision:     precision,
		PrecisionMean: precision * mean,
	}
}

func (z *GaussDist) String() string {
	return fmt.Sprintf("{μ:%.6g σ:%.6g}", z.Mean, z.Stddev)
}

// Mul sets z to the product x*y and returns z.
func (z *GaussDist) Mul(x, y *GaussDist) *GaussDist {
	z.Precision = x.Precision + y.Precision
	z.PrecisionMean = x.PrecisionMean + y.PrecisionMean
	z.fromPrecisionMean()
	return z
}

// Div sets z to the product x/y and returns z.
func (z *GaussDist) Div(x, y *GaussDist) *GaussDist {
	z.Precision = x.Precision - y.Precision
	z.PrecisionMean = x.PrecisionMean - y.PrecisionMean
	z.fromPrecisionMean()
	return z
}

func (z *GaussDist) fromPrecisionMean() {
	z.Variance = 1 / z.Precision
	z.Stddev = math.Sqrt(z.Variance)
	z.Mean = z.PrecisionMean / z.Precision
}

// Returns the LogProductNormalization of x and y.
func LogProdNorm(x, y *GaussDist) float64 {
	if x.Precision == 0 || y.Precision == 0 {
		return 0
	}

	varSum := x.Variance + y.Variance
	meanDiff := x.Mean - y.Mean
	meanDiff2 := meanDiff * meanDiff

	return -logSqrt2Pi - (math.Log(varSum)+meanDiff2/varSum)/2.0
}

// Returns the LogRatioNormalization of x and y.
func LogRatioNorm(x, y *GaussDist) float64 {
	if x.Precision == 0 || y.Precision == 0 {
		return 0
	}

	varDiff := x.Variance - y.Variance
	meanDiff := x.Mean - y.Mean
	meanDiff2 := meanDiff * meanDiff

	return math.Log(y.Variance) + logSqrt2Pi - (math.Log(varDiff)-meanDiff2/varDiff)/2.0
}

// Computes the absolute difference between two Gaussians
func AbsDiff(x, y *GaussDist) float64 {
	return math.Max(math.Abs(x.PrecisionMean-y.PrecisionMean), math.Sqrt(math.Abs(x.Precision-y.Precision)))
}
