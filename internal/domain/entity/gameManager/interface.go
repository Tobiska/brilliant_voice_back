package gameManager

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

type IRoundProvider interface {
	CreateRound(context.Context) (game.Round, error)
}
