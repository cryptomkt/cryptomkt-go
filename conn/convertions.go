package conn

import (
	"strconv"
)

// cleanMap deletes keys with empty string
func cleanMap(aMap *map[string]string) {
	for k, v := range *aMap {
		if v == "" {
			delete(*aMap, k)
		}
	}
}

// ToMap converts a ticket object to a map
func (ticker *Ticker) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["high"] = ticker.High
	asMap["volume"] = ticker.Volume
	asMap["low"] = ticker.Low
	asMap["ask"] = ticker.Ask
	asMap["timestamp"] = ticker.Timestamp
	asMap["bid"] = ticker.Bid
	asMap["last_price"] = ticker.LastPrice
	asMap["market"] = ticker.Market
	cleanMap(&asMap)
	return asMap
}

// ToMap converts a balance object to a map
func (balance *Balance) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["wallet"] = balance.Wallet
	asMap["available"] = balance.Available
	asMap["balance"] = balance.Balance
	cleanMap(&asMap)
	return asMap
}

// ToMap converts an order object to a map
func (order *Order) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["id"] = order.Id
	asMap["status"] = order.Status
	asMap["type"] = order.Type
	asMap["price"] = order.Price
	asMap["amount_original"] = order.Amount.Original
	asMap["amount_remaining"] = order.Amount.Remaining
	asMap["amount_executed"] = order.Amount.Executed
	asMap["execution_price"] = order.ExecutionPrice
	asMap["avg_execution_price"] = strconv.Itoa(order.AvgExecutionPrice)
	asMap["market"] = order.Market
	asMap["created_at"] = order.CreatedAt
	asMap["updated_at"] = order.UpdatedAt
	asMap["executed_at"] = order.ExecutedAt
	cleanMap(&asMap)
	return asMap
}

func (transaction *Transaction) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["id"] = transaction.Id
	asMap["type"] = strconv.Itoa(transaction.Type)
	asMap["amount"] = transaction.Amount
	asMap["fee_percent"] = transaction.FeePercent
	asMap["fee_amount"] = transaction.FeeAmount
	asMap["balance"] = transaction.Balance
	asMap["date"] = transaction.Date
	asMap["hash"] = transaction.Hash
	asMap["address"] = transaction.Address
	asMap["memo"] = transaction.Memo
	cleanMap(&asMap)
	return asMap
}

func (instant *Instant) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["obtained"] = strconv.FormatFloat(instant.Obtained, 'f', -1, 64)
	asMap["required"] = strconv.FormatFloat(instant.Required, 'f', -1, 64)
	cleanMap(&asMap)
	return asMap
}

// ToMap converts a *TradeData object to a map
func (trade *TradeData) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["market_taker"] = trade.MarketTaker
	asMap["price"] = trade.Price
	asMap["amount"] = trade.Amount
	asMap["tid"] = trade.Tid
	asMap["timestamp"] = trade.Timestamp
	asMap["market"] = trade.Market
	cleanMap(&asMap)
	return asMap
}
