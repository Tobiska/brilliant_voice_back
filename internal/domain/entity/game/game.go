package game

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"context"
	"errors"
)

var (
	ErrUserAlreadyExist  = errors.New("game with id is already exist")
	ErrDeleteNoExist     = errors.New("game doesn't exist(deleting)")
	ErrOwnerLeave        = errors.New("owner leave")
	ErrOwnerYetNotJoined = errors.New("owner yet not joined")
)

type Users struct {
	cnt         map[string]*User
	countAnswer int
}

func NewUsersContainer() *Users {
	return &Users{
		cnt:         make(map[string]*User),
		countAnswer: 0,
	}
}

func (us Users) ToSlice() (sl []*User) {
	sl = make([]*User, len(us.cnt))
	var i uint
	for _, u := range us.cnt {
		sl[i] = u
		i++
	}
	return
}

func (us Users) Len() int {
	return len(us.cnt)
}

func (us Users) Delete(id string) error {
	if u, ex := us.cnt[id]; ex {
		closeCtx, _ := context.WithTimeout(context.Background(), 1)
		go u.Conn.Close(closeCtx)

		delete(us.cnt, id)
		return nil
	} else {
		return ErrDeleteNoExist
	}
}

func (us Users) Add(u *User) error {
	if _, ex := us.cnt[u.ID]; !ex {
		us.cnt[u.ID] = u
		return nil
	} else {
		return ErrUserAlreadyExist
	}
}

func (us Users) Answer(id, answer string) error {
	if u, ok := us.cnt[id]; ok {
		u.Answer = answer
		us.countAnswer++
		return nil
	} else {
		return ErrOwnerYetNotJoined
	}
}

func (us Users) CheckAnswers() bool {
	return len(us.cnt) == us.countAnswer
}

func (us Users) Clear() {
	for id := range us.cnt {
		_ = us.Delete(id)
	}
}

type Status string

type Game struct {
	Descriptor
	Users   *Users
	OwnerId string
	status  Status

	NumberOfRound int

	Timer ITimer

	Rounds []*Round
}

func (g *Game) GetOwner() (*User, error) {
	if u, ok := g.Users.cnt[g.OwnerId]; ok {
		return u, nil
	} else {
		return nil, ErrOwnerYetNotJoined
	}
}

func (g *Game) StartTimer() {
	g.Timer.Start(context.Background(), TimerInfo{
		TickerPeriod:  15000, //todo refactor
		TimeOutPeriod: g.Properties.TimerDuration,
	})
}

func (g *Game) StopTimer() {
	g.Timer.Stop(context.Background())
}

func (g *Game) DeleteUser(u *User) error {
	if u.ID == g.OwnerId {
		g.Users.Clear()
		return ErrOwnerLeave
	} else {
		if pu, ok := g.Users.cnt[u.ID]; ok {
			return g.Users.Delete(pu.ID)
		}
		return nil
	}
}

func (g *Game) AddUser(u *User) error {
	if u == nil {
		return errors.New("user is nil")
	}
	if _, err := g.GetOwner(); err != nil && u.ID != g.OwnerId {
		return err
	}
	return g.Users.Add(u)
}

type Descriptor struct {
	Code       string
	IsFully    bool
	Properties properties.Properties
	Status     string
}
