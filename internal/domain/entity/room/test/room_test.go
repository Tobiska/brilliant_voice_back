package test

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/logicTimer"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/room"
	"brillian_voice_back/internal/infrustucture/conn"
	"brillian_voice_back/internal/infrustucture/roundsProvider/inmemory"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	testCases := []struct {
		name                 string
		users                []*game.User
		expectedCountPlayers uint
		valid                bool
	}{
		{
			name: "JoinBeforeOwner",
			users: []*game.User{
				{ID: "test1", Conn: &conn.MockConn{}},
			},
			expectedCountPlayers: 0,
			valid:                false,
		},
		{
			name: "JoinOwner",
			users: []*game.User{
				{ID: "admin_code", Conn: &conn.MockConn{}},
				{ID: "test1", Conn: &conn.MockConn{}},
			},
			expectedCountPlayers: 2,
			valid:                true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rs, err := inmemory.NewRoundProvider().PrepareRounds(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			m := logicTimer.NewManager()
			r := room.NewRoom("code", "admin_code", properties.Properties{
				CountPlayers:  2,
				TimerDuration: 1,
			}, rs, m)
			for _, u := range tc.users {
				err := r.JoinToRoom(u)
				if tc.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestAsyncOneUser(t *testing.T) {
	c := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}
	u := game.NewUser("admin_code", "admin", c)

	rs, err := inmemory.NewRoundProvider().PrepareRounds(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	m := logicTimer.NewManager()
	r := room.NewRoom("code", "admin_code", properties.Properties{
		CountPlayers:  2,
		TimerDuration: 5,
	}, rs, m)
	r.Run()
	err = r.JoinToRoom(u)
	assert.NoError(t, err)

	r.ActionChannel() <- actions.AddUserAction(u)

	uErr := <-c.ErrCh
	assert.NotNil(t, uErr)
}

func TestAsyncManyUsers(t *testing.T) {
	cf := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}

	cs := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}
	uf := game.NewUser("admin_code", "admin1", cf)
	us := game.NewUser("admin_code", "admin2", cs)

	rs, err := inmemory.NewRoundProvider().PrepareRounds(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	m := logicTimer.NewManager()
	r := room.NewRoom("code", "admin_code", properties.Properties{
		CountPlayers:  2,
		TimerDuration: 5,
	}, rs, m)
	r.Run()
	err = r.JoinToRoom(uf)
	assert.NoError(t, err)

	r.ActionChannel() <- actions.AddUserAction(uf)
	r.ActionChannel() <- actions.AddUserAction(us)

	uErrF := <-cf.ErrCh
	uErrS := <-cs.ErrCh
	assert.NotNil(t, uErrF, uErrS)
}

func TestJoinUpdate(t *testing.T) {
	cf := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}

	cs := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}
	uf := game.NewUser("admin_code", "admin1", cf)
	us := game.NewUser("test", "test1", cs)

	rs, err := inmemory.NewRoundProvider().PrepareRounds(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	m := logicTimer.NewManager()
	r := room.NewRoom("code", "admin_code", properties.Properties{
		CountPlayers:  2,
		TimerDuration: 5,
	}, rs, m)
	r.Run()
	err = r.JoinToRoom(uf)
	assert.NoError(t, err)

	r.ActionChannel() <- actions.AddUserAction(us)
	assert.NotNil(t, <-cf.UpdateCh)

	r.ActionChannel() <- actions.LeaveUserAction(us)
	assert.Equal(t, game.ErrDone, <-cs.ErrCh)
	assert.NotNil(t, <-cf.UpdateCh)
}

func TestLogicTimer(t *testing.T) {
	cf := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}
	uf := game.NewUser("admin_code", "admin1", cf)

	rs, err := inmemory.NewRoundProvider().PrepareRounds(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	m := logicTimer.NewManager()
	r := room.NewRoom("code", "admin_code", properties.Properties{
		CountPlayers:  1,
		TimerDuration: 5,
	}, rs, m)
	r.Run()

	if err := r.JoinToRoom(uf); err != nil {
		t.Fatal(err)
	}

	r.ActionChannel() <- actions.StartAction(uf)

	assert.NotNil(t, <-cf.UpdateCh) //after start round
	assert.NotNil(t, <-cf.UpdateCh) //after timeout
}

func TestAnswers(t *testing.T) {
	cf := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}
	cs := &Conn{
		UpdateCh: make(chan game.Game),
		ErrCh:    make(chan error),
	}
	uf := game.NewUser("admin_code", "admin1", cf)
	us := game.NewUser("test", "admin1", cs)

	rs, err := inmemory.NewRoundProvider().PrepareRounds(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	m := logicTimer.NewManager()
	r := room.NewRoom("code", "admin_code", properties.Properties{
		CountPlayers:  2,
		TimerDuration: 5000000000,
	}, rs, m)
	r.Run()

	if err := r.JoinToRoom(uf); err != nil {
		t.Fatal(err)
	}

	if err := r.JoinToRoom(us); err != nil {
		t.Fatal(err)
	}

	r.ActionChannel() <- actions.StartAction(uf)
	r.ActionChannel() <- actions.AnswerAction(uf, "TextAdmin")
	r.ActionChannel() <- actions.AnswerAction(us, "TextTest")

	assert.NotNil(t, <-cf.UpdateCh) //after start round
	g1 := <-cf.UpdateCh
	assert.Equal(t, 1, len(g1.Rounds[g1.NumberOfRound].Answers))

	g2 := <-cf.UpdateCh
	assert.Equal(t, 2, len(g2.Rounds[g2.NumberOfRound].Answers))
	assert.Equal(t, "wait users", g2.Status)
}
