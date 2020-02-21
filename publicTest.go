package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
	"github.com/cryptomkt/cryptomkt-go/conn"
)

var argus [8]args.Argument = [8]args.Argument{args.Market("ETHCLP"), args.Type("buy"), args.Type("sell"), args.Page(0), args.Limit(50), args.Start("03/03/2017"), args.End("03/03/2018"), args.Timeframe("60")}

var client *conn.Client = conn.NewClient("YourApiKey", "YourSecretKey")

func find(slice []int, value int) bool {
	for i := 0; i < len(slice); i++ {
		if value == slice[i] {
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

func testBook(times int) error {
	rand.Seed(time.Now().UnixNano())
	var optional [2]args.Argument = [2]args.Argument{argus[3], argus[4]}
	for i := 0; i < times; i++ { //here you can change the number of repetitions
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

func testTrades(times int) error {
	rand.Seed(time.Now().UnixNano())
	var optional [4]args.Argument = [4]args.Argument{argus[3], argus[4], argus[5], argus[6]}
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

func testPrices(times int) error {
	rand.Seed(time.Now().UnixNano())
	var optional [2]args.Argument = [2]args.Argument{argus[3], argus[4]}
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

func main() {
	fmt.Println("Testing trades . . .")
	err := testTrades(500)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Trades test successfull")
	}
	/*fmt.Println("Testing prices . . .")
	err = testPrices(500)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Prices test successfull")
	}
	*/
}
