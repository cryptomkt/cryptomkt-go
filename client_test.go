package main
import (
    "fmt"
    "os"
    "log"
	"bufio"
	"testing"
)
func TestGetAutenticated(t *testing.T) {
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
		fmt.Println("error making the client")
	}

	t.Run("account", func(t *testing.T){
		response := client.getAccount()
		fmt.Println(response)
		if len(response) == 0 {
			t.Errorf("Return msg of len 0")
		}
	})
}
