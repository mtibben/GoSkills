package trueskill

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
	"math"
	"runtime"
	"testing"
)

const (
	errorTolerance = 0.085
)

func AllTwoPlayerScenarios(t *testing.T, calc skills.Calc) {
	TwoPlayerTestNotDrawn(t, calc)
	TwoPlayerTestDrawn(t, calc)
	OneOnOneMassiveUpsetDrawTest(t, calc)

	TwoPlayerChessTestNotDrawn(t, calc)
}

func AllTwoTeamScenarios(t *testing.T, calc skills.Calc) {
	OneOnTwoSimpleTest(t, calc)
	OneOnTwoSomewhatBalanced(t, calc)
	OneOnThreeSimpleTest(t, calc)
	OneOnTwoDrawTest(t, calc)
	OneOnThreeDrawTest(t, calc)
	OneOnSevenSimpleTest(t, calc)

	TwoOnTwoSimpleTest(t, calc)
	TwoOnTwoDrawTest(t, calc)
	TwoOnTwoUnbalancedDrawTest(t, calc)
	TwoOnTwoUpsetTest(t, calc)

	ThreeOnTwoTests(t, calc)

	FourOnFourSimpleTest(t, calc)
}

func AllMultipleTeamScenarios(t *testing.T, calc skills.Calc) {
	ThreeTeamsOfOneNotDrawn(t, calc)
	ThreeTeamsOfOneDrawn(t, calc)
	FourTeamsOfOneNotDrawn(t, calc)
	FiveTeamsOfOneNotDrawn(t, calc)
	EightTeamsOfOneDrawn(t, calc)
	EightTeamsOfOneUpset(t, calc)
	SixteenTeamsOfOneNotDrawn(t, calc)

	TwoOnFourOnTwoWinDraw(t, calc)
}

func PartialPlayScenarios(t *testing.T, calc skills.Calc) {
	OneOnTwoBalancedPartialPlay(t, calc)
}

//------------------- Actual Tests ---------------------------
// If you see more than 3 digits of precision in the decimal point, then the expected values calculated from 
// F# RalfH's implementation with the same input. It didn't support teams, so team values all came from the 
// online calculator at http://atom.research.microsoft.com/trueskill/rankcalculator.aspx
//
// All match quality expected values came from the online calculator

// In both cases, there may be some discrepancy after the first decimal point. I think this is due to my implementation
// using slightly higher precision in GaussianDistribution.

//------------------------------------------------------------------------------
// Two Player Tests
//------------------------------------------------------------------------------

func TwoPlayerTestNotDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())

	teams := []skills.Team{team2, team1}
	ranks := []int{2, 1}

	newRatings := calc.CalcNewRatings(gameInfo, teams, ranks...)

	if teams[0].Players()[0] != *player2 {
		t.Errorf("client teams slice reordered")
	}
	if ranks[0] != 2 {
		t.Errorf("client ranks slice reordered")
	}

	player1NewRating := newRatings[*player1]
	AssertRating(t, 29.39583201999924, 7.171475587326186, player1NewRating)

	player2NewRating := newRatings[*player2]
	AssertRating(t, 20.60416798000076, 7.171475587326186, player2NewRating)
	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoPlayerTestDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	player1NewRating := newRatings[*player1]
	AssertRating(t, 25.0, 6.4575196623173081, player1NewRating)

	player2NewRating := newRatings[*player2]
	AssertRating(t, 25.0, 6.4575196623173081, player2NewRating)
	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoPlayerChessTestNotDrawn(t *testing.T, calc skills.Calc) {
	// Inspired by a real bug :-)
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := &skills.GameInfo{
		InitialMean:     1200,
		InitialStddev:   1200 / 3,
		Beta:            200,
		DynamicsFactor:  1200 / 300,
		DrawProbability: 0.03,
	}

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, skills.NewRating(1301.0007, 42.9232))

	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, skills.NewRating(1188.7560, 42.5570))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	player1NewRating := newRatings[*player1]
	AssertRating(t, 1304.7820836053318, 42.843513887848658, player1NewRating)

	player2NewRating := newRatings[*player2]
	AssertRating(t, 1185.0383099003536, 42.485604606897752, player2NewRating)
}

func OneOnOneMassiveUpsetDrawTest(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, skills.NewRating(50, 12.5))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	player1NewRating := newRatings[*player1]
	AssertRating(t, 31.662, 7.137, player1NewRating)

	player2NewRating := newRatings[*player2]
	AssertRating(t, 35.010, 7.910, player2NewRating)
	AssertMatchQuality(t, 0.110, calc.CalcMatchQual(gameInfo, teams))
}

//------------------------------------------------------------------------------
// Two Team Tests
//------------------------------------------------------------------------------

func TwoOnTwoSimpleTest(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())
	team1.AddPlayer(*player2, gameInfo.DefaultRating())

	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player3, gameInfo.DefaultRating())
	team2.AddPlayer(*player4, gameInfo.DefaultRating())

	teams := []skills.Team{team2, team1}

	ranks := []int{2, 1}
	newRatings := calc.CalcNewRatings(gameInfo, teams, ranks...)
	if ranks[0] != 2 {
		t.Errorf("client ranks slice reordered")
	}

	// Winners
	AssertRating(t, 28.108, 7.774, newRatings[*player1])
	AssertRating(t, 28.108, 7.774, newRatings[*player2])

	// Losers
	AssertRating(t, 21.892, 7.774, newRatings[*player3])
	AssertRating(t, 21.892, 7.774, newRatings[*player4])

	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoOnTwoDrawTest(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())
	team1.AddPlayer(*player2, gameInfo.DefaultRating())

	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player3, gameInfo.DefaultRating())
	team2.AddPlayer(*player4, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 25, 7.455, newRatings[*player1])
	AssertRating(t, 25, 7.455, newRatings[*player2])

	// Losers
	AssertRating(t, 25, 7.445, newRatings[*player3])
	AssertRating(t, 25, 7.445, newRatings[*player4])

	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func TwoOnTwoUnbalancedDrawTest(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, skills.NewRating(15, 8))
	team1.AddPlayer(*player2, skills.NewRating(20, 6))

	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player3, skills.NewRating(25, 4))
	team2.AddPlayer(*player4, skills.NewRating(30, 3))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 21.570, 6.556, newRatings[*player1])
	AssertRating(t, 23.696, 5.418, newRatings[*player2])

	// Losers
	AssertRating(t, 23.357, 3.833, newRatings[*player3])
	AssertRating(t, 29.075, 2.931, newRatings[*player4])

	AssertMatchQuality(t, 0.214, calc.CalcMatchQual(gameInfo, teams))
}

func TwoOnTwoUpsetTest(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	gameInfo := skills.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, skills.NewRating(20, 8))
	team1.AddPlayer(*player2, skills.NewRating(25, 6))

	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player3, skills.NewRating(35, 7))
	team2.AddPlayer(*player4, skills.NewRating(40, 5))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 29.698, 7.008, newRatings[*player1])
	AssertRating(t, 30.455, 5.594, newRatings[*player2])

	// Losers
	AssertRating(t, 27.575, 6.346, newRatings[*player3])
	AssertRating(t, 36.211, 4.768, newRatings[*player4])

	AssertMatchQuality(t, 0.084, calc.CalcMatchQual(gameInfo, teams))
}

func FourOnFourSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)

	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())
	team1.AddPlayer(*player2, gameInfo.DefaultRating())
	team1.AddPlayer(*player3, gameInfo.DefaultRating())
	team1.AddPlayer(*player4, gameInfo.DefaultRating())

	player5 := skills.NewPlayer(5)
	player6 := skills.NewPlayer(6)
	player7 := skills.NewPlayer(7)
	player8 := skills.NewPlayer(8)

	team2 := skills.NewTeam()
	team2.AddPlayer(*player5, gameInfo.DefaultRating())
	team2.AddPlayer(*player6, gameInfo.DefaultRating())
	team2.AddPlayer(*player7, gameInfo.DefaultRating())
	team2.AddPlayer(*player8, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 27.198, 8.059, newRatings[*player1])
	AssertRating(t, 27.198, 8.059, newRatings[*player2])
	AssertRating(t, 27.198, 8.059, newRatings[*player3])
	AssertRating(t, 27.198, 8.059, newRatings[*player4])

	// Losers
	AssertRating(t, 22.802, 8.059, newRatings[*player5])
	AssertRating(t, 22.802, 8.059, newRatings[*player6])
	AssertRating(t, 22.802, 8.059, newRatings[*player7])
	AssertRating(t, 22.802, 8.059, newRatings[*player8])

	AssertMatchQuality(t, 0.447, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnTwoSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())
	team2.AddPlayer(*player3, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 33.730, 7.317, newRatings[*player1])

	// Losers
	AssertRating(t, 16.270, 7.317, newRatings[*player2])
	AssertRating(t, 16.270, 7.317, newRatings[*player3])

	AssertMatchQuality(t, 0.135, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnTwoSomewhatBalanced(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, skills.NewRating(40, 6))

	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, skills.NewRating(20, 7))
	team2.AddPlayer(*player3, skills.NewRating(25, 8))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 42.744, 5.602, newRatings[*player1])

	// Losers
	AssertRating(t, 16.266, 6.359, newRatings[*player2])
	AssertRating(t, 20.123, 7.028, newRatings[*player3])

	AssertMatchQuality(t, 0.478, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnThreeSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())
	team2.AddPlayer(*player3, gameInfo.DefaultRating())
	team2.AddPlayer(*player4, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 36.337, 7.527, newRatings[*player1])

	// Losers
	AssertRating(t, 13.663, 7.527, newRatings[*player2])
	AssertRating(t, 13.663, 7.527, newRatings[*player3])
	AssertRating(t, 13.663, 7.527, newRatings[*player4])

	AssertMatchQuality(t, 0.012, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnTwoDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())
	team2.AddPlayer(*player3, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 31.660, 7.138, newRatings[*player1])

	// Losers
	AssertRating(t, 18.340, 7.138, newRatings[*player2])
	AssertRating(t, 18.340, 7.138, newRatings[*player3])

	AssertMatchQuality(t, 0.135, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnThreeDrawTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())
	team2.AddPlayer(*player3, gameInfo.DefaultRating())
	team2.AddPlayer(*player4, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	// Winners
	AssertRating(t, 34.990, 7.455, newRatings[*player1])

	// Losers
	AssertRating(t, 15.010, 7.455, newRatings[*player2])
	AssertRating(t, 15.010, 7.455, newRatings[*player3])
	AssertRating(t, 15.010, 7.455, newRatings[*player4])

	AssertMatchQuality(t, 0.012, calc.CalcMatchQual(gameInfo, teams))
}

func OneOnSevenSimpleTest(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, gameInfo.DefaultRating())

	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	player6 := skills.NewPlayer(6)
	player7 := skills.NewPlayer(7)
	player8 := skills.NewPlayer(8)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player2, gameInfo.DefaultRating())
	team2.AddPlayer(*player3, gameInfo.DefaultRating())
	team2.AddPlayer(*player4, gameInfo.DefaultRating())
	team2.AddPlayer(*player5, gameInfo.DefaultRating())
	team2.AddPlayer(*player6, gameInfo.DefaultRating())
	team2.AddPlayer(*player7, gameInfo.DefaultRating())
	team2.AddPlayer(*player8, gameInfo.DefaultRating())

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 40.582, 7.917, newRatings[*player1])

	// Losers
	AssertRating(t, 9.418, 7.917, newRatings[*player2])
	AssertRating(t, 9.418, 7.917, newRatings[*player3])
	AssertRating(t, 9.418, 7.917, newRatings[*player4])
	AssertRating(t, 9.418, 7.917, newRatings[*player5])
	AssertRating(t, 9.418, 7.917, newRatings[*player6])
	AssertRating(t, 9.418, 7.917, newRatings[*player7])
	AssertRating(t, 9.418, 7.917, newRatings[*player8])

	AssertMatchQuality(t, 0.000, calc.CalcMatchQual(gameInfo, teams))
}

func ThreeOnTwoTests(t *testing.T, calc skills.Calc) {
	gameInfo := skills.DefaultGameInfo

	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	team1 := skills.NewTeam()
	team1.AddPlayer(*player1, skills.NewRating(28, 7))
	team1.AddPlayer(*player2, skills.NewRating(27, 6))
	team1.AddPlayer(*player3, skills.NewRating(26, 5))

	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	team2 := skills.NewTeam()
	team2.AddPlayer(*player4, skills.NewRating(30, 4))
	team2.AddPlayer(*player5, skills.NewRating(31, 3))

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

	// Winners
	AssertRating(t, 28.658, 6.770, newRatings[*player1])
	AssertRating(t, 27.484, 5.856, newRatings[*player2])
	AssertRating(t, 26.336, 4.917, newRatings[*player3])

	// Losers
	AssertRating(t, 29.785, 3.958, newRatings[*player4])
	AssertRating(t, 30.879, 2.983, newRatings[*player5])

	newRatings = calc.CalcNewRatings(gameInfo, teams, 2, 1)

	// Losers
	AssertRating(t, 21.840, 6.314, newRatings[*player1])
	AssertRating(t, 22.474, 5.575, newRatings[*player2])
	AssertRating(t, 22.857, 4.757, newRatings[*player3])

	// Winners
	AssertRating(t, 32.012, 3.877, newRatings[*player4])
	AssertRating(t, 32.132, 2.949, newRatings[*player5])

	AssertMatchQuality(t, 0.254, calc.CalcMatchQual(gameInfo, teams))
}

//------------------------------------------------------------------------------
// Multiple Teams Tests
//------------------------------------------------------------------------------
func TwoOnFourOnTwoWinDraw(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)

	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeam()
	team1.AddPlayer(player1, skills.NewRating(40, 4))
	team1.AddPlayer(player2, skills.NewRating(45, 3))

	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	player6 := skills.NewPlayer(6)

	team2 := skills.NewTeam()
	team2.AddPlayer(player3, skills.NewRating(20, 7))
	team2.AddPlayer(player4, skills.NewRating(19, 6))
	team2.AddPlayer(player5, skills.NewRating(30, 9))
	team2.AddPlayer(player6, skills.NewRating(10, 4))

	player7 := skills.NewPlayer(7)
	player8 := skills.NewPlayer(8)

	team3 := skills.NewTeam()
	team3.AddPlayer(player7, skills.NewRating(50, 5))
	team3.AddPlayer(player8, skills.NewRating(30, 2))

	teams := []skills.Team{team1, team2, team3}
	newRatingsWinLose = calc.CalculateNewRatings(gameInfo, teams, 1, 2, 2)

	// Winners
	AssertRating(40.877, 3.840, newRatingsWinLose[player1])
	AssertRating(45.493, 2.934, newRatingsWinLose[player2])
	AssertRating(19.609, 6.396, newRatingsWinLose[player3])
	AssertRating(18.712, 5.625, newRatingsWinLose[player4])
	AssertRating(29.353, 7.673, newRatingsWinLose[player5])
	AssertRating(9.872, 3.891, newRatingsWinLose[player6])
	AssertRating(48.830, 4.590, newRatingsWinLose[player7])
	AssertRating(29.813, 1.976, newRatingsWinLose[player8])

	AssertMatchQuality(0.367, calc.CalculateMatchQuality(gameInfo, teams))
}

func ThreeTeamsOfOneNotDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, gameInfo.DefaultRating)
	team2 := skills.NewTeamWithPlayer(player2, gameInfo.DefaultRating)
	team3 := skills.NewTeamWithPlayer(player3, gameInfo.DefaultRating)

	teams := []skills.Team{team1, team2, team3}
	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 2, 3)

	player1NewRating := newRatings[player1]
	AssertRating(31.675352419172107, 6.6559853776206905, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(25.000000000003912, 6.2078966412243233, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(18.324647580823971, 6.6559853776218318, player3NewRating)

	AssertMatchQuality(0.200, calc.CalculateMatchQuality(gameInfo, teams))
}

func ThreeTeamsOfOneDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, gameInfo.DefaultRating)
	team2 := skills.NewTeamWithPlayer(player2, gameInfo.DefaultRating)
	team3 := skills.NewTeamWithPlayer(player3, gameInfo.DefaultRating)

	teams := []skills.Team{team1, team2, team3}
	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 1, 1)

	player1NewRating := newRatings[player1]
	AssertRating(25.000, 5.698, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(25.000, 5.695, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(25.000, 5.698, player3NewRating)

	AssertMatchQuality(0.200, calc.CalculateMatchQuality(gameInfo, teams))
}

func FourTeamsOfOneNotDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, gameInfo.DefaultRating)
	team2 := skills.NewTeamWithPlayer(player2, gameInfo.DefaultRating)
	team3 := skills.NewTeamWithPlayer(player3, gameInfo.DefaultRating)
	team4 := skills.NewTeamWithPlayer(player4, gameInfo.DefaultRating)

	teams := []skills.Team{team1, team2, team3, team4}

	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 2, 3, 4)

	player1NewRating := newRatings[player1]
	AssertRating(33.206680965631264, 6.3481091698077057, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(27.401454693843323, 5.7871629348447584, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(22.598545306188374, 5.7871629348413451, player3NewRating)

	player4NewRating := newRatings[player4]
	AssertRating(16.793319034361271, 6.3481091698144967, player4NewRating)

	AssertMatchQuality(0.089, calc.CalculateMatchQuality(gameInfo, teams))
}

func FiveTeamsOfOneNotDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, gameInfo.DefaultRating)
	team2 := skills.NewTeamWithPlayer(player2, gameInfo.DefaultRating)
	team3 := skills.NewTeamWithPlayer(player3, gameInfo.DefaultRating)
	team4 := skills.NewTeamWithPlayer(player4, gameInfo.DefaultRating)
	team5 := skills.NewTeamWithPlayer(player5, gameInfo.DefaultRating)

	teams := []skills.Team{team1, team2, team3, team4, team5}
	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 2, 3, 4, 5)

	player1NewRating := newRatings[player1]
	AssertRating(34.363135705841188, 6.1361528798112692, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(29.058448805636779, 5.5358352402833413, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(25.000000000031758, 5.4200805474429847, player3NewRating)

	player4NewRating := newRatings[player4]
	AssertRating(20.941551194426314, 5.5358352402709672, player4NewRating)

	player5NewRating := newRatings[player5]
	AssertRating(15.636864294158848, 6.136152879829349, player5NewRating)

	AssertMatchQuality(0.040, calc.CalculateMatchQuality(gameInfo, teams))
}

func EightTeamsOfOneDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	player6 := skills.NewPlayer(6)
	player7 := skills.NewPlayer(7)
	player8 := skills.NewPlayer(8)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, gameInfo.DefaultRating)
	team2 := skills.NewTeamWithPlayer(player2, gameInfo.DefaultRating)
	team3 := skills.NewTeamWithPlayer(player3, gameInfo.DefaultRating)
	team4 := skills.NewTeamWithPlayer(player4, gameInfo.DefaultRating)
	team5 := skills.NewTeamWithPlayer(player5, gameInfo.DefaultRating)
	team6 := skills.NewTeamWithPlayer(player6, gameInfo.DefaultRating)
	team7 := skills.NewTeamWithPlayer(player7, gameInfo.DefaultRating)
	team8 := skills.NewTeamWithPlayer(player8, gameInfo.DefaultRating)

	teams := []skills.Team{team1, team2, team3, team4, team5, team6, team7, team8}
	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 1, 1, 1, 1, 1, 1, 1)

	player1NewRating := newRatings[player1]
	AssertRating(25.000, 4.592, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(25.000, 4.583, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(25.000, 4.576, player3NewRating)

	player4NewRating := newRatings[player4]
	AssertRating(25.000, 4.573, player4NewRating)

	player5NewRating := newRatings[player5]
	AssertRating(25.000, 4.573, player5NewRating)

	player6NewRating := newRatings[player6]
	AssertRating(25.000, 4.576, player6NewRating)

	player7NewRating := newRatings[player7]
	AssertRating(25.000, 4.583, player7NewRating)

	player8NewRating := newRatings[player8]
	AssertRating(25.000, 4.592, player8NewRating)

	AssertMatchQuality(0.004, calc.CalculateMatchQuality(gameInfo, teams))
}

func EightTeamsOfOneUpset(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	player6 := skills.NewPlayer(6)
	player7 := skills.NewPlayer(7)
	player8 := skills.NewPlayer(8)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, skills.NewRating(10, 8))
	team2 := skills.NewTeamWithPlayer(player2, skills.NewRating(15, 7))
	team3 := skills.NewTeamWithPlayer(player3, skills.NewRating(20, 6))
	team4 := skills.NewTeamWithPlayer(player4, skills.NewRating(25, 5))
	team5 := skills.NewTeamWithPlayer(player5, skills.NewRating(30, 4))
	team6 := skills.NewTeamWithPlayer(player6, skills.NewRating(35, 3))
	team7 := skills.NewTeamWithPlayer(player7, skills.NewRating(40, 2))
	team8 := skills.NewTeamWithPlayer(player8, skills.NewRating(45, 1))

	teams := []skills.Team{team1, team2, team3, team4, team5, team6, team7, team8}
	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 2, 3, 4, 5, 6, 7, 8)

	player1NewRating := newRatings[player1]
	AssertRating(35.135, 4.506, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(32.585, 4.037, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(31.329, 3.756, player3NewRating)

	player4NewRating := newRatings[player4]
	AssertRating(30.984, 3.453, player4NewRating)

	player5NewRating := newRatings[player5]
	AssertRating(31.751, 3.064, player5NewRating)

	player6NewRating := newRatings[player6]
	AssertRating(34.051, 2.541, player6NewRating)

	player7NewRating := newRatings[player7]
	AssertRating(38.263, 1.849, player7NewRating)

	player8NewRating := newRatings[player8]
	AssertRating(44.118, 0.983, player8NewRating)

	AssertMatchQuality(0.000, calc.CalculateMatchQuality(gameInfo, teams))
}

func SixteenTeamsOfOneNotDrawn(t *testing.T, calc skills.Calc) {
	player1 := skills.NewPlayer(1)
	player2 := skills.NewPlayer(2)
	player3 := skills.NewPlayer(3)
	player4 := skills.NewPlayer(4)
	player5 := skills.NewPlayer(5)
	player6 := skills.NewPlayer(6)
	player7 := skills.NewPlayer(7)
	player8 := skills.NewPlayer(8)
	player9 := skills.NewPlayer(9)
	player10 := skills.NewPlayer(10)
	player11 := skills.NewPlayer(11)
	player12 := skills.NewPlayer(12)
	player13 := skills.NewPlayer(13)
	player14 := skills.NewPlayer(14)
	player15 := skills.NewPlayer(15)
	player16 := skills.NewPlayer(16)
	gameInfo := GameInfo.DefaultGameInfo

	team1 := skills.NewTeamWithPlayer(player1, gameInfo.DefaultRating)
	team2 := skills.NewTeamWithPlayer(player2, gameInfo.DefaultRating)
	team3 := skills.NewTeamWithPlayer(player3, gameInfo.DefaultRating)
	team4 := skills.NewTeamWithPlayer(player4, gameInfo.DefaultRating)
	team5 := skills.NewTeamWithPlayer(player5, gameInfo.DefaultRating)
	team6 := skills.NewTeamWithPlayer(player6, gameInfo.DefaultRating)
	team7 := skills.NewTeamWithPlayer(player7, gameInfo.DefaultRating)
	team8 := skills.NewTeamWithPlayer(player8, gameInfo.DefaultRating)
	team9 := skills.NewTeamWithPlayer(player9, gameInfo.DefaultRating)
	team10 := skills.NewTeamWithPlayer(player10, gameInfo.DefaultRating)
	team11 := skills.NewTeamWithPlayer(player11, gameInfo.DefaultRating)
	team12 := skills.NewTeamWithPlayer(player12, gameInfo.DefaultRating)
	team13 := skills.NewTeamWithPlayer(player13, gameInfo.DefaultRating)
	team14 := skills.NewTeamWithPlayer(player14, gameInfo.DefaultRating)
	team15 := skills.NewTeamWithPlayer(player15, gameInfo.DefaultRating)
	team16 := skills.NewTeamWithPlayer(player16, gameInfo.DefaultRating)

	newRatings =
		calc.CalculateNewRatings(
			gameInfo,
			Teams.Concat(
				team1, team2, team3, team4, team5,
				team6, team7, team8, team9, team10,
				team11, team12, team13, team14, team15,
				team16),
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)

	player1NewRating := newRatings[player1]
	AssertRating(40.53945776946920, 5.27581643889050, player1NewRating)

	player2NewRating := newRatings[player2]
	AssertRating(36.80951229454210, 4.71121217610266, player2NewRating)

	player3NewRating := newRatings[player3]
	AssertRating(34.34726355544460, 4.52440328139991, player3NewRating)

	player4NewRating := newRatings[player4]
	AssertRating(32.33614722608720, 4.43258628279632, player4NewRating)

	player5NewRating := newRatings[player5]
	AssertRating(30.55048814671730, 4.38010805034365, player5NewRating)

	player6NewRating := newRatings[player6]
	AssertRating(28.89277312234790, 4.34859291776483, player6NewRating)

	player7NewRating := newRatings[player7]
	AssertRating(27.30952161972210, 4.33037679041216, player7NewRating)

	player8NewRating := newRatings[player8]
	AssertRating(25.76571046519540, 4.32197078088701, player8NewRating)

	player9NewRating := newRatings[player9]
	AssertRating(24.23428953480470, 4.32197078088703, player9NewRating)

	player10NewRating := newRatings[player10]
	AssertRating(22.69047838027800, 4.33037679041219, player10NewRating)

	player11NewRating := newRatings[player11]
	AssertRating(21.10722687765220, 4.34859291776488, player11NewRating)

	player12NewRating := newRatings[player12]
	AssertRating(19.44951185328290, 4.38010805034375, player12NewRating)

	player13NewRating := newRatings[player13]
	AssertRating(17.66385277391300, 4.43258628279643, player13NewRating)

	player14NewRating := newRatings[player14]
	AssertRating(15.65273644455550, 4.52440328139996, player14NewRating)

	player15NewRating := newRatings[player15]
	AssertRating(13.19048770545810, 4.71121217610273, player15NewRating)

	player16NewRating := newRatings[player16]
	AssertRating(9.46054223053080, 5.27581643889032, player16NewRating)
}

//------------------------------------------------------------------------------
// Partial Play Tests
//------------------------------------------------------------------------------

func OneOnTwoBalancedPartialPlay(t *testing.T, calc skills.Calc) {
	gameInfo := GameInfo.DefaultGameInfo

	p1 := skills.NewPlayer(1)
	team1 := skills.NewTeamWithPlayer(p1, gameInfo.DefaultRating)

	p2 := skills.NewPlayer(2, 0.0)
	p3 := skills.NewPlayer(3, 1.00)

	team2 := skills.NewTeam()
	team2.AddPlayer(p2, gameInfo.DefaultRating)
	team2.AddPlayer(p3, gameInfo.DefaultRating)

	teams := []skills.Team{team1, team2}
	newRatings := calc.CalculateNewRatings(gameInfo, teams, 1, 2)
	matchQuality = calc.CalculateMatchQuality(gameInfo, teams)

}

func testLoc() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%v:%v", file, line)
	}
	return ""
}

func AssertRating(t *testing.T, expectedMean, expectedStddev float64, actual skills.Rating) {
	if r := actual.Mean(); math.Abs(r-expectedMean) > errorTolerance {
		t.Errorf("actual.Mean = %v, want %v\n%v", r, expectedMean, testLoc())
	}
	if r := actual.Stddev(); math.Abs(r-expectedStddev) > errorTolerance {
		t.Errorf("actual.Stddev = %v, want %v\n%v", r, expectedStddev, testLoc())
	}
}
func AssertMatchQuality(t *testing.T, expectedMatchQual, actualMatchQual float64) {
	if r := actualMatchQual; math.Abs(r-expectedMatchQual) > errorTolerance {
		t.Errorf("actualMatchQual = %v, want %v\n%v", r, expectedMatchQual, testLoc())
	}
}
