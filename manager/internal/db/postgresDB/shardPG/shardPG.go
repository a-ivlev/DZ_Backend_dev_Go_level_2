package shardPG

import (
	"errors"
	"sync"
)

var (
	ErrorShardNotFound = errors.New("shard not found")
)

type Shard struct {
	Address string
	Number  int
}

type Manager struct {
	size int
	ss   *sync.Map
}

func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
	}
}

func (m *Manager) Add(s *Shard) {
	m.ss.Store(s.Number, s)
}

func (m *Manager) ShardById(entityId int) (*Shard, error) {
	if entityId < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityId % m.size
	if s, ok := m.ss.Load(n); ok {
		return s.(*Shard), nil
	}

	return nil, ErrorShardNotFound
}
