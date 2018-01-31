package dibs

import (
	"log"
	"sync"
)

type pool struct {
	freeBuilders []*builder
	mutex        *sync.Mutex
}

func (p *pool) selectBuilder(t tag) *builder {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var b *builder
	log.Printf("builders: %v\n", p.freeBuilders)
	b, p.freeBuilders = p.freeBuilders[0], p.freeBuilders[1:]
	return b
}

func (p *pool) recycleBuilder(b *builder) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.freeBuilders = append(p.freeBuilders, b)
}
