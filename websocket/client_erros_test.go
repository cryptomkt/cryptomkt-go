package websocket

import (
	"testing"
)

func TestRequestWithDeadline(t *testing.T) {
	// client, err := NewMarketDataClient()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*10))
	// defer cancelFunc()
	// result, err := client.GetCurrency(ctx, args.Currency("EOS"))
	// if err.Error() != ctx.Err().Error() { // should be err = "context deadline exeeded"
	// 	t.Fatal("context failed to stop execution")
	// }
	// if err := checkCurrency(result); err == nil {
	// 	t.Log(result)
	// 	t.Fatal("should be an invalid currency")
	// }
	// result, err = client.GetCurrency(context.Background(), args.Currency("EOS"))
	// if err != nil {
	// 	t.Log(err)
	// 	t.Fatal("failed to do a good request after a bad one")
	// }
	// if err := checkCurrency(result); err != nil {
	// 	t.Fatal("should be a valid currency")
	// }
}

func TestAPIErrorHandling(t *testing.T) {
	// client, err := NewPublicClient()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ctx := context.Background()
	// _, err = client.GetCurrency(ctx, args.Currency("eosas"))
	// if err == nil {
	// 	t.Fatal("should return an error")
	// }
}

func TestErrorOnClosedConnection(t *testing.T) {
	// client, err := NewPublicClient()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ctx := context.Background()
	// if _, err = client.GetCurrency(ctx, args.Currency("EOS")); err != nil {
	// 	t.Fatal(err)
	// }
	// client.Close()
	// if _, err = client.GetCurrency(ctx, args.Currency("EOS")); err == nil {
	// 	t.Log("should err in closed connection")
	// 	t.Fatal(err)
	// }
}
