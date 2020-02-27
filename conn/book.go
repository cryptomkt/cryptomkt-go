package conn

import (
	"bytes"
	"fmt"

	"github.com/cryptomkt/cryptomkt-go/args"
)

type BookResponse struct {
	Status     string
	Message    string
	Pagination Pagination
	Data       []BookData
}

type Book struct {
	args       map[string]string
	pagination Pagination
	client     *Client
	Data       []BookData
}

func (book *Book) String() string {
	var b bytes.Buffer
	b.WriteString("Book{")
	b.WriteString("\n\tclient:")
	b.WriteString(book.client.String())
	b.WriteString("\n\tpagination:")
	b.WriteString(book.pagination.String())
	b.WriteString("\n\tdata:")
	for _, bookData := range book.Data {
		b.WriteString("\n")
		b.WriteString(bookData.String())
	}
	b.WriteString("\n}")
	return b.String()
}

type BookData struct {
	Price     string
	Amount    string
	Timestamp string
}

func (bookData *BookData) String() string {
	return "{price:" +
		bookData.Price +
		" amount:" +
		bookData.Amount +
		" Timestamp:" +
		bookData.Timestamp +
		"}"
}

// GetPrevious lets you go to the previous page if it exists, returns (*Book, nil) if
// it is successfull and (nil, error) otherwise
func (b *Book) GetPrevious() (*Book, error) {
	if b.pagination.Previous == nil {
		return nil, fmt.Errorf("Previous page does not exist")
	}
	return b.client.GetBook(
		args.Market(b.args["market"]),
		args.Type(b.args["type"]),
		args.Page(int(b.pagination.Previous.(float64))),
		args.Limit(b.pagination.Limit))
}

// GetNext lets you go to the next page if it exists, returns (*Book, nil) if
// it is successfull and (nil, error) otherwise
func (b *Book) GetNext() (*Book, error) {
	if b.pagination.Next == nil {
		return nil, fmt.Errorf("Next page does not exist")
	}
	return b.client.GetBook(
		args.Market(b.args["market"]),
		args.Type(b.args["type"]),
		args.Page(int(b.pagination.Next.(float64))),
		args.Limit(b.pagination.Limit))
}

// GetPage returns the actual page of the request.
func (b *Book) GetPage() int {
	return b.pagination.Page
}

// GetLimit returns the number of trades per page
func (b *Book) GetLimit() int {
	return b.pagination.Limit
}
