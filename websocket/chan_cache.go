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
		reqChans := val.(*requestsChans)
		close(reqChans.ch)
		close(reqChans.errCh)
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

func (cache *chanCache) store(chans *requestsChans) int64 {
	id := cache.nextID()
	cache.chans.Store(id, chans)
	return id
}

func (cache *chanCache) pop(id int64) (*requestsChans, bool) {
	if val, ok := cache.chans.LoadAndDelete(id); ok {
		return val.(*requestsChans), ok
	}
	return nil, false
}
