package provider

import (
	"brillian_voice_back/internal/domain/entity/properties"
	"brillian_voice_back/internal/domain/entity/room"
	"context"
	"errors"
	"sync"
)

var (
	ErrRoomNotFound = errors.New("room not found")
)

type Provider struct {
	storage map[string]*room.Room
	mu      *sync.RWMutex
}

func NewProvider() *Provider {
	return &Provider{
		storage: make(map[string]*room.Room, 0),
		mu:      &sync.RWMutex{},
	}
}

func (p *Provider) CreateRoom(ctx context.Context, ownerID string, properties properties.Properties) (*room.Room, error) {
	//todo generate code
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.storage[ownerID]; !ok {
		p.storage[ownerID] = room.NewRoom(ownerID, ownerID, properties)
	}
	return p.storage[ownerID], nil
}

func (p *Provider) FindRoom(ctx context.Context, code string) (*room.Room, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	r, ok := p.storage[code]
	if !ok {
		return nil, ErrRoomNotFound
	}
	return r, nil
}

func (p *Provider) gc() {
	panic("implement storage dead rooms collect")
}
