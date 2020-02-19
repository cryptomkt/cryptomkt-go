package conn

import (
	"fmt"
	"math/rand"
	"time"

	"../args"
)

var client *Client = NewClient("RandomKey", "RandomSecret")

var argus [8]args.Argument = [8]args.Argument{args.Market("ETHCLP"), args.Type("buy"), args.Type("sell"), args.Page(0), args.Limit(50), args.Start("2017-03-03"), args.End("2018-03-03"), args.TimeFrame("60")}

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

func TestMarket() error {
	if _, err := client.GetMarket(); err != nil {
		return fmt.Errorf("Market Test failed because %s", err)
	}
	return nil
}
func TestTicker() error {
	if _, err := client.GetTicker(); err != nil {
		return fmt.Errorf("Ticker with no optional arguments failed because %s", err)
	}
	if _, err := client.GetTicker(argus[0]); err != nil {
		return fmt.Errorf("Ticker with market argument failed because %s", err)
	}
	return nil
}

func TestBook(times int) error {
	rand.Seed(time.Now().UnixNano())
	var optional [2]args.Argument = [2]args.Argument{argus[3], argus[4]}
	for i := 0; i < times; i++ {
		var numArgs int = rand.Intn(3)
		switch numArgs {
		case 0:
			if _, err := client.GetBook(argus[0], argus[1]); err != nil {
				return fmt.Errorf("Book with cero optional args failed: %s", err)
			}
		case 1:
			var random int = rand.Intn(2)
			if _, err := client.GetBook(argus[0], argus[1], optional[random]); err != nil {
				return fmt.Errorf("Book with %v optional args failed: %s", 1, err)
			}
		case 2:
			if _, err := client.GetBook(argus[0], argus[1], optional[0], optional[1]); err != nil {
				return fmt.Errorf("Book with 2 optional arguments failed because %s ", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

func TestTrades(times int) error {
	rand.Seed(time.Now().UnixNano())
	var optional [4]Argument = [4]Argument{argus[3], argus[4], argus[5], argus[6]}
	for i := 0; i < times; i++ {
		var numArgs int = rand.Intn(5)
		switch numArgs {
		case 0:
			if _, err := client.GetTrades(argus[0]); err != nil {
				return fmt.Errorf("Trades with cero optional arguments failed because %s", err)
			}
		case 1:
			var randomIndex int = rand.Intn(4)
			if _, err := client.GetTrades(argus[0], argus[randomIndex]); err != nil {
				return fmt.Errorf("Trades with one optional argument failed")
			}
		case 2:
			var randomIndexes []int = generateIndexes(2, 4)
			if _, err := client.GetTrades(argus[0], optional[randomIndexes[0]], optional[randomIndexes[1]]); err != nil {
				return fmt.Errorf("Trades with 2 optional arguments failed, %s", err)
			}
		case 3:
			var randomIndexes []int = generateIndexes(3, 4)
			if _, err := client.GetTrades(argus[0], optional[randomIndexes[0]], optional[randomIndexes[1]], optional[randomIndexes[2]]); err != nil {
				return fmt.Errorf("Trades with 3 optional arguments failed, %s", err)
			}
		case 4:
			if _, err := client.GetTrades(argus[0], optional[0], optional[1], optional[2], optional[3]); err != nil {
				return fmt.Errorf("Trades with 4 optional args failed %s", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

func TestPrices(times int) error {
	rand.Seed(time.Now().UnixNano())
	var optional [2]Argument = [2]Argument{argus[3], argus[4]}
	for i := 0; i < times; i++ {
		var numArgs int = rand.Intn(3)
		switch numArgs {
		case 0:
			if _, err := client.GetPrices(argus[0], argus[7]); err != nil {
				return fmt.Errorf("Prices with zero optional args failed, %s", err)
			}
		case 1:
			var randomIndex int = rand.Intn(2)
			if _, err := client.GetPrices(argus[0], argus[7], optional[randomIndex]); err != nil {
				return fmt.Errorf("Prices with one optional argument failed, %s", err)
			}
		case 2:
			if _, err := client.GetPrices(argus[0], argus[7], optional[0], optional[1]); err != nil {
				return fmt.Errorf("Prices with 2 optional args failed, %s", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

func TestPrevious() {
	book, err := client.GetBook(args.Market("ETHCLP"), args.Type("buy"), args.Page(1))
	if err == nil {
		book, err = book.GetPrevious()
		if err != nil {
			fmt.Println("Error: ", err, " in previous book")
		} else {
			fmt.Println("Previous book finished succesfully")
		}
	} else {
		fmt.Println("Error in book: ", err)
	}
	trades, err2 := client.GetTrades(args.Market("ETHCLP"), args.Start("2019-12-12"), args.End("2020-01-01"), args.Page(1))
	if err2 == nil {
		_, err2 = trades.GetPrevious()
		if err2 != nil {
			fmt.Println("Error: ", err, " in previous trades")
		} else {
			fmt.Println("Previous trades finished succesfully")
		}
	} else {
		fmt.Println("Error Trades: ", err)
	}
	_, err3 := client.GetPrices(args.Market("ETHCLP"), args.TimeFrame("60"), args.Page(1), args.Limit(40))
	if err3 == nil {
		_, err3 = trades.GetPrevious()
		if err3 != nil {
			fmt.Println("Error: ", err, " in previous prices")
		} else {
			fmt.Println("Previous prices finished successfully")
		}
	}
}

func TestNext() {
	book, err := client.GetBook(args.Market("ETHCLP"), args.Type("buy"), args.Page(0))
	if err == nil {
		book, err = book.GetNext()
		if err != nil {
			fmt.Println("Error: ", err, " in Next book")
		} else {
			fmt.Println("Next book finished succesfully")
		}
	} else {
		fmt.Println("Error in book: ", err)
	}
	trades, err2 := client.GetTrades(args.Market("ETHCLP"), args.Start("2019-12-12"), args.Page(0))
	if err2 == nil {
		_, err2 = trades.GetNext()
		if err2 != nil {
			fmt.Println("Error: ", err, " in Next trades")
		} else {
			fmt.Println("Next trades finished succesfully")
		}
	} else {
		fmt.Println("Error Trades: ", err)
	}
	_, err3 := GetPrices(args.Market("ETHCLP"), args.TimeFrame("10080"), args.Page(0), args.Limit(40))
	if err3 == nil {
		_, err3 = trades.GetNext()
		if err3 != nil {
			fmt.Println("Error: ", err, " in Next prices")
		} else {
			fmt.Println("Next prices finished successfully")
		}
	}
}
