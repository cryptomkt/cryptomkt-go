// package args holds the argument logic for the requests the client
// will use to comunicate with Cryptomarket. Separate package between
// args and the request is prefered as its crearer for the user to use
// them if they are called args.AnArgument, instead of request.AnArgument.
package args

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cryptomkt/cryptomkt-go/requests"
)

type DateError struct {
	caller     string
	dateFormat string
}

func (de *DateError) Error() string {
	return de.caller + "format error: must be " + de.dateFormat
}

// An Argument is a function that servers the porpouse of arguments for a
// requests.
// Works by modifying the given request, creating the corresponding data there.
type Argument func(*requests.Request) error

// assertDateFormatV2 assert the format yyyy-mm-dd of a date string
// and returns an error if fails.
func assertDateFormatV2(val, caller string) error {
	dateFormat := "yyyy-mm-dd"
	if len(val) != 10 {
		return &DateError{caller, dateFormat}
	}
	day, err := strconv.Atoi(val[8:10])
	if err != nil {
		return &DateError{caller, dateFormat}
	}
	if day < 1 || 31 < day {
		return fmt.Errorf("%s format error: invalid day value", caller)
	}
	month, err := strconv.Atoi(val[5:7])
	if err != nil {
		return &DateError{caller, dateFormat}
	}
	if month < 1 || 12 < month {
		return fmt.Errorf("%s format error: invalid month value", caller)
	}
	_, err = strconv.Atoi(val[0:4])
	if err != nil {
		return &DateError{caller, dateFormat}
	}
	return nil
}

// assertDateFormatV1 assert the format dd/mm/yyyy of a date string
// and returns an error if fails.
func assertDateFormatV1(val, caller string) error {
	dateFormat := "dd/mm/yyyy"
	if len(val) != 10 {
		return &DateError{caller, dateFormat}
	}
	day, err := strconv.Atoi(val[0:2])
	if err != nil {
		return &DateError{caller, dateFormat}
	}
	if day < 1 || 31 < day {
		return fmt.Errorf("%s format error: invalid day value", caller)
	}
	month, err := strconv.Atoi(val[3:5])
	if err != nil {
		return &DateError{caller, dateFormat}
	}
	if month < 1 || 12 < month {
		return fmt.Errorf("%s format error: invalid month value", caller)
	}
	_, err = strconv.Atoi(val[6:10])
	if err != nil {
		return &DateError{caller, dateFormat}
	}
	return nil
}

// Amount is an argument of a request.
// Represents numbers as strings.
//
// Number format: without thousand separator, and . (a dot) as decimal point.
func Amount(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("amount", val)
		return nil
	}
}

// Market is an argument of a request.
//
// Accepts a par of currencies. e.g. "ETHCLP" or "BTCARS".
func Market(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("market", val)
		return nil
	}
}

// Type is an argument of a request.
//
// Accepts either "buy" or "sell".
func Type(val string) Argument {
	return func(request *requests.Request) error {
		if !(val == "buy" || val == "sell") {
			return errors.New("type must be either \"buy\" or \"sell\"")
		}
		request.AddArgument("type", val)
		return nil
	}
}

// Page is an argument of a request.
//
// Accepts an integer greater or equal to 0, asumed to be 0 by the server if not given.
func Page(val int) Argument {
	return func(request *requests.Request) error {
		if val < 0 {
			return errors.New("page must be an integer greater or equal to 0")
		}
		request.AddArgument("page", strconv.Itoa(val))
		return nil
	}
}

// Limit is an argument of a request. It accepts an integer greater
// or equal to 20 and lesser or equal to 100.
//
// Asumed to be 20 by the server if the argument is not given.
func Limit(val int) Argument {
	return func(request *requests.Request) error {
		if val < 20 || 100 < val {
			return errors.New("limit must be an integer greater or equal to 20 and lesser or equal to 100")
		}
		request.AddArgument("limit", strconv.Itoa(val))
		return nil
	}
}

// Start is an argument of a request.
//
// Date format: dd/mm/yyyy.
func Start(val string) Argument {
	return func(request *requests.Request) error {
		err := assertDateFormatV1(val, "Start")
		if err != nil {
			return err
		}
		request.AddArgument("start", val)
		return nil
	}
}

// End is an argument of a request.
//
// Date format: dd/mm/yyyy.
func End(val string) Argument {
	return func(request *requests.Request) error {
		err := assertDateFormatV1(val, "End")
		if err != nil {
			return err
		}
		request.AddArgument("end", val)
		return nil
	}
}

// Timeframe is an argument of a request. Its the lapse between two candles.
//
// Accepts: 1, 5, 15, 60, 240, 1440 or 10080 as strings.
func Timeframe(val string) Argument {
	return func(request *requests.Request) error {
		if !(val == "1" || val == "5" || val == "15" || val == "60" || val == "240" || val == "1440" || val == "10080") {
			return errors.New("timeframe must be one of the following numbers as string: 1, 5, 15, 60, 240, 1440 or 10080")
		}
		request.AddArgument("timeframe", val)
		return nil
	}
}

// Price is an argument of a request.
func Price(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("price", val)
		return nil
	}
}

// Currency is an argument of a request. Its a currency as "EUR" or "XLM".
func Currency(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("currency", val)
		return nil
	}
}

// Id is an argument of a request.
func Id(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("id", val)
		return nil
	}
}

// Date is an argument of a request.
//
// Needed in deposit requests for México.
//
// Date format: dd/mm/yyyy.
func Date(val string) Argument {
	return func(request *requests.Request) error {
		err := assertDateFormatV1(val, "Date")
		if err != nil {
			return err
		}
		request.AddArgument("date", val)
		return nil
	}
}

// TrackingCode is an argument of a request.
//
// Its needed in deposit requests for México.
func TrackingCode(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("tracking_code", val)
		return nil
	}
}

// Voucher is an argument of a request.
//
// Its needed in deposit requests for México, Brasil and European Union.
func Voucher(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("voucher", val)
		return nil
	}
}

// BancAccount is an argument of a request.
func BankAccount(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("bank_account", val)
		return nil
	}
}

// Address is an argument of a request.
func Address(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("address", val)
		return nil
	}
}

// Memo is an argument of a request.
func Memo(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("memo", val)
		return nil
	}
}

// CallbackUrl is an argument of a request.
//
// Max 256 caracteres.
func CallbackUrl(val string) Argument {
	return func(request *requests.Request) error {
		if len(val) > 256 {
			return errors.New("callback url too long, max 256 caracteres")
		}
		request.AddArgument("callback_url", val)
		return nil
	}
}

// ErrorUrl is an argument of a request.
//
// Max 256 caracteres.
func ErrorUrl(val string) Argument {
	return func(request *requests.Request) error {
		if len(val) > 256 {
			return errors.New("Error url too long, max 256 caracteres")
		}
		request.AddArgument("error_url", val)
		return nil
	}
}

// ErrorUrl is an argument of a request.
//
// Max 64 caracteres.
func ExternalId(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("external_id", val)
		if len(val) > 64 {
			return errors.New("callback url too long, max 64 caracteres")
		}
		return nil
	}
}

// PaymentReceiver is an argument of a request.
func PaymentReceiver(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("payment_receiver", val)
		return nil
	}
}

// SuccessUrl is an argument of a request.
//
// Max 256 caracteres
func SuccessUrl(val string) Argument {
	return func(request *requests.Request) error {
		if len(val) > 256 {
			return errors.New("callback url too long, max 256 caracteres")
		}
		request.AddArgument("success_url", val)
		return nil
	}
}

// ToReceive is an argument of a request.
func ToReceive(val float64) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("to_receive", strconv.FormatFloat(val, 'f', 2, 64))
		return nil
	}
}

// ToReceiveCurrency is an argument of a request.
func ToReceiveCurrency(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("to_receive_currency", val)
		return nil
	}
}

// Language is an argument of a request.
//
// Supported languages are "es", "en" and "pt".
func Language(val string) Argument {
	return func(request *requests.Request) error {
		if !(val == "es" || val == "en" || val == "pt") {
			return errors.New("language not supported. Supported languages are \"es\", \"en\" and \"pt\"")
		}
		request.AddArgument("language", val)
		return nil
	}
}

// Token is an argument of a request.
func Token(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("token", val)
		return nil
	}
}

// Wallet is an argument of a request.
//
// examples: "ETH", "XLM" or "BTC".
func Wallet(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("wallet", val)
		return nil
	}
}

// StartDate is an argument of a request.
//
// Date format: dd/mm/yyyy
func StartDate(val string) Argument {
	return func(request *requests.Request) error {
		err := assertDateFormatV1(val, "StartDate")
		if err != nil {
			return err
		}
		request.AddArgument("start_date", val)
		return nil
	}
}

// EndDate is an argument of a request.
//
// Date format: dd/mm/yyyy
func EndDate(val string) Argument {
	return func(request *requests.Request) error {
		err := assertDateFormatV1(val, "EndDate")
		if err != nil {
			return err
		}
		request.AddArgument("end_date", val)
		return nil
	}
}

// RefundMail is an argument of a request.
func RefundEmail(val string) Argument {
	return func(request *requests.Request) error {
		request.AddArgument("refund_email", val)
		return nil
	}
}
