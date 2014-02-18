package trueskill

import (
	"testing"
)

func TestFullFactorGraphCalculatorTests(t *testing.T) {
	calculator := &FactorGraphTrueSkillCalc{}

	// We can test all classes
	AllTwoPlayerScenarios(calculator)
	AllTwoTeamScenarios(calculator)
	AllMultipleTeamScenarios(calculator)

	PartialPlayScenarios(calculator)
}
