# cryptomkt-go

<!--- badges -->
[![GitHub version](https://badge.fury.io/gh/cryptomkt%2Fcryptomkt-go.svg)](https://badge.fury.io/gh/cryptomkt%2Fcryptomkt-go)

cryptomkt-go is the SDK for cryptomkt in the GO programing language

## Installation
To install the sdk, run the `go get` command

`go get github.com/cryptomkt/cryptomkt-go`

## Documentation

For further information about the sdk see the godoc

The base api for this sdk can be found in [here](https://developers.cryptomkt.com/)

## API Key

To make use of the sdk, you need to [enable an API key](https://www.cryptomkt.com/platform/account#api_tab) in the cryptomkt account you'll be using.

If you don't have an account yet, sign up [here](https://www.cryptomkt.com/account/register)

Once you enable an API key, you'll get two keys that are needed to make a client to connect with cryptomkt. All calls are done with this client.

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
)
client := conn.NewClient(apiKey, apiSecret)
```

## Configuring Calls
Arguments are needed for most of the calls you can make. For each new call, you'll pass a different set of configuration arguments. All arguments are in the `args` package

Each call specifies which arguments are required and which ones are optional. you can find this information in the documentation or in the [api page](https://developers.cryptomkt.com/) of Cryptomkt. Also, an error mentioning the unmeeted required arguments is given when an incomplete call is done.

As an example, to create a buy order in the ETHCLP market, we can use CreateOrder, which requires the Amount, Market, Price and Type arguments and have no optional arguments.

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

order, err := client.CreateOrder(
    args.Amount(0.3),
    args.Market("ETHCLP"),
    args.Price(1000),
    args.Type("buy"),
)
// if we forgot a required argument, this error will tell us
if err != nil {
    fmt.Errorf("Error making an order: %s", err)
}
```

## Calls to the API

Calls have multiple return formats.
All calls return at least one informative error if something goes wrong, as an unmeeted argument, an invalid apiKey or a "not_enough_balance" as a replay from the server if you try to buy more than your money can take.

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

//account is a pointer to a struct that matches the account information
account, err := client.GetAccount()
if err != nil {
    fmt.Errorf("Error getting the account: %s", err)
}
```

If we want to go over a long range of trade data of a market, we can call `client.GetTrades` to get a list of `Trades`, this list can be one page of many. When we read the data of one single page, to get the rest of the pages, we can call over and over `GetNext()` over the struct, until an `Next page does not exist` error is raised. Replace `GetObject` with the appropriate method. The structs that support this functionality so far are Trades, Book, Prices and Orders. Here is in code:

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/args"
    "github.com/cryptomkt/cryptomkt-go/conn"
)

var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey,apiSecret)

response, err := client.GetObject(args.Argument1(value1), args.Argument2(value2), ...)

nextPage, err := reponse.GetNext()
previousPage, err := response.GetPrevious()

// You can call these methods from its response if the page exists
nextPage2, err := nextPage.GetNext()
previousPage2, err := previousPage.GetPrevious()

```

Also you can close and refresh and close orders directly from their structures

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)

var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey,apiSecret)

order, err := client.OrderStatus(args.Id("M103966"))

// to update the order
order, err = order.Refresh()
if err != nil {
    fmt.Errorf("error while updating the order: %v", err)
}

// and to close it
order, err = order.Close()
if err != nil {
    fmt.Errorf("error while closing the order: %v", err)
}
```

## API Calls Examples


### Public endpoints

Responses from client methods are pointers to its structures.

**Listing available markets**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
)
var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey, apiSecret)

// marketList is a list of enabled markets
marketList, err := client.GetMarkets()
if err != nil {
    fmt.Errorf("Error getting the market list: %s", err)
}
```

**Getting tickers of active markets**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)

var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey, apiSecret)

// Here you get the ticker list for the ethereum chilean pesos market. It is 
// your choice to give the Market argument
ticker, err := client.GetTicker(args.Market("ETHCLP"))

if err != nil {
    fmt.Errorf("Error getting the ticker, %s", err)
}else{
    // here you have the data
    fmt.Println(ticker.Data)
}

// or, if you prefer, you can get all markets tickers
allTickers, err := client.GetTicker()
if err != nil{
    fmt.Errorf("Error getting all tickers, %s", err)
}else{
    fmt.Println(allTickers.Data)
}
```

**Getting active orders book**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)

var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey,apiSecret)

// Here you call with the requiered (Market and Type) arguments. See here: https://developers.cryptomkt.com/es/#ordenes
// or the documentation for more info 
book,err := client.GetBook(args.Market("ETHCLP"), args.Type("buy"))
if err != nil{
    fmt.Errorf("Error getting orders book, %s", err)
}else{
    fmt.Println(book.Data)
}

```

**Getting trades list**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)

var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey,apiSecret)

// Here you call trades from bitcoin argentinean pesos market. 
// You can see the optional arguments here: https://developers.cryptomkt.com/es/#trades 
// or in the documentation
trades,err:= client.GetTrades(args.Market("BTCARS"))
if err != nil {
     fmt.Errorf("Error getting trades, %s", err)
}
```

**Getting prices list**
```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)

var apiKey string = "YourApiKey"
var apiSecret string = "YourApiSecretKey"

client := conn.NewClient(apiKey,apiSecret)

// Here you call prices from ethereum chilean pesos market and 
// a timeframe of 60 minutes. Optional args here: https://developers.cryptomkt.com/es/#precios
// or in the documentation
prices,err := client.GetPrices(args.Market("ETHCLP"),args.TimeFrame("60"))
if err != nil{
    fmt.Errorf("Error getting prices, %s", err)
}else{
    fmt.Println(prices.Data)
}
```


### Authenticated endpoints

**Get account info**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
)
client := conn.NewClient(apiKey, apiSecret)

//account is pointer to a struct with the account info
account, err := client.Account()
if err != nil {
    fmt.Errorf("Error getting account: %s", err)
}
```


**Create order**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

order, err := client.CreateOrder(
    args.Amount(0.3),
    args.Market("ETHCLP"),
    args.Price(1000),
    args.Type("buy"))
if err != nil {
    fmt.Errorf("Error making an order: %s", err)
}
```

**Active Orders**
```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

// See the optional args here https://developers.cryptomkt.com/es/#ordenes-de-mercado
// or in the documentation
orders,err := client.GetActiveOrders(args.Market("BTCARS")) 
if err != nil{
    fmt.Errorf("Error getting active orders, %s", err)
}
```

**Executed Orders**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

// See the optional args here https://developers.cryptomkt.com/es/#ordenes-de-mercado
// or in the documentation
orders,err := client.GetExecutedOrders(args.Market("BTCARS")) 
if err != nil{
    fmt.Errorf("Error getting executed orders, %s", err)
}
```

**Order Status**
```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

// See the optional args here https://developers.cryptomkt.com/es/#ordenes-de-mercado
// or in the documentation
orders,err := client.GetOrderStatus(args.Id("YourId")) 
if err != nil{
    fmt.Errorf("Error getting order status, %s", err)
}
```
**Cancel Order**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

order,err := client.CancelOrder(args.Id("YourId"))

if err != nil{
    fmt.Errorf("Error canceling order, %s", err)
}

```