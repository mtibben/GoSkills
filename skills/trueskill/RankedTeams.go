package trueskill

import (
	"fmt"
	"github.com/ChrisHines/GoSkills/skills"
)

type rankedTeams struct {
	teams []skills.Team
	ranks []int
}

func newRankedTeams(teams []skills.Team, ranks []int) *rankedTeams {
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
