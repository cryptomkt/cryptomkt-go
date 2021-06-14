package websocket

import "sync"

type chanCache struct {
	currentID int64
	chans     *sync.Map
	idLock    *sync.Mutex
}

func newChanCache() *chanCache {
	return &chanCache{
		currentID: 1,
		chans:     new(sync.Map),
		idLock:    new(sync.Mutex),
	}
}

func (cache *chanCache) close() {
	cache.chans.Range(func(key interface{}, val interface{}) bool {
		ch := val.(chan []byte)
		close(ch)
		return true
	})
}

func (cache *chanCache) nextID() int64 {
	cache.idLock.Lock()
	defer cache.idLock.Unlock()
	if cache.currentID < 1 {
		cache.currentID = 1
	}
	id := cache.currentID
	cache.currentID++
	return id
}

func (cache *chanCache) store(ch chan []byte) int64 {
	id := cache.nextID()
	cache.chans.Store(id, ch)
	return id
}

func (cache *chanCache) pop(id int64) (chan []byte, bool) {
	if val, ok := cache.chans.LoadAndDelete(id); ok {
		return val.(chan []byte), ok
	}
	return nil, false
}

func (cache *chanCache) getSubcriptionCh(key string) (chan []byte, bool) {
	if val, ok := cache.chans.Load(key); ok {
		return val.(chan []byte), ok
	}
	return nil, false
}

func (cache *chanCache) storeSubscriptionCh(key string, ch chan []byte) {
	if val, ok := cache.chans.Load(key); ok {
		close(val.(chan []byte)) // close old channel
	}
	cache.chans.Store(key, ch)
}

func (cache *chanCache) deleteSubscriptionCh(key string) {
	cache.chans.Delete(key)
}