package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
	"time"
)

type Tick struct { //todo timeout has no any user
	Action
	RestOfTime time.Duration
}

func (t Tick) String() string {
	return fmt.Sprintf("Tick %s", t.U)
}

func TickAction(u *game.User, rest time.Duration) Tick {
	return Tick{
		Action{
			U: u,
		},
		rest,
	}
}
