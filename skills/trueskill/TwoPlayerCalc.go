package trueskill

import (
	"fmt"
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
func (calc *TwoPlayerCalc) CalcNewRatings(gi *skills.GameInfo, teams []*skills.Team, ranks ...int) map[skills.Player]skills.Rating {
	newSkills := make(map[skills.Player]skills.Rating)

	// Basic argument checking
	ValidateTeamCountAndPlayersCountPerTeam(teams, twoPlayerTeamRange, twoPlayerPlayerRange)

	// Make sure things are in order
	sort.Sort(newRankedTeams(teams, ranks))

	// Since we verified that each team has one player, we know the player is the first one
	winningTeam := teams[0]
	winner := winningTeam.Players()[0]
	winnerPrevRating := winningTeam.PlayerRating(winner)

	losingTeam := teams[1]
	loser := losingTeam.Players()[0]
	loserPrevRating := losingTeam.PlayerRating(loser)

	wasDraw := ranks[0] == ranks[1]

	newSkills[winner] = CalculateNewRating(gi, winnerPrevRating, loserPrevRating, cond(wasDraw, skills.Draw, skills.Win))
	newSkills[loser] = CalculateNewRating(gi, loserPrevRating, winnerPrevRating, cond(wasDraw, skills.Draw, skills.Lose))

	return newSkills
}

func CalculateNewRating(gi *skills.GameInfo, selfRating, oppRating skills.Rating, comparison int) skills.Rating {
	drawMargin := DrawMarginFromDrawProbability(gi.DrawProbability, gi.Beta)

	c := math.Sqrt(sqr(selfRating.Stddev) + sqr(oppRating.Stddev) + 2*sqr(gi.Beta))

	winningMean := selfRating.Mean
	losingMean := oppRating.Mean

	if comparison == skills.Lose {
		winningMean = oppRating.Mean
		losingMean = selfRating.Mean
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

	meanMultiplier := (sqr(selfRating.Stddev) + sqr(gi.DynamicsFactor)) / c

	varianceWithDynamics := sqr(selfRating.Stddev) + sqr(gi.DynamicsFactor)
	stdDevMultiplier := varianceWithDynamics / sqr(c)

	newMean := selfRating.Mean + (rankMultiplier * meanMultiplier * v)
	newStdDev := math.Sqrt(varianceWithDynamics * (1 - w*stdDevMultiplier))

	return skills.NewRating(newMean, newStdDev)
}

// Calculates the match quality as the likelihood of all teams drawing (0% = bad, 100% = well matched).
func (calc *TwoPlayerCalc) CalcMatchQual(gi *skills.GameInfo, teams []*skills.Team) float64 {
	ValidateTeamCountAndPlayersCountPerTeam(teams, twoPlayerTeamRange, twoPlayerPlayerRange)

	team1 := teams[0]
	player1 := team1.Players()[0]
	player1Rating := team1.PlayerRating(player1)

	team2 := teams[1]
	player2 := team2.Players()[0]
	player2Rating := team2.PlayerRating(player2)

	// We just use equation 4.1 found on page 8 of the TrueSkill 2006 paper:
	betaSquared := sqr(gi.Beta)
	player1SigmaSquared := sqr(player1Rating.Stddev)
	player2SigmaSquared := sqr(player2Rating.Stddev)

	// This is the square root part of the equation:
	sqrtPart := math.Sqrt(2 * betaSquared / (2*betaSquared + player1SigmaSquared + player2SigmaSquared))

	// This is the exponent part of the equation:
	expPart := math.Exp((-1 * sqr(player1Rating.Mean-player2Rating.Mean)) / (2 * (2*betaSquared + player1SigmaSquared + player2SigmaSquared)))

	return sqrtPart * expPart
}

func ValidateTeamCountAndPlayersCountPerTeam(teams []*skills.Team, teamsAllowed, playersAllowed numerics.Range) {
	if n := len(teams); !teamsAllowed.In(n) {
		panic(fmt.Errorf("len(teams) [%v] outside of expected range [%v]", n, teamsAllowed))
	}
	for _, t := range teams {
		if n := t.PlayerCount(); !playersAllowed.In(n) {
			panic(fmt.Errorf("PlayerCount [%v] outside of expected range [%v]", n, playersAllowed))
		}
	}
}

var (
	twoPlayerTeamRange   = numerics.Exactly(2)
	twoPlayerPlayerRange = numerics.Exactly(1)
)

type rankedTeams struct {
	teams []*skills.Team
	ranks []int
}

func newRankedTeams(teams []*skills.Team, ranks []int) *rankedTeams {
	if len(teams) != len(ranks) {
		panic(fmt.Errorf("Number of teams [%v] does not match number of ranks [%v]", len(teams), len(ranks)))
	}
	return &rankedTeams{teams, ranks}
}

func (rt *rankedTeams) Len() int           { return len(rt.teams) }
func (rt *rankedTeams) Less(i, j int) bool { return rt.ranks[i] < rt.ranks[j] }

func (rt *rankedTeams) Swap(i, j int) {
	rt.teams[i], rt.teams[j] = rt.teams[j], rt.teams[i]
	rt.ranks[i], rt.ranks[j] = rt.ranks[j], rt.ranks[i]
}

func sqr(x float64) float64 { return x * x }

func cond(c bool, t, f int) int {
	if c {
		return t
	}
	return f
}
