package conn

import (
	"strconv"
)

// parseState returns messages based in the code *PaymentOrder provides
func (paymentOrder *PaymentOrder) parseState() string {
	switch paymentOrder.Status {
	case -4:
		return "Pago Múltiple"
	case -3:
		return "Monto pagado no concuerda"
	case -2:
		return "Falló conversión"
	case -1:
		return "Expiró orden de pago"
	case 0:
		return "Esperando pago"
	case 1:
		return "Esperando bloque"
	case 2:
		return "Esperando procesamiento"
	case 3:
		return "Pago exitoso"
	}
	return "Invalid state"
}

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

// ToMap converts a PaymentOrder object to a map
func (payment *PaymentOrder) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["id"] = payment.Id
	asMap["external_id"] = payment.ExternalId
	asMap["status"] = strconv.Itoa(payment.Status)
	asMap["to_receive"] = payment.ToReceive
	asMap["to_receive_currency"] = payment.ToReceiveCurrency
	asMap["expected_amount"] = payment.ExpectedAmount
	asMap["expected_currency"] = payment.ExpectedCurrency
	asMap["deposit_address"] = payment.DepositAddress
	asMap["refund_email"] = payment.RefundEmail
	asMap["qr"] = payment.Qr
	asMap["obs"] = payment.Obs
	asMap["callback_url"] = payment.CallbackUrl
	asMap["error_url"] = payment.ErrorUrl
	asMap["success_url"] = payment.SuccessUrl
	asMap["payment_url"] = payment.PaymentUrl
	asMap["remaining"] = strconv.Itoa(payment.Remaining)
	asMap["language"] = payment.Language
	asMap["created_at"] = payment.CreatedAt
	asMap["updated_at"] = payment.UpdatedAt
	asMap["server_at"] = payment.ServerAt
	cleanMap(&asMap)
	return asMap
}

// ToMap converts a *Transaction object to a map
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

// ToMap converts a *Quantity object to a map
func (instant *Quantity) ToMap() map[string]string {
	asMap := make(map[string]string)
	asMap["obtained"] = instant.Obtained
	asMap["required"] = instant.Required
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
