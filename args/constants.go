package args

// SideType is an order side or a trade side of an order
type SideType string

const (
	SideSell SideType = "sell"
	SideBuy  SideType = "buy"
)

// OrderType is a type of order
type OrderType string

// OrderTypes
const (
	OrderLimit            OrderType = "limit"
	OrderMarket           OrderType = "market"
	OrderStopLimit        OrderType = "stopLimit"
	OrderStopMarket       OrderType = "stopMarket"
	OrderTakeProfitLimit  OrderType = "takeProfitLimit"
	OrderTakeProfitMarket OrderType = "takeProfitMarket"
)

// TimeInForceType is the time in force of an order
type TimeInForceType string

const (
	TimeInForceGTC TimeInForceType = "GTC" // Good Till Cancel
	TimeInForceIOC TimeInForceType = "IOC" // Immediate or Cancel
	TimeInForceFOK TimeInForceType = "FOK" // Fill or Kill
	TimeInForceDAY TimeInForceType = "Day" // valid during Day
	TimeInForceGTD TimeInForceType = "GTD" // Good Till Date
)

// SortType is the sorting direction of a query
type SortType string

const (
	SortASC  SortType = "ASC"
	SortDESC SortType = "DESC"
)

// SortBy is the field used for sorting lists
type SortByType string

const (
	SortByTimestamp SortByType = "timestamp"
	SortByID        SortByType = "id"
	SortByCreatedAt SortByType = "created_at"
)

// PeriodType is the period of a candle
type PeriodType string

const (
	Period1Minute   PeriodType = "M1"
	Period3Minutes  PeriodType = "M3"
	Period5Minutes  PeriodType = "M5"
	Period15Minutes PeriodType = "M15"
	Period30Minutes PeriodType = "M30"
	Period1Hour     PeriodType = "H1"
	Period4Hours    PeriodType = "H4"
	Period1Day      PeriodType = "D1"
	Period7Days     PeriodType = "D7"
	Period1Month    PeriodType = "1M"
)

// IdentifyByType is the type of identifier used for transafers between users
type IdentifyByType string

const (
	IdentifyByEmail    IdentifyByType = "email"
	IdentifyByUsername IdentifyByType = "username"
)

// AccountType is the type of account used for internal user tranfers
type AccountType string

const (
	AccountWallet AccountType = "wallet"
	AccountSpot   AccountType = "spot"
)

type UseOffchainType string

const (
	UseOffchainNever     UseOffchainType = "never"
	UseOffchainOptionaly UseOffchainType = "optionaly"
	UseOffChainRequired  UseOffchainType = "required"
)

type TransactionTypeType string

const (
	TransactionTypeDeposit  TransactionTypeType = "DEPOSIT"
	TransactionTypeWithdraw TransactionTypeType = "WITHDRAW"
	TransactionTypeTransfer TransactionTypeType = "TRANSFER"
	TransactionTypeSwap     TransactionTypeType = "SWAP"
)

type TransactionSubTypeType string

const (
	TransactionSubTypeUnclassified    TransactionSubTypeType = "UNCLASSIFIED"
	TransactionSubTypeBlockchain      TransactionSubTypeType = "BLOCKCHAIN"
	TransactionSubTypeAffiliate       TransactionSubTypeType = "AFFILIATE"
	TransactionSubTypeOffchain        TransactionSubTypeType = "OFFCHAIN"
	TransactionSubTypeFiat            TransactionSubTypeType = "FIAT"
	TransactionSubTypeSubAccount      TransactionSubTypeType = "SUB_ACCOUNT"
	TransactionSubTypeWalletToSpot    TransactionSubTypeType = "WALLET_TO_SPOT"
	TransactionSubTypeSpotToWallet    TransactionSubTypeType = "SPOT_TO_WALLET"
	TransactionSubTypeChainSwitchFrom TransactionSubTypeType = "CHAIN_SWITCH_FROM"
	TransactionSubTypeChainSwitchTo   TransactionSubTypeType = "CHAIN_SWITCH_TO"
)

type TransactionStatusType string

const (
	TransactionStatusCreated    TransactionStatusType = "CREATED"
	TransactionStatusPending    TransactionStatusType = "PENDING"
	TransactionStatusFailed     TransactionStatusType = "FAILED"
	TransactionStatusSuccess    TransactionStatusType = "SUCCESS"
	TransactionStatusRolledBack TransactionStatusType = "ROLLED_BACK"
)

type TickerSpeedType string

const (
	TickerSpeed1s TickerSpeedType = "1s"
	TickerSpeed3s TickerSpeedType = "3s"
)

type OrderBookSpeedType string

const (
	OrderBookSpeed100ms  OrderBookSpeedType = "100ms"
	OrderBookSpeed500ms  OrderBookSpeedType = "500ms"
	OrderBookSpeed1000ms OrderBookSpeedType = "1000ms"
)

type WSDepthType string

const (
	WSDepth5  WSDepthType = "D5"
	WSDepth10 WSDepthType = "D10"
	WSDepth20 WSDepthType = "D20"
)

type NotificationType string

const (
	NotificationSnapshot NotificationType = "snapshot"
	NotificationUpdate   NotificationType = "update"
	NotificationData     NotificationType = "data"
)

type ContingencyType string

const (
	ContingencyAllOrNone                ContingencyType = "allOrNone"
	ContingencyAON                      ContingencyType = "allOrNone"
	ContingencyOneCancelOther           ContingencyType = "oneCancelOther"
	ContingencyOCO                      ContingencyType = "oneCancelOther"
	ContingencyOneTriggerOneCancelOther ContingencyType = "oneTriggerOneCancelOther"
	ContingencyOTOCO                    ContingencyType = "oneTriggerOneCancelOther"
)

type OrderStatusType string

const (
	OrderStatusNew             OrderStatusType = "new"
	OrderStatusSuspended       OrderStatusType = "suspended"
	OrderStatusPartiallyFilled OrderStatusType = "partiallyFilled"
	OrderStatusFilled          OrderStatusType = "filled"
	OrderStatusCanceled        OrderStatusType = "canceled"
	OrderStatusExpired         OrderStatusType = "expired"
)

type ReportType string

const (
	ReportStatus    ReportType = "status"
	ReportNew       ReportType = "new"
	ReportCanceled  ReportType = "canceled"
	ReportRejected  ReportType = "rejected"
	ReportExpired   ReportType = "expired"
	ReportSuspended ReportType = "suspended"
	ReportTrade     ReportType = "trade"
	ReportReplaced  ReportType = "replaced"
)

type SymbolStatusType string

const (
	SymbolStatusWorking   SymbolStatusType = "working"
	SymbolStatusSuspended SymbolStatusType = "suspended"
)

type MetaTransactionStatusType string

const (
	MetaTransactionStatusActive   MetaTransactionStatusType = "ACTIVE"
	MetaTransactionStatusInactive MetaTransactionStatusType = "INACTIVE"
)

type TransferTypeType string

const (
	TransferToSubAccount   TransferTypeType = "to_sub_account"
	TransferFromSubAccount TransferTypeType = "from_sub_account"
)
