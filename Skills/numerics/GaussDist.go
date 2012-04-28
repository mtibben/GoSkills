package numerics

import (
	"math"
)

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

// Mul sets z to the product x*y and returns z.
func (z *GaussDist) Mul(x, y *GaussDist) *GaussDist {
	z.Precision = x.Precision + y.Precision
	z.PrecisionMean = x.PrecisionMean + y.PrecisionMean
	z.fromPrecisionMean()
	return z
}

func (z *GaussDist) fromPrecisionMean() {
	z.Variance = 1 / z.Precision
	z.Stddev = math.Sqrt(z.Variance)
	z.Mean = z.PrecisionMean / z.Precision
}
