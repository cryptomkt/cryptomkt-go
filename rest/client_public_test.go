package rest

import (
	"context"
	"testing"

	"github.com/cryptomkt/go-api/args"
)

func TestGetCurrencies(t *testing.T) {
	client := NewClient("", "")
	t.Run("all currencies", func(t *testing.T) {
		result, err := client.getCurrencies(context.Background())
		if err != nil {
			t.Error(err)
		}
		for _, currency := range result {
			if err = checkCurrency(&currency); err != nil {
				t.Error(err)
			}
		}
	})
	t.Run("some currencies", func(t *testing.T) {
		result, err := client.getCurrencies(context.Background(), args.Currencies([]string{"EOS", "ETH"}))
		if err != nil {
			t.Error(err)
		}
		for _, currency := range result {
			if err = checkCurrency(&currency); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestGetCurrency(t *testing.T) {
	client := NewClient("", "")
	t.Run("valid currency", func(t *testing.T) {
		result, err := client.getCurrency(context.Background(), args.Currency("EOS"))
		if err != nil {
			t.Error(err)
			return
		}
		if err = checkCurrency(result); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid currency", func(t *testing.T) {
		result, err := client.getCurrency(context.Background(), args.Currency("classic"))
		if err != nil {
			return
		}
		if err = checkCurrency(result); err != nil {
			t.Error(err)
		}
	})

	t.Run("no currency", func(t *testing.T) {
		result, err := client.getCurrency(context.Background(), args.Offset(30))
		if err != nil {
			return
		}
		t.Error(result)
	})
}

func TestGetSymbols(t *testing.T) {
	client := NewClient("", "")
	t.Run("all symbols", func(t *testing.T) {
		result, err := client.getSymbols(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		for _, symbol := range result {
			if err = checkSymbol(&symbol); err != nil {
				t.Error(err)
			}
		}
	})
	t.Run("some symbols", func(t *testing.T) {
		result, err := client.getSymbols(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}))
		if err != nil {
			t.Error(err)
			return
		}
		for _, symbol := range result {
			if err = checkSymbol(&symbol); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestGetSymbol(t *testing.T) {
	client := NewClient("", "")
	t.Run("valid symbol", func(t *testing.T) {
		result, err := client.getSymbol(context.Background(), args.Symbol("EOSETH"))
		if err != nil {
			t.Error(err)
			return
		}
		if err = checkSymbol(result); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid symbol", func(t *testing.T) {
		result, err := client.getSymbol(context.Background(), args.Symbol("tree"))
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("no symbol", func(t *testing.T) {
		result, err := client.getSymbol(context.Background())
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("junk arguments", func(t *testing.T) {
		// should still work
		result, err := client.getSymbol(context.Background(), args.Symbol("ETHBTC"), args.Currency("EOS"), args.Offset(19))
		if err != nil {
			t.Error(err)
			return
		}
		if err = checkSymbol(result); err != nil {
			t.Error(err)
		}
	})
}

func TestGetTickers(t *testing.T) {
	client := NewClient("", "")
	t.Run("all tickers", func(t *testing.T) {
		result, err := client.getTickers(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		nullTickers := map[string]bool{
			"BTCEUR": true,
			"ETHCLP": true,
			"ETHMXN": true,
			"ETHUYU": true,
			"ETHBRL": true,
			"BTCARS": true,
			"ETHCOP": true,
			"BTCUYU": true,
			"BTCCOP": true,
			"ETHVEF": true,
			"ETHARS": true,
			"ETHPEN": true,
			"BTCCLP": true,
			"BTCMXN": true,
			"BTCBRL": true,
			"ETHEUR": true,
			"BTCPEN": true,
			"BTCVEF": true,
		}
		for _, ticker := range result {
			if _, ok := nullTickers[ticker.Symbol]; ok {
				continue
			}
			if err = checkTicker(&ticker); err != nil {
				t.Error(err)
			}
		}
	})
	t.Run("some tickers", func(t *testing.T) {
		result, err := client.getTickers(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}))
		if err != nil {
			t.Error(err)
			return
		}
		for _, ticker := range result {
			if err = checkTicker(&ticker); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestGetTicker(t *testing.T) {
	client := NewClient("", "")
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.getTicker(context.Background(), args.Symbol("EOSETH"))
		if err != nil {
			t.Error(err)
			return
		}
		if err = checkTicker(result); err != nil {
			t.Error(err)
		}

	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.getTicker(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
}

func TestGetTrades(t *testing.T) {
	client := NewClient("", "")
	t.Run("from all symbols, no arguments", func(t *testing.T) {
		result, err := client.getTrades(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		for _, trades := range result {
			for _, trade := range trades {
				if err = checkPublicTrade(&trade); err != nil {
					t.Error(err)
				}
			}
		}
	})
	t.Run("from some symbols, limit at 2", func(t *testing.T) {
		result, err := client.getTrades(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(2))
		if err != nil {
			t.Error(err)
			return
		}
		for _, trades := range result {
			for _, trade := range trades {
				if err = checkPublicTrade(&trade); err != nil {
					t.Error(err)
				}
			}
		}
	})
}

func TestGetTradesOfSymbol(t *testing.T) {
	client := NewClient("", "")
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.getTradesOfSymbol(context.Background(), args.Symbol("EOSETH"), args.Limit(2))
		if err != nil {
			t.Error(err)
			return
		}
		for _, trade := range result {
			if err = checkPublicTrade(&trade); err != nil {
				t.Error(err)
			}
		}

	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.getTradesOfSymbol(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("from one symbol, limit 10, filter by id", func(t *testing.T) {
		result, err := client.getTradesOfSymbol(context.Background(), args.Symbol("ETHBTC"), args.Limit(50), args.From("1085615118"), args.By(args.SortByID))
		if err != nil {
			t.Error(err)
			return
		}
		for _, trade := range result {
			if err = checkPublicTrade(&trade); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestGetOrderbooks(t *testing.T) {
	client := NewClient("", "")
	t.Run("from all symbols, no arguments", func(t *testing.T) {
		result, err := client.getOrderbooks(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		for _, orderbook := range result {
			if err = checkOrderbook(&orderbook); err != nil {
				t.Error(err)
			}
		}
	})
	t.Run("from some symbols, limit at 2", func(t *testing.T) {
		result, err := client.getOrderbooks(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(2))
		if err != nil {
			t.Error(err)
			return
		}
		for _, orderbook := range result {
			if err = checkOrderbook(&orderbook); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestGetOrderbook(t *testing.T) {
	client := NewClient("", "")
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.getOrderbook(context.Background(), args.Symbol("EOSETH"), args.Limit(2))
		if err != nil {
			t.Error(err)
			return
		}
		if err = checkOrderbook(result); err != nil {
			t.Error(err)
		}
	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.getOrderbook(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("with volume", func(t *testing.T) {
		result, err := client.getOrderbook(context.Background(), args.Symbol("EOSETH"), args.Volume("3"))
		if err != nil {
			t.Error(err)
			return
		}
		if err = checkOrderbook(result); err != nil {
			t.Error(err)
		}
	})
}

func TestGetCandles(t *testing.T) {
	client := NewClient("", "")
	t.Run("from all symbols, no arguments", func(t *testing.T) {
		result, err := client.getCandles(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		for _, candles := range result {
			for _, candle := range candles {
				if err = checkCandle(&candle); err != nil {
					t.Error(err)
				}
			}
		}
	})
	t.Run("from some symbols, limit at 2", func(t *testing.T) {
		result, err := client.getCandles(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(2))
		if err != nil {
			t.Error(err)
			return
		}
		for _, candles := range result {
			for _, candle := range candles {
				if err = checkCandle(&candle); err != nil {
					t.Error(err)
				}
			}
		}
	})
}

func TestGetCandlesOfSymbol(t *testing.T) {
	client := NewClient("", "")
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.getCandlesOfSymbol(context.Background(), args.Symbol("EOSETH"), args.Limit(2))
		if err != nil {
			t.Error(err)
			return
		}
		for _, candle := range result {
			if err = checkCandle(&candle); err != nil {
				t.Error(err)
			}
		}
	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.getOrderbook(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
}
