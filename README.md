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
let apiSecret="21b12401"
// default window
let window = 0
// or some other window in miliseconds
window = 25000

client := rest.NewClient(apiKey, apiSecret, window)
ctx := context.Background()

// get currencies
currencies, err := client.GetCurrencies(ctx)

// get order books
orderBook, err := client.GetOrderBook(ctx, args.Symbol("EOSETH"))

// get your wallet balances
accountBalanceList, err := client.GetWalletBalances(ctx)

// get your trading balances
tradingBalanceList, err := client.GetSpotTradingBalances(ctx)

// move balance from wallet to spot trading
result, err := client.TransferMoneyFromAccountBalanceToTradingBalance(
  ctx,
  args.Currency("ETH"),
  args.Amount("3.2"),
  args.Source(args.AccountWallet),
  args.Destination(args.AccountSpot),
)

// get your active orders
ordersList, _ := client.GetAllActiveSpotOrders(ctx, args.Symbol("EOSETH"))

// create a new order
order, err := client.CreateSpotOrder(ctx, args.Symbol("EOSETH"), args.Side(args.SideTypeBuy), args.Quantity("10"), args.Price("10"))
```

## websocket clients

there are three websocket clients, `MarketDataClient`, the `SpotTradingClient` and the `WalletManagementClient`. The `MarketDataClient` is public, while the others require authentication to be used.

websocket subscriptions make use of notification channels. To close the notification channel of a subscription, remember to call the corrensponding Unsubscription.

### MarketDataClient

Unsubscription for the MarketDataClient are called from the subscription structure, as seen in the examples. Keep in mind that this stop the client from processing the messages, but the server will continue to send them. To completely stop recieving messages is recomended to close the MarketDataClient.

```go
// instance a client
client, err := NewMarketDataClient()
// close the client
defer client.Close()
// subscribe to public trades
subscription, err := client.SubscribeToTrades(
  args.Symbols([]string{"EOSETH", "ETHBTC"}),
  args.Limit(10),
)
subscribedSymbols := subscription.Symbols
go func() {
  for notification := range subscription.NotificationCh {
      notificationType := notification.NotificationType
      if notificationType == args.NotificationSnapshot {
        fmt.Println('is a snapshot')
      }
      if notificationType == args.NotificationUpdate {
        fmt.Println('is an update')
      }
      for _, tradeList := range notification.Data {
        for _, trade := range tradeList {
          fmt.Println(trade)
      }
    }
  }
}()
// unsubscribe
UnsubscribeTo(subscription.NotificationChannel)

// subscribe to symbol tickers
subscription, err = client.SubscribeToTicker(
  args.Symbols([]string{"EOSETH"}),
  args.TickerSpeed(args.TickerSpeed1s),
)
go func() {
  for notification := range subscription.NotificationCh {
    notificationType := notification.NotificationType
    if notificationType == args.NotificationData {
      fmt.Println('is always a data notification')
    }
    fmt.Println("tickers")
    for symbol, ticker := range notification.Data {
      fmt.Println("["+symbol+"]="+ ticker)
    }
  }
}()
// unsubscribe
UnsubscribeTo(subscription.NotificationChannel)
```

### SpotTradingClient

```go
// instance a client with default window of 10 seconds
client, err := NewSpotTradingClient(APIKey, APISecret, 0)
// close the client
defer client.Close()

// subscribe to order reports
notificationCh, err := client.SubscribeToReports()
go func() {
  for notification := range notificationCh {
    for _, report := range notification.Data {
      fmt.Println(report)
    }
  }
}()
// unsubscribe from order reports
client.UnsubscribeToReports()

clientOrderID := fmt.Sprint(time.Now().Unix())

// create an order
client.CreateSpotOrder(
  context.Background(),
  args.Symbol("EOSETH"),
  args.Side(args.SideSell),
  args.Price("1000"),
  args.Quantity("0.01"),
  args.ClientOrderID(clientOrderID),
)

// candel an order
client.CancelSpotOrder(
  context.Background(),
  args.ClientOrderID(clientOrderID),
)

```

### WalletManagementClient

```go
// instance a client with 20 seconds of window
client, err := NewWalletManagementClient(APIKey, APISecret, 20_000)
// close the client
defer client.Close()

// subscribe to wallet transactions
notificationCh, err := client.SubscribeToTransactions()
go func() {
  for notification := range notificationCh {
    transaction := notification.Data
    fmt.Println(transaction)
  }
}()

// unsubscribe from wallet transactions
err = client.UnsubscribeToTransactions()

// get wallet balances
balances, err := client.GetWalletBalances(context.Background())
for _, balance := range balances {
  fmt.Println(balance)
}

// transfer assets from wallet account and spot account
restClient.TransferBetweenWalletAndExchange(
  context.Background(),
  args.Amount("0.2"),
  args.Currency("EOS"),
  args.Source(args.AccountWallet),
  args.Destination(args.AccountSpot),
)
```

## arguments and constants of interest

all the arguments for the clients are in the args package, as well as the custom types for the arguments. check the package documentation, and the method documentation of the clients for more info.

# Checkout our other SDKs

<!-- agregar links -->

python sdk

nodejs sdk

java sdk

ruby sdk
