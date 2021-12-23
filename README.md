# CryptoMarket-go
[main page](https://www.cryptomkt.com/)


[sign up in CryptoMarket](https://www.cryptomkt.com/account/register).

# Installation
To install the Cryptomarket client
```
 go get github.com/cryptomarket/cryptomarket-go
```
# Documentation
This sdk makes use of the [api version 2](https://api.exchange.cryptomkt.com/v2) of cryptomarket


# Quick Start

## rest client
```go
import (
	"context"

    "github.com/cryptomarket/cryptomarket-go/args"
    "github.com/cryptomarket/cryptomarket-go/rest"
)

// instance a client
let apiKey="AB32B3201"
let api_secret="21b12401"
client := rest.NewClient(apiKey, api_secret)
ctx := context.Background()

// get currencies
currencies, err := client.GetCurrencies(ctx)

// get order books
orderBook, err := client.GetOrderBook(ctx, args.Symbol("EOSETH"))

// get your account balances
accountBalanceList, err := client.GetAccountBalance(ctx)

// get your trading balances
tradingBalanceList, err := client.GetTradingBalance(ctx)

// move balance from account to trading
result, err := client.TransferMoneyFromAccountBalanceToTradingBalance(ctx, args.Currency("ETH"), args.Amount("3.2"))

// get your active orders
ordersList, _ := client.GetActiveOrders(ctx, args.Symbol("EOSETH"))

// create a new order
order, err := client.CreateOrder(ctx, args.Symbol("EOSETH"), args.Side(args.SideTypeBuy), args.Quantity("10"), args.Price("10"))
```

## websocket client

There are three diferent websocket clients, the public client, the trading client and the account client.

websocket requests also accept contexts, except subscriptions and unsubscriptions.

subscriptions returns a channel for the subscription feed

```go
import (
	"context"

    "github.com/cryptomarket/cryptomarket-go/args"
    "github.com/cryptomarket/cryptomarket-go/websocket"
)

let apiKey="AB32B3201"
let api_secret="21b12401"
publicClient, err := websocket.NewPublicClient()
tradingClient, err := websocket.NewTradingClient(apiKey, api_secret)
accountClient, err := websocket.NewAccountClient(apiKey, api_secret)

// get currencies
currencies, err := publicClient.GetCurrencies(ctx)


// get your account balances
accountBalanceList, err := accountClient.GetAccountBalance(ctx)


// get your trading balances
tradingBalanceList, err := tradingClient.GetTradingBalance(ctx)

// get your active orders
ordersList := tradingClient.GetActiveOrders(ctx, args.Symbol("EOSETH"))

// create a new order
order, err := tradingClient.CreateOrder(ctx, args.ClientOrderID("aBcDeFgHi"), args.Symbol("EOSETH"), args.Side(args.SideTypeBuy), args.Quantity("10"), args.Price("10"))



// a subscription
feedChannel, err := publicClient.SubscribeToTicker(args.Symbol("ETHBTC"))
if err != nil {
    fmt.Println(err)
}
go func() {
    for ticker := range feedChannel {
        fmt.Println(ticker)
    }
}
time.Sleep(10)
if err := publicClient.SubscribeToTicker(args.Symbol("ETHBTC")); err != nil {
    fmt.Prtintln(err)
    // if the error is from the exchange, and not from the sdk, then the feedChannel is still efectively closed
    // and all subsequent updates are dropped.
}
```

## error handling
for the rest client and the three websocket clients, all requests accepts (not subcriptions or unsubscriptions) context for cancelation.

an informative error is also forwarded from the exchange server if an error is present in a request.

```go
client, err := websocket.NewPublicClient()
if err != nil {
    t.Fatal(err)
}
ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*10))
defer cancelFunc()
if _, err := client.GetCurrency(ctx, args.Currency("EOS")); err.Error() == ctx.Err().Error() { 
    // good, this error is expected
}
```
## arguments and constants of interest
all the arguments for the clients are in the args package, as well as the custom types for the arguments. check the package documentation, and the method documentation of the clients for more info.

# Checkout our other SDKs
<!-- agregar links -->
python sdk

nodejs sdk

java sdk

ruby sdk