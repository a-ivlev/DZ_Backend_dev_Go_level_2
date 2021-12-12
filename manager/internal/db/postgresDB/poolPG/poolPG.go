package poolPG

import (
	"database/sql"
	"sync"
)

type Pool struct {
	sync.RWMutex
	cc map[string]*sql.DB
}

func NewPool() *Pool {
	return &Pool{
		cc: map[string]*sql.DB{},
	}
}

func (p *Pool) Connection(addr string) (*sql.DB, error) {
	p.RLock()
	if c, ok := p.cc[addr]; ok {
		defer p.RUnlock()
		return c, nil
	}
	p.RUnlock()
	p.Lock()
	defer p.Unlock()
	if c, ok := p.cc[addr]; ok {
		return c, nil
	}
	var err error
	p.cc[addr], err = sql.Open("postgres", addr)
	return p.cc[addr], err
}
