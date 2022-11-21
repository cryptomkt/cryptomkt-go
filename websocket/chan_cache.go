package websocket

import (
	"sync"
)

type chanCache struct {
	currentID             int64
	idLock                *sync.Mutex
	notificationsChans    map[int64]chan []byte
	notificationsChanLock *sync.Mutex
	subscriptionChans     map[string]chan []byte
	subscriptionChansLock *sync.RWMutex
}

func newChanCache() *chanCache {
	return &chanCache{
		currentID:             1,
		idLock:                new(sync.Mutex),
		notificationsChans:    make(map[int64]chan []byte),
		notificationsChanLock: new(sync.Mutex),
		subscriptionChans:     make(map[string]chan []byte),
		subscriptionChansLock: new(sync.RWMutex),
	}
}

func (cache *chanCache) close() {
	cache.subscriptionChansLock.Lock()
	defer cache.subscriptionChansLock.Unlock()
	for key, ch := range cache.subscriptionChans {
		close(ch)
		delete(cache.subscriptionChans, key)
	}
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
	cache.notificationsChanLock.Lock()
	defer cache.notificationsChanLock.Unlock()
	cache.notificationsChans[id] = ch
	return id
}

func (cache *chanCache) pop(id int64) (chan []byte, bool) {
	cache.notificationsChanLock.Lock()
	defer cache.notificationsChanLock.Unlock()
	if ch, ok := cache.notificationsChans[id]; ok {
		delete(cache.notificationsChans, id)
		return ch, true
	}
	return nil, false
}

func (cache *chanCache) sendViaSubscriptionCh(key string, data []byte) {
	cache.subscriptionChansLock.RLock()
	defer cache.subscriptionChansLock.RUnlock()
	if ch, ok := cache.subscriptionChans[key]; ok {
		ch <- data
	}
}

func (cache *chanCache) storeSubscriptionCh(key string, ch chan []byte) {
	cache.subscriptionChansLock.Lock()
	defer cache.subscriptionChansLock.Unlock()
	cache.subscriptionChans[key] = ch
}

func (cache *chanCache) deleteSubscriptionCh(key string) {
	cache.subscriptionChansLock.Lock()
	defer cache.subscriptionChansLock.Unlock()
	if ch, ok := cache.subscriptionChans[key]; ok {
		close(ch)
		delete(cache.subscriptionChans, key)
	}
}

func (cache *chanCache) getCh(key string) (chan []byte, bool) {
	cache.subscriptionChansLock.Lock()
	defer cache.subscriptionChansLock.Unlock()
	if ch, ok := cache.subscriptionChans[key]; ok {
		return ch, ok
	}
	return nil, false
}
