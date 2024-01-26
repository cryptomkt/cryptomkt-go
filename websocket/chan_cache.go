package websocket

import (
	"sync"
)

type chanCache struct {
	currentID             int64
	idLock                *sync.Mutex
	notificationsChans    map[int64]*reusableChan
	notificationsChanLock *sync.Mutex
	subscriptionChans     map[string]chan []byte
	subscriptionChansLock *sync.RWMutex
}

func newChanCache() *chanCache {
	return &chanCache{
		currentID:             1,
		idLock:                new(sync.Mutex),
		notificationsChans:    make(map[int64]*reusableChan),
		notificationsChanLock: new(sync.Mutex),
		subscriptionChans:     make(map[string]chan []byte),
		subscriptionChansLock: new(sync.RWMutex),
	}
}

func (cache *chanCache) close() {
	cache.closeSubscriptionChs()
	cache.closeNotificationChs()
}

func (cache *chanCache) closeSubscriptionChs() {
	cache.subscriptionChansLock.Lock()
	defer cache.subscriptionChansLock.Unlock()
	for key, ch := range cache.subscriptionChans {
		close(ch)
		delete(cache.subscriptionChans, key)
	}
}

func (cache *chanCache) closeNotificationChs() {
	cache.notificationsChanLock.Lock()
	defer cache.notificationsChanLock.Unlock()
	for key, ch := range cache.notificationsChans {
		close(ch.channel)
		delete(cache.notificationsChans, key)
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

func (cache *chanCache) saveCh(ch chan []byte, callCount int) int64 {
	id := cache.nextID()
	cache.notificationsChanLock.Lock()
	defer cache.notificationsChanLock.Unlock()
	cache.notificationsChans[id] = newReusableChan(ch, callCount)
	return id
}

func (cache *chanCache) getChan(id int64) (*reusableChan, bool) {
	cache.notificationsChanLock.Lock()
	defer cache.notificationsChanLock.Unlock()
	reusableCh, ok := cache.notificationsChans[id]
	if !ok {
		return nil, false
	}
	if reusableCh.isDone() {
		delete(cache.notificationsChans, id)
	}
	return reusableCh, true
}

func (cache *chanCache) closeAndRemoveCh(id int64) {
	cache.notificationsChanLock.Lock()
	defer cache.notificationsChanLock.Unlock()
	reusableCh, ok := cache.notificationsChans[id]
	if !ok {
		return
	}
	delete(cache.notificationsChans, id)
	reusableCh.close()
}

func (cache *chanCache) sendViaSubscriptionCh(key string, data []byte) {
	cache.subscriptionChansLock.RLock()
	defer cache.subscriptionChansLock.RUnlock()
	if ch, ok := cache.subscriptionChans[key]; ok {
		ch <- data
	}
}

func (cache *chanCache) saveSubscriptionCh(key string, ch chan []byte) {
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

func (cache *chanCache) getSubscriptionChan(key string) (chan []byte, bool) {
	cache.subscriptionChansLock.Lock()
	defer cache.subscriptionChansLock.Unlock()
	ch, ok := cache.subscriptionChans[key]
	if !ok {
		return nil, false
	}
	return ch, ok
}
