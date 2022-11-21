package models

import "github.com/cryptomarket/cryptomarket-go/args"

type WSTradeFeed map[string][]WSTrade
type WSTrade struct {
	Timestamp int64  `json:"t"`
	ID        int64  `json:"i"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	Side      string `json:"s"`
}

type WSCandleFeed map[string][]WSCandle
type WSCandle struct {
	Timestamp   int64  `json:"t"`
	Open        string `json:"o"`
	Close       string `json:"c"`
	High        string `json:"h"`
	Low         string `json:"l"`
	Volume      string `json:"v"`
	VolumeQuote string `json:"q"`
}

type MiniTickerFeed map[string]MiniTicker
type MiniTicker struct {
	Timestamp   int64  `json:"t"`
	Open        string `json:"o"`
	Last        string `json:"c"`
	High        string `json:"h"`
	Low         string `json:"l"`
	VolumeBase  string `json:"v"`
	VolumeQuote string `json:"q"`
}

type WSTickerFeed map[string]WSTicker
type WSTicker struct {
	Timestamp          int64  `json:"t"`
	BestAsk            string `json:"a"`
	BestAskQuantity    string `json:"A"`
	BestBid            string `json:"b"`
	BestBidQuantity    string `json:"B"`
	Last               string `json:"c"`
	Open               string `json:"o"`
	High               string `json:"h"`
	Low                string `json:"l"`
	VolumeBase         string `json:"v"`
	VolumeQuote        string `json:"q"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	LastTradeID        int64  `json:"L"`
}

type WSOrderbookFeed map[string]WSOrderbook
type WSOrderbook struct {
	Timestamp      int64      `json:"t"`
	SequenceNumber int64      `json:"s"`
	Ask            [][]string `json:"a"`
	Bid            [][]string `json:"b"`
}

type OrderbookTopFeed map[string]OrderbookTop
type OrderbookTop struct {
	Timestamp       int64  `json:"t"`
	BestAsk         string `json:"a"`
	BestAskQuantity string `json:"A"`
	BestBid         string `json:"b"`
	BestBidQuantity string `json:"B"`
}

type FeedType interface {
	WSTradeFeed |
		WSCandleFeed |
		WSOrderbookFeed |
		OrderbookTopFeed |
		MiniTickerFeed |
		WSTickerFeed |
		[]Report |
		Transaction |
		[]Balance
}

type Notification[ft FeedType] struct {
	Data             ft
	NotificationType args.NotificationType
}

type Subscription[ft FeedType] struct {
	NotificationChannel string
	NotificationCh      chan Notification[ft]
	Symbols             []string
}
