package numerics

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTwoByTwoDeterminantTests(t *testing.T) {

	a := NewSquareMatrix(
		1, 2,
		3, 4,
	)
	AssertAreEqual(t, -2, a.Determinant())

	b := NewSquareMatrix(
		3, 4,
		5, 6,
	)
	AssertAreEqual(t, -2, b.Determinant())

	c := NewSquareMatrix(
		1, 1,
		1, 1,
	)
	AssertAreEqual(t, 0, c.Determinant())

	d := NewSquareMatrix(
		12, 15,
		17, 21,
	)
	AssertAreEqual(t, 12*21-15*17, d.Determinant())
}

func TestThreeByThreeDeterminantTests(t *testing.T) {

	a := NewSquareMatrix(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)
	AssertAreEqual(t, 0, a.Determinant())

	π := NewSquareMatrix(
		3, 1, 4,
		1, 5, 9,
		2, 6, 5,
	)
	// Verified against http://www.wolframalpha.com/input/?i=determinant+%7B%7B3%2C1%2C4%7D%2C%7B1%2C5%2C9%7D%2C%7B2%2C6%2C5%7D%7D
	AssertAreEqual(t, -90, π.Determinant())
}
func TestFourByFourDeterminantTests(t *testing.T) {

	a := NewSquareMatrix(
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	)
	in := a.Determinant()
	out := 0
	Convey(fmt.Sprintf("Determinant %v should equal %v", in, out), t, func() {
		So(in, ShouldEqual, out)
	})

	π := NewSquareMatrix(
		3, 1, 4, 1,
		5, 9, 2, 6,
		5, 3, 5, 8,
		9, 7, 9, 3,
	)

	// Verified against http://www.wolframalpha.com/input/?i=determinant+%7B+%7B3%2C1%2C4%2C1%7D%2C+%7B5%2C9%2C2%2C6%7D%2C+%7B5%2C3%2C5%2C8%7D%2C+%7B9%2C7%2C9%2C3%7D%7D
	in = π.Determinant()
	out = 98
	// Verified against http://www.wolframalpha.com/input/?i=determinant+%7B%7B3%2C1%2C4%7D%2C%7B1%2C5%2C9%7D%2C%7B2%2C6%2C5%7D%7D
	Convey(fmt.Sprintf("Determinant %v should equal %v", in, out), t, func() {
		So(in, ShouldEqual, out)
	})
}

func TestEightByEightDeterminantTests(t *testing.T) {
	a := NewSquareMatrix(
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32,
		33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 32, 44, 45, 46, 47, 48,
		49, 50, 51, 52, 53, 54, 55, 56,
		57, 58, 59, 60, 61, 62, 63, 64,
	)

	AssertAreEqual(t, 0, a.Determinant())

	π := NewSquareMatrix(
		3, 1, 4, 1, 5, 9, 2, 6,
		5, 3, 5, 8, 9, 7, 9, 3,
		2, 3, 8, 4, 6, 2, 6, 4,
		3, 3, 8, 3, 2, 7, 9, 5,
		0, 2, 8, 8, 4, 1, 9, 7,
		1, 6, 9, 3, 9, 9, 3, 7,
		5, 1, 0, 5, 8, 2, 0, 9,
		7, 4, 9, 4, 4, 5, 9, 2,
	)

	// Verified against http://www.wolframalpha.com/input/?i=det+%7B%7B3%2C1%2C4%2C1%2C5%2C9%2C2%2C6%7D%2C%7B5%2C3%2C5%2C8%2C9%2C7%2C9%2C3%7D%2C%7B2%2C3%2C8%2C4%2C6%2C2%2C6%2C4%7D%2C%7B3%2C3%2C8%2C3%2C2%2C7%2C9%2C5%7D%2C%7B0%2C2%2C8%2C8%2C4%2C1%2C9%2C7%7D%2C%7B1%2C6%2C9%2C3%2C9%2C9%2C3%2C7%7D%2C%7B5%2C1%2C0%2C5%2C8%2C2%2C0%2C9%7D%2C%7B7%2C4%2C9%2C4%2C4%2C5%2C9%2C2%7D%7D
	AssertAreEqual(t, 1378143, π.Determinant())
}

func TestEqualsTest(t *testing.T) {
	a := NewSquareMatrix(
		1, 2,
		3, 4,
	)

	b := NewSquareMatrix(
		1, 2,
		3, 4,
	)

	AssertIsTrue(t, a.Equals(b))
	AssertShouldResemble(t, a, b)

	c := NewMatrixFromValues(2, 3,
		1, 2, 3,
		4, 5, 6,
	)

	d := NewMatrixFromValues(2, 3,
		1, 2, 3,
		4, 5, 6,
	)

	AssertIsTrue(t, c.Equals(d))
	AssertShouldResemble(t, c, d)

	e := NewMatrixFromValues(3, 2,
		1, 4,
		2, 5,
		3, 6,
	)

	f := e.Transpose()
	AssertIsTrue(t, d.Equals(f))
	AssertShouldResemble(t, d, f)

	AssertAreEqual(t, d.GetHashCode(), f.GetHashCode())

	// Test rounding (thanks to nsp on GitHub for finding this case)
	g := NewSquareMatrix(
		1, 2.00000000000001,
		3, 4,
	)

	h := NewSquareMatrix(
		1, 2,
		3, 4,
	)

	AssertIsTrue(t, g.Equals(h))
	AssertAreEqual(t, g, h)
	AssertAreEqual(t, g.GetHashCode(), h.GetHashCode())
}

func TestAdjugateTests(t *testing.T) {
	// From Wikipedia: http://en.wikipedia.org/wiki/Adjugate_matrix

	a := NewSquareMatrix(
		1, 2,
		3, 4,
	)

	b := NewSquareMatrix(
		4, -2,
		-3, 1,
	)

	AssertShouldResemble(t, b, a.Adjugate())

	c := NewSquareMatrix(
		-3, 2, -5,
		-1, 0, -2,
		3, -4, 1,
	)

	d := NewSquareMatrix(
		-8, 18, -4,
		-5, 12, -1,
		4, -6, 2,
	)

	AssertShouldResemble(t, d, c.Adjugate())
}

func TestInverseTests(t *testing.T) {
	// see http://www.mathwords.com/i/inverse_of_a_matrix.htm
	a := NewSquareMatrix(
		4, 3,
		3, 2,
	)

	b := NewSquareMatrix(
		-2, 3,
		3, -4,
	)

	aInverse := a.Inverse()
	AssertShouldResemble(t, b, aInverse)

	identity2x2 := NewIdentityMatrix(
		2,
	)

	aaInverse := Multiply(a, aInverse)
	AssertIsTrue(t, identity2x2.Equals(aaInverse))
	AssertShouldResemble(t, identity2x2, aaInverse)

	c := NewSquareMatrix(
		1, 2, 3,
		0, 4, 5,
		1, 0, 6,
	)

	cInverse := c.Inverse()
	d := MultiplyBy(1.0/22, NewSquareMatrix(
		24, -12, -2,
		5, 3, -5,
		-4, 2, 4,
	))

	AssertIsTrue(t, d.Equals(cInverse))
	AssertShouldResemble(t, d, cInverse)
	identity3x3 := NewIdentityMatrix(
		3,
	)

	ccInverse := Multiply(c, cInverse)
	AssertIsTrue(t, identity3x3.Equals(ccInverse))
	AssertShouldResemble(t, identity3x3, ccInverse)
}

func AssertAreEqual(t *testing.T, a, b interface{}) {
	Convey(fmt.Sprintf("%v should equal %v", a, b), t, func() {
		So(a, ShouldEqual, b)
	})
}

func AssertShouldResemble(t *testing.T, a, b interface{}) {
	Convey(fmt.Sprintf("%v should equal %v", a, b), t, func() {
		So(a, ShouldResemble, b)
	})
}

func AssertIsTrue(t *testing.T, a bool) {
	AssertAreEqual(t, a, true)
}
