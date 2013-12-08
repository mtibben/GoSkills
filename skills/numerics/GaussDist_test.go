package numerics

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

const (
	errorTolerance = 0.000001

	Sqrt10 = 3.1622776601683793319988935444327
)

func TestGaussAt(t *testing.T) {
	const in, out = 0.5, 0.352065
	Convey(fmt.Sprintf("GaussAt(%v) should equal %v", in, out), t, func() {
		So(GaussAt(in), ShouldAlmostEqual, out, errorTolerance)
	})
}

func TestGaussCumulativeTo(t *testing.T) {
	const in, out = 0.5, 0.69146246
	Convey(fmt.Sprintf("GaussCumulativeTo(%v) should equal %v", in, out), t, func() {
		So(GaussCumulativeTo(in), ShouldAlmostEqual, out, errorTolerance)
	})
}

func TestGaussInvCumulativeTo(t *testing.T) {
	const mu, sig, in, out = 0, 1, 0.69146246, 0.5
	Convey(fmt.Sprintf("GaussInvCumulativeTo(%v, %v, %v) should equal %v", in, mu, sig, out), t, func() {
		So(GaussInvCumulativeTo(in, mu, sig), ShouldAlmostEqual, out, errorTolerance)
	})
}

func TestInvErfc(t *testing.T) {
	// Verified with WolframAlpha
	// (e.g. http://www.wolframalpha.com/input/?i=CDF%5BNormalDistribution%5B0%2C1%5D%2C+0.5%5D )
	const in, out = 0.4794999836952529, 0.5
	Convey(fmt.Sprintf("InvErfc(%v) should equal %v", in, out), t, func() {
		So(InvErfc(in), ShouldAlmostEqual, out, errorTolerance)
	})
}

func TestSub(t *testing.T) {
	Convey("Given two Gaussian distributions", t, func() {
		stdNormal := NewGaussDist(0, 1)
		shiftedGaussian := NewGaussDist(2, 3)

		Convey("After subtracting one from the other", func() {
			diff := new(GaussDist).Sub(stdNormal, shiftedGaussian)

			Convey("The new mean should equal the difference of the means", func() {
				So(diff.Mean, ShouldEqual, -2)
			})

			Convey("The new standard deviation should equal the square root of the sum of the squares of the standard deviations", func() {
				So(diff.Stddev, ShouldAlmostEqual, math.Sqrt(10.0), errorTolerance)
			})
		})
	})
}

func TestMul(t *testing.T) {
	// Verified against the formula at http://www.tina-vision.net/tina-knoppix/tina-memo/2003-003.pdf

	Convey("Given two Gausian distributions", t, func() {
		stdNormal := NewGaussDist(0, 1)
		shiftedGaussian := NewGaussDist(2, 3)

		Convey("After multiplying one by the other", func() {
			product := new(GaussDist).Mul(stdNormal, shiftedGaussian)

			Convey("The new mean should equal the mean of the product", func() {
				So(product.Mean, ShouldAlmostEqual, 0.2, errorTolerance)
			})

			Convey("The new standard deviation should equal the standard deviation of the product", func() {
				So(product.Stddev, ShouldAlmostEqual, 3.0/Sqrt10, errorTolerance)
			})
		})
	})

	Convey("Given two Gausian distributions", t, func() {
		m4s5 := NewGaussDist(4, 5)
		m6s7 := NewGaussDist(6, 7)

		Convey("After multiplying one by the other", func() {
			product := new(GaussDist).Mul(m4s5, m6s7)

			Convey("The new mean should equal the mean of the product", func() {
				const expectedMean = (4.*7.*7. + 6.*5.*5.) / (5.*5. + 7.*7.)
				So(product.Mean, ShouldAlmostEqual, expectedMean, errorTolerance)
			})

			Convey("The new standard deviation should equal the standard deviation of the product", func() {
				expectedStddev := math.Sqrt((5. * 5. * 7. * 7.) / (5.*5. + 7.*7.))
				So(product.Stddev, ShouldAlmostEqual, expectedStddev, errorTolerance)
			})
		})
	})
}

func TestDiv(t *testing.T) {
	// Since the multiplication was worked out by hand, we use the same numbers but work backwards
	Convey("Given two Gausian distributions", t, func() {
		product := NewGaussDist(0.2, 3.0/Sqrt10)
		stdNormal := NewGaussDist(0, 1)

		Convey("After dividing one by the other", func() {
			quotient := new(GaussDist).Div(product, stdNormal)

			Convey("The new mean should equal the mean of the quotient", func() {
				const expectedMean = 2.0
				So(quotient.Mean, ShouldAlmostEqual, expectedMean, errorTolerance)
			})

			Convey("The new standard deviation should equal the standard deviation of the quotient", func() {
				const expectedStddev = 3.0
				So(quotient.Stddev, ShouldAlmostEqual, expectedStddev, errorTolerance)
			})
		})
	})

	Convey("Given two Gausian distributions", t, func() {
		const productMean = (4.*7.*7. + 6.*5.*5.) / (5.*5. + 7.*7.)
		productStddev := math.Sqrt((5. * 5. * 7. * 7.) / (5.*5. + 7.*7.))
		product := NewGaussDist(productMean, productStddev)
		m4s5 := NewGaussDist(4, 5)

		Convey("After dividing one by the other", func() {
			quotient := new(GaussDist).Div(product, m4s5)

			Convey("The new mean should equal the mean of the quotient", func() {
				const expectedMean = 6.0
				So(quotient.Mean, ShouldAlmostEqual, expectedMean, errorTolerance)
			})

			Convey("The new standard deviation should equal the standard deviation of the quotient", func() {
				const expectedStddev = 7.0
				So(quotient.Stddev, ShouldAlmostEqual, expectedStddev, errorTolerance)
			})
		})
	})
}

func TestCumulativeTo(t *testing.T) {
	Convey("Given a Gaussian distribution", t, func() {
		stdNormal := NewGaussDist(0, 1)
		const expected = 0.15865525
		Convey(fmt.Sprintf("g.CumulativeTo(-1) should equal %v", expected), func() {
			So(stdNormal.CumulativeTo(-1), ShouldAlmostEqual, expected, errorTolerance)
		})
	})
}

func TestLogProdNorm(t *testing.T) {
	// Verified with Ralf Herbrich's F# implementation
	Convey("Given a Gaussian distribution", t, func() {
		stdNormal := NewGaussDist(0, 1)
		Convey("LogProdNorm(g, g) returns the log product normalization", func() {
			const expected = -1.2655121234846454
			So(LogProdNorm(stdNormal, stdNormal), ShouldAlmostEqual, expected, errorTolerance)
		})
	})

	Convey("Given two Gaussian distributions", t, func() {
		m1s2 := NewGaussDist(1, 2)
		m3s4 := NewGaussDist(3, 4)
		Convey("LogProdNorm(g1, g2) returns the log product normalization", func() {
			const expected = -2.5168046699816684
			So(LogProdNorm(m1s2, m3s4), ShouldAlmostEqual, expected, errorTolerance)
		})
	})
}

func TestLogRatioNorm(t *testing.T) {
	// Verified with Ralf Herbrich's F# implementation
	Convey("Given two Gaussian distributions", t, func() {
		m1s2 := NewGaussDist(1, 2)
		m3s4 := NewGaussDist(3, 4)
		Convey("LogRatioNorm(g1, g2) returns the log ratio normalization", func() {
			const expected = 2.6157405972171204
			So(LogRatioNorm(m1s2, m3s4), ShouldAlmostEqual, expected, errorTolerance)
		})
	})
}

func TestAbsDiff(t *testing.T) {
	// Verified with Ralf Herbrich's F# implementation
	Convey("Given a Gaussian distribution", t, func() {
		stdNormal := NewGaussDist(0, 1)
		Convey("AbsDiff(g, g) returns the absolute difference", func() {
			So(AbsDiff(stdNormal, stdNormal), ShouldEqual, 0)
		})
	})

	Convey("Given two Gaussian distributions", t, func() {
		m1s2 := NewGaussDist(1, 2)
		m3s4 := NewGaussDist(3, 4)
		Convey("AbsDiff(g1, g2) returns the log product normalization", func() {
			const expected = 0.4330127018922193
			So(AbsDiff(m1s2, m3s4), ShouldAlmostEqual, expected, errorTolerance)
		})
	})
}
