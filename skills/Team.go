package skills

type Team struct {
	PlayerRatings
}

func NewTeam() Team {
	return Team{make(PlayerRatings)}
}

func (t Team) AddPlayer(p Player, r Rating) {
	t.PlayerRatings[p] = r
}

func (t Team) PlayerCount() int {
	return len(t.PlayerRatings)
}

func (t Team) Players() []Player {
	ps := []Player{}
	for p := range t.PlayerRatings {
		ps = append(ps, p)
	}
	return ps
}

func (t Team) PlayerRating(p Player) Rating {
	return t.PlayerRatings[p]
}
