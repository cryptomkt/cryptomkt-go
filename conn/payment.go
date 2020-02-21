package conn

import (
	"fmt"
	"github.com/cryptomkt/cryptomkt-go/args"
)

type PaymentOrder struct {
	client *Client
	Id                string
	ExternalId        string `json:"external_id"`
	Status            int
	ToReceive         string `json:"to_receive"`
	ToReceiveCurrency string `json:"to_receive_currency"`
	ExpectedAmount    string `json:"expected_amount"`
	ExpectedCurrency  string `json:"expected_currency"`
	DepositAddress    string `json:"deposit_address"`
	RefundEmail       string `json:"refund_email"`
	Qr                string
	Obs               string
	CallbackUrl       string `json:"callback_url"`
	ErrorUrl          string `json:"error_url"`
	SuccessUrl        string `json:"success_url"`
	PaymentUrl        string `json:"payment_url"`
	Remaining         int    `json:"remanining"`
	Language          string
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ServerAt          string `json:"server_at"`
}

type PaymentResponse struct {
	Status  string
	Message string
	Data    PaymentOrder
}

type PaymentOrderList struct {
	pagination Pagination
	startDate string
	endDate string
	client *Client
	Data []PaymentOrder
}

type PaymentOrdersResponse struct {
	Status     string
	Message    string
	Pagination Pagination
	Data       []PaymentOrder
}

func (po *PaymentOrderList) GetPrevious() (*PaymentOrderList, error) {
	if po.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	return po.client.PaymentOrders(
		args.StartDate(po.startDate),
		args.EndDate(po.endDate),
		args.Page(int(po.pagination.Previous.(float64))),
		args.Limit(po.pagination.Limit))
}

func (po *PaymentOrderList) GetNext() (*PaymentOrderList, error) {
	if po.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	return po.client.PaymentOrders(
		args.StartDate(po.startDate),
		args.EndDate(po.endDate),
		args.Page(int(po.pagination.Next.(float64))),
		args.Limit(po.pagination.Limit))
}

// GetAllPaymentOrders get all the payment orders between the two given dates.
// Returns an array of PaymentOrder
//
// List of accepted Arguments:
//   - required: StartDate, EndDate
//   - optional: none
func (client *Client) GetAllPaymentOrders(arguments... args.Argument) (*[]PaymentOrder, error) {
	req, err := makeReq([]string{"start_date", "end_date"}, arguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllPaymentOrders: %s", err)
	}
	neededArguments := []args.Argument{args.Page(0), args.Limit(100)}
	argsMap := req.GetArguments()
	val := argsMap["start_date"]
	neededArguments = append(neededArguments, args.StartDate(val))
	val = argsMap["end_date"]
	neededArguments = append(neededArguments, args.EndDate(val))
	
	poList, err := client.PaymentOrders(neededArguments...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetAllPaymentOrders: %s", err)
	}
	allpo := make([]PaymentOrder, len(poList.Data))
	copy(allpo, poList.Data)
	for poList, err = poList.GetNext(); err == nil; poList, err = poList.GetNext() {
		allpo = append(allpo, poList.Data...)
	}
	return &allpo, nil
}
