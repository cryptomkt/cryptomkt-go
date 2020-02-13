package client

import (
	"errors"
	"fmt"
	"strconv"
)

type Argument func(*Request) error

func assertDayFormat(val, date string) error {
	day, err := strconv.Atoi(val[0:2])
	if err != nil {
		return fmt.Errorf("%s format error: must be dd/mm/yyyy", date)
	}
	if day < 1 || 31 < day {
		return fmt.Errorf("%s format error: invalid day value", date)
	}
	month, err := strconv.Atoi(val[3:5])
	if err != nil {
		return fmt.Errorf("%s format error: must be dd/mm/yyyy", date)
	}
	if month < 1 || 12 < month {
		return fmt.Errorf("%s format error: invalid month value", date)
	}
	_, err = strconv.Atoi(val[6:10])
	if err != nil {
		return fmt.Errorf("%s format error: must be dd/mm/yyyy", date)
	}
	return nil
}

// Amount is an argument of a request
// accepts positive numbers
// number format: without thousand separator, and '.' as decimal point
func Amount(val string) Argument {
	return func(request *Request) error {
		request.addArgument("amount", val)
		return nil
	}
}

// Market is an argument of a request,
// accepts a par of currencies. e.g. "ETHCLP" or "BTCARS"
func Market(val string) Argument {
	return func(request *Request) error {
		request.addArgument("market", val)
		return nil
	}
}

// Type is an argument of a request,
// accepts either "buy" or "sell"
func Type(val string) Argument {
	return func(request *Request) error {
		if !(val == "buy" || val == "sell") {
			return errors.New("type must be either \"buy\" or \"sell\"")
		}
		request.addArgument("type", val)
		return nil
	}
}

// Page is an argument of a request,
// accepts an integer greater or equal to 0
// default value is 0
func Page(val int) Argument {
	return func(request *Request) error {
		if val < 0 {
			return errors.New("page must be an integer greater or equal to 0")
		}
		request.addArgument("page", string(val))
		return nil
	}
}

// Limit is an argument of a request,
// accepts an integer greater or equal to 0 and lesser or equal to 200.
// if not given, its default value is 20.
func Limit(val int) Argument {
	return func(request *Request) error {
		if val < 20 || 100 < val {
			return errors.New("limit must be an integer greater or equal to 20 and lesser or equal to 100")
		}
		request.addArgument("limit", string(val))
		return nil
	}
}

// Timeframe is an argument of a request,
// its the lapse between two candles,
// accepts 1, 5, 15, 60, 240, 1440 or 10080 as integers
func Timeframe(val string) Argument {
	return func(request *Request) error {
		if !(val == "1" || val == "5" || val == "15" || val == "60" || val == "240" || val == "1440" || val == "10080") {
			return errors.New("timeframe must be one of the following numbers: 1, 5, 15, 60, 240, 1440 or 10080, as string")
		}
		request.addArgument("limit", string(val))
		return nil
	}
}

// Price is an argument of a request
// must be a positive number
func Price(val string) Argument {
	return func(request *Request) error {
		request.addArgument("price", val)
		return nil
	}
}

// Currency is an argument of a request,
// its a currency as "EUR" or "XLM"
func Currency(val string) Argument {
	return func(request *Request) error {
		request.addArgument("currency", val)
		return nil
	}
}

// Id is an argument of a request.
func Id(val string) Argument {
	return func(request *Request) error {
		request.addArgument("id", val)
		return nil
	}
}

// Date is an argument of a request.
// needed in deposit for México
// date format: dd/mm/yyyy
func Date(val string) Argument {
	return func(request *Request) error {
		err := assertDayFormat(val, "date")
		if err != nil {
			return err
		}
		request.addArgument("date", val)
		return nil
	}
}

// TrackingCode is an argument of a request
// its needed in deposits for México
func TrackingCode(val string) Argument {
	return func(request *Request) error {
		request.addArgument("tracking_code", val)
		return nil
	}
}

// Voucher is an argument of a request
// its needed in deposits for México, Brasil and European Union.
func Voucher(val string) Argument {
	return func(request *Request) error {
		request.addArgument("voucher", val)
		return nil
	}
}

// BancAccount is an argument of a request
func BankAccount(val string) Argument {
	return func(request *Request) error {
		request.addArgument("bank_account", val)
		return nil
	}
}

// Address is an argument of a request
func Address(val string) Argument {
	return func(request *Request) error {
		request.addArgument("address", val)
		return nil
	}
}

// Memo is an argument of a request
func Memo(val string) Argument {
	return func(request *Request) error {
		request.addArgument("memo", val)
		return nil
	}
}

// CallbackUrl is an argument of a request
// max 256 caracteres.
func CallbackUrl(val string) Argument {
	return func(request *Request) error {
		if len(val) > 256 {
			return errors.New("callback url too long, max 256 caracteres")
		}
		request.addArgument("callback_url", val)
		return nil
	}
}

// ErrorUrl is an argument of a request
// max 256 caracteres.
func ErrorUrl(val string) Argument {
	return func(request *Request) error {
		if len(val) > 256 {
			return errors.New("Error url too long, max 256 caracteres")
		}
		request.addArgument("error_url", val)
		return nil
	}
}

// ErrorUrl is an argument of a request
// max 64 caracteres.
func ExternalId(val string) Argument {
	return func(request *Request) error {
		request.addArgument("external_id", val)
		if len(val) > 64 {
			return errors.New("callback url too long, max 64 caracteres")
		}
		return nil
	}
}

// PaymentReceiver is an argument of a request
func PaymentReceiver(val string) Argument {
	return func(request *Request) error {
		request.addArgument("payment_receiver", val)
		return nil
	}
}

// SuccessUrl is an argument of a request
// max 256 caracteres
func SuccessUrl(val string) Argument {
	return func(request *Request) error {
		if len(val) > 256 {
			return errors.New("callback url too long, max 256 caracteres")
		}
		request.addArgument("success_url", val)
		return nil
	}
}

// ToReceive is an argument of a request.
func ToReceive(val float64) Argument {
	return func(request *Request) error {
		request.addArgument("to_receive", strconv.FormatFloat(val, 'f', 2, 64))
		return nil
	}
}

// ToReceiveCurrency is an argument of a request.
func ToReceiveCurrency(val string) Argument {
	return func(request *Request) error {
		request.addArgument("to_receive_currency", val)
		return nil
	}
}

// Language is an argument of a request.
// supported languages are "es", "en" and "pt"
func Language(val string) Argument {
	return func(request *Request) error {
		if !(val == "es" || val == "en" || val == "pt") {
			return errors.New("language not supported. Supported languages are \"es\", \"en\" and \"pt\"")
		}
		request.addArgument("language", val)
		return nil
	}
}

// Token is an argument of a request.
func Token(val string) Argument {
	return func(request *Request) error {
		request.addArgument("token", val)
		return nil
	}
}

// Wallet is an argument of a request.
// enabled wallets are "ETH", "XLM" and "BTC"
func Wallet(val string) Argument {
	return func(request *Request) error {
		if !(val == "ETH" || val == "XLM" || val == "BTC") {
			return errors.New("wallet not supported. Wallets enabled are \"ETH\", \"XLM\" and \"BTC\"")
		}
		request.addArgument("wallet", val)
		return nil
	}
}

// StartDate is an argument of a request.
// date format: dd/mm/yyyy
func StartDate(val string) Argument {
	return func(request *Request) error {
		err := assertDayFormat(val, "start date")
		if err != nil {
			return err
		}
		request.addArgument("start_date", val)
		return nil
	}
}

// EndDate is an argument of a request.
// date format: dd/mm/yyyy
func EndDate(val string) Argument {
	return func(request *Request) error {
		err := assertDayFormat(val, "end date")
		if err != nil {
			return err
		}
		request.addArgument("end_date", val)
		return nil
	}
}

// RefundMail is an argument of a request.
func RefundEmail(val string) Argument {
	return func(request *Request) error {
		request.addArgument("refund_email", val)
		return nil
	}
}
