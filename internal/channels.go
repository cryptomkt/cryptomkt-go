package internal

const (
	ChannelTrades  = "trades"
	ChannelCandles = "candles/%s"

	ChannelMiniTickerInBatch = "ticker/price/%s/batch"
	ChannelTicker            = "ticker/%s"
	ChannelMiniTicker        = "ticker/price/%s"
	ChannelTickerInBatch     = "ticker/%s/batch"

	ChannelOrderBookFull           = "orderbook/full"
	ChannelOrderbookPartial        = "orderbook/%s/%s"
	ChannelOrderbookPartialInBatch = "orderbook/%s/%s/batch"
	ChannelOrderbookTop            = "orderbook/top/%s"
	ChannelOrderbookTopInBatch     = "orderbook/top/%s/batch"
)
