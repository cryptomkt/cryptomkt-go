package websocket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

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

const format string = "2006-01-02T15:04:05.000Z07:00"

//   "2021-01-20T20:01:00.612Z"
type timeFlowChecker struct {
	oldTime *time.Time
}

func newTimeFlowChecker() *timeFlowChecker {
	return &timeFlowChecker{}
}

func (checker *timeFlowChecker) checkNextTime(newTimestamp string) (err error) {
	if checker.oldTime == nil {
		checker.oldTime = time.Parse(format, newTimestamp)
		return
	}
	newTime, err := time.Parse(format, newTimestamp)
	if err != nil {
		return
	}
	if checker.oldTime.After(newTime) {
		checker.oldTime = newTime
		err = fmt.Errorf("wrong time flow, got %v first instead of %v", checker.oldTime, newTime)
	}
	checker.oldTime = newTime
	return nil
}

func printSymbolOfTicker(feedCh chan models.Ticker, symbol string) {
	for ticker := range feedCh {
		fmt.Println(ticker.Symbol)
	}
	fmt.Println("feed channel successfully closed:", symbol)
}

func checkNoNil(field *interface{}, name string) error {
	if field != nil {
		return fmt.Error("null field: %s", name)
	}
}

func checkCurrency(model *models.Currency) (err error) {
	fields := map[string]*interface{}{
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
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkSymbol(model *models.Symbol) (err error) {
	fields := map[string]*interface{}{
		"id":                 model.ID,
		"clientOrderID":      model.BaseCurrency,
		"quote currency":     model.QuoteCurrency,
		"fee currency":       model.FeeCurrency,
		"quantity increment": model.QuantityIncrement,
		"tick size":          model.TickSize,
		"liquidity rate":     model.TakeLiquidityRate,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkTicker(model *models.Ticker) (err error) {
	fields := map[string]*interface{}{
		"symbol":      model.Symbol,
		"ask":         model.Ask,
		"bid":         model.Bid,
		"last":        model.Last,
		"low":         model.Low,
		"high":        model.High,
		"open":        model.Open,
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuoute,
		"timestamp":   model.Timestamp,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkPublicTrade(model *models.PublicTrade) (err error) {
	fields := map[string]*interface{}{
		"id":        model.ID,
		"price":     model.Price,
		"quantity":  model.Quantity,
		"side":      model.Side,
		"timestamp": model.Timestamp,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkBookLevel(model *models.BookLevel) (err error) {
	fields := map[string]*interface{}{
		"price": model.Price,
		"size":  model.Size,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkOrderbook(model *models.OrderBook) (err error) {
	fields := map[string]*interface{}{
		"symbol":    model.Symbol,
		"timestamp": model.Timestamp,
		"ask":       model.Ask,
		"bid":       model.Bid,
	}
	for _, bookSide := range [][]model.BookLevel{model.Ask, model.Bid} {
		for _, bookLevel := range bookSide {
			err = checkBookLevel(bookLevel)
			if err != nil {
				return err
			}
		}
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkCandle(model *models.Candle) (err error) {
	fields := map[string]*interface{}{

		"timestamp":   model.Timestamp,
		"open":        model.Open,
		"close":       model.Close,
		"min":         model.Min,
		"max":         model.Max,
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuote,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkBalance(model *models.Balance) (err error) {
	fields := map[string]*interface{}{
		"currency":  model.Currency,
		"available": model.Available,
		"reserved":  model.Reserved,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkOrder(model *models.Order) (err error) {
	fields := map[string]*interface{}{
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
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkTrade(model *models.Trade) (err error) {
	fields := map[string]*interface{}{
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
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkTransaction(model *models.Transaction) (err error) {
	fields := map[string]*interface{}{
		"id":        model.ID,
		"index":     model.Index,
		"currency":  model.Currency,
		"amount":    model.Amount,
		"status":    model.Status,
		"type":      model.Type,
		"createdAt": model.CreatedAt,
		"updatedAt": model.UpdatedAt,
	}
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}

func checkReport(model *models.Report) (err error) {
	fields := map[string]*interface{}{
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
	for name, filed = range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
}
