package idgen

import "sync"

type IDGenerator struct {
	mtx sync.Mutex
	id  int64
}

func (g *IDGenerator) GenerateID() int64 {
	g.mtx.Lock()

	g.id += 1
	id := g.id

	g.mtx.Unlock()
	return id
}
