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
    "github.com/cryptomkt/cryptomkt-go/args"
)

```

and then create a ``Client`` object for interacting with the API:


```golang
Client, err := client.New(apiKey, apiSecret)
if err != nil{
    fmt.Errorf("Problems creating the client: %s", err)
}
```

## Making API Calls

With a `client`, you can now make API calls. We've included some examples below.  Each API method returns a ``struct`` representing the JSON response from the API.

### Public endpoints

**Listing available markets**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
)
Client, err := cryptomarket.NewClient(apiKey, apiSecret)
if err != nil{
    fmt.Errorf("Problems creating the client %s:", err)
}

marketList, err := client.Markets()
if err != nil {
    fmt.Errorf("Error while getting market list: %s", err)
}
// if we are here, then marketList has the requested data
```

### Authenticated endpoints

**Get account info**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
)
Client, err := NewClient(apiKey, apiSecret)
if err != nil{
    fmt.Errorf("Problems creating the client %s:", err)
}

account, err := client.Account()
if err != nil {
    fmt.Errorf("Error while getting account: %s", err)
}
// if we are here, then account has the requested data
```

**Create order**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/client"
    "github.com/cryptomkt/cryptomkt-go/args"
)
Client, err := cryptomarket.NewClient(apiKey, apiSecret)
if err != nil{
    fmt.Errorf("Problems creating the client %s:", err)
}

account, err := client.CreateOrder(
    args.Amount(0.3),
    args.Market("ETHCLP"),
    args.Price(1000),
    args.Type("buy"),
)
if err != nil {
    fmt.Errorf("Error while creating an order: %s", err)
}
// if we are here, then account has the requested data
```

