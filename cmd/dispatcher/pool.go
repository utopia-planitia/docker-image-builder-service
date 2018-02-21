package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type clientID string

const reservation = 10 * time.Second

func (s *dispatcher) reselect(c clientID) (*builder, bool) {
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

func (s *dispatcher) findScheduleable(c clientID, t, cf string) (*builder, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var selectableBuilders []*builder
	for _, i := range r.Perm(len(s.builders)) {
		b := s.builders[i]
		if !scheduleable(b) {
			continue
		}
		selectableBuilders = append(selectableBuilders, b)
	}
	if len(selectableBuilders) == 0 {
		return nil, nil
	}

	b, err := selectByUncachedBytes(selectableBuilders, t, cf)
	if err != nil {
		return nil, fmt.Errorf("failed to select builder on uncached bytes size: %s", err)
	}

	b.dedicatedTo = c
	b.lastestUse = time.Now().Unix()
	atomic.AddInt32(&b.openConnections, 1)
	return b, nil
}

func selectByUncachedBytes(bs []*builder, t, cf string) (*builder, error) {
	if t == "" {
		return bs[0], nil
	}

	var smallestSize int64
	smallestSize = math.MaxInt64
	var indexOfSmallestSize = -1
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	for i, b := range bs {
		resp, err := client.Get(b.name.String() + "/uncachedBytes?t=" + t + "&cachefrom=" + cf)
		if err != nil {
			log.Printf("rpc for uncached bytes failed: %s", err)
			continue
		}
		if resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("reading response for uncached bytes failed: %s", err)
				continue
			}
			size, err := strconv.ParseInt(string(body), 10, 64)
			if err != nil {
				log.Printf("failed to parse uncached size bytes %s: %s", string(body), err)
				continue
			}
			if size < smallestSize {
				smallestSize = size
				indexOfSmallestSize = i
			}
		}
	}
	if indexOfSmallestSize == -1 {
		return nil, errors.New("all builders failed to report uncached bytes size")
	}
	return bs[indexOfSmallestSize], nil
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

func (s *dispatcher) recycle(b *builder) {
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

func (s *dispatcher) selectWorker(c clientID, t string, cf string) *builder {

	for {

		b, ok := s.reselect(c)
		if ok {
			log.Printf("reselected worker %s for client %s\n", b.name, c)
			return b
		}

		s.mutex.Lock()

		b, err := s.findScheduleable(c, t, cf)
		if err != nil {
			log.Printf("failed to select builder: %s\n", err)
		}
		if b != nil {
			log.Printf("selected worker %s for client %s\n", b.name, c)
			s.mutex.Unlock()
			return b
		}

		log.Printf("waiting for worker to become free\n")
		s.cond.Wait()
	}
}
