package args

import "fmt"

// Argument are functions that serves as arguments for the diferent
// requests to the server, either rest request, or websocket requests.
// Argument works by extending a given map with a particular key.
// For example Symbol("EOSETH") returns an Argument that extends
// a map assigning "EOSETH" to the key "symbol".
// usually, a function returning an Argument is what's used as a
// request parameter. just like the Symbol(val string) function
// in the example above.
type Argument func(map[string]interface{})

const (
	// SideSell is the "sell" side, for an order or for a trade
	SideSell = "sell"
	// SideBuy is the "buy" side, for an order or for a trade
	SideBuy = "buy"

	// TypeLimit is the limit order type.
	TypeLimit = "limit"
	// TypeMarket is the market order type.
	TypeMarket = "market"
	// TypeStopLimit is the stop limit order type.
	TypeStopLimit = "stopLimit"
	// TypeStopMarket is the stop market order type.
	TypeStopMarket = "stopMarket"

	// TimeInForceGTC is Good Till Canceled
	TimeInForceGTC = "GTC"
	// TimeInForceIOC Inmediate or Cancel
	TimeInForceIOC = "IOC"
	// TimeInForceFOK is Fill Or Kill
	TimeInForceFOK = "FOK"
	// TimeInForceDAY is Day
	TimeInForceDAY = "DAY"
	// TimeInForceGTD is Good Till Date
	TimeInForceGTD = "GTD"

	// SortAsc is sorting for pagination, ascending order
	SortAsc = "ASC"
	// SortDesc is sorting for pagination, descending order
	SortDesc = "DESC"

	// SortByTimestamp is the sorting field for pagination, sorting by timestamp
	SortByTimestamp = "timestamp"
	// SortByID is the sorting field for pagination, sorting by id
	SortByID = "id"

	// TransferByEmail signals the identifier type to transfer by.
	TransferByEmail = "email"
	// TransferByUsername signals the identifier type to transfer by.
	TransferByUsername = "username"
)

// BuildParams makes a map with the Arguments functions,
// and check for the presence of "requireds" keys in the map,
// raising an error if some required keys are not present.
func BuildParams(arguments []Argument, requireds ...string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	for _, argFunc := range arguments {
		argFunc(params)
	}
	missing := []string{}
	for _, required := range requireds {
		if _, ok := params[required]; !ok {
			missing = append(missing, required)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("CryptomarketSDKError: missing arguments: %v", missing)
	}
	return params, nil
}

// Currencies returns a "currencies" Argument
func Currencies(val []string) Argument {
	return func(params map[string]interface{}) {
		params["currencies"] = val
	}
}

// Currency returns a "currency" Argument
func Currency(val string) Argument {
	return func(params map[string]interface{}) {
		params["currency"] = val
	}
}

// Symbols returns a "symbols" Argument
func Symbols(val []string) Argument {
	return func(params map[string]interface{}) {
		params["symbols"] = val
	}
}

// Symbol returns a "symbol" Argument
func Symbol(val string) Argument {
	return func(params map[string]interface{}) {
		params["symbol"] = val
	}
}

// Sort returns a "sort" Argument
func Sort(val string) Argument {
	return func(params map[string]interface{}) {
		params["side"] = val
	}
}

// By returns a "by" Argument
func By(val string) Argument {
	return func(params map[string]interface{}) {
		params["by"] = val
	}
}

// From returns a "from" Argument
func From(val string) Argument {
	return func(params map[string]interface{}) {
		params["from"] = val
	}
}

// Till returns a "till" Argument
func Till(val string) Argument {
	return func(params map[string]interface{}) {
		params["till"] = val
	}
}

// Limit returns a "limit" Argument
func Limit(val int) Argument {
	return func(params map[string]interface{}) {
		params["limit"] = val
	}
}

// Offset returns a "offset" Argument
func Offset(val int) Argument {
	return func(params map[string]interface{}) {
		params["offset"] = val
	}
}

// Volume returns a "volume" Argument
func Volume(val string) Argument {
	return func(params map[string]interface{}) {
		params["volume"] = val
	}
}

// Period returns a "period" Argument
func Period(val string) Argument {
	return func(params map[string]interface{}) {
		params["period"] = val
	}
}

// ClientOrderID returns a "clientOrderId" Argument
func ClientOrderID(val string) Argument {
	return func(params map[string]interface{}) {
		params["clientOrderId"] = val
	}
}

// Side returns a "side" Argument
func Side(val string) Argument {
	return func(params map[string]interface{}) {
		params["side"] = val
	}
}

// Quantity returns a "quantity" Argument
func Quantity(val string) Argument {
	return func(params map[string]interface{}) {
		params["quantity"] = val
	}
}

// Price returns a "price" Argument
func Price(val string) Argument {
	return func(params map[string]interface{}) {
		params["price"] = val
	}
}

// StopPrice returns a "stopPrice" Argument
func StopPrice(val string) Argument {
	return func(params map[string]interface{}) {
		params["stopPrice"] = val
	}
}

// TimeInForce returns a "timeInForce" Argument
func TimeInForce(val string) Argument {
	return func(params map[string]interface{}) {
		params["timeInForce"] = val
	}
}

// ExpireTime returns a "expireTime" Argument
func ExpireTime(val string) Argument {
	return func(params map[string]interface{}) {
		params["expireTime"] = val
	}
}

// StrictValidate returns a "strictValidate" Argument
func StrictValidate(val bool) Argument {
	return func(params map[string]interface{}) {
		params["strictValidate"] = val
	}
}

// PostOnly returns a "postOnly" Argument
func PostOnly(val bool) Argument {
	return func(params map[string]interface{}) {
		params["postOnly"] = val
	}
}

// Margin returns a "margin" Argument
func Margin(val string) Argument {
	return func(params map[string]interface{}) {
		params["margin"] = val
	}
}

// OrderID returns a "orderId" Argument
func OrderID(val int64) Argument {
	return func(params map[string]interface{}) {
		params["orderId"] = val
	}
}

// Amount returns a "amount" Argument
func Amount(val string) Argument {
	return func(params map[string]interface{}) {
		params["amount"] = val
	}
}

// Address returns a "address" Argument
func Address(val string) Argument {
	return func(params map[string]interface{}) {
		params["address"] = val
	}
}

// PaymentID returns a "paymentId" Argument
func PaymentID(val string) Argument {
	return func(params map[string]interface{}) {
		params["paymentId"] = val
	}
}

// IncludeFee returns a "includeFee" Argument
func IncludeFee(val bool) Argument {
	return func(params map[string]interface{}) {
		params["includeFee"] = val
	}
}

// AutoCommit returns a "autoCommit" Argument
func AutoCommit(val bool) Argument {
	return func(params map[string]interface{}) {
		params["autoCommit"] = val
	}
}

// PublicComment returns a "publicComment" Argument
func PublicComment(val string) Argument {
	return func(params map[string]interface{}) {
		params["publicComment"] = val
	}
}

// FromCurrency returns a "fromCurrency" Argument
func FromCurrency(val string) Argument {
	return func(params map[string]interface{}) {
		params["fromCurrency"] = val
	}
}

// ToCurrency returns a "toCurrency" Argument
func ToCurrency(val string) Argument {
	return func(params map[string]interface{}) {
		params["toCurrency"] = val
	}
}

// TransferType returns a "type" Argument
func TransferType(val string) Argument {
	return func(params map[string]interface{}) {
		params["type"] = val
	}
}

// OrderType returns a "type" Argument
func OrderType(val string) Argument {
	return func(params map[string]interface{}) {
		params["type"] = val
	}
}

// Type returns a "type" Argument
func Type(val string) Argument {
	return func(params map[string]interface{}) {
		params["type"] = val
	}
}

// Identifier returns a "identifier" Argument
func Identifier(val string) Argument {
	return func(params map[string]interface{}) {
		params["identifier"] = val
	}
}

// IdentifyBy returns a "by" Argument
func IdentifyBy(val string) Argument {
	return func(params map[string]interface{}) {
		params["by"] = val
	}
}

// ShowSenders returns a "showSenders" Argument
func ShowSenders(val string) Argument {
	return func(params map[string]interface{}) {
		params["showSenders"] = val
	}
}

// RequestClientID returns a "requestClientId" Argument
func RequestClientID(val string) Argument {
	return func(params map[string]interface{}) {
		params["requestClientId"] = val
	}
}

// ID returns a "id" Argument
func ID(val string) Argument {
	return func(params map[string]interface{}) {
		params["id"] = val
	}
}
