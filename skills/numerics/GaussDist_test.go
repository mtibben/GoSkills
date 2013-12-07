package numerics

import (
	"math"
	"testing"
)

const (
	errorTolerance = 0.000001

	Sqrt10 = 3.1622776601683793319988935444327
)

func TestGaussAt(t *testing.T) {
	const in, out = 0.5, 0.352065
	if r := GaussAt(in); math.Abs(r-out) > errorTolerance {
		t.Errorf("GaussAt(%v) = %v, want %v", in, r, out)
	}
}

func TestGaussCumulativeTo(t *testing.T) {
	const in, out = 0.5, 0.69146246
	if r := GaussCumulativeTo(in); math.Abs(r-out) > errorTolerance {
		t.Errorf("GaussCumulativeTo(%v) = %v, want %v", in, r, out)
	}
}

func TestGaussInvCumulativeTo(t *testing.T) {
	const μ, σ, in, out = 0, 1, 0.69146246, 0.5
	if r := GaussInvCumulativeTo(in, μ, σ); math.Abs(r-out) > errorTolerance {
		t.Errorf("GaussCumulativeTo(%v) = %v, want %v", in, r, out)
	}
}

func TestInvErfc(t *testing.T) {
	// Verified with WolframAlpha
	// (e.g. http://www.wolframalpha.com/input/?i=CDF%5BNormalDistribution%5B0%2C1%5D%2C+0.5%5D )
	const in, out = 0.4794999836952529, 0.5
	if r := InvErfc(in); math.Abs(r-out) > errorTolerance {
		t.Errorf("InvErfc(%v) = %v, want %v", in, r, out)
	}
}

func TestSub(t *testing.T) {
	{
		stdNormal := NewGaussDist(0, 1)
		shiftedGaussian := NewGaussDist(2, 3)

		diff := new(GaussDist).Sub(stdNormal, shiftedGaussian)

		const expectedMean = -2.0
		if r := diff.Mean; math.Abs(r-expectedMean) > errorTolerance {
			t.Errorf("diff.Mean = %v, want %v", r, expectedMean)
		}

		var expectedStddev = math.Sqrt(10.0)
		if r := diff.Stddev; math.Abs(r-expectedStddev) > errorTolerance {
			t.Errorf("diff.Stddev = %v, want %v", r, expectedStddev)
		}
	}
}

func TestMul(t *testing.T) {
	// Verified against the formula at http://www.tina-vision.net/tina-knoppix/tina-memo/2003-003.pdf
	{
		stdNormal := NewGaussDist(0, 1)
		shiftedGaussian := NewGaussDist(2, 3)

		product := new(GaussDist).Mul(stdNormal, shiftedGaussian)

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

func TestDiv(t *testing.T) {
	// Since the multiplication was worked out by hand, we use the same numbers but work backwards
	{
		product := NewGaussDist(0.2, 3.0/Sqrt10)
		stdNormal := NewGaussDist(0, 1)

		quotient := new(GaussDist).Div(product, stdNormal)

		const expectedMean = 2.0
		if r := quotient.Mean; math.Abs(r-expectedMean) > errorTolerance {
			t.Errorf("quotient.Mean = %v, want %v", r, expectedMean)
		}

		const expectedStddev = 3.0
		if r := quotient.Stddev; math.Abs(r-expectedStddev) > errorTolerance {
			t.Errorf("quotient.Stddev = %v, want %v", r, expectedStddev)
		}
	}

	{
		const productMean = (4.*7.*7. + 6.*5.*5.) / (5.*5. + 7.*7.)
		productStddev := math.Sqrt((5. * 5. * 7. * 7.) / (5.*5. + 7.*7.))
		product := NewGaussDist(productMean, productStddev)
		m4s5 := NewGaussDist(4, 5)

		quotient := new(GaussDist).Div(product, m4s5)

		const expectedMean = 6.0
		if r := quotient.Mean; math.Abs(r-expectedMean) > errorTolerance {
			t.Errorf("quotient.Mean = %v, want %v", r, expectedMean)
		}

		expectedStddev := 7.0
		if r := quotient.Stddev; math.Abs(r-expectedStddev) > errorTolerance {
			t.Errorf("quotient.Stddev = %v, want %v", r, expectedStddev)
		}
	}
}

func TestCumulativeTo(t *testing.T) {
	{
		stdNormal := NewGaussDist(0, 1)

		const expected = 0.15865525
		if r := stdNormal.CumulativeTo(-1); math.Abs(r-expected) > errorTolerance {
			t.Errorf("stdNormal.CumulativeTo(-1) = %v, want %v", r, expected)
		}
	}
}

func TestLogProdNorm(t *testing.T) {
	// Verified with Ralf Herbrich's F# implementation
	{
		stdNormal := NewGaussDist(0, 1)
		const expected = -1.2655121234846454
		if r := LogProdNorm(stdNormal, stdNormal); math.Abs(r-expected) > errorTolerance {
			t.Errorf("LogProdNorm(%v, %v) = %v, want %v", stdNormal, stdNormal, r, expected)
		}
	}

	{
		m1s2 := NewGaussDist(1, 2)
		m3s4 := NewGaussDist(3, 4)
		const expected = -2.5168046699816684
		if r := LogProdNorm(m1s2, m3s4); math.Abs(r-expected) > errorTolerance {
			t.Errorf("LogProdNorm(%v, %v) = %v, want %v", m1s2, m3s4, r, expected)
		}
	}
}

func TestLogRatioNorm(t *testing.T) {
	// Verified with Ralf Herbrich's F# implementation
	m1s2 := NewGaussDist(1, 2)
	m3s4 := NewGaussDist(3, 4)
	const expected = 2.6157405972171204
	if r := LogRatioNorm(m1s2, m3s4); math.Abs(r-expected) > errorTolerance {
		t.Errorf("LogProdNorm(%v, %v) = %v, want %v", m1s2, m3s4, r, expected)
	}
}

func TestAbsDiff(t *testing.T) {
	// Verified with Ralf Herbrich's F# implementation
	{
		stdNormal := NewGaussDist(0, 1)
		const expected = 0.0
		if r := AbsDiff(stdNormal, stdNormal); math.Abs(r-expected) > errorTolerance {
			t.Errorf("AbsDiff(%v, %v) = %v, want %v", stdNormal, stdNormal, r, expected)
		}
	}

	{
		m1s2 := NewGaussDist(1, 2)
		m3s4 := NewGaussDist(3, 4)
		const expected = 0.4330127018922193
		if r := AbsDiff(m1s2, m3s4); math.Abs(r-expected) > errorTolerance {
			t.Errorf("AbsDiff(%v, %v) = %v, want %v", m1s2, m3s4, r, expected)
		}
	}
}
