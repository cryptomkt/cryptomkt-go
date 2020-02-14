# cryptomkt-go

Golang API of CryptoMarket

## Installation

`go get cryptomarket`

## Quick Start

The first thing you'll need to do is [sign up for cryptomkt](https://www.cryptomkt.com).

## API Key

If you're writing code for your own CryptoMarket account, [enable an API key](https://www.cryptomkt.com/platform/account#api_tab). Next, import the package:


```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
)

```

and then create a ``Client`` struct for interacting with the API:


```golang
Client := client.New(apiKey, apiSecret)

```

## Making API Calls

With a `client`, you can now make API calls. We've included some examples below.  Each API method returns a ``struct`` representing the JSON response from the API.

### Public endpoints

**Listing available markets**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
)
Client := conn.NewClient(apiKey, apiSecret)

marketList, err := client.Markets()
if err != nil {
    fmt.Errorf("Error while getting market list: %s", err)
}
```

### Authenticated endpoints

**Get account info**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
)
Client := conn.NewClient(apiKey, apiSecret)

account, err := client.Account()
if err != nil {
    fmt.Errorf("Error while getting account: %s", err)
}
```

**Create order**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
    "github.com/cryptomkt/cryptomkt-go/args"
)
Client := conn.NewClient(apiKey, apiSecret)

order, err := client.CreateOrder(
    args.Amount(0.3),
    args.Market("ETHCLP"),
    args.Price(1000),
    args.Type("buy"),
)
if err != nil {
    fmt.Errorf("Error while making an order: %s", err)
}
```

**Create Wallet**


```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

walletStatus, err := client.CreateWallet(
    args.Id("P2023132"),
    args.Token("xToY232aheSt8F"),
    args.Wallet("ETH"),
)
if err != nil {
    fmt.Errorf("Error while creating the Wallet: %s", err)
}
```

