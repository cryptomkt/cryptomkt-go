package client

import (
	"fmt"
	"testing"
)

func TestTransactions(t *testing.T) {
	client, err := newDebugClient("keys.txt")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Run("missing=Currency", func(t *testing.T) {
		_, err := client.Transactions()
		if err == nil {
			t.Errorf("no error rised, should rise missing Currency arg")
		}
		fmt.Println("should miss required argument of transaction: currency")
		fmt.Println(err)
		fmt.Println()
	})
	/*
		fails, currency has no restrictions, although it can be added
		to be CLP, ARS, BRL, MXN, EUR, ETH, XLM or BTC.
		t.Run("unsupported=Currency", func(t *testing.T) {
			_, err := client.Transactions(
				Currency("USD"),
			)
			if err == nil {
				t.Errorf("no error rised, should rise not supported currency")
			}
			fmt.Println("should inform unsupported currency")
			fmt.Println(err)
			fmt.Println()
		})
	*/
}

func TestActiveOrders(t *testing.T) {
	client, err := newDebugClient("keys.txt")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Run("missing=market", func(t *testing.T) {
		_, err := client.ActiveOrders()
		if err == nil {
			t.Errorf("no error rised, shold rise missing market arg")
		}
		fmt.Println("should miss required argument of ActiveOrders: market")
		fmt.Println(err)
		fmt.Println()
	})
}

func ExecutedOrders(t *testing.T) {
	client, err := newDebugClient("keys.txt")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Run("missing=market", func(t *testing.T) {
		_, err := client.ExecutedOrders()
		if err == nil {
			t.Errorf("no error rised, shold rise missing market arg")
		}
		fmt.Println("should miss required argument of ExecutedOrders: market")
		fmt.Println(err)
		fmt.Println()
	})
}

func OrderStatus(t *testing.T) {
	client, err := newDebugClient("keys.txt")
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Run("missing=market", func(t *testing.T) {
		_, err := client.OrderStatus()
		if err == nil {
			t.Errorf("no error rised, shold rise missing id arg")
		}
		fmt.Println("should miss required argument of OrderStatus: id")
		fmt.Println(err)
		fmt.Println()
	})
}

/*
#TODO: cuando se tengan los enpoints no autorizados,
	   probar los valores que acepta el parametro limit
	   en el listado de trades, dice que el mínimo es 20 y el
	   máximo es 100, pero puede que acepte menos o más que estos
	   valores, recuerdo haber usado 10 y que me aceptara la request
	   *Mientras*, limit se restringira a ser mayor a cero y menor a 200.
func TestLimitArgumentRestrictions(t *testing.T) {
//load the keys for the connection to crypto market.
	//apikey first, apisecret after,  every in its own line.
	file, err := os.Open("keys.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	apiKey := scanner.Text()
	scanner.Scan()
	apiSecret := scanner.Text()

	client, err := New(apiKey, apiSecret)
	if err != nil {
		t.Errorf("error making the client")
	}
}
*/
