package rest

import (
	"context"
	"testing"

	"github.com/cryptomarket/cryptomarket-go/args"
)

func TestGetCurrencies(t *testing.T) {
	client := NewPublicClient()
	t.Run("all currencies", func(t *testing.T) {
		result, err := client.GetCurrencies(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		for _, currency := range result {
			if err = checkCurrency(&currency); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("some currencies", func(t *testing.T) {
		result, err := client.GetCurrencies(context.Background(), args.Currencies([]string{"EOS", "ETH"}))
		if err != nil {
			t.Fatal(err)
		}
		for _, currency := range result {
			if err = checkCurrency(&currency); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetCurrency(t *testing.T) {
	client := NewPublicClient()
	t.Run("valid currency", func(t *testing.T) {
		result, err := client.GetCurrency(context.Background(), args.Currency("EOS"))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkCurrency(result); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid currency", func(t *testing.T) {
		result, err := client.GetCurrency(context.Background(), args.Currency("classic"))
		if err != nil {
			return
		}
		if err = checkCurrency(result); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("no currency", func(t *testing.T) {
		result, err := client.GetCurrency(context.Background(), args.Offset(30))
		if err != nil {
			return
		}
		t.Error(result)
	})
}

func TestGetSymbols(t *testing.T) {
	client := NewPublicClient()
	t.Run("all symbols", func(t *testing.T) {
		result, err := client.GetSymbols(context.Background())
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, symbol := range result {
			if err = checkSymbol(&symbol); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("some symbols", func(t *testing.T) {
		result, err := client.GetSymbols(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, symbol := range result {
			if err = checkSymbol(&symbol); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetSymbol(t *testing.T) {
	client := NewPublicClient()
	t.Run("valid symbol", func(t *testing.T) {
		result, err := client.GetSymbol(context.Background(), args.Symbol("EOSETH"))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkSymbol(result); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("invalid symbol", func(t *testing.T) {
		result, err := client.GetSymbol(context.Background(), args.Symbol("tree"))
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("no symbol", func(t *testing.T) {
		result, err := client.GetSymbol(context.Background())
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("junk arguments", func(t *testing.T) {
		// should still work
		result, err := client.GetSymbol(context.Background(), args.Symbol("ETHBTC"), args.Currency("EOS"), args.Offset(19))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkSymbol(result); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetTickers(t *testing.T) {
	client := NewPublicClient()
	t.Run("all tickers", func(t *testing.T) {
		result, err := client.GetTickers(context.Background())
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, ticker := range result {
			if err = checkTicker(&ticker); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("some tickers", func(t *testing.T) {
		result, err := client.GetTickers(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, ticker := range result {
			if err = checkTicker(&ticker); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetTicker(t *testing.T) {
	client := NewPublicClient()
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.GetTickerOfSymbol(context.Background(), args.Symbol("EOSETH"))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkTicker(result); err != nil {
			t.Fatal(err)
		}

	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.GetTickerOfSymbol(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
}

func TestGetPrice(t *testing.T) {
	client := NewPublicClient()
	t.Run("all prices", func(t *testing.T) {
		result, err := client.GetPrices(context.Background(), args.To("ETH"))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, price := range result {
			if err = checkQuotationPrice(&price); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("one price", func(t *testing.T) {
		result, err := client.GetPrices(context.Background(), args.From("BTC"), args.To("ETH"))
		if err != nil {
			t.Fatal(err)
			return
		}
		if len(result) != 1 {
			t.Fatal("should only have one price")
		}
		for _, price := range result {
			if err = checkQuotationPrice(&price); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetPriceHistory(t *testing.T) {
	client := NewPublicClient()
	t.Run("all prices histories", func(t *testing.T) {
		result, err := client.GetPricesHistory(context.Background(), args.To("ETH"), args.Limit(3))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, priceHistory := range result {
			if err = checkQuotationPriceHistory(&priceHistory); err != nil {
				t.Fatal(err)
			}
			if len(priceHistory.History) != 3 {
				t.Fatal("should have 3 history points")
			}
		}
	})
}

func TestGetTickerLastPrices(t *testing.T) {
	client := NewPublicClient()
	t.Run("all prices histories", func(t *testing.T) {
		result, err := client.GetTickerLastPrices(context.Background(), args.Symbols([]string{"ETHBTC", "EOSETH"}))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, tickerPrice := range result {
			if err = checkTickerPrice(&tickerPrice); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetTickerLastPricesOfSymbol(t *testing.T) {
	client := NewPublicClient()
	t.Run("all prices histories", func(t *testing.T) {
		result, err := client.GetTickerLastPricesOfSymbol(context.Background(), args.Symbol("EOSETH"))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkTickerPrice(result); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetTrades(t *testing.T) {
	client := NewPublicClient()
	t.Run("from all symbols, no arguments", func(t *testing.T) {
		result, err := client.GetTrades(context.Background())
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, trades := range result {
			for _, trade := range trades {
				if err = checkPublicTrade(&trade); err != nil {
					t.Fatal(err)
				}
			}
		}
	})
	t.Run("from some symbols, limit at 2", func(t *testing.T) {
		result, err := client.GetTrades(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(2), args.Offset(10))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, trades := range result {
			for _, trade := range trades {
				if err = checkPublicTrade(&trade); err != nil {
					t.Fatal(err)
				}
			}
		}
	})
}

func TestGetTradesOfSymbol(t *testing.T) {
	client := NewPublicClient()
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.GetTradesOfSymbol(context.Background(), args.Symbol("EOSETH"), args.Limit(2))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, trade := range result {
			if err = checkPublicTrade(&trade); err != nil {
				t.Fatal(err)
			}
		}

	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.GetTradesOfSymbol(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("from one symbol, limit 10, filter by id", func(t *testing.T) {
		result, err := client.GetTradesOfSymbol(context.Background(), args.Symbol("ETHBTC"), args.Limit(50), args.From("1085615118"), args.SortBy(args.SortByID))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, trade := range result {
			if err = checkPublicTrade(&trade); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetOrderbooks(t *testing.T) {
	client := NewPublicClient()
	t.Run("from all symbols, no arguments", func(t *testing.T) {
		result, err := client.GetOrderbooks(context.Background())
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, orderbook := range result {
			if err = checkOrderbook(&orderbook); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("from some symbols, limit at 2", func(t *testing.T) {
		result, err := client.GetOrderbooks(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(2))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, orderbook := range result {
			if err = checkOrderbook(&orderbook); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestGetOrderbook(t *testing.T) {
	client := NewPublicClient()
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.GetOrderBookOfSymbol(context.Background(), args.Symbol("EOSETH"), args.Limit(2))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkOrderbook(result); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.GetOrderBookOfSymbol(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
	t.Run("with volume", func(t *testing.T) {
		result, err := client.GetOrderBookOfSymbol(context.Background(), args.Symbol("EOSETH"), args.Volume("3"))
		if err != nil {
			t.Fatal(err)
			return
		}
		if err = checkOrderbook(result); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetCandles(t *testing.T) {
	client := NewPublicClient()
	t.Run("from all symbols, no arguments", func(t *testing.T) {
		result, err := client.GetCandles(context.Background())
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, candles := range result {
			for _, candle := range candles {
				if err = checkCandle(&candle); err != nil {
					t.Fatal(err)
				}
			}
		}
	})
	t.Run("from some symbols, limit at 2", func(t *testing.T) {
		result, err := client.GetCandles(context.Background(), args.Symbols([]string{"EOSETH", "ETHBTC"}), args.Limit(2))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, candles := range result {
			for _, candle := range candles {
				if err = checkCandle(&candle); err != nil {
					t.Fatal(err)
				}
			}
		}
	})
}

func TestGetCandlesOfSymbol(t *testing.T) {
	client := NewPublicClient()
	t.Run("from valid symbol", func(t *testing.T) {
		result, err := client.GetCandlesOfSymbol(context.Background(), args.Symbol("EOSETH"), args.Limit(2))
		if err != nil {
			t.Fatal(err)
			return
		}
		for _, candle := range result {
			if err = checkCandle(&candle); err != nil {
				t.Fatal(err)
			}
		}
	})
	t.Run("from invalid symbol", func(t *testing.T) {
		result, err := client.GetOrderBookOfSymbol(context.Background(), args.Symbol("orange"))
		if err != nil {
			return
		}
		t.Error(result)
	})
}
