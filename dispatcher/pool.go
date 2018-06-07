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

	"github.com/damoon/ttlcache"
)

type clientID string

type builderBytesize struct {
	builder  *builder
	bytesize int64
	err      error
}

const reservation = 10 * time.Second

type cache struct {
	*ttlcache.Cache
}

func (cache *cache) Set(k clientID, v *builder) {
	cache.SetUnsafe(k, v)
}

func (cache *cache) Get(k clientID) (*builder, bool) {
	cached, found := cache.GetUnsafe(k)
	if !found {
		return &builder{}, false
	}
	return cached.(*builder), true
}

func (s *dispatcher) returnWorker(b *builder) {
	atomic.AddInt32(&b.openConnections, -1)
}

func (s *dispatcher) selectWorker(cxt context.Context, cID clientID, v url.Values, h http.Header) (*builder, error) {

	builder, found := s.cache.Get(cID)
	if found {
		atomic.AddInt32(&builder.openConnections, 1)
		log.Printf("reselected worker %s for client %s\n", builder.name, cID)
		return builder, nil
	}

	builder, err := s.findScheduleable(cID, v, h)
	if err != nil {
		log.Printf("failed to select builder: %s\n", err)
	}
	if builder != nil {
		s.cache.Set(cID, builder)
		atomic.AddInt32(&builder.openConnections, 1)
		log.Printf("selected worker %s for client %s\n", builder.name, cID)
		return builder, nil
	}

	builder, found = s.findLeastConnected()
	if found {
		s.cache.Set(cID, builder)
		atomic.AddInt32(&builder.openConnections, 1)
		log.Printf("selected worker %s for client %s\n", builder.name, cID)
		return builder, nil
	}

	return nil, fmt.Errorf("failed to find a worker")
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

func scheduleable(b *builder) bool {
	if b.dedicatedTo == "" {
		return true
	}
	notConnected := b.openConnections == 0
	reservationEnded := b.lastestUse <= time.Now().Add(reservation).Unix()
	return notConnected && reservationEnded
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

func (s *dispatcher) findLeastConnected() (*builder, bool) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var selected *builder
	for _, i := range r.Perm(len(s.builders)) {
		builder := s.builders[i]
		if selected == nil {
			selected = builder
			continue
		}
		if atomic.LoadInt32(&selected.openConnections) > atomic.LoadInt32(&builder.openConnections) {
			selected = builder
		}
	}
	return selected, selected != nil
}
