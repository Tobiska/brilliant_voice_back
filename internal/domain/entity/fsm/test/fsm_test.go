package test

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/fsm/states"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/infrustucture/conn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedState(t *testing.T) {
	testCases := []struct {
		name          string
		expectedState fsm.IState
		initState     fsm.IState
	}{
		{
			name:          "DefaultTransition",
			initState:     &states.Created{},
			expectedState: &states.WaitStart{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := fsm.InitFsm(tc.initState, &game.Game{})
			assert.Equal(t, f.CurrentState(), tc.expectedState)
		})
	}
}

func TestRoundRunningState(t *testing.T) {
	testCases := []struct {
		name          string
		initState     fsm.IState
		expectedState fsm.IState
		actions       []fsm.IUserAction
		gameState     *game.Game
	}{
		{
			name:          "AllAnswers",
			initState:     &states.Ready{},
			expectedState: &states.WaitUsers{},
			gameState: TestGame("code",
				"admin_code",
				properties.Properties{CountPlayers: 2, TimerDuration: 5},
				&game.User{ID: "admin_code", Conn: &conn.MockConn{}},
				&game.User{ID: "test", Conn: &conn.MockConn{}},
			),
			actions: []fsm.IUserAction{
				actions.StartAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}), //for start timer
				actions.AnswerAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}, "text_admin"),
				actions.AnswerAction(&game.User{ID: "test", Conn: &conn.MockConn{}}, "text_test"),
			},
		},
		{
			name:          "NotAllAnswers",
			initState:     &states.Ready{},
			expectedState: &states.RoundRunning{},
			gameState: TestGame("code",
				"admin_code",
				properties.Properties{CountPlayers: 2, TimerDuration: 5},
				&game.User{ID: "admin_code", Conn: &conn.MockConn{}},
				&game.User{ID: "test", Conn: &conn.MockConn{}},
			),
			actions: []fsm.IUserAction{
				actions.StartAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
				actions.AnswerAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}, "text_admin"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := fsm.InitFsm(tc.initState, tc.gameState)
			for _, a := range tc.actions {
				assert.NoError(t, f.SendAction(a))
			}
			assert.Equal(t, tc.expectedState, f.CurrentState())
		})
	}
}

func TestReadyState(t *testing.T) {
	testCases := []struct {
		name          string
		initState     fsm.IState
		expectedState fsm.IState
		actions       []fsm.IUserAction
		gameState     *game.Game
		err           error
	}{
		{
			name:          "Start",
			initState:     &states.Ready{},
			expectedState: &states.RoundRunning{},
			actions: []fsm.IUserAction{
				actions.StartAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
		{
			name:          "Start",
			initState:     &states.Ready{},
			expectedState: &states.Ready{}, //todo RoundRunning
			actions: []fsm.IUserAction{
				actions.StartAction(&game.User{ID: "not_admin", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       states.ErrStartNotOwner,
		},
		{
			name:          "Sustain",
			initState:     &states.Ready{},
			expectedState: &states.WaitStart{},
			actions: []fsm.IUserAction{
				actions.LeaveUserAction(&game.User{ID: "test", Conn: &conn.MockConn{}}),
				actions.StartAction(&game.User{ID: "test", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := fsm.InitFsm(tc.initState, tc.gameState)
			for _, a := range tc.actions {
				if tc.err == nil {
					assert.NoError(t, f.SendAction(a))
				} else {
					assert.Error(t, f.SendAction(a))
				}
			}
			assert.Equal(t, tc.expectedState, f.CurrentState())
		})
	}
}

func TestWaitStartState(t *testing.T) {
	testCases := []struct {
		name          string
		initState     fsm.IState
		expectedState fsm.IState
		actions       []fsm.IUserAction
		gameState     *game.Game
		err           error
	}{
		{
			name:          "Sustain",
			initState:     &states.WaitStart{},
			expectedState: &states.WaitStart{},
			actions: []fsm.IUserAction{
				actions.StartAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
		{
			name:          "ContinueWait",
			initState:     &states.WaitStart{},
			expectedState: &states.WaitStart{},
			actions: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
		{
			name:          "WaitToReady",
			initState:     &states.WaitStart{},
			expectedState: &states.Ready{},
			actions: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
		{
			name:          "JoinJoinLeave",
			initState:     &states.WaitStart{},
			expectedState: &states.WaitStart{},
			actions: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),

				actions.LeaveUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
		{
			name:          "JoinJoinLeaveJoin",
			initState:     &states.WaitStart{},
			expectedState: &states.Ready{},
			actions: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),

				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),

				actions.LeaveUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),

				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       nil,
		},
		{
			name:          "AdminLeave",
			initState:     &states.WaitStart{},
			expectedState: &states.Dead{},
			actions: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),

				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),

				actions.LeaveUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       game.ErrOwnerLeave,
		},
		{
			name:          "OwnerNotJoined",
			initState:     &states.WaitStart{},
			expectedState: &states.WaitStart{},
			actions: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "not_admin", Conn: &conn.MockConn{}}),

				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			gameState: Game,
			err:       game.ErrOwnerYetNotJoined,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := fsm.InitFsm(tc.initState, tc.gameState)
			for _, a := range tc.actions {
				_ = f.SendAction(a)
			}
			assert.Equal(t, tc.expectedState, f.CurrentState())
		})
	}
}
