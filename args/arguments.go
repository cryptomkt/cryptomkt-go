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

// SideType is an order side or a trade side of an order
type SideType string

const (
	// SideTypeSell is the sell side for an order or a trade
	SideTypeSell SideType = "sell"
	// SideTypeBuy is the buy side for an order or a trade
	SideTypeBuy SideType = "buy"
)

// OrderType is a type of order
type OrderType string

// OrderTypes
const (
	OrderTypeLimit      OrderType = "limit"
	OrderTypeMarket     OrderType = "market"
	OrderTypeStopLimit  OrderType = "stopLimit"
	OrderTypeStopMarket OrderType = "stopMarket"
)

// TimeInForceType is the time in force of an order
type TimeInForceType string

// types of time in force
const (
	TimeInForceTypeGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceTypeIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceTypeFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceTypeDAY TimeInForceType = "Day" // valid during Day
	TimeInForceTypeGTD TimeInForceType = "GTD" // Good Till Date
)

// SortType is the sorting direction of a query
type SortType string

const (
	// SortTypeASC is the ascending sorting direction of a query
	SortTypeASC SortType = "ASC"
	// SortTypeDESC is the descending sorting direction of a query
	SortTypeDESC SortType = "DESC"
)

// SortByType is the parameter for sorting
type SortByType string

const (
	// SortByTypeTimestamp is the sorting field for pagination, sorting by timestamp
	SortByTypeTimestamp SortByType = "timestamp"
	// SortByTypeID is the sorting field for pagination, sorting by id
	SortByTypeID SortByType = "id"
)

const (
	// TransferByEmail signals the identifier type to transfer by.
	TransferByEmail = "email"
	// TransferByUsername signals the identifier type to transfer by.
	TransferByUsername = "username"
)

// PeriodType is the period of a candle
type PeriodType string

// candle periods
const (
	PeriodType1Minutes  PeriodType = "M1"
	PeriodType3Minutes  PeriodType = "M3"
	PeriodType5Minutes  PeriodType = "M5"
	PeriodType15Minutes PeriodType = "M15"
	PeriodType30Minutes PeriodType = "M30"
	PeriodType1Hours    PeriodType = "H1"
	PeriodType4Hours    PeriodType = "H4"
	PeriodType1Day      PeriodType = "D1"
	PeriodType7Days     PeriodType = "D7"
	PeriodType1Month    PeriodType = "1M"
)

// MarginType is the type of margin of a trade
type MarginType string

// types of margin
const (
	MarginTypeInclude MarginType = "include"
	MarginTypeOnly    MarginType = "only"
	MarginTypeIgnore  MarginType = "ignore"
)

// IdentifyByType for transfers
type IdentifyByType string

// identify by types
const (
	IdentifyByTypeEmail    IdentifyByType = "email"
	IdentifyByTypeUsername IdentifyByType = "username"
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
func Sort(val SortType) Argument {
	return func(params map[string]interface{}) {
		params["sort"] = val
	}
}

// SortBy returns a "by" Argument
func SortBy(val SortByType) Argument {
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
func Period(val PeriodType) Argument {
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

// Wait returns a "wait" Argument
func Wait(val int) Argument {
	return func(params map[string]interface{}) {
		params["wait"] = val
	}
}

// Side returns a "side" Argument
func Side(val SideType) Argument {
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
func TimeInForce(val TimeInForceType) Argument {
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
func Margin(val MarginType) Argument {
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

// Type returns a "type" Argument
func Type(val OrderType) Argument {
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
func IdentifyBy(val IdentifyByType) Argument {
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
