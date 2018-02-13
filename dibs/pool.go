package dibs

import (
	"sync/atomic"
	"time"
)

type tag string
type clientID string

func (s *scheduler) selectBuilder(t tag, c clientID) *builder {

	defer s.mutex.Unlock()
	for {
		s.mutex.Lock()

		for _, b := range s.builders {
			if b.dedicatedTo == c {
				atomic.AddInt32(&b.openConnections, 1)
				b.lastestUse = time.Now().Unix()
				return b
			}
		}

		for _, b := range s.builders {
			if b.dedicatedTo == "" {
				b.dedicatedTo = c
				b.lastestUse = time.Now().Unix()
				atomic.AddInt32(&b.openConnections, 1)
				return b
			}
		}

		for _, b := range s.builders {
			if b.openConnections == 0 && b.lastestUse < time.Now().Unix()-10 {
				b.dedicatedTo = c
				b.lastestUse = time.Now().Unix()
				atomic.AddInt32(&b.openConnections, 1)
				return b
			}
		}

		s.mutex.Unlock()
		time.Sleep(time.Second)
	}
}

func (b *builder) Close() {
	b.lastestUse = time.Now().Unix()
	atomic.AddInt32(&b.openConnections, -1)
}
