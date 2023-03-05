package test

import (
	"brillian_voice_back/internal/domain/entity/actions"
	"brillian_voice_back/internal/domain/entity/fsm"
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/gameManager"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/infrustucture/conn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFsmStress(t *testing.T) {
	testDependCases := []struct {
		name          string
		sendAction    []fsm.IUserAction
		exceptedState string
	}{
		{
			name: "Created/WaitStart",
			sendAction: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			exceptedState: "ready",
		},
		{
			name: "ReadyLeave",
			sendAction: []fsm.IUserAction{
				actions.LeaveUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			exceptedState: "wait start",
		},
		{
			name: "Ready",
			sendAction: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "test1", Conn: &conn.MockConn{}}),
			},
			exceptedState: "ready",
		},
		{
			name: "Ready",
			sendAction: []fsm.IUserAction{
				actions.AddUserAction(&game.User{ID: "admin_code", Conn: &conn.MockConn{}}),
			},
			exceptedState: "ready", //todo RunningRound
		},
	}
	m := gameManager.NewManager("test", "admin_code", properties.Properties{
		CountPlayers:  2,
		TimerDuration: 2,
	})
	for _, tc := range testDependCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, a := range tc.sendAction {
				assert.NoError(t, m.DoSync(a))
			}
			assert.Equal(t, tc.exceptedState, m.Game().Status)
		})
	}
}
