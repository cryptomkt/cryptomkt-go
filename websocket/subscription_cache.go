package websocket

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/cryptomkt/go-api/models"
)

type orderbookState int
type sortOrder int

const (
	// // orderbook states
	orderbookStateUpdating orderbookState = iota
	orderbookStateWaiting
	orderbookStateBroken

	// // ordering of side of orderbook
	sortOrderAscending sortOrder = iota
	sortOrderDescending
)

type orderbookSink struct {
	orderbook orderbookSnapshot
	in        chan []byte
	out       chan models.OrderBook
	wg        *sync.WaitGroup
	state     orderbookState
}

func newOrderbookSink() orderbookSink {
	return orderbookSink{
		in:    make(chan []byte, 10),
		out:   make(chan models.OrderBook),
		wg:    new(sync.WaitGroup),
		state: orderbookStateWaiting,
	}
}

func (obSink orderbookSink) sinkLoop(key string) {
	defer close(obSink.out)
	for data := range obSink.in {
		var result struct {
			Params orderbookSnapshot
		}
		json.Unmarshal(data, &result)
		fmt.Println(result.Params.Sequence)
		switch obSink.state {
		case orderbookStateWaiting:
			snapshot := result.Params
			fmt.Println("snapshot:", snapshot.Sequence)
			obSink.state = orderbookStateUpdating
			obSink.orderbook = snapshot
			obSink.out <- models.OrderBook{
				Symbol:    snapshot.Symbol,
				Timestamp: snapshot.Timestamp,
				Ask:       snapshot.Ask,
				Bid:       snapshot.Bid,
			}
		case orderbookStateUpdating:
			update := result.Params
			fmt.Println("update:", update.Sequence)
			snapshot := obSink.orderbook

			if update.Sequence-snapshot.Sequence != 1 {
				obSink.state = orderbookStateBroken
				continue
			}
			snapshot.Sequence = update.Sequence
			snapshot.Timestamp = update.Timestamp
			if ask := update.Ask; ask != nil {
				updateOrderbookSide(snapshot.Ask, ask, sortOrderAscending)
			}
			if bid := update.Bid; bid != nil {
				updateOrderbookSide(snapshot.Bid, bid, sortOrderDescending)
			}
			obSink.out <- models.OrderBook{
				Symbol:    snapshot.Symbol,
				Timestamp: snapshot.Timestamp,
				Ask:       snapshot.Ask,
				Bid:       snapshot.Bid,
			}
		case orderbookStateBroken:
			continue
		}
	}
	fmt.Println("closing orderbook feed: ", key)
}

type subscriptionCache struct {
	reportsChan  chan models.Report
	tickersChans map[string]chan models.Ticker
	orderbooks   map[string]orderbookSink
	tradesChans  map[string]chan models.PublicTrade
	candlesChans map[string]chan models.Candle
}

func newSubscriptionCache() subscriptionCache {
	return subscriptionCache{
		tickersChans: make(map[string]chan models.Ticker),
		orderbooks:   make(map[string]orderbookSink),
		tradesChans:  make(map[string]chan models.PublicTrade),
		candlesChans: make(map[string]chan models.Candle),
	}
}

func (cache subscriptionCache) close() {
	if cache.reportsChan != nil {
		close(cache.reportsChan)
		cache.reportsChan = nil
	}
	for key, ch := range cache.tickersChans {
		close(ch)
		delete(cache.tickersChans, key)
	}
	for key, sink := range cache.orderbooks {
		close(sink.in)
		delete(cache.orderbooks, key)
	}
	for key, ch := range cache.tradesChans {
		close(ch)
		delete(cache.tradesChans, key)
	}
	for key, ch := range cache.candlesChans {
		close(ch)
		delete(cache.candlesChans, key)
	}
}

func (cache subscriptionCache) closeFeedReport(key string) {
	if cache.reportsChan != nil {
		close(cache.reportsChan)
		cache.reportsChan = nil
	}
}
func (cache subscriptionCache) closeFeedTicker(key string) {
	if ch, ok := cache.tickersChans[key]; ok {
		close(ch)
		delete(cache.tickersChans, key)
	}
}
func (cache subscriptionCache) closeFeedOrderbook(key string) {
	if sink, ok := cache.orderbooks[key]; ok {
		close(sink.in)
		delete(cache.orderbooks, key)
	}
}
func (cache subscriptionCache) closeFeedTrade(key string) {
	if ch, ok := cache.tradesChans[key]; ok {
		close(ch)
		delete(cache.tradesChans, key)
	}
}
func (cache subscriptionCache) closeFeedCandle(key string) {
	if ch, ok := cache.candlesChans[key]; ok {
		close(ch)
		delete(cache.candlesChans, key)
	}

}

func (cache subscriptionCache) makeFeedReport(key string, buffSize int) chan models.Report {
	cache.closeFeedReport(key)
	ch := make(chan models.Report, buffSize)
	cache.reportsChan = ch
	return ch
}

func (cache subscriptionCache) makeFeedTicker(key string, buffSize int) chan models.Ticker {
	if ch, ok := cache.tickersChans[key]; ok {
		close(ch)
	}
	ch := make(chan models.Ticker, buffSize)
	cache.tickersChans[key] = ch
	return ch
}

func (cache subscriptionCache) makeFeedOrderbook(key string, buffSize int) chan models.OrderBook {
	if sink, ok := cache.orderbooks[key]; ok {
		close(sink.in)
	}
	sink := newOrderbookSink()
	go sink.sinkLoop(key)
	cache.orderbooks[key] = sink
	return sink.out
}

func (cache subscriptionCache) makeFeedTrade(key string, buffSize int) chan models.PublicTrade {
	if ch, ok := cache.tradesChans[key]; ok {
		close(ch)
	}
	ch := make(chan models.PublicTrade, buffSize)
	cache.tradesChans[key] = ch
	return ch
}

func (cache subscriptionCache) makeFeedCandle(key string, buffSize int) chan models.Candle {
	if ch, ok := cache.candlesChans[key]; ok {
		close(ch)
	}
	ch := make(chan models.Candle, buffSize)
	cache.candlesChans[key] = ch
	return ch
}

func (cache subscriptionCache) update(method, key string, data []byte) {
	fmt.Println("arriving method ", method)
	switch method {
	case methodActiveOrders:
		var result struct {
			Params []models.Report
		}
		json.Unmarshal(data, &result)
		for _, report := range result.Params {
			cache.reportsChan <- report
		}
	case methodReport:
		var result struct {
			Params models.Report
		}
		json.Unmarshal(data, &result)
		cache.reportsChan <- result.Params
	case methodTicker:
		go func() {

			var result struct {
				Params models.Ticker
			}
			json.Unmarshal(data, &result)
			if ch, ok := cache.tickersChans[key]; ok {
				ch <- result.Params
			}
		}()
	case methodSnapshotOrderbook, methodUpdateOrderbook:
		if obSink, ok := cache.orderbooks[key]; ok {
			obSink.in <- data
		}
	case methodSnapshotTrades, methodUpdateTrades:
		var result struct {
			Params []models.PublicTrade
		}
		json.Unmarshal(data, &result)
		if ch, ok := cache.tradesChans[key]; ok {
			for _, trade := range result.Params {
				ch <- trade
			}
		}
	case methodSnapshotCandles, methodUpdateCandles:
		var result struct {
			Params []models.Candle
		}
		json.Unmarshal(data, &result)
		if ch, ok := cache.candlesChans[key]; ok {
			for _, candle := range result.Params {
				ch <- candle
			}
		}

	}
}

func updateOrderbookSide(oldSide []models.BookLevel, updateSide []models.BookLevel, sortOrdering sortOrder) {
	oldIdx := 0
	updateIdx := 0
	newSide := make([]models.BookLevel, len(oldSide))
	for oldIdx < len(oldSide) && updateIdx < len(updateSide) {
		updateEntry := updateSide[updateIdx]
		oldEntry := oldSide[oldIdx]
		order := priceOrder(oldEntry, updateEntry, sortOrdering)
		if order == 0 {
			if !zeroSize(updateEntry) {
				newSide = append(newSide, updateEntry)
				updateIdx++
				oldIdx++
			}
		} else if order == 1 {
			newSide = append(newSide, oldEntry)
			oldIdx++
		} else {
			newSide = append(newSide, updateEntry)
			updateIdx++
		}
	}
	if updateIdx == len(updateSide) {
		for idx := oldIdx; idx < len(oldSide); idx++ {
			newSide = append(newSide, oldSide[idx])
		}
	}
	if oldIdx == len(oldSide) {
		for idx := updateIdx; idx < len(updateSide); idx++ {
			if !zeroSize(updateSide[idx]) {
				newSide = append(newSide, updateSide[idx])
			}
		}
	}
	oldSide = newSide
}

func zeroSize(entry models.BookLevel) bool {
	size, _ := new(big.Float).SetString(entry.Size)
	return size.Cmp(new(big.Float)) == 0
}

func priceOrder(oldEntry, updateEntry models.BookLevel, sortOrdering sortOrder) int {
	oldPrice, _ := new(big.Float).SetString(updateEntry.Price)
	updatePrice, _ := new(big.Float).SetString(updateEntry.Price)
	direction := oldPrice.Cmp(updatePrice)
	if sortOrdering == sortOrderAscending {
		return -direction
	} // sort Ordering == sortOrderDescending
	return direction

}

func (cache subscriptionCache) requireOrderbookSnapshot(key string) bool {
	if obSink, ok := cache.orderbooks[key]; ok {
		if obSink.state == orderbookStateWaiting {
			return true
		}
	}
	return false
}

func (cache subscriptionCache) waitingOrderbookSnapshot(key string) bool {
	if obSink, ok := cache.orderbooks[key]; ok {
		return obSink.state == orderbookStateWaiting
	}
	return false
}
