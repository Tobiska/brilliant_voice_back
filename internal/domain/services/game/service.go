package game

import (
	"brillian_voice_back/internal/domain/dto"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/room"
	"context"
)

type Service struct {
	provider IGameProvider
}

func NewGameService(provider IGameProvider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) CreateRoom(ctx context.Context, input *dto.InputCreateGameDto) (*room.Room, error) {
	return s.provider.CreateRoom(ctx, input.ID, input.Username, properties.Properties{
		CountPlayers:  input.CountPlayers,
		TimerDuration: input.TimeDurationRound,
	})
}

func (s *Service) JoinToRoom(ctx context.Context, input *dto.InputJoinGameDto) (*room.Room, error) {
	return s.provider.FindRoom(ctx, input.Code)
}
