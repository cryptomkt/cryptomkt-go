package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cryptomkt/go-api/models"
)

type APIKeys struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

func LoadKeys() (apiKeys APIKeys) {
	data, err := ioutil.ReadFile("../../keys.json")
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(data, &apiKeys)
	return
}

func checkNoNil(field interface{}, name string) (err error) {
	errMsg := fmt.Errorf("null field: %s", name)
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

func checkCurrency(model *models.Currency) (err error) {
	fields := map[string]interface{}{
		"id":                 model.ID,
		"fullName":           model.FullName,
		"crypto":             model.Crypto,
		"payinEnabled":       model.PayinEnabled,
		"payinPaymentId":     model.PayinPaymentID,
		"payinConfirmations": model.PayinConfirmations,
		"payoutEnabled":      model.PayoutEnabled,
		"payoutIsPaymentId":  model.PayinPaymentID,
		"transferEnabled":    model.TransferEnabled,
		"delisted":           model.Delisted,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkSymbol(model *models.Symbol) (err error) {
	fields := map[string]interface{}{
		"id":                 model.ID,
		"clientOrderID":      model.BaseCurrency,
		"quote currency":     model.QuoteCurrency,
		"fee currency":       model.FeeCurrency,
		"quantity increment": model.QuantityIncrement,
		"tick size":          model.TickSize,
		"liquidity rate":     model.TakeLiquidityRate,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkTicker(model *models.Ticker) (err error) {
	fields := map[string]interface{}{
		"symbol":      model.Symbol,
		"ask":         model.Ask,
		"bid":         model.Bid,
		"last":        model.Last,
		"low":         model.Low,
		"high":        model.High,
		"open":        model.Open,
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuote,
		"timestamp":   model.Timestamp,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			fmt.Println(model.Symbol)
			return err
		}
	}
	return
}

func checkPublicTrade(model *models.PublicTrade) (err error) {
	fields := map[string]interface{}{
		"id":        model.ID,
		"price":     model.Price,
		"quantity":  model.Quantity,
		"side":      model.Side,
		"timestamp": model.Timestamp,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkBookLevel(model *models.BookLevel) (err error) {
	fields := map[string]interface{}{
		"price": model.Price,
		"size":  model.Size,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkOrderbook(model *models.OrderBook) (err error) {
	fields := map[string]interface{}{
		"symbol":    model.Symbol,
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
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkCandle(model *models.Candle) (err error) {
	fields := map[string]interface{}{

		"timestamp":   model.Timestamp,
		"open":        model.Open,
		"close":       model.Close,
		"min":         model.Min,
		"max":         model.Max,
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuote,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkBalance(model *models.Balance) (err error) {
	fields := map[string]interface{}{
		"currency":  model.Currency,
		"available": model.Available,
		"reserved":  model.Reserved,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkOrder(model *models.Order) (err error) {
	fields := map[string]interface{}{
		"id":            model.ID,
		"clientOrderId": model.ClientOrderID,
		"symbol":        model.Symbol,
		"side":          model.Side,
		"status":        model.Status,
		"type":          model.Type,
		"timeInForce":   model.TimeInForce,
		"quantity":      model.Quantity,
		"price":         model.Price,
		"cumQuantity":   model.CumQuantity,
		"createdAt":     model.CreatedAt,
		"updatedAt":     model.UpdatedAt,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
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
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkTransaction(model *models.Transaction) (err error) {
	fields := map[string]interface{}{
		"id":        model.ID,
		"index":     model.Index,
		"currency":  model.Currency,
		"amount":    model.Amount,
		"status":    model.Status,
		"type":      model.Type,
		"createdAt": model.CreatedAt,
		"updatedAt": model.UpdatedAt,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}

func checkReport(model *models.Report) (err error) {
	fields := map[string]interface{}{
		"id":            model.ID,
		"clientOrderId": model.ClientOrderID,
		"symbol":        model.Symbol,
		"side":          model.Side,
		"status":        model.Status,
		"type":          model.Type,
		"timeInForce":   model.TimeInForce,
		"quantity":      model.Quantity,
		"price":         model.Price,
		"cumQuantity":   model.CumQuantity,
		"createdAt":     model.CreatedAt,
		"updatedAt":     model.UpdatedAt,
		"report type":   model.ReportType,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return
}
