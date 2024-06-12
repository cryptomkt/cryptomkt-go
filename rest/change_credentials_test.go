package rest

import (
	"testing"
)

func TestChangeCredentials(t *testing.T) {
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
	client.ChangeCredentials("", "")
	result, err = client.GetWalletBallances(bg)
	if err == nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatal(err)
	}
}
