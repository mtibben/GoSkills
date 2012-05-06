package trueskill

import (
	"testing"
)

func TestTwoTeamCalc(t *testing.T) {
	AllTwoPlayerScenarios(t, &TwoPlayerCalc{})
	AllTwoTeamScenarios(t, &TwoTeamCalc{})
}
