package conn

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
)

var client *Client = NewClient("RandomKey", "RandomSecret")

var argus [8]args.Argument = [8]args.Argument{args.Market("ETHCLP"), args.Type("buy"), args.Type("sell"), args.Page(0), args.Limit(50), args.Start("2017-03-03"), args.End("2018-03-03"), args.Timeframe("60")}

func find(slice []int, value int) bool {
	for i := 0; i < len(slice); i++ {
		if value == slice[i] {
			time.Sleep(time.Second)
			return true
		}
	}
	return false
}

func generateIndexes(cantidad int, rango int) []int {
	rand.Seed(time.Now().UnixNano())
	var resp []int = make([]int, cantidad)
	for i := 1; i < cantidad; i++ {
		var value int = rand.Intn(rango)
		var b bool = find(resp, value)
		for b {
			value = rand.Intn(rango)
			b = find(resp, value)
		}
		resp[i] = value
	}
	return resp
}

func TestMarket(t *testing.T) {
	if _, err := client.GetMarkets(); err != nil {
		t.Errorf("Market Test failed because %s", err)
	}
}
func TestTicker(t *testing.T) {
	if _, err := client.GetTicker(); err != nil {
		t.Errorf("Ticker with no optional arguments failed because %s", err)
	}
	if _, err := client.GetTicker(argus[0]); err != nil {
		t.Errorf("Ticker with market argument failed because %s", err)
	}
}

func TestBook(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var optional [2]args.Argument = [2]args.Argument{argus[3], argus[4]}
	for i := 0; i < 100; i++ { //here you can change the number of repetitions
		var numArgs int = rand.Intn(3)
		switch numArgs {
		case 0:
			if _, err := client.GetBook(argus[0], argus[1]); err != nil {
				t.Errorf("Book with cero optional args failed: %s", err)
			}
		case 1:
			var random int = rand.Intn(2)
			if _, err := client.GetBook(argus[0], argus[1], optional[random]); err != nil {
				t.Errorf("Book with %v optional args failed: %s", 1, err)
			}
		case 2:
			if _, err := client.GetBook(argus[0], argus[1], optional[0], optional[1]); err != nil {
				t.Errorf("Book with 2 optional arguments failed because %s ", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func TestTrades(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var optional [4]args.Argument = [4]args.Argument{argus[3], argus[4], argus[5], argus[6]}
	for i := 0; i < 100; i++ {
		var numArgs int = rand.Intn(5)
		switch numArgs {
		case 0:
			if _, err := client.GetTrades(argus[0]); err != nil {
				t.Errorf("Trades with cero optional arguments failed because %s", err)
			}
		case 1:
			var randomIndex int = rand.Intn(4)
			if _, err := client.GetTrades(argus[0], argus[randomIndex]); err != nil {
				t.Errorf("Trades with one optional argument failed")
			}
		case 2:
			var randomIndexes []int = generateIndexes(2, 4)
			if _, err := client.GetTrades(argus[0], optional[randomIndexes[0]], optional[randomIndexes[1]]); err != nil {
				t.Errorf("Trades with 2 optional arguments failed, %s", err)
			}
		case 3:
			var randomIndexes []int = generateIndexes(3, 4)
			if _, err := client.GetTrades(argus[0], optional[randomIndexes[0]], optional[randomIndexes[1]], optional[randomIndexes[2]]); err != nil {
				t.Errorf("Trades with 3 optional arguments failed, %s", err)
			}
		case 4:
			if _, err := client.GetTrades(argus[0], optional[0], optional[1], optional[2], optional[3]); err != nil {
				t.Errorf("Trades with 4 optional args failed %s", err)
			}
		}
		time.Sleep(3 * time.Second)
	}

}

func TestPrices(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var optional [2]args.Argument = [2]args.Argument{argus[3], argus[4]}
	for i := 0; i < 100; i++ {
		var numArgs int = rand.Intn(3)
		switch numArgs {
		case 0:
			if _, err := client.GetPrices(argus[0], argus[7]); err != nil {
				t.Errorf("Prices with zero optional args failed, %s", err)
			}
		case 1:
			var randomIndex int = rand.Intn(2)
			if _, err := client.GetPrices(argus[0], argus[7], optional[randomIndex]); err != nil {
				t.Errorf("Prices with one optional argument failed, %s", err)
			}
		case 2:
			if _, err := client.GetPrices(argus[0], argus[7], optional[0], optional[1]); err != nil {
				t.Errorf("Prices with 2 optional args failed, %s", err)
			}
		}
		time.Sleep(3 * time.Second)
	}

}

func TestPrevious(t *testing.T) {
	book, err := client.GetBook(args.Market("ETHCLP"), args.Type("buy"), args.Page(1))
	if err == nil {
		book, err = book.GetPrevious()
		if err != nil {
			t.Errorf("Error: %s in previous book", err)
		} else {
			fmt.Println("Previous book finished succesfully")
		}
	} else {
		t.Errorf("Error in book: %s", err)
	}
	trades, err2 := client.GetTrades(args.Market("ETHCLP"), args.Start("2019-12-12"), args.End("2020-01-01"), args.Page(1))
	if err2 == nil {
		_, err2 = trades.GetPrevious()
		if err2 != nil {
			t.Errorf("Error: %s in previous trades", err2)
		}
	} else {
		t.Errorf("Error Trades: %s", err)
	}
	_, err3 := client.GetPrices(args.Market("ETHCLP"), args.Timeframe("60"), args.Page(1), args.Limit(40))
	if err3 == nil {
		_, err3 = trades.GetPrevious()
		if err3 != nil {
			t.Errorf("Error: %s in previous prices", err3)
		}
	}
}

func TestNext(t *testing.T) {
	book, err := client.GetBook(args.Market("ETHCLP"), args.Type("buy"), args.Page(0))
	if err == nil {
		book, err = book.GetNext()
		if err != nil {
			t.Errorf("Error: %s in Next book", err)
		}
	} else {
		t.Errorf("Error in book: %s", err)
	}
	trades, err2 := client.GetTrades(args.Market("ETHCLP"), args.Start("2019-12-12"), args.Page(0))
	if err2 == nil {
		_, err2 = trades.GetNext()
		if err2 != nil {
			t.Errorf("Error: %s in Next trades", err2)
		}
	} else {
		t.Errorf("Error Trades: %s", err)
	}
	_, err3 := client.GetPrices(args.Market("ETHCLP"), args.Timeframe("10080"), args.Page(0), args.Limit(40))
	if err3 == nil {
		_, err3 = trades.GetNext()
		if err3 != nil {
			t.Errorf("Error: %s in Next prices", err3)
		}
	}
}
