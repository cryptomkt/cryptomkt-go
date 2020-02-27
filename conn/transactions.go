package conn

import (
	"fmt"

	"github.com/cryptomkt/cryptomkt-go/args"
)

// Transactions structs needed for this sdk

type TransactionsResponse struct {
	client     *Client
	Status     string
	Message    string
	Pagination Pagination
	Data       []Transaction
}

type Transaction struct {
	Id         string
	Type       int
	Amount     string
	FeePercent string `json:"fee_percent"`
	FeeAmount  string `json:"fee_amount"`
	Balance    string
	Date       string
	Hash       string
	Address    string
	Memo       string
}

type TransactionList struct {
	currency   string
	client     *Client
	pagination Pagination
	Data       []Transaction
}

// GetPrevious get the previous page of the Transaction list
// If there is no previous page, rise an error.
func (tList *TransactionList) GetPrevious() (*TransactionList, error) {
	if tList.pagination.Previous == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	tList, err := tList.client.GetTransactions(
		args.Currency(tList.currency),
		args.Page(int(tList.pagination.Previous.(float64))),
		args.Limit(tList.pagination.Limit))
	if err != nil {
		return nil, fmt.Errorf("error getting the previous page: %s", err)
	}
	return tList, nil
}

// GetNext lets you go to the next page if it exists, returns (*Prices, nil) if
// it is successfull and (nil, error) otherwise
func (tList *TransactionList) GetNext() (*TransactionList, error) {
	if tList.pagination.Next == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	tList, err := tList.client.GetTransactions(
		args.Currency(tList.currency),
		args.Page(int(tList.pagination.Next.(float64))),
		args.Limit(tList.pagination.Limit))
	if err != nil {
		return nil, fmt.Errorf("error getting the next page: %s", err)
	}
	return tList, nil
}

// GetPage returns the actual page.
func (tList *TransactionList) GetPage() int {
	return tList.pagination.Page
}

// GetLimit returns the limit number of elements per page
func (tList *TransactionList) GetLimit() int {
	return tList.pagination.Limit
}
