package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cryptomarket/cryptomarket-go/models"
)

type APIKeys struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

func LoadKeys() (apiKeys APIKeys) {
	data, err := ioutil.ReadFile("/home/ismael/cryptomarket/keys-v3.json")
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(data, &apiKeys)
	return
}

func checkNoNil(field interface{}, name string) (err error) {
	errMsg := fmt.Errorf("nil field: %s", name)
	switch v := field.(type) {
	case int, int64:
		if v == 0 {
			err = errMsg
		}
	case string:
		if v == "" {
			err = errMsg
		}
	default:
		if field == nil {
			err = errMsg
		}
	}
	return
}

func checkFields(fields map[string]interface{}) (err error) {
	for name, field := range fields {
		if err = checkNoNil(field, name); err != nil {
			return
		}
	}
	return nil
}

func checkCurrency(model *models.Currency) (err error) {
	fields := map[string]interface{}{
		"fullName":        model.FullName,
		"payinEnabled":    model.PayinEnabled,
		"payoutEnabled":   model.PayoutEnabled,
		"transferEnabled": model.TransferEnabled,
	}
	err = checkFields(fields)
	return
}

func checkSymbol(model *models.Symbol) (err error) {
	fields := map[string]interface{}{
		"type":               model.Type,
		"clientOrderID":      model.BaseCurrency,
		"quote currency":     model.QuoteCurrency,
		"status":             model.Status,
		"quantity increment": model.QuantityIncrement,
		"tick size":          model.TickSize,
		"take rate":          model.TakeRate,
		"make rate":          model.MakeRate,
		"fee currency":       model.FeeCurrency,
	}
	err = checkFields(fields)
	return
}

func checkTicker(model *models.Ticker) (err error) {
	fields := map[string]interface{}{
		// "ask":         model.Ask, // can be nil
		// "bid":         model.Bid, // can be nil
		// "last":        model.Last, // can be nil
		"low":  model.Low,
		"high": model.High,
		// "open":        model.Open, // can be nil
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuote,
		"timestamp":   model.Timestamp,
	}
	err = checkFields(fields)
	return
}

func checkQuotationPrice(model *models.Price) (err error) {
	fields := map[string]interface{}{
		"currency":  model.Currency,
		"price":     model.Price,
		"timestamp": model.Timestamp,
	}
	err = checkFields(fields)
	return
}

func checkTickerPrice(model *models.Price) (err error) {
	fields := map[string]interface{}{
		"price":     model.Price,
		"timestamp": model.Timestamp,
	}
	err = checkFields(fields)
	return
}

func checkQuotationPriceHistory(model *models.PriceHistory) (err error) {
	fields := map[string]interface{}{
		"currency": model.Currency,
		"history":  model.History,
	}
	err = checkFields(fields)
	if err != nil {
		return
	}
	for _, historyPoint := range model.History {
		err = checkHistoryPoint(&historyPoint)
		if err != nil {
			return
		}
	}
	return
}

func checkHistoryPoint(model *models.HistoryPoint) (err error) {
	fields := map[string]interface{}{
		"timestamp": model.Timestamp,
		"open":      model.Open,
		"close":     model.Close,
		"min":       model.Min,
		"max":       model.Max,
	}
	err = checkFields(fields)
	return
}

func checkPublicTrade(model *models.PublicTrade) (err error) {
	fields := map[string]interface{}{
		"id":        model.ID,
		"price":     model.Price,
		"qty":       model.Quantity,
		"side":      model.Side,
		"timestamp": model.Timestamp,
	}
	err = checkFields(fields)
	return
}

func checkBookLevel(model *models.BookLevel) (err error) {
	fields := map[string]interface{}{
		"price": model.Price,
		"size":  model.Amount,
	}
	err = checkFields(fields)
	return
}

func checkOrderbook(model *models.OrderBook) (err error) {
	fields := map[string]interface{}{
		"timestamp": model.Timestamp,
		"ask":       model.Ask,
		"bid":       model.Bid,
	}
	for _, bookSide := range [][]models.BookLevel{model.Ask, model.Bid} {
		for _, bookLevel := range bookSide {
			err = checkBookLevel(&bookLevel)
			if err != nil {
				return err
			}
		}
	}
	err = checkFields(fields)
	return
}

func checkCandle(model *models.Candle) (err error) {
	fields := map[string]interface{}{
		"timestamp":   model.Timestamp,
		"open":        model.Open,
		"close":       model.Close,
		"min":         model.Low,
		"max":         model.High,
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuote,
	}
	err = checkFields(fields)
	return
}

func checkBalance(model *models.Balance) (err error) {
	fields := map[string]interface{}{
		// "currency":  model.Currency,  only present when querying a list of balances
		"available": model.Available,
		"reserved":  model.Reserved,
	}
	err = checkFields(fields)
	return
}

func checkOrder(model *models.Order) (err error) {
	fields := map[string]interface{}{
		"id":                  model.ID,
		"client_order_id":     model.ClientOrderID,
		"symbol":              model.Symbol,
		"side":                model.Side,
		"status":              model.Status,
		"type":                model.Type,
		"time_in_force":       model.TimeInForce,
		"quantity":            model.Quantity,
		"price":               model.Price,
		"quantity_cumulative": model.QuantityCumulative,
		"created_at":          model.CreatedAt,
		"updated_at":          model.UpdatedAt,
	}
	err = checkFields(fields)
	return
}

func checkTrade(model *models.Trade) (err error) {
	fields := map[string]interface{}{
		"id":            model.ID,
		"orderId":       model.OrderID,
		"clientOrderId": model.ClientOrderID,
		"symbol":        model.Symbol,
		"side":          model.Side,
		"quantity":      model.Quantity,
		"price":         model.Price,
		"fee":           model.Fee,
		"timestamp":     model.Timestamp,
	}
	err = checkFields(fields)
	return
}

func checkTransaction(model *models.Transaction) (err error) {
	fields := map[string]interface{}{
		"id":        model.ID,
		"status":    model.Status,
		"type":      model.Type,
		"createdAt": model.CreatedAt,
		"updatedAt": model.UpdatedAt,
	}
	err = checkFields(fields)
	return
}

func checkReport(model *models.Report) (err error) {
	fields := map[string]interface{}{
		"id":            model.ID,
		"clientOrderId": model.ClientOrderID,
		"symbol":        model.Symbol,
		"side":          model.Side,
		"status":        model.Status,
		"type":          model.OrderType,
		"timeInForce":   model.TimeInForce,
		"quantity":      model.Quantity,
		"price":         model.Price,
		"cumQuantity":   model.QuantityCumulative,
		"createdAt":     model.CreatedAt,
		"updatedAt":     model.UpdatedAt,
		"report type":   model.ReportType,
	}
	err = checkFields(fields)
	return
}

func checkCryptoAddress(model *models.CryptoAddress) (err error) {
	fields := map[string]interface{}{
		"address":  model.Address,
		"currency": model.Currency,
		// "paymentID":   model.PaymentID, // optional
		// "publicKey":     model.PublicKey, // optional
	}
	err = checkFields(fields)
	return
}

func checkTradingCommission(model *models.TradingCommission) (err error) {
	fields := map[string]interface{}{
		"symbol":    model.Symbol,
		"make rate": model.MakeRate,
		"take rate": model.TakeRate,
	}
	err = checkFields(fields)
	return
}

type checkable interface {
	models.Currency |
		models.Symbol |
		models.Price |
		models.Order |
		models.Balance |
		models.Trade |
		models.CryptoAddress |
		models.Transaction |
		models.TradingCommission
}

func checkList[C checkable](checkerFn func(*C) error, list []C) error {
	for _, aCheckable := range list {
		err := checkerFn(&aCheckable)
		if err != nil {
			return err
		}
	}
	return nil
}
