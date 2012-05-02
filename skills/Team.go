package skills

import (
	"fmt"
)

type Team struct {
	ratingMap map[Player]Rating
}

func NewTeam(p Player, r Rating) (t *Team) {
	t = &Team{}
	t.ratingMap = make(map[Player]Rating)
	t.ratingMap[p] = r
	return
}

func (t Team) String() string {
	return fmt.Sprintf("%v", t.ratingMap)
}

func (t *Team) PlayerCount() int {
	return len(t.ratingMap)
}

func (t *Team) Players() (ps []Player) {
	for p := range t.ratingMap {
		ps = append(ps, p)
	}
	return
}

func (t *Team) PlayerRating(p Player) Rating {
	return t.ratingMap[p]
}
