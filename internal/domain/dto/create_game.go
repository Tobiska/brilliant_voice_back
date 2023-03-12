package dto

import "brillian_voice_back/internal/domain/entity/room"

type InputCreateGameDto struct {
	Username          string
	ID                string
	CountPlayers      int
	TimeDurationRound int
}

type InputJoinGameDto struct {
	Username string
	ID       string
	Code     string
	Room     *room.Room
}
