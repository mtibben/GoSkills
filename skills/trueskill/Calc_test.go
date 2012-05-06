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
	/*
		OneOnTwoSimpleTest(t, calc)
		OneOnTwoDrawTest(t, calc)
		OneOnTwoSomewhatBalanced(t, calc)
		OneOnThreeDrawTest(t, calc)
		OneOnThreeSimpleTest(t, calc)
		OneOnSevenSimpleTest(t, calc)
	*/

	TwoOnTwoSimpleTest(t, calc)
	TwoOnTwoDrawTest(t, calc)
	TwoOnTwoUnbalancedDrawTest(t, calc)
	TwoOnTwoUpsetTest(t, calc)

	/*
		ThreeOnTwoTests(t, calc)

		FourOnFourSimpleTest(t, calc)
	*/
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

	teams := []skills.Team{team1, team2}

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

	teams := []skills.Team{team1, team2}

	newRatings := calc.CalcNewRatings(gameInfo, teams, 1, 2)

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

/*
   func FourOnFourSimpleTest(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);
       var player2 = new Player(2);
       var player3 = new Player(3);
       var player4 = new Player(4);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, gameInfo.DefaultRating)
           .AddPlayer(player2, gameInfo.DefaultRating)
           .AddPlayer(player3, gameInfo.DefaultRating)
           .AddPlayer(player4, gameInfo.DefaultRating);

       var player5 = new Player(5);
       var player6 = new Player(6);
       var player7 = new Player(7);
       var player8 = new Player(8);

       var team2 = new Team()
                   .AddPlayer(player5, gameInfo.DefaultRating)
                   .AddPlayer(player6, gameInfo.DefaultRating)
                   .AddPlayer(player7, gameInfo.DefaultRating)
                   .AddPlayer(player8, gameInfo.DefaultRating);


       var teams = Teams.Concat(team1, team2);

       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 2);

       // Winners
       AssertRating(27.198, 8.059, newRatingsWinLose[player1]);
       AssertRating(27.198, 8.059, newRatingsWinLose[player2]);
       AssertRating(27.198, 8.059, newRatingsWinLose[player3]);
       AssertRating(27.198, 8.059, newRatingsWinLose[player4]);            

       // Losers
       AssertRating(22.802, 8.059, newRatingsWinLose[player5]);
       AssertRating(22.802, 8.059, newRatingsWinLose[player6]);
       AssertRating(22.802, 8.059, newRatingsWinLose[player7]);
       AssertRating(22.802, 8.059, newRatingsWinLose[player8]);

       AssertMatchQuality(0.447, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func OneOnTwoSimpleTest(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, gameInfo.DefaultRating);

       var player2 = new Player(2);
       var player3 = new Player(3);

       var team2 = new Team()
                   .AddPlayer(player2, gameInfo.DefaultRating)
                   .AddPlayer(player3, gameInfo.DefaultRating);

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 2);

       // Winners
       AssertRating(33.730, 7.317, newRatingsWinLose[player1]);

       // Losers
       AssertRating(16.270, 7.317, newRatingsWinLose[player2]);
       AssertRating(16.270, 7.317, newRatingsWinLose[player3]);

       AssertMatchQuality(0.135, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func OneOnTwoSomewhatBalanced(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, new Rating(40, 6));

       var player2 = new Player(2);
       var player3 = new Player(3);

       var team2 = new Team()
                   .AddPlayer(player2, new Rating(20, 7))
                   .AddPlayer(player3, new Rating(25, 8));

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 2);

       // Winners
       AssertRating(42.744, 5.602, newRatingsWinLose[player1]);

       // Losers
       AssertRating(16.266, 6.359, newRatingsWinLose[player2]);
       AssertRating(20.123, 7.028, newRatingsWinLose[player3]);

       AssertMatchQuality(0.478, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func OneOnThreeSimpleTest(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, gameInfo.DefaultRating);

       var player2 = new Player(2);
       var player3 = new Player(3);
       var player4 = new Player(4);

       var team2 = new Team()
                   .AddPlayer(player2, gameInfo.DefaultRating)
                   .AddPlayer(player3, gameInfo.DefaultRating)
                   .AddPlayer(player4, gameInfo.DefaultRating);

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 2);

       // Winners
       AssertRating(36.337, 7.527, newRatingsWinLose[player1]);

       // Losers
       AssertRating(13.663, 7.527, newRatingsWinLose[player2]);
       AssertRating(13.663, 7.527, newRatingsWinLose[player3]);
       AssertRating(13.663, 7.527, newRatingsWinLose[player4]);

       AssertMatchQuality(0.012, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func OneOnTwoDrawTest(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, gameInfo.DefaultRating);

       var player2 = new Player(2);
       var player3 = new Player(3);

       var team2 = new Team()
                   .AddPlayer(player2, gameInfo.DefaultRating)
                   .AddPlayer(player3, gameInfo.DefaultRating);

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 1);

       // Winners
       AssertRating(31.660, 7.138, newRatingsWinLose[player1]);

       // Losers
       AssertRating(18.340, 7.138, newRatingsWinLose[player2]);
       AssertRating(18.340, 7.138, newRatingsWinLose[player3]);

       AssertMatchQuality(0.135, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func OneOnThreeDrawTest(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, gameInfo.DefaultRating);

       var player2 = new Player(2);
       var player3 = new Player(3);
       var player4 = new Player(4);

       var team2 = new Team()
                   .AddPlayer(player2, gameInfo.DefaultRating)
                   .AddPlayer(player3, gameInfo.DefaultRating)
                   .AddPlayer(player4, gameInfo.DefaultRating);

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 1);

       // Winners
       AssertRating(34.990, 7.455, newRatingsWinLose[player1]);

       // Losers
       AssertRating(15.010, 7.455, newRatingsWinLose[player2]);
       AssertRating(15.010, 7.455, newRatingsWinLose[player3]);
       AssertRating(15.010, 7.455, newRatingsWinLose[player4]);

       AssertMatchQuality(0.012, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func OneOnSevenSimpleTest(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);

       var gameInfo = GameInfo.DefaultGameInfo;

       var team1 = new Team()
           .AddPlayer(player1, gameInfo.DefaultRating);

       var player2 = new Player(2);
       var player3 = new Player(3);
       var player4 = new Player(4);
       var player5 = new Player(5);
       var player6 = new Player(6);
       var player7 = new Player(7);
       var player8 = new Player(8);

       var team2 = new Team()
                   .AddPlayer(player2, gameInfo.DefaultRating)
                   .AddPlayer(player3, gameInfo.DefaultRating)
                   .AddPlayer(player4, gameInfo.DefaultRating)
                   .AddPlayer(player5, gameInfo.DefaultRating)
                   .AddPlayer(player6, gameInfo.DefaultRating)
                   .AddPlayer(player7, gameInfo.DefaultRating)
                   .AddPlayer(player8, gameInfo.DefaultRating);

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLose = calculator.CalculateNewRatings(gameInfo, teams, 1, 2);

       // Winners
       AssertRating(40.582, 7.917, newRatingsWinLose[player1]);

       // Losers
       AssertRating(9.418, 7.917, newRatingsWinLose[player2]);
       AssertRating(9.418, 7.917, newRatingsWinLose[player3]);
       AssertRating(9.418, 7.917, newRatingsWinLose[player4]);
       AssertRating(9.418, 7.917, newRatingsWinLose[player5]);
       AssertRating(9.418, 7.917, newRatingsWinLose[player6]);
       AssertRating(9.418, 7.917, newRatingsWinLose[player7]);
       AssertRating(9.418, 7.917, newRatingsWinLose[player8]);

       AssertMatchQuality(0.000, calculator.CalculateMatchQuality(gameInfo, teams));
   }

   func ThreeOnTwoTests(t *testing.T, calc skills.Calc)
   {
       var player1 = new Player(1);
       var player2 = new Player(2);
       var player3 = new Player(3);

       var team1 = new Team()
                   .AddPlayer(player1, new Rating(28, 7))
                   .AddPlayer(player2, new Rating(27, 6))
                   .AddPlayer(player3, new Rating(26, 5));


       var player4 = new Player(4);
       var player5 = new Player(5);

       var team2 = new Team()
                   .AddPlayer(player4, new Rating(30, 4))
                   .AddPlayer(player5, new Rating(31, 3));

       var gameInfo = GameInfo.DefaultGameInfo;

       var teams = Teams.Concat(team1, team2);
       var newRatingsWinLoseExpected = calculator.CalculateNewRatings(gameInfo, teams, 1, 2);

       // Winners
       AssertRating(28.658, 6.770, newRatingsWinLoseExpected[player1]);
       AssertRating(27.484, 5.856, newRatingsWinLoseExpected[player2]);
       AssertRating(26.336, 4.917, newRatingsWinLoseExpected[player3]);

       // Losers
       AssertRating(29.785, 3.958, newRatingsWinLoseExpected[player4]);
       AssertRating(30.879, 2.983, newRatingsWinLoseExpected[player5]);

       var newRatingsWinLoseUpset = calculator.CalculateNewRatings(gameInfo, Teams.Concat(team1, team2), 2, 1);

       // Winners
       AssertRating(32.012, 3.877, newRatingsWinLoseUpset[player4]);
       AssertRating(32.132, 2.949, newRatingsWinLoseUpset[player5]);

       // Losers
       AssertRating(21.840, 6.314, newRatingsWinLoseUpset[player1]);
       AssertRating(22.474, 5.575, newRatingsWinLoseUpset[player2]);
       AssertRating(22.857, 4.757, newRatingsWinLoseUpset[player3]);

       AssertMatchQuality(0.254, calculator.CalculateMatchQuality(gameInfo, teams));
   }
*/

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
