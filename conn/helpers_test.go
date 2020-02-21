package conn

import (
	"math/rand"
	"time"
	"log"
	"os"
	"bufio"
)

var keysfile = "../keys.txt"

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

