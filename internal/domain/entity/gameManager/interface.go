package gameManager

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

type IRoundProvider interface { //todo может быть убрать в другое место
	PrepareRounds(context.Context) ([]*game.Round, error)
}
