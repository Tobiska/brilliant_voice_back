package gameManager

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/user"
)

type GameStatus string

var (
	WaitStart GameStatus = "WAIT_START"
	WaitUsers GameStatus = "WAIT_USERS"
	Running   GameStatus = "RUNNING"
	Pause     GameStatus = "PAUSE"
	Dead      GameStatus = "DEAD"
)

type GameState struct {
	Descriptor
	Users           map[string]*user.User
	OwnerId         string
	CurrentQuestion string
	status          GameStatus
	NumberOfRound   int
}

func (s *GameState) GetOwner() (*user.User, error) {
	if u, ok := s.Users[s.OwnerId]; ok {
		return u, nil
	} else {
		return nil, ErrOwnerYetNotJoined
	}
}

func (s *GameState) Status() GameStatus {
	return s.status
}

type Descriptor struct {
	Code       string
	IsFully    bool
	Properties properties.Properties
	Status     string
}
