package models

func FromOrderbookJsonToOrderbook(obJson OrderBookJson) OrderBook {
	asks := make([]BookLevel, 0, len(obJson.Ask))
	for _, level := range obJson.Ask {
		asks = append(asks, BookLevel{
			Price:  level[0],
			Amount: level[1],
		})
	}
	bids := make([]BookLevel, 0, len(obJson.Bid))
	for _, level := range obJson.Bid {
		bids = append(asks, BookLevel{
			Price:  level[0],
			Amount: level[1],
		})
	}
	return OrderBook{
		Timestamp: obJson.Timestamp,
		Ask:       asks,
		Bid:       bids,
	}
}
