package trueskill

import (
	"math"
	"sort"

	"github.com/ChrisHines/GoSkills/skills"
	"github.com/ChrisHines/GoSkills/skills/numerics"
)

var (
	factorGraphTeamRange   = numerics.AtLeast(2)
	factorGraphPlayerRange = numerics.AtLeast(1)
)

// FactorGraphTrueSkillCalc calculates TrueSkill using a full factor graph.
type FactorGraphTrueSkillCalc struct{}

func (calc *FactorGraphTrueSkillCalc) CalcNewRatings(gi *skills.GameInfo, teams []skills.Team, ranks ...int) skills.PlayerRatings {
	newSkills := make(map[skills.Player]skills.Rating)

	// Basic argument checking
	validateTeamCount(teams, factorGraphTeamRange)
	validatePlayersPerTeam(teams, factorGraphPlayerRange)

	// Copy slices so we don't confuse the client code
	steams := append([]skills.Team{}, teams...)
	sranks := append([]int{}, ranks...)

	// Make sure things are in order
	sort.Sort(skills.NewRankedTeams(steams, sranks))

	factorGraph := &TrueSkillFactorGraph{gameInfo, steams, sranks}
	factorGraph.BuildGraph()
	factorGraph.RunSchedule()

	probabilityOfOutcome := factorGraph.GetProbabilityOfRanking()

	return factorGraph.GetUpdatedRatings()
}

func (calc *FactorGraphTrueSkillCalc) CalculateMatchQuality(gi *skills.GameInfo, teams []skills.Team) float64 {

	// We need to create the A matrix which is the player team assigments.
	teamAssignmentsList := teams.ToList()
	skillsMatrix := GetPlayerCovarianceMatrix(teamAssignmentsList)
	meanVector := GetPlayerMeansVector(teamAssignmentsList)
	meanVectorTranspose := meanVector.Transpose

	playerTeamAssignmentsMatrix := CreatePlayerTeamAssignmentMatrix(teamAssignmentsList, meanVector.Rows)
	playerTeamAssignmentsMatrixTranspose := playerTeamAssignmentsMatrix.Transpose

	betaSquared := Square(gameInfo.Beta)

	start := meanVectorTranspose * playerTeamAssignmentsMatrix
	aTa := (betaSquared * playerTeamAssignmentsMatrixTranspose) * playerTeamAssignmentsMatrix
	aTSA := playerTeamAssignmentsMatrixTranspose * skillsMatrix * playerTeamAssignmentsMatrix
	middle := aTa + aTSA

	middleInverse := middle.Inverse

	end := playerTeamAssignmentsMatrixTranspose * meanVector

	expPartMatrix := -0.5 * (start * middleInverse * end)
	expPart := expPartMatrix.Determinant

	sqrtPartNumerator := aTa.Determinant
	sqrtPartDenominator := middle.Determinant
	sqrtPart := sqrtPartNumerator / sqrtPartDenominator

	result := math.Exp(expPart) * math.Sqrt(sqrtPart)

	return result
}

func GetPlayerMeansVector(teamAssignmentsList skills.PlayerRatings) Vector {
	// A simple vector of all the player means.
	return Vector(GetPlayerRatingValues(teamAssignmentsList, func(rating skills.Rating) float64 { return rating.Mean() }))
}

func GetPlayerCovarianceMatrix(teamAssignmentsList skills.PlayerRatings) Matrix {
	// This is a square matrix whose diagonal values represent the variance (square of standard deviation) of all players.
	return DiagonalMatrix(GetPlayerRatingValues(teamAssignmentsList, func(rating skills.Rating) float64 { Square(rating.StandardDeviation) }))
}

// Helper function that gets a list of values for all player ratings
func GetPlayerRatingValues(teamAssignmentsList skills.PlayerRatings, playerRatingFunction func(skills.Rating) float64) []float64 {

	playerRatingValues := []float64{}

	for currentTeam := range teamAssignmentsList {
		for currentRating := range currentTeam.Values {
			playerRatingValues = append(playerRatingValues, playerRatingFunction(currentRating))
		}
	}

	return playerRatingValues
}

func CreatePlayerTeamAssignmentMatrix(teamAssignmentsList skills.PlayerRatings, totalPlayers int) Matrix {

	// The team assignment matrix is often referred to as the "A" matrix. It's a matrix whose rows represent the players
	// and the columns represent teams. At Matrix[row, column] represents that player[row] is on team[col]
	// Positive values represent an assignment and a negative value means that we subtract the value of the next
	// team since we're dealing with pairs. This means that this matrix always has teams - 1 columns.
	// The only other tricky thing is that values represent the play percentage.

	// For example, consider a 3 team game where team1 is just player1, team 2 is player 2 and player 3, and
	// team3 is just player 4. Furthermore, player 2 and player 3 on team 2 played 25% and 75% of the time
	// (e.g. partial play), the A matrix would be:

	// A = this 4x2 matrix:
	// |  1.00  0.00 |
	// | -0.25  0.25 |
	// | -0.75  0.75 |
	// |  0.00 -1.00 |

	playerAssignments := map[float64]float64{}
	totalPreviousPlayers := 0

	for i := 0; i < teamAssignmentsList.Count-1; i++ {

		currentTeam := teamAssignmentsList[i]

		// Need to add in 0's for all the previous players, since they're not
		// on this team
		currentRowValues := make([]float64{}, totalPreviousPlayers)
		for currentRating := range currentTeam {
			currentRowValues = append(currentRowValues,
				PartialPlay.GetPartialPlayPercentage(currentRating.Key))
			// indicates the player is on the team
			totalPreviousPlayers++
		}

		nextTeam := teamAssignmentsList[i+1]

		for nextTeamPlayerPair := range nextTeam {
			// Add a -1 * playing time to represent the difference
			currentRowValues = append(currentRowValues,
				-1*PartialPlay.GetPartialPlayPercentage(nextTeamPlayerPair.Key))
		}

		playerAssignments = append(playerAssignments, currentRowValues)
	}

	playerTeamAssignmentsMatrix := Matrix(totalPlayers, teamAssignmentsList.Count-1, playerAssignments)

	return playerTeamAssignmentsMatrix
}
