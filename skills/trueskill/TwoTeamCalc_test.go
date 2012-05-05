package trueskill

import (
	"testing"
)

func TestTwoTeamCalc(t *testing.T) {
	AllTwoTeamScenarios(t, &TwoTeamCalc{})
}
