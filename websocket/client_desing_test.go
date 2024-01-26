package websocket

import (
	"testing"
)

func TestClientConcurrency(t *testing.T) {
	// client, err := NewPublicClient()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// bg := context.Background()
	// result, err := client.GetCurrencies(bg)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// errCh := make(chan error)
	// done := make(chan struct{})
	// // for _, currency := range result[:5] {
	// // 	go func(client *PublicClient, id string) {
	// // 		defer func() {
	// // 			done <- struct{}{}
	// // 		}()
	// // 		resp, err := client.GetCurrency(bg, args.Currency(id))
	// // 		if err != nil {
	// // 			errCh <- err
	// // 			return
	// // 		}
	// // 		if err = checkCurrency(resp); err != nil {
	// // 			errCh <- err
	// // 			return
	// // 		}
	// // 	}(client, currency.ID)
	// // }
	// go func() {
	// 	for range result[:5] {
	// 		<-done
	// 	}
	// 	close(errCh)
	// }()
	// for err := range errCh {
	// 	t.Fatal(err)
	// }
	// close(done)
}
