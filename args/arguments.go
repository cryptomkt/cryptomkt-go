package args

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cryptomarket/cryptomarket-go/internal"
)

// Argument are functions that serves as arguments for the diferent
// requests to the server, either rest request, or websocket requests.
// Argument works by extending a given map with a particular key.
// For example Symbol("EOSETH") returns an Argument that extends
// a map assigning "EOSETH" to the key "symbol".
// usually, a function returning an Argument is what's used as a
// request parameter. just like the Symbol(val string) function
// in the example above.
type Argument func(map[string]interface{})

func fromSnakeCaseToCamelCase(s string) string {
	snakeParts := strings.Split(s, "_")
	camelParts := make([]string, 0)
	for _, snakePart := range snakeParts {
		camelParts = append(camelParts, strings.Title(snakePart))
	}
	return strings.Join(camelParts, "")
}

// BuildParams makes a map with the Arguments functions,
// and check for the presence of "requireds" keys in the map,
// raising an error if some required keys are not present.
func BuildParams(
	arguments []Argument,
	requireds ...string,
) (map[string]interface{}, error) {
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
		missingAsCamelCase := make([]string, 0)
		for _, miss := range missing {
			missingAsCamelCase = append(
				missingAsCamelCase,
				fromSnakeCaseToCamelCase(miss),
			)
		}
		return nil, fmt.Errorf(
			"CryptomarketSDKError: missing arguments: %v", missingAsCamelCase,
		)
	}
	return params, nil
}

type stringable interface {
	TransactionTypeType |
		TransactionSubTypeType |
		TransactionStatusType
}

func toString[str stringable](list []str) string {
	asStrs := make([]string, len(list))
	for i, val := range list {
		asStrs[i] = string(val)
	}
	return strings.Join(asStrs, ",")
}

func BuildQuery(params map[string]interface{}) string {
	query := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case int:
			query.Add(key, strconv.Itoa(v))
		case []string:
			strs := strings.Join(v, ",")
			query.Add(key, strs)
		case []TransactionTypeType:
			query.Add(key, toString(v))
		case []TransactionSubTypeType:
			query.Add(key, toString(v))
		case []TransactionStatusType:
			query.Add(key, toString(v))
		default:
			query.Add(key, fmt.Sprint(v))
		}
	}
	return query.Encode()
}

func Currencies(val []string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameCurrencies] = val
	}
}

func Currency(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameCurrency] = val
	}
}

func Symbols(val []string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSymbols] = val
	}
}

func Symbol(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSymbol] = val
	}
}

func Sort(val SortType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSort] = val
	}
}

func SortBy(val SortByType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSortBy] = val
	}
}

func From(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameFrom] = val
	}
}

func To(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameTo] = val
	}
}

func Till(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameTill] = val
	}
}

func Limit(val int) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameLimit] = val
	}
}

func Offset(val int) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameOffset] = val
	}
}

func Volume(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameVolume] = val
	}
}

func Period(val PeriodType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNamePeriod] = val
	}
}

func ClientOrderID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameClientOrderID] = val
	}
}

func Side(val SideType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSide] = val
	}
}

func Quantity(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameQuantity] = val
	}
}

func Price(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNamePrice] = val
	}
}

func StopPrice(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameStopPrice] = val
	}
}

func TimeInForce(val TimeInForceType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameTimeInForce] = val
	}
}

func ExpireTime(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameExpireTime] = val
	}
}

func StrictValidate(val bool) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameStrictValidate] = val
	}
}

func PostOnly(val bool) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNamePostOnly] = val
	}
}

func OrderID(val int) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameOrderID] = val
	}
}

func Amount(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameAmount] = val
	}
}

func Address(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameAddress] = val
	}
}

func PaymentID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNamePaymentID] = val
	}
}

func IncludeFee(val bool) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameIncludeFee] = val
	}
}

func AutoCommit(val bool) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameAutoCommit] = val
	}
}

func PublicComment(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNamePublicComment] = val
	}
}

func FromCurrency(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameFromCurrency] = val
	}
}

func ToCurrency(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameToCurrency] = val
	}
}

func TransferType(val TransferTypeType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameTransferType] = val
	}
}

func Identifier(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameIdentifier] = val
	}
}

func IdentifyBy(val IdentifyByType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameIdentifyBy] = val
	}
}

func ShowSenders(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameShowSenders] = val
	}
}

func ID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameID] = val
	}
}

func Source(val AccountType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSource] = val
	}
}

func Destination(val AccountType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameDestination] = val
	}
}

func Since(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSince] = val
	}
}

func Untill(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameUntil] = val
	}
}

func Depth(val int) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameDepth] = val
	}
}

func WSDepth(val WSDepthType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameDepth] = val
	}
}

func TakeRate(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameTakeRate] = val
	}
}

func MakeRate(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameMakeRate] = val
	}
}

func NewClientOrderID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameNewClientOrderID] = val
	}
}

func UseOffchain(val UseOffchainType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameNewClientOrderID] = val
	}
}

func RequestClientOrderID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameRequestClientOrderID] = val
	}
}

func OrderBookSpeed(val OrderBookSpeedType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSpeed] = val
	}
}

func TickerSpeed(val TickerSpeedType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSpeed] = val
	}
}

func OrderListID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameOrderListID] = val
	}
}

func Contingency(val ContingencyType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameContingencyType] = val
	}
}
func IDTill(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameIDTill] = val
	}
}

func IDFrom(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameIDFrom] = val
	}
}

type OrderRequest struct {
	ClientOrderID  string          `json:"client_order_id"`
	Symbol         string          `json:"symbol"`
	Side           SideType        `json:"side"`
	Type           OrderType       `json:"type"`
	TimeInForce    TimeInForceType `json:"time_in_force"`
	Quantity       string          `json:"quantity"`
	Price          string          `json:"price"`
	StopPrice      string          `json:"stop_price"`
	ExpireTime     string          `json:"expire_time"`
	StrictValidate bool            `json:"strict_validate"`
	PostOnly       bool            `json:"post_only"`
	MakeRate       string          `json:"make_rate"`
	TakeRate       string          `json:"take_rate"`
}

func Orders(val []OrderRequest) Argument {
	return func(params map[string]interface{}) {
		data, _ := json.Marshal(val)
		params[internal.ArgNameOrders] = string(data)
	}
}

func TransactionType(val TransactionTypeType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameType] = val
	}
}

func TransactionTypes(list ...TransactionTypeType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameTypes] = list
	}
}

func TransactionSubTypes(list ...TransactionSubTypeType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSubtypes] = list
	}
}

func TransactionStatuses(list ...TransactionStatusType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameStatuses] = list
	}
}

func SubAccountID(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSubAccountID] = val
	}
}

func SubAccountIDs(list ...string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSubAccountIDs] = list
	}
}

func DepositAddressGenerationEnabled(val bool) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameDepositAddressGenerationEnabled] = val
	}
}

func WithdrawEnabled(val bool) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameWithdrawEnabled] = val
	}
}

func Description(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameDescription] = val
	}
}

func CreatedAt(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameCreatedAt] = val
	}
}

func UpdatedAt(val string) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameUpdatedAt] = val
	}
}

type SubscriptionType string

func SubscriptionTypeTrades() SubscriptionType {
	return SubscriptionType(internal.ChannelTrades)
}

func SubscriptionsTypeCandles(period PeriodType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelCandles, period))
}

func SubscriptionTypeMiniTicker(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelMiniTicker, speed))
}

func SubscriptionTypeMiniTickerInBatch(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelMiniTickerInBatch, speed))
}

func SubscriptionTypeTicker(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelTicker, speed))
}

func SubscriptionTypeTickerInBatch(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelTickerInBatch, speed))
}

func SubscriptionTypeFullOrderbook(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(internal.ChannelOrderBookFull)
}

func SubscriptionTypePartialOrderbook(depth WSDepthType, speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelOrderbookPartial, depth, speed))
}

func SubscriptionTypePartialOrderbookInBatch(depth WSDepthType, speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelOrderbookPartialInBatch, depth, speed))
}

func SubscriptionTypeOrderbookTop(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelOrderbookTop, speed))
}

func SubscriptionTypeOrderbookTopInBatch(speed OrderBookSpeedType) SubscriptionType {
	return SubscriptionType(fmt.Sprintf(internal.ChannelOrderbookTopInBatch, speed))
}

func Subscription(val SubscriptionType) Argument {
	return func(params map[string]interface{}) {
		params[internal.ArgNameSubscription] = val
	}
}
