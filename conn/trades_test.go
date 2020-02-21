package conn

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
)

func TestTrades(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	rand.Seed(time.Now().UnixNano())
	var optional [4]args.Argument = [4]args.Argument{args.Page(0), args.Limit(50), args.Start("2017-03-03"), args.Timeframe("60")}
	for i := 0; i < 100; i++ {
		var numArgs int = rand.Intn(5)
		switch numArgs {
		case 0:
			if _, err := client.GetTrades(args.Market("ETHCLP")); err != nil {
				t.Errorf("Trades with cero optional arguments failed because %s", err)
			}
		case 1:
			var randomIndex int = rand.Intn(4)
			if _, err := client.GetTrades(args.Market("ETHCLP"), argus[randomIndex]); err != nil {
				t.Errorf("Trades with one optional argument failed")
			}
		case 2:
			var randomIndexes []int = generateIndexes(2, 4)
			if _, err := client.GetTrades(args.Market("ETHCLP"), optional[randomIndexes[0]], optional[randomIndexes[1]]); err != nil {
				t.Errorf("Trades with 2 optional arguments failed, %s", err)
			}
		case 3:
			var randomIndexes []int = generateIndexes(3, 4)
			if _, err := client.GetTrades(args.Market("ETHCLP"), optional[randomIndexes[0]], optional[randomIndexes[1]], optional[randomIndexes[2]]); err != nil {
				t.Errorf("Trades with 3 optional arguments failed, %s", err)
			}
		case 4:
			if _, err := client.GetTrades(args.Market("ETHCLP"), optional[0], optional[1], optional[2], optional[3]); err != nil {
				t.Errorf("Trades with 4 optional args failed %s", err)
			}
		}
		time.Sleep(3 * time.Second)
	}

}

func TestTradesGetPrevious(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	trades, err := client.GetTrades(args.Market("ETHCLP"), args.Start("2019-12-12"), args.End("2020-01-01"), args.Page(1))
	if err != nil {
		t.Errorf("Error Trades: %s", err)
	}
	_, err = trades.GetPrevious()
	if err != nil {
		t.Errorf("Error in previous trades: %s", err)
	}
}

func TestTradesGetNext(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	trades, err := client.GetTrades(args.Market("ETHCLP"), args.Start("2019-12-12"), args.End("2020-01-01"), args.Page(1))
	if err != nil {
		t.Errorf("Error Trades: %s", err)
	}
	_, err = trades.GetNext()
	if err != nil {
		t.Errorf("Error in next trades: %s", err)
	}
}


func TestGetAllTrades(t *testing.T) {
	client := NewClient("NoKey", "NoSecret")
	rand.Seed(time.Now().UnixNano())
	var optional [4]args.Argument = [4]args.Argument{argus[3], argus[4], argus[5], argus[6]}
	for i := 0; i < 100; i++ {
		var numArgs int = rand.Intn(5)
		switch numArgs {
		case 0:
			if trades, err := client.GetTrades(args.Market("ETHCLP")); err != nil {
				t.Errorf("Trades with zero optional arguments failed because %s", err)
				_, err2 := trades.GetAllTrades()
				if err2 != nil {
					t.Errorf("Error getting all trades with zero optional args, %s", err)
				}

			}
		case 1:
			var randomIndex int = rand.Intn(4)
			if trades, err := client.GetTrades(args.Market("ETHCLP"), argus[randomIndex]); err != nil {
				t.Errorf("Trades with one optional argument failed")
				_, err2 := trades.GetAllTrades()
				if err2 != nil {
					t.Errorf("Error getting all trades with one optional arg, %s", err)
				}
			}
		case 2:
			var randomIndexes []int = generateIndexes(2, 4)
			if trades, err := client.GetTrades(args.Market("ETHCLP"), optional[randomIndexes[0]], optional[randomIndexes[1]]); err != nil {
				t.Errorf("Trades with 2 optional arguments failed, %s", err)
				_, err2 := trades.GetAllTrades()
				if err2 != nil {
					t.Errorf("Error getting all trades with two optional args, %s", err)
				}
			}
		case 3:
			var randomIndexes []int = generateIndexes(3, 4)
			if trades, err := client.GetTrades(args.Market("ETHCLP"), optional[randomIndexes[0]], optional[randomIndexes[1]], optional[randomIndexes[2]]); err != nil {
				t.Errorf("Trades with 3 optional arguments failed, %s", err)
				_, err2 := trades.GetAllTrades()
				if err2 != nil {
					t.Errorf("Error getting all trades with three optional args, %s", err)
				}
			}
		case 4:
			if trades, err := client.GetTrades(args.Market("ETHCLP"), optional[0], optional[1], optional[2], optional[3]); err != nil {
				t.Errorf("Trades with 4 optional args failed %s", err)
				_, err2 := trades.GetAllTrades()
				if err2 != nil {
					t.Errorf("Error getting all trades with four optional args, %s", err)
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}