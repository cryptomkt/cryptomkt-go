package rest

import (
	"testing"
)

func TestChangeWindow(t *testing.T) {
	client, bg := beforeEach()
	result, err := client.GetWalletBallances(bg)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("should have results")
	}
	if err = checkList(checkBalance, result); err != nil {
		t.Fatal(err)
	}
	client.ChangeWindow(100)
	result, err = client.GetWalletBallances(bg)
	if err == nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatal(err)
	}
}
