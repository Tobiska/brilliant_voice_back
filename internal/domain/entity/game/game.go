package game

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"errors"
)

var (
	ErrUserAlreadyExist  = errors.New("game with id is already exist")
	ErrDeleteNoExist     = errors.New("game doesn't exist(deleting)")
	ErrOwnerLeave        = errors.New("owner leave")
	ErrOwnerYetNotJoined = errors.New("owner yet not joined")
)

type Users map[string]*User

func (us Users) Delete(id string) error {
	if u, ex := us[id]; ex {
		u.Conn.Close()
		delete(us, id)
		return nil
	} else {
		return ErrDeleteNoExist
	}
}

func (us Users) Add(u *User) error {
	if _, ex := us[u.ID]; !ex {
		us[u.ID] = u
		return nil
	} else {
		return ErrUserAlreadyExist
	}
}

func (us Users) Clear() {
	for id := range us {
		_ = us.Delete(id)
	}
}

type GameStatus string

type Game struct {
	Descriptor
	Users           Users
	OwnerId         string
	CurrentQuestion string
	status          GameStatus
	NumberOfRound   int
}

func (g *Game) GetOwner() (*User, error) {
	if u, ok := g.Users[g.OwnerId]; ok {
		return u, nil
	} else {
		return nil, ErrOwnerYetNotJoined
	}
}

func (g *Game) DeleteUser(u *User) error {
	if u.ID == g.OwnerId {
		g.Users.Clear()
		return ErrOwnerLeave
	} else {
		if pu, ok := g.Users[u.ID]; ok {
			return g.Users.Delete(pu.ID)
		}
		return nil
	}
}

func (g *Game) AddUser(u *User) error {
	if _, err := g.GetOwner(); err != nil && u.ID != g.OwnerId {
		return err
	}
	return g.Users.Add(u)
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
