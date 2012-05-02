package trueskill

import (
	"testing"
)

func TestTwoPlayerCalc(t *testing.T) {
	// We only support two players
	AllTwoPlayerScenarios(t, &TwoPlayerCalc{})

	// TODO: Assert failures for larger teams
}
