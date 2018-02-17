package main

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

type tag string
type clientID string

const reservation = 10 * time.Second

func (s *scheduler) reselect(t tag, c clientID) (*builder, bool) {
	for _, b := range s.builders {
		if b.dedicatedTo != c {
			continue
		}
		atomic.AddInt32(&b.openConnections, 1)
		if b.dedicatedTo != c {
			atomic.AddInt32(&b.openConnections, -1)
			continue
		}
		b.lastestUse = time.Now().Unix()
		return b, true
	}
	return nil, false
}

func (s *scheduler) findScheduleable(t tag, c clientID) (*builder, bool) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(s.builders)) {
		b := s.builders[i]
		if !scheduleable(b) {
			continue
		}
		b.dedicatedTo = c
		b.lastestUse = time.Now().Unix()
		atomic.AddInt32(&b.openConnections, 1)
		return b, true
	}
	return nil, false
}

func scheduleable(b *builder) bool {
	if b.dedicatedTo == "" {
		return true
	}
	if b.openConnections == 0 && b.lastestUse <= time.Now().Add(reservation).Unix() {
		return true
	}
	return false
}

func (s *scheduler) recycle(b *builder) {
	t := time.Now().Unix()
	b.lastestUse = t
	o := atomic.AddInt32(&b.openConnections, -1)
	if o != 0 {
		return
	}
	go func(b *builder, t int64) {
		time.Sleep(reservation)
		if b.lastestUse != t {
			return
		}
		log.Printf("recycled worker %s\n", b.name)
		s.mutex.Lock()
		s.cond.Broadcast()
		s.mutex.Unlock()
	}(b, t)
}

func (s *scheduler) selectWorker(t tag, c clientID) *builder {

	for {

		b, found := s.reselect(t, c)
		if found {
			log.Printf("reselected worker %s for client %s (tag %s)\n", b.name, c, t)
			return b
		}

		s.mutex.Lock()

		b, found = s.findScheduleable(t, c)
		if found {
			log.Printf("scheduled worker %s for client %s (tag %s)\n", b.name, c, t)
			s.mutex.Unlock()
			return b
		}

		log.Printf("waiting for worker to become free\n")
		s.cond.Wait()
	}
}
