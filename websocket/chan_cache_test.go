package websocket

import (
	"testing"
)

func TestOneUseNotificationChLifetime(t *testing.T) {
	chanCache := newChanCache()
	ch := make(chan []byte, 1)
	id := chanCache.saveCh(ch, 1)
	go dropAllMessagesFromCh(ch)
	reusableCh, ok := chanCache.getChan(id)
	if !ok {
		t.Fatal("should be present")
	}
	if reusableCh.isDone() {
		t.Fatal("should not be done")
	}
	reusableCh.send([]byte{0101})
	if !reusableCh.isDone() {
		t.Fatal("should be done")
	}
	chanCache.closeAndRemoveCh(id)
	_, ok = chanCache.getChan(id)
	if ok {
		t.Fatal("should not be present")
	}
}

func TestTwoUsesNotificationChLifetime(t *testing.T) {
	chanCache := newChanCache()
	ch := make(chan []byte, 1)
	id := chanCache.saveCh(ch, 2)
	go dropAllMessagesFromCh(ch)
	reusableCh, ok := chanCache.getChan(id)
	if !ok {
		t.Fatal("should be present")
	}
	if reusableCh.isDone() {
		t.Fatal("should not be done")
	}
	reusableCh.send([]byte{0101})
	if reusableCh.isDone() {
		t.Fatal("should not be done")
	}
	reusableCh, ok = chanCache.getChan(id)
	if !ok {
		t.Fatal("should be present")
	}
	if reusableCh.isDone() {
		t.Fatal("should not be done")
	}
	reusableCh.send([]byte{0101})
	if !reusableCh.isDone() {
		t.Fatal("should be done")
	}
	chanCache.closeAndRemoveCh(id)
	_, ok = chanCache.getChan(id)
	if ok {
		t.Fatal("should not be present")
	}
}

func TestCacheAutoRemovesNotificationChWhenDone(t *testing.T) {
	chanCache := newChanCache()
	ch := make(chan []byte, 1)
	id := chanCache.saveCh(ch, 1)
	go dropAllMessagesFromCh(ch)
	reusableCh, ok := chanCache.getChan(id)
	if !ok {
		t.Fatal("should be present")
	}
	if reusableCh.isDone() {
		t.Fatal("should not be done")
	}
	reusableCh.send([]byte{0101})
	if !reusableCh.isDone() {
		t.Fatal("should be done")
	}
	_, ok = chanCache.getChan(id)
	if !ok {
		t.Fatal("should not be eliminated")
	}
	_, ok = chanCache.getChan(id)
	if ok {
		t.Fatal("should be eliminated")
	}
}

func dropAllMessagesFromCh(ch chan []byte) {
	for range <-ch {
	}
}
