package fsm

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/user"
	"errors"
)

var (
	ErrUserAlreadyExist = errors.New("user with id is already exist")
)

type GameStatus string

type Game struct {
	Descriptor
	Users           map[string]*user.User
	OwnerId         string
	CurrentQuestion string
	status          GameStatus
	NumberOfRound   int
}

func (g *Game) GetOwner() (*user.User, error) {
	if u, ok := g.Users[g.OwnerId]; ok {
		return u, nil
	} else {
		return nil, ErrOwnerYetNotJoined
	}
}

func (g *Game) DeleteUser(u *user.User) error {
	if u.ID == g.OwnerId {
		g.Clear()
		return ErrOwnerLeave
	} else {
		if pu, ok := g.Users[u.ID]; ok {
			delete(g.Users, u.ID)
			_ = pu.DeleteAndClose()
		}
		return nil
	}
}

func (g *Game) Clear() {
	for _, u := range g.Users {
		_ = u.DeleteAndClose()
	}
}

func (g *Game) AddUser(u *user.User) error {
	if _, err := g.GetOwner(); err != nil && u.ID != g.OwnerId {
		return err
	}
	if _, ex := g.Users[u.ID]; ex {
		return ErrUserAlreadyExist
	} else {
		g.Users[u.ID] = u
		return nil
	}
}

func (g *Game) Status() GameStatus {
	return g.status
}

type Descriptor struct {
	Code       string
	IsFully    bool
	Properties properties.Properties
	Status     string
}
