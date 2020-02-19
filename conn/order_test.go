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

func TestUnmarshalingAnOrderList(t *testing.T) {
	dummyJSONOrder := `{
		"status": "success",
		"pagination": {
			"previous": 3,
			"limit": 20,
			"page": 0,
			"next": "null"		   
		},
		"data": [
		   {
			  "status": "executed",
			  "created_at": "2017-08-31T21:37:42.282102",
			  "amount": {
				 "executed": "0.6",
				 "original": "3.75"
			  },
			  "execution_price": "8000",
			  "executed_at": "2017-08-31T22:01:19.481403",
			  "price": "8000",
			  "type": "buy",
			  "id": "M103959",
			  "market": "ETHCLP"
		   },
		   {
			  "status": "executed",
			  "created_at": "2017-08-31T21:37:42.282102",
			  "amount": {
				 "executed": "0.5",
				 "original": "3.75"
			  },
			  "execution_price": "8000",
			  "executed_at": "2017-08-31T22:00:13.805482",
			  "price": "8000",
			  "type": "buy",
			  "id": "M103959",
			  "market": "ETHCLP"
		   },
		   {
			  "status": "executed",
			  "created_at": "2016-11-26T23:27:54.502024",
			  "amount": {
				 "executed": "1.5772",
				 "original": "1.5772"
			  },
			  "execution_price": "6340",
			  "executed_at": "2017-01-02T22:56:03.897534",
			  "price": "6340",
			  "type": "buy",
			  "id": "M103260",
			  "market": "ETHCLP"
		   }
		]
	 }`
	orders := &OrderList{}
	err := json.Unmarshal([]byte(dummyJSONOrder), orders)
	if err != nil {
		t.Errorf("error unmarshling %s", err)
	}
	if orders.Pagination.NextHolder == "null" {
		orders.Pagination.Next = -1
	} else {
		orders.Pagination.Next = int(orders.Pagination.NextHolder.(float64))
	}

	if orders.Pagination.PreviousHolder == "null" {
		orders.Pagination.Previous = -1
	} else {
		orders.Pagination.Previous = int(orders.Pagination.PreviousHolder.(float64))
	}


	if orders.Status != "success" {
		t.Errorf("status %v", orders)
	}

	if orders.Pagination.Previous != 3 {
		t.Errorf("previous page should be 3 %v", orders)
	}

	if orders.Pagination.Next != -1 {
		t.Errorf("next page should be -1 %v", orders)
	}
}

/* TODO
func TestOrderFlow(t *testing.T) {
	client, err := newDebugClient(keysfile)
	if err != nil {
		t.Errorf("%s", err)
	}
	
}
*/