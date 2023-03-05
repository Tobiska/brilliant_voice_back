package game

import (
	"brillian_voice_back/internal/domain/dto"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/room"
	"context"
	"github.com/rs/zerolog/log"
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
	r, err := s.provider.CreateRoom(ctx, input.ID, properties.Properties{
		CountPlayers:  input.CountPlayers,
		TimerDuration: input.TimeDurationRound,
	})
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	log.Info().Str("room_code", r.Desc().Code).Msg("room has created")
	return r, nil
}

func (s *Service) JoinToRoom(ctx context.Context, input *dto.InputJoinGameDto) (*room.Room, error) {
	return s.provider.FindRoom(ctx, input.Code)
}
