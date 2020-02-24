package conn

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalingOrder(t *testing.T) {
	dummyJSONOrder := `{
		"status": "active",
		"created_at": "2017-09-01T14:02:36.386967",
		"amount": {
		   "original": "1.25",
		   "remaining": "1.25"
		},
		"execution_price": null,
		"price": "8000",
		"type": "buy",
		"id": "M103967",
		"market": "ETHCLP",
		"updated_at": "2017-09-01T14:02:36.386967"
	 }`
	order := &Order{}
	json.Unmarshal([]byte(dummyJSONOrder), &order)
	if order.Status != "active" {
		t.Errorf("status %v", order)
	}
	if order.CreatedAt != "2017-09-01T14:02:36.386967" {
		t.Errorf("created at %v", order.CreatedAt)
	}
	if order.Amount.Original != "1.25" {
		t.Errorf("amount original %v", order)
	}
	if order.Amount.Remaining != "1.25" {
		t.Errorf("amount remaining %v", order)
	}
	if order.ExecutionPrice != "" {
		t.Errorf("execution price %v", order.ExecutionPrice)
	}
	if order.Price != "8000" {
		t.Errorf("price %v", order)
	}
	if order.Type != "buy" {
		t.Errorf("type %v", order)
	}
	if order.Id != "M103967" {
		t.Errorf("id %v", order)
	}
	if order.Market != "ETHCLP" {
		t.Errorf("market %v", order)
	}
	if order.UpdatedAt != "2017-09-01T14:02:36.386967" {
		t.Errorf("updated at %v", order)
	}
}
