package conn

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/cryptomkt/cryptomkt-go/args"
)

var keysfile = "../keys.txt"

var delay = 6 * time.Second
var argus [8]args.Argument = [8]args.Argument{args.Market("ETHCLP"), args.Type("buy"), args.Type("sell"), args.Page(0), args.Limit(50), args.Start("2017-03-03"), args.End("2018-03-03"), args.Timeframe("60")}

// newDebugClient initializes a client to run request,
// so its recomendable to not have money in the account for testing
// reads the first two lines of a file, the first one is the api key,
// the second one is the api secret
// **WARNING** DO NOT SHARE YOUR KEYS, KEEP THEM OUT OF THE REPOSITORY
// (with .gitignore for example)
func newDebugClient(filePath string) (*Client, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	apiKey := scanner.Text()
	scanner.Scan()
	apiSecret := scanner.Text()

	client := NewClient(apiKey, apiSecret)
	return client, nil
}

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
