package numerics

import (
	"math"
	"testing"
)

const (
	errorTolerance = 0.000001

	Sqrt10 = 3.1622776601683793319988935444327
)

func TestCumulativeTo(t *testing.T) {
	// Verified with WolframAlpha
	// (e.g. http://www.wolframalpha.com/input/?i=CDF%5BNormalDistribution%5B0%2C1%5D%2C+0.5%5D )
	const in, out = 0.5, 0.691462
	if r := GaussCumulativeTo(in); math.Abs(r-out) > errorTolerance {
		t.Errorf("GaussCumulativeTo(%v) = %v, want %v", in, r, out)
	}
}

func TestAt(t *testing.T) {
	// Verified with WolframAlpha
	// (e.g. http://www.wolframalpha.com/input/?i=PDF%5BNormalDistribution%5B0%2C1%5D%2C+0.5%5D )
	const in, out = 0.5, 0.352065
	if r := GaussAt(in); math.Abs(r-out) > errorTolerance {
		t.Errorf("GaussAt(%v) = %v, want %v", in, r, out)
	}
}

func TestMultiplication(t *testing.T) {
	// Verified against the formula at http://www.tina-vision.net/tina-knoppix/tina-memo/2003-003.pdf
	{
		standardNormal := NewGaussDist(0, 1)
		shiftedGaussian := NewGaussDist(2, 3)

		product := new(GaussDist).Mul(standardNormal, shiftedGaussian)

		const expectedMean = 0.2
		if r := product.Mean; math.Abs(r-expectedMean) > errorTolerance {
			t.Errorf("product.Mean = %v, want %v", r, expectedMean)
		}

		const expectedStddev = 3.0 / Sqrt10
		if r := product.Stddev; math.Abs(r-expectedStddev) > errorTolerance {
			t.Errorf("product.Stddev = %v, want %v", r, expectedStddev)
		}
	}

	{
		m4s5 := NewGaussDist(4, 5)
		m6s7 := NewGaussDist(6, 7)

		product := new(GaussDist).Mul(m4s5, m6s7)

		const expectedMean = (4.*7.*7. + 6.*5.*5.) / (5.*5. + 7.*7.)
		if r := product.Mean; math.Abs(r-expectedMean) > errorTolerance {
			t.Errorf("product.Mean = %v, want %v", r, expectedMean)
		}

		expectedStddev := math.Sqrt((5. * 5. * 7. * 7.) / (5.*5. + 7.*7.))
		if r := product.Stddev; math.Abs(r-expectedStddev) > errorTolerance {
			t.Errorf("product.Stddev = %v, want %v", r, expectedStddev)
		}
	}
}
