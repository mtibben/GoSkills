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

	team1 := skills.NewTeam(*player1, gameInfo.DefaultRating())
	team2 := skills.NewTeam(*player2, gameInfo.DefaultRating())
	teams := []*skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

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

	team1 := skills.NewTeam(*player1, gameInfo.DefaultRating())
	team2 := skills.NewTeam(*player2, gameInfo.DefaultRating())
	teams := []*skills.Team{team1, team2}

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

	team1 := skills.NewTeam(*player1, skills.NewRating(1301.0007, 42.9232))
	team2 := skills.NewTeam(*player2, skills.NewRating(1188.7560, 42.5570))
	teams := []*skills.Team{team1, team2}

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

	team1 := skills.NewTeam(*player1, gameInfo.DefaultRating())
	team2 := skills.NewTeam(*player2, skills.NewRating(50, 12.5))
	teams := []*skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 1)

	player1NewRating := newRatings[*player1]
	AssertRating(t, 31.662, 7.137, player1NewRating)

	player2NewRating := newRatings[*player2]
	AssertRating(t, 35.010, 7.910, player2NewRating)
	AssertMatchQuality(t, 0.110, calc.CalcMatchQual(gameInfo, teams))
}

func testLoc() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%v:%v", file, line)
	}
	return ""
}

func AssertRating(t *testing.T, expectedMean, expectedStddev float64, actual skills.Rating) {
	if r := actual.Mean; math.Abs(r-expectedMean) > errorTolerance {
		t.Errorf("actual.Mean = %v, want %v\n%v", r, expectedMean, testLoc())
	}
	if r := actual.Stddev; math.Abs(r-expectedStddev) > errorTolerance {
		t.Errorf("actual.Stddev = %v, want %v\n%v", r, expectedStddev, testLoc())
	}
}
func AssertMatchQuality(t *testing.T, expectedMatchQual, actualMatchQual float64) {
	if r := actualMatchQual; math.Abs(r-expectedMatchQual) > errorTolerance {
		t.Errorf("actualMatchQual = %v, want %v\n%v", r, expectedMatchQual, testLoc())
	}
}
