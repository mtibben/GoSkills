package trueskill

import (
	"testing"
	"math"
)

func TestDrawMarginFromDrawProbability(t *testing.T) {
	const beta = 25.0 / 6.0

	// The expected values were compared against Ralf Herbrich's implementation in F#
	AssertDrawMargin(t, 0.10, beta, 0.74046637542690541)
	AssertDrawMargin(t, 0.25, beta, 1.87760059883033)
	AssertDrawMargin(t, 0.33, beta, 2.5111010132487492)
}

func AssertDrawMargin(t *testing.T, drawProb, beta, expected float64) {
	const errorTolerance = 0.000001
	actual := DrawMarginFromDrawProbability(drawProb, beta)
	if r := actual; math.Abs(r-expected) > errorTolerance {
		t.Errorf("draw margin = %v, want %v\n%v", r, expected, testLoc())
	}
}
