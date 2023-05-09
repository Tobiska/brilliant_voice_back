package roomProvider

import "context"

type ICodeRoomProvider interface {
	Generate(ctx context.Context) string
}
