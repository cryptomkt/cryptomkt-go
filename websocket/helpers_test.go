package websocket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/cryptomarket/cryptomarket-go/models"
)

type APIKeys struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

func LoadKeys() (apiKeys APIKeys) {
	data, err := ioutil.ReadFile("/home/ismael/cryptomarket/apis/keys.json")
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
		newTime, err := time.Parse(format, newTimestamp)
		if err != nil {
			return err
		}
		checker.oldTime = &newTime
		return nil
	}
	newTime, err := time.Parse(format, newTimestamp)
	if err != nil {
		return err
	}
	if checker.oldTime.After(newTime) {
		checker.oldTime = &newTime
		err = fmt.Errorf("wrong time flow, got %v first instead of %v", checker.oldTime, newTime)
	}
	checker.oldTime = &newTime
	return nil
}

type sequenceFlowChecker struct {
	sequence int64
}

func newSequenceFlowChecker() *sequenceFlowChecker {
	return &sequenceFlowChecker{}
}

func (checker *sequenceFlowChecker) checkNextSequence(nextSequence int64) (err error) {
	if checker.sequence != 0 && nextSequence <= checker.sequence {
		err = fmt.Errorf("wrong sequence, old:%v\tnew:%v", checker.sequence, nextSequence)
	}
	checker.sequence = nextSequence
	return
}

func checkNoNil(field interface{}, name string) error {
	if field == nil {
		return fmt.Errorf("null field: %s", name)
	}
	switch v := field.(type) {
	case string:
		if v == "" {
			return fmt.Errorf("null string: %v", name)
		}
	case int, int64:
		if v == 0 {
			return fmt.Errorf("null number: %v", name)
		}
	}
	return nil
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
	return nil
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
	return nil
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
			return err
		}
	}
	return nil
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
	return nil
}

func checkBookLevel(model *models.BookLevel) (err error) {
	fields := map[string]interface{}{
		"price": model.Price,
		"size":  model.Size,
	}
	for name, field := range fields {
		if err = checkNoNil(field, name); err != nil {
			return err
		}
	}
	size, _ := new(big.Float).SetString(model.Size)
	zero, _ := new(big.Float).SetString("0.00")
	if size.Cmp(zero) == 0 {
		fmt.Println(model)
		return fmt.Errorf("zero level")
	}
	return nil
}

func checkOrderbook(model *models.OrderBook) (err error) {
	fields := map[string]interface{}{
		"symbol":    model.Symbol,
		"timestamp": model.Timestamp,
		"ask":       model.Ask,
		"bid":       model.Bid,
	}
	for name, field := range fields {
		if err = checkNoNil(field, name); err != nil {
			return err
		}
	}
	sides := make([][]models.BookLevel, 2)
	sides[0] = model.Ask
	sides[1] = model.Bid
	for _, bookSide := range sides {
		for _, bookLevel := range bookSide {
			if err = checkBookLevel(&bookLevel); err != nil {
				return err
			}
		}
	}
	return nil
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
	return nil
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
	return nil
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
	return nil
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
	return nil
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
	return nil
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
		"reportType":    model.ReportType,
	}
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return nil
}
