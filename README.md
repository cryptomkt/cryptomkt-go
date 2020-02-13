# cryptomkt-go

## Installation

`go get cryptomarket`

## Quick Start

The first thing you'll need to do is [sign up for cryptomkt](https://www.cryptomkt.com).

## API Key

If you're writing code for your own CryptoMarket account, [enable an API key](https://www.cryptomkt.com/platform/account#api_tab). Next, import the package:


```golang
import (
    "cryptomarket"
)

```

and then create a ``Client`` object for interacting with the API:


```golang
Client, err := cryptomarket.NewClient(apiKey, apiSecret)
if err != nil{
    fmt.Errorf("Problems creating the client %s:", err)
}
```

## Making API Calls

With a `client`, you can now make API calls. We've included some examples below.  Each API method returns a ``struct`` representing the JSON response from the API.

### Public endpoints

**Listing available markets**

```golang
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
Client, err := cryptomarket.NewClient(apiKey, apiSecret)
if err != nil{
    fmt.Errorf("Problems creating the client %s:", err)
}

account, err := client.CreateOrder(
    Amount(0.3),
    Market("ETHCLP"),
    Price(1000),
    Type("buy"),
)
if err != nil {
    fmt.Errorf("Error while creating an order: %s", err)
}
// if we are here, then account has the requested data
```