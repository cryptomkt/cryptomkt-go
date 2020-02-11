package client
import (
    "os"
    "log"
	"bufio"
	"testing"
	"encoding/json"
)
func assertSuccess(response string, t *testing.T) {
	var resp map[string] interface{}
		json.Unmarshal([]byte(response), &resp)
		if val, ok := resp["status"]; ok {
			switch val {
			case "error":
				t.Errorf(response)
			case "success":
				//all good
			default:
				t.Errorf("unexpected response")
			}
		} else {
			t.Errorf("error in the response")
		}
}

func TestAutenticated(t *testing.T) {
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
	
	client, err := NewClient(apiKey, apiSecret)
	if err != nil {
		t.Errorf("error making the client")
	}
	t.Run("account", func(t *testing.T){
		response := client.getAccount()
		assertSuccess(response, t)	
	})
	t.Run("wallet", func(t *testing.T){
		response := client.getBalance()
		assertSuccess(response, t)	
	})
}
