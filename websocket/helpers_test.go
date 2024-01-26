package websocket

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cryptomkt/cryptomkt-go/models"
)

type APIKeys struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

func LoadKeys() (apiKeys APIKeys) {
	data, err := os.ReadFile("/home/ismael/cryptomarket/keys.json")
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(data, &apiKeys)
	if err != nil {
		fmt.Print(err)
	}
	return
}

type stringSaver struct {
	save   chan string
	take   chan string
	memory []string
}

type errorSaver struct {
	save   chan error
	take   chan error
	memory []error
}

type saver struct {
	strSaver      stringSaver
	errSaver      errorSaver
	errorsPrinted bool
}

func runSaver() *saver {
	strSaver := stringSaver{
		save:   make(chan string, 1),
		take:   make(chan string, 1),
		memory: make([]string, 0),
	}
	errSaver := errorSaver{
		save:   make(chan error, 1),
		take:   make(chan error, 1),
		memory: make([]error, 0),
	}
	go func() {
		defer close(strSaver.take)
		for str := range strSaver.save {
			strSaver.memory = append(strSaver.memory, str)
		}
		for _, str := range strSaver.memory {
			strSaver.take <- str
		}
	}()
	go func() {
		defer close(errSaver.take)
		for err := range errSaver.save {
			errSaver.memory = append(errSaver.memory, err)
		}
		for _, err := range errSaver.memory {
			errSaver.take <- err
		}
	}()
	return &saver{
		strSaver:      strSaver,
		errSaver:      errSaver,
		errorsPrinted: false,
	}
}

func (saver *saver) strSaveCh() chan<- string {
	return saver.strSaver.save
}

func (saver *saver) errSaveCh() chan<- error {
	return saver.errSaver.save
}

func (saver *saver) strTakeCh() <-chan string {
	return saver.strSaver.take
}

func (saver *saver) errTakeCh() <-chan error {
	return saver.errSaver.take
}

func (saver *saver) printSavedStrings() {
	for str := range saver.strTakeCh() {
		fmt.Println("str:", str)
	}
}

func (saver *saver) printSavedErrors() {
	for err := range saver.errTakeCh() {
		saver.errorsPrinted = true
		fmt.Println("error:", err)
	}
}

func (saver *saver) wereErrorsPrinted() bool {
	return saver.errorsPrinted
}

func (saver *saver) close() {
	close(saver.errSaver.save)
	close(saver.strSaver.save)
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

func checkFields(fields map[string]interface{}) (err error) {
	for name, field := range fields {
		err = checkNoNil(field, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkMiniTicker(model *models.MiniTicker) (err error) {
	fields := map[string]interface{}{
		"timestamp":    model.Timestamp,
		"open":         model.Open,
		"last":         model.Last,
		"high":         model.High,
		"low":          model.Low,
		"volume base":  model.VolumeBase,
		"volume quote": model.VolumeQuote,
	}
	return checkFields(fields)
}

func checkWSTicker(model *models.WSTicker) (err error) {
	fields := map[string]interface{}{
		"timestamp":            model.Timestamp,
		"best ask":             model.BestAsk,
		"best ask quantity":    model.BestAskQuantity,
		"best bid":             model.BestBid,
		"best bid quantity":    model.BestBidQuantity,
		"open":                 model.Open,
		"last":                 model.Last,
		"high":                 model.High,
		"low":                  model.Low,
		"volume base":          model.VolumeBase,
		"volume quote":         model.VolumeQuote,
		"price change":         model.PriceChange,
		"price change percent": model.PriceChangePercent,
		"last trade ID":        model.LastTradeID,
	}
	return checkFields(fields)
}

func checkWSTrade(model *models.WSTrade) (err error) {
	fields := map[string]interface{}{
		"id":        model.ID,
		"price":     model.Price,
		"quantity":  model.Quantity,
		"side":      model.Side,
		"timestamp": model.Timestamp,
	}
	return checkFields(fields)
}

type OBChecker struct {
	lastSequence int64
}

func (obchecker *OBChecker) checkOrderbook(model *models.WSOrderbook) (err error) {
	fields := map[string]interface{}{
		"timestamp": model.Timestamp,
		"ask":       model.Ask,
		"bid":       model.Bid,
		"sequence":  model.SequenceNumber,
	}
	err = checkFields(fields)
	if err != nil {
		return
	}
	if model.SequenceNumber != obchecker.lastSequence+1 &&
		obchecker.lastSequence != 0 {
		err = fmt.Errorf(
			"wrong sequence: wanted:%v. recieved:%v",
			obchecker.lastSequence+1,
			model.SequenceNumber,
		)
	}
	obchecker.lastSequence = model.SequenceNumber
	if err != nil {
		return err
	}
	sides := make([][][]string, 2)
	sides[0] = model.Ask
	sides[1] = model.Bid
	for _, bookSide := range sides {
		for _, bookLevel := range bookSide {
			if len(bookLevel) != 2 {
				return fmt.Errorf("wrong book level size, should be 2, and its:%v", len(bookLevel))
			}
		}
	}
	return nil
}

func checkOrderbookTop(model *models.OrderbookTop) (err error) {
	fields := map[string]interface{}{
		"timestamp":         model.Timestamp,
		"best ask":          model.BestAsk,
		"best ask quantity": model.BestAskQuantity,
		"best bid":          model.BestBid,
		"best bid quantity": model.BestBidQuantity,
	}
	return checkFields(fields)
}

func checkCandle(model *models.WSCandle) (err error) {
	fields := map[string]interface{}{
		"timestamp":   model.Timestamp,
		"open":        model.Open,
		"close":       model.Close,
		"min":         model.Low,
		"max":         model.High,
		"volume":      model.Volume,
		"volumeQuote": model.VolumeQuote,
	}
	return checkFields(fields)
}

func checkBalance(model *models.Balance) (err error) {
	fields := map[string]interface{}{
		"currency":  model.Currency,
		"available": model.Available,
		"reserved":  model.Reserved,
	}
	return checkFields(fields)
}

func checkCommission(model *models.TradingCommission) (err error) {
	fields := map[string]interface{}{
		"symbol":    model.Symbol,
		"make rate": model.MakeRate,
		"take rate": model.TakeRate,
	}
	return checkFields(fields)
}

func checkTransaction(model *models.Transaction) (err error) {
	fields := map[string]interface{}{
		"id":        model.ID,
		"status":    model.Status,
		"type":      model.Type,
		"createdAt": model.CreatedAt,
		"updatedAt": model.UpdatedAt,
	}
	return checkFields(fields)
}

func checkReport(model *models.Report) (err error) {
	fields := map[string]interface{}{
		"id":                 model.ID,
		"clientOrderId":      model.ClientOrderID,
		"symbol":             model.Symbol,
		"side":               model.Side,
		"status":             model.Status,
		"type":               model.OrderType,
		"timeInForce":        model.TimeInForce,
		"quantity":           model.Quantity,
		"price":              model.Price,
		"quantityCumulative": model.QuantityCumulative,
		"createdAt":          model.CreatedAt,
		"updatedAt":          model.UpdatedAt,
		"reportType":         model.ReportType,
	}
	return checkFields(fields)
}

func checkCandleFeed(saver *saver, feed *models.WSCandleFeed) {
	saver.strSaveCh() <- fmt.Sprint(feed)
	for _, candles := range *feed {
		for _, candle := range candles {
			if err := checkCandle(&candle); err != nil {
				saver.errSaveCh() <- err
			}
		}
	}
}

func checkMiniTickerFeed(saver *saver, feed *models.MiniTickerFeed) {
	saver.strSaveCh() <- fmt.Sprint(feed)
	for _, miniTicker := range *feed {
		if err := checkMiniTicker(&miniTicker); err != nil {
			saver.errSaveCh() <- err
		}
	}
}

func checkWSTickerFeed(saver *saver, feed *models.WSTickerFeed) {
	saver.strSaveCh() <- fmt.Sprint(feed)
	for _, ticker := range *feed {
		if err := checkWSTicker(&ticker); err != nil {
			saver.errSaveCh() <- err
		}
	}
}

func checkOrderbookFeed(obchecker *OBChecker, saver *saver, feed *models.WSOrderbookFeed) {
	saver.strSaveCh() <- fmt.Sprint(feed)
	for _, ob := range *feed {
		if err := obchecker.checkOrderbook(&ob); err != nil {
			saver.errSaveCh() <- err
		}
	}
}

func checkOrderbookTopFeed(saver *saver, feed *models.OrderbookTopFeed) {
	saver.strSaveCh() <- fmt.Sprint(feed)
	for _, ob := range *feed {
		if err := checkOrderbookTop(&ob); err != nil {
			saver.errSaveCh() <- err

		}
	}
}
