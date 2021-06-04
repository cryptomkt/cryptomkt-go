package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/cryptomkt/go-api/args"
)

func TestRequestWithDeadline(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*10))
	defer cancelFunc()
	result, err := client.getCurrency(ctx, args.Currency("EOS"))
	if err.Error() != ctx.Err().Error() { // should be err = "context deadline exeeded"
		t.Fatal("context failed to stop execution")
	}
	time.Sleep(2 *time.Second)
	if result != nil && result.ID != "" {
		t.Fatal("should be an empty result")
	}
	ctx = context.Background()
	result, err = client.getCurrency(ctx, args.Currency("EOS"))
	if err != nil {
		t.Log(err)
		t.Fatal("failed to do a good request after a bad one")
	}
	if result.ID == "" {
		t.Fatal("currency should have id")
	}
}

func TestAPIErrorHandling(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, err = client.getCurrency(ctx, args.Currency("eosas"))
	if err == nil {
		t.Fatal("should return an error")
	}
}

func TestErrorOnClosedConnection(t *testing.T) {
	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, err = client.getCurrency(ctx, args.Currency("EOS"))
	if err != nil {
		t.Fatal(err)
	}
	client.Close()
	_, err = client.getCurrency(ctx, args.Currency("EOS"))
	if err == nil {
		t.Log("should err in closed connection")
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}
