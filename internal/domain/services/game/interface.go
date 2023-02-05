package game

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/room"
	"context"
)

type IGameProvider interface {
	FindRoom(ctx context.Context, code string) (*room.Room, error)
	CreateRoom(ctx context.Context, ownerName, ownerID string, properties properties.Properties) (*room.Room, error)
}
