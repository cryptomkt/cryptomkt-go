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