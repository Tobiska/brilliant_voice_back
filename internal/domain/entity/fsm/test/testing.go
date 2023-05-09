package test

import (
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/infrustucture/conn"
	"brillian_voice_back/internal/infrustucture/roundsProvider/inmemory"
	"context"
)

var (
	Game = TestGame("code",
		"admin_code",
		game.Properties{CountPlayers: 2, TimerDuration: 5},
		&game.User{ID: "admin_code", Conn: &conn.MockConn{}},
		&game.User{ID: "test", Conn: &conn.MockConn{}},
	)
)

type MockTimeAdapter struct {
	ch           chan game.TimerInfo
	timerRunning bool
}

func (a *MockTimeAdapter) Send(_ context.Context, _ game.TimerInfo) error {
	a.timerRunning = true
	return nil
}

func (a *MockTimeAdapter) CheckTimer() bool {
	return a.timerRunning
}

func TestGame(code, ownerId string, prop game.Properties, us ...*game.User) *game.Game {
	usersContainer := game.NewUsersContainer()
	for _, u := range us {
		if err := usersContainer.Add(u); err != nil {
			return nil
		}
	}
	rounds, _ := inmemory.NewRoundProvider().PrepareRounds(context.Background())
	return &game.Game{
		Users:   usersContainer,
		OwnerId: ownerId,
		Timer:   &MockTimeAdapter{},
		Rounds:  rounds,
		Descriptor: game.Descriptor{
			Code:       code,
			Properties: prop,
		},
	}
}
