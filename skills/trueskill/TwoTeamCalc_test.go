package trueskill

import (
	"testing"
)

func TestTwoTeamCalc(t *testing.T) {
	AllTwoPlayerScenarios(t, &TwoTeamCalc{})
	AllTwoTeamScenarios(t, &TwoTeamCalc{})
}
