package websocket

import (
	"encoding/json"
	"math/big"

	"github.com/cryptomarket/cryptomarket-go/models"
)

type orderbookSnapshot struct {
	Symbol    string             `json:"symbol"`
	Sequence  int64              `json:"sequence"`
	Timestamp string             `json:"timestamp"`
	Ask       []models.BookLevel `json:"ask"`
	Bid       []models.BookLevel `json:"bid"`
}

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

type orderbookCache struct {
	orderbook      *orderbookSnapshot
	orderbookState orderbookState
}

func newOrderbookCache() *orderbookCache {
	return &orderbookCache{
		orderbook:      nil,
		orderbookState: orderbookStateWaiting,
	}
}

func parseOB(data []byte) *orderbookSnapshot {
	var resp struct {
		Params orderbookSnapshot
	}
	json.Unmarshal(data, &resp)
	return &resp.Params

}

func (cache *orderbookCache) obSnapshot(data []byte) {
	cache.orderbook = parseOB(data)
	cache.orderbookState = orderbookStateUpdating
}

func (cache *orderbookCache) obUpdate(data []byte) {
	if cache.orderbookState != orderbookStateUpdating {
		return
	}
	updateOB := parseOB(data)
	oldOB := cache.orderbook
	if updateOB.Sequence-oldOB.Sequence != 1 {
		cache.orderbookState = orderbookStateBroken
		return
	}
	oldOB.Sequence = updateOB.Sequence
	oldOB.Timestamp = updateOB.Timestamp

	oldOB.Ask = updateOrderbookSide(oldOB.Ask, updateOB.Ask, sortOrderAscending)
	oldOB.Bid = updateOrderbookSide(oldOB.Bid, updateOB.Bid, sortOrderDescending)
}

func updateOrderbookSide(oldSide []models.BookLevel, updateSide []models.BookLevel, sortOrdering sortOrder) []models.BookLevel {
	oldIdx := 0
	updateIdx := 0
	newSide := make([]models.BookLevel, 0)
	for oldIdx < len(oldSide) && updateIdx < len(updateSide) {
		updateEntry := updateSide[updateIdx]
		oldEntry := oldSide[oldIdx]
		order := priceOrder(oldEntry, updateEntry, sortOrdering)
		if order == 0 {
			if !zeroSize(updateEntry) {
				newSide = append(newSide, updateEntry)
			}
			updateIdx++
			oldIdx++
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
	return newSide
}

func zeroSize(entry models.BookLevel) bool {
	size, _ := new(big.Float).SetString(entry.Size)
	return size.Cmp(new(big.Float)) == 0
}

func priceOrder(oldEntry, updateEntry models.BookLevel, sortOrdering sortOrder) int {
	oldPrice, _ := new(big.Float).SetString(oldEntry.Price)
	updatePrice, _ := new(big.Float).SetString(updateEntry.Price)
	direction := oldPrice.Cmp(updatePrice)
	if sortOrdering == sortOrderAscending {
		return -direction
	}
	return direction

}

func (cache *orderbookCache) getOrderbookData() []byte {
	data, err := json.Marshal(cache.orderbook)
	if err != nil {
		return []byte{}
	}
	return data
}

func (cache *orderbookCache) obWaiting() bool {
	return cache.orderbookState == orderbookStateWaiting
}

func (cache *orderbookCache) obBroken() bool {
	return cache.orderbookState == orderbookStateBroken
}

func (cache *orderbookCache) waitOBSnapshot() {
	cache.orderbookState = orderbookStateWaiting
}
