package conn

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

type MockConn struct{}

func (mc *MockConn) UpdateGame(_ context.Context, _ game.Game) {}
func (mc *MockConn) SendError(_ context.Context, _ error)      {}
func (mc *MockConn) Close(_ context.Context)                   {}
