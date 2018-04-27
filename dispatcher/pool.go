package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"
)

type clientID string

type builderBytesize struct {
	builder  *builder
	bytesize int64
	err      error
}

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

func (s *dispatcher) findScheduleable(c clientID, v url.Values, h http.Header) (*builder, error) {
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

	b, err := selectByUncachedSize(selectableBuilders, v, h)
	if err != nil {
		return nil, fmt.Errorf("failed to select builder on uncached bytes size: %s", err)
	}

	b.dedicatedTo = c
	b.lastestUse = time.Now().Unix()
	atomic.AddInt32(&b.openConnections, 1)
	return b, nil
}

func selectByUncachedSize(bs []*builder, v url.Values, h http.Header) (*builder, error) {
	if len(v["t"]) == 0 && len(v["cachefrom"]) == 0 {
		return bs[0], nil
	}

	c := &http.Client{
		Timeout: 1 * time.Second,
	}
	chs := make([]chan builderBytesize, len(bs))
	for i, b := range bs {
		ch := make(chan builderBytesize)
		chs[i] = ch
		go func(b *builder, ch chan builderBytesize) {
			bytes, err := uncachedSize(c, b, v, h)
			ch <- builderBytesize{
				builder:  b,
				bytesize: bytes,
				err:      err,
			}
			close(ch)
		}(b, ch)
	}

	var smallestSize int64 = math.MaxInt64
	var b *builder
	for _, ch := range chs {
		bb := <-ch
		if bb.err != nil {
			log.Printf("builder failed to determine uncached size: %s", bb.err)
			continue
		}
		if bb.bytesize <= smallestSize {
			smallestSize = bb.bytesize
			b = bb.builder
		}
	}
	if b == nil {
		return nil, errors.New("all builders failed to report uncached bytes size")
	}
	return b, nil
}

func uncachedSize(c *http.Client, b *builder, v url.Values, h http.Header) (int64, error) {
	req, err := http.NewRequest("GET", b.name.String()+"/uncachedSize?"+v.Encode(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request to ask for uncached bytes: %s", err)
	}
	req.Header = h

	resp, err := c.Do(req)
	if err != nil {
		return 0, fmt.Errorf("rpc for uncached bytes failed: %s", err)
	}
	defer func(r *http.Response) {
		closingErr := resp.Body.Close()
		if closingErr != nil {
			log.Printf("closing response body failed: %s", err)
		}
	}(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response for uncached bytes failed: %s", err)
	}
	size, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uncached size bytes %s: %s", string(body), err)
	}
	return size, nil
}

func scheduleable(b *builder) bool {
	if b.dedicatedTo == "" {
		return true
	}
	notConnected := b.openConnections == 0
	reservationEnded := b.lastestUse <= time.Now().Add(reservation).Unix()
	return notConnected && reservationEnded
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
		log.Printf("worker can be recycled: %s\n", b.name)
		s.cond.Signal()
		log.Printf("broadcasted 'recycled worker'\n")
	}(b, t)
}

func (s *dispatcher) selectWorker(cxt context.Context, c clientID, v url.Values, h http.Header) (*builder, error) {

	b, ok := s.reselect(c)
	if ok {
		log.Printf("reselected worker %s for client %s\n", b.name, c)
		return b, nil
	}

	s.cond.L.Lock()
	for {
		b, ok := s.reselect(c)
		if ok {
			log.Printf("reselected worker %s for client %s\n", b.name, c)
			s.cond.L.Unlock()
			return b, nil
		}

		b, err := s.findScheduleable(c, v, h)
		if err != nil {
			log.Printf("failed to select builder: %s\n", err)
		}
		if b != nil {
			log.Printf("selected worker %s for client %s\n", b.name, c)
			s.cond.L.Unlock()
			return b, nil
		}

		log.Printf("waiting for worker to become free\n")
		s.cond.Wait()
		log.Printf("worker became free\n")

		if cxt.Err() != nil {
			s.cond.L.Unlock()
			s.cond.Signal()
			return nil, fmt.Errorf("stoped waiting for a worker: %s", cxt.Err())
		}
	}
}
