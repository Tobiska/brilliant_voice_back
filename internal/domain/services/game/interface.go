package game

import (
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/room"
	"context"
)

type IGameProvider interface {
	FindRoom(ctx context.Context, code string) (*room.Room, error)
	CreateRoom(ctx context.Context, ownerID string, properties game.Properties) (*room.Room, error)
}
