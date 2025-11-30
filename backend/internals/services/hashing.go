package services

import (
	"sync"
	"time"
)

type HASH[V any] struct {
	mu   sync.RWMutex
	Data V
	TTL  time.Time
}

type HashMap[V any] struct {
	mu  sync.RWMutex
	Map map[string]*HASH[V]
}

func NewHashMap[V any]() *HashMap[V] {
	hm := &HashMap[V]{
		Map: make(map[string]*HASH[V]),
	}
	go hm.TimeoutLoop()
	return hm
}

func (hm *HashMap[V]) TimeoutLoop() {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		now := time.Now()

		hm.mu.Lock()
		for key, h := range hm.Map {
			if now.After(h.TTL) {
				delete(hm.Map, key)
			}
		}
		hm.mu.Unlock()
	}
}

func (hm *HashMap[V]) Get(k string) (V, bool) {
	hm.mu.RLock()
	item, ok := hm.Map[k]
	hm.mu.RUnlock()

	if !ok {
		var zero V
		return zero, false
	}

	item.mu.RLock()
	defer item.mu.RUnlock()
	return item.Data, true
}

func (hm *HashMap[V]) Delete(k string) {
	hm.mu.Lock()
	delete(hm.Map, k)
	hm.mu.Unlock()
}

func (hm *HashMap[V]) Set(k string, v V, ttl time.Duration) {
	hm.mu.Lock()
	item, exists := hm.Map[k]

	// create new entry
	if !exists {
		item = &HASH[V]{
			Data: v,
			TTL:  time.Now().Add(ttl),
		}
		hm.Map[k] = item
		hm.mu.Unlock()
		return
	}

	// update existing entry
	hm.mu.Unlock()

	item.mu.Lock()
	item.Data = v
	item.TTL = time.Now().Add(ttl)
	item.mu.Unlock()
}
