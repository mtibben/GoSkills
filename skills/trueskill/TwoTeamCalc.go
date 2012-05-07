package trueskill

import (
	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/numerics"
	"math"
	"sort"
)

// Calculates new ratings for only two teams where each team has 1 or more players.
// When you only have two teams, the math is still simple: no factor graphs are used yet.
type TwoTeamCalc struct{}

// Calculates new ratings based on the prior ratings and team ranks use 1 for first place, repeat the number for a tie (e.g. 1, 2, 2).
func (calc *TwoTeamCalc) CalcNewRatings(gi *skills.GameInfo, teams []skills.Team, ranks ...int) skills.PlayerRatings {
	newSkills := make(map[skills.Player]skills.Rating)

	// Basic argument checking
	ValidateTeamCountAndPlayersCountPerTeam(teams, twoTeamTeamRange, twoTeamPlayerRange)

	// Make sure things are in order
	sort.Sort(skills.NewRankedTeams(teams, ranks))

	winningTeam := teams[0]
	losingTeam := teams[1]

	wasDraw := ranks[0] == ranks[1]

	twoTeamUpdateRatings(gi, newSkills, winningTeam, losingTeam, cond(wasDraw, skills.Draw, skills.Win))
	twoTeamUpdateRatings(gi, newSkills, losingTeam, winningTeam, cond(wasDraw, skills.Draw, skills.Lose))

	return newSkills
}

func twoTeamUpdateRatings(gi *skills.GameInfo, newSkills skills.PlayerRatings, selfTeam, otherTeam skills.Team, comparison int) {
	drawMargin := DrawMarginFromDrawProbability(gi.DrawProbability, gi.Beta)
	betaSqr := numerics.Sqr(gi.Beta)
	tauSqr := numerics.Sqr(gi.DynamicsFactor)

	totalPlayers := selfTeam.PlayerCount() + otherTeam.PlayerCount()

	selfMeanSum := selfTeam.Accum(skills.MeanSum)
	otherMeanSum := otherTeam.Accum(skills.MeanSum)

	c := math.Sqrt(selfTeam.Accum(skills.VarianceSum) + otherTeam.Accum(skills.VarianceSum) + float64(totalPlayers)*betaSqr)

	winningMean := selfMeanSum
	losingMean := otherMeanSum

	if comparison == skills.Lose {
		winningMean, losingMean = losingMean, winningMean
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

	for p, r := range selfTeam.PlayerRatings {
		prevPlayerRating := r

		meanMultiplier := (prevPlayerRating.Variance() + tauSqr) / c
		stdDevMultiplier := (prevPlayerRating.Variance() + tauSqr) / numerics.Sqr(c)

		playerMeanDelta := rankMultiplier * meanMultiplier * v
		newMean := prevPlayerRating.Mean() + playerMeanDelta

		newStdDev := math.Sqrt((prevPlayerRating.Variance() + tauSqr) * (1 - w*stdDevMultiplier))

		newSkills[p] = skills.NewRating(newMean, newStdDev)
	}
}

// Calculates the match quality as the likelihood of all teams drawing (0% = bad, 100% = well matched).
func (calc *TwoTeamCalc) CalcMatchQual(gi *skills.GameInfo, teams []skills.Team) float64 {
	// Basic argument checking
	ValidateTeamCountAndPlayersCountPerTeam(teams, twoTeamTeamRange, twoTeamPlayerRange)

	// We've verified that there's just two teams
	team1 := teams[0]
	team1Count := team1.PlayerCount()

	team2 := teams[1]
	team2Count := team2.PlayerCount()

	totalPlayers := team1Count + team2Count

	betaSqr := numerics.Sqr(gi.Beta)

	team1MeanSum := team1.Accum(skills.MeanSum)
	team1VarSum := team1.Accum(skills.VarianceSum)

	team2MeanSum := team2.Accum(skills.MeanSum)
	team2VarSum := team2.Accum(skills.VarianceSum)

	// This comes from equation 4.1 in the TrueSkill paper on page 8            
	// The equation was broken up into the part under the square root sign and 
	// the exponential part to make the code easier to read.

	betaSqrPlayers := betaSqr * float64(totalPlayers)

	sqrtPart := math.Sqrt(betaSqrPlayers / (betaSqrPlayers + team1VarSum + team2VarSum))
	expPart := math.Exp(-.5 * numerics.Sqr(team1MeanSum-team2MeanSum) / (betaSqrPlayers + team1VarSum + team2VarSum))

	return expPart * sqrtPart
}

var (
	twoTeamTeamRange   = numerics.Exactly(2)
	twoTeamPlayerRange = numerics.AtLeast(1)
)
