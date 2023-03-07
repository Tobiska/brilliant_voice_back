package roomProvider

import (
	"brillian_voice_back/internal/domain/entity/gameManager"
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/room"
	"context"
	"errors"
	"sync"
)

var (
	ErrRoomNotFound       = errors.New("room not found")
	ErrTooManyRooms       = errors.New("too many rooms")
	ErrRoomIsNotAvailable = errors.New("room is not available")
)

type Provider struct {
	roundProvider gameManager.IRoundProvider
	codeProvider  ICodeRoomProvider
	storage       map[string]*room.Room
	mu            *sync.RWMutex

	limit int
}

func NewProvider(cp ICodeRoomProvider, limit int, roundProvider gameManager.IRoundProvider) *Provider {
	return &Provider{
		storage:       make(map[string]*room.Room, 0),
		codeProvider:  cp,
		mu:            &sync.RWMutex{},
		limit:         limit,
		roundProvider: roundProvider,
	}
}

func (p *Provider) CreateRoom(ctx context.Context, ownerID string, properties properties.Properties) (*room.Room, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.storage) >= p.limit {
		return nil, ErrTooManyRooms
	}
	for {
		code := p.codeProvider.Generate(ctx)
		if _, ok := p.storage[code]; !ok {
			p.storage[code] = room.NewRoom(code, ownerID, properties, p.roundProvider)
			return p.storage[code], nil
		}
	}
}

func (p *Provider) FindRoom(_ context.Context, code string) (*room.Room, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	r, ok := p.storage[code]
	if !ok {
		return nil, ErrRoomNotFound
	}
	if !filterRoom(r) {
		return nil, ErrRoomIsNotAvailable
	}
	return r, nil
}

func filterRoom(r *room.Room) bool {
	if r.Desc().Status == "dead" {
		return false
	}
	return true
}

func (p *Provider) gc() {
	panic("implement storage dead rooms collect")
}
