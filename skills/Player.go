package skills

import (
	"fmt"
)

type Identifier interface {
}

type Player struct {
	id Identifier
}

func NewPlayer(id Identifier) *Player {
	return &Player{
		id: id,
	}
}

func (p Player) String() string {
	return fmt.Sprintf("%v", p.id)
}
