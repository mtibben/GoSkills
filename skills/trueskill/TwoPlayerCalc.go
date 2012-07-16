package trueskill

import (
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/numerics"
	"math"
	"sort"
)

// Calculates the new ratings for only two players.
// When you only have two players, a lot of the math simplifies. The main purpose of this type
// is to show the bare minimum of what a TrueSkill implementation should have.
type TwoPlayerCalc struct{}

// Calculates new ratings based on the prior ratings and team ranks use 1 for first place, repeat the number for a tie (e.g. 1, 2, 2).
func (calc *TwoPlayerCalc) CalcNewRatings(gi *skills.GameInfo, teams []skills.Team, ranks ...int) skills.PlayerRatings {
	newSkills := make(map[skills.Player]skills.Rating)

	// Basic argument checking
	ValidateTeamCount(teams, twoPlayerTeamRange)
	ValidatePlayersPerTeam(teams, twoPlayerPlayerRange)

	// Copy ranks slice so we don't confuse the client code
	sranks := append([]int{}, ranks...)

	// Make sure things are in order
	sort.Sort(skills.NewRankedTeams(teams, sranks))

	// Since we verified that each team has one player, we know the player is the first one
	winningTeam := teams[0]
	winner := winningTeam.Players()[0]
	winnerPrevRating := winningTeam.PlayerRating(winner)

	losingTeam := teams[1]
	loser := losingTeam.Players()[0]
	loserPrevRating := losingTeam.PlayerRating(loser)

	wasDraw := sranks[0] == sranks[1]

	newSkills[winner] = twoPlayerCalcNewRating(gi, winnerPrevRating, loserPrevRating, cond(wasDraw, skills.Draw, skills.Win))
	newSkills[loser] = twoPlayerCalcNewRating(gi, loserPrevRating, winnerPrevRating, cond(wasDraw, skills.Draw, skills.Lose))

	return newSkills
}

func twoPlayerCalcNewRating(gi *skills.GameInfo, selfRating, oppRating skills.Rating, comparison int) skills.Rating {
	drawMargin := DrawMarginFromDrawProbability(gi.DrawProbability, gi.Beta)

	c := math.Sqrt(numerics.Sqr(selfRating.Stddev()) + numerics.Sqr(oppRating.Stddev()) + 2*numerics.Sqr(gi.Beta))

	winningMean := selfRating.Mean()
	losingMean := oppRating.Mean()

	if comparison == skills.Lose {
		winningMean = oppRating.Mean()
		losingMean = selfRating.Mean()
	}

	meanDelta := winningMean - losingMean

	var v, w, rankMultiplier float64

	if comparison != skills.Draw {
		v = VExceedsMarginC(meanDelta, drawMargin, c)
		w = WExceedsMarginC(meanDelta, drawMargin, c)
		rankMultiplier = float64(comparison)
	} else {
		v = VWithinMarginC(meanDelta, drawMargin, c)
		w = WWithinMarginC(meanDelta, drawMargin, c)
		rankMultiplier = 1
	}

	meanMultiplier := (numerics.Sqr(selfRating.Stddev()) + numerics.Sqr(gi.DynamicsFactor)) / c

	varianceWithDynamics := numerics.Sqr(selfRating.Stddev()) + numerics.Sqr(gi.DynamicsFactor)
	stdDevMultiplier := varianceWithDynamics / numerics.Sqr(c)

	newMean := selfRating.Mean() + (rankMultiplier * meanMultiplier * v)
	newStdDev := math.Sqrt(varianceWithDynamics * (1 - w*stdDevMultiplier))

	return skills.NewRating(newMean, newStdDev)
}

// Calculates the match quality as the likelihood of all teams drawing (0% = bad, 100% = well matched).
func (calc *TwoPlayerCalc) CalcMatchQual(gi *skills.GameInfo, teams []skills.Team) float64 {
	ValidateTeamCount(teams, twoPlayerTeamRange)
	ValidatePlayersPerTeam(teams, twoPlayerPlayerRange)

	team1 := teams[0]
	player1 := team1.Players()[0]
	player1Rating := team1.PlayerRating(player1)

	team2 := teams[1]
	player2 := team2.Players()[0]
	player2Rating := team2.PlayerRating(player2)

	// We just use equation 4.1 found on page 8 of the TrueSkill 2006 paper:
	betaSquared := numerics.Sqr(gi.Beta)
	player1SigmaSquared := numerics.Sqr(player1Rating.Stddev())
	player2SigmaSquared := numerics.Sqr(player2Rating.Stddev())

	// This is the square root part of the equation:
	sqrtPart := math.Sqrt(2 * betaSquared / (2*betaSquared + player1SigmaSquared + player2SigmaSquared))

	// This is the exponent part of the equation:
	expPart := math.Exp((-1 * numerics.Sqr(player1Rating.Mean()-player2Rating.Mean())) / (2 * (2*betaSquared + player1SigmaSquared + player2SigmaSquared)))

	return sqrtPart * expPart
}

var (
	twoPlayerTeamRange   = numerics.Exactly(2)
	twoPlayerPlayerRange = numerics.Exactly(1)
)
