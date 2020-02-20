# cryptomkt-go

cryptomkt-go is the SDK for cryptomkt in the GO programing language

## Installing
To install the sdk, run the `go get` command

`go get github.com/cryptomkt/cryptomkt-go`

## Documentation

For further information about the skd see the [godoc documentation](link to godoc once the repo is available) of the module.

The base api for this skd can be foun in [here](https://developers.cryptomkt.com/)

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
All calls return at least one informative error value as an unmeeted argument, an invalid apiKey or a "not_enough_balance" as a replay from the server if you try to buy more than your money can take.
Some calls return structs (or slice of structs).
And lastly we have some calls that return the same information as before, buy instead of using structs, they use map\[string\]string (or slices of maps) to store all values.

for example, there are two calls that can give you the account information:

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

//accountAsMap is the same information of account, but stored in a map[string]string
accountAsMap, err := client.GetAccountAsMap()
if err != nil {
    fmt.Errorf("Error getting the account: %s", err)
}
```

The calls that returns maps end with 'AsMap' or 'AsMapList' in contrast with the ones that return structs, that have the same name, but without the suffix.

The advantage of the map format is its simplicity and ease of use, while using structs gives aditional functionality over the recieved data. 

for example, if we want to go over a long range of trade data of a market, we can call `client.getTrades` to get a list of `Trades`, this list can be one page of many, so once we read the data on the page, to get the rest of the pages, we can call over an over `GetNext()` over the struct, until an empty list is returned. Here is in code:

/* por implementar (big requests)*/
To protect from attacks, Cryptomarket only accepts a maximum amount of message per minute. If you go over this number, your ip is blocked so you can't keep making request using neither the skd nor the api. In order to keep your ip usable, big requests, as getting all trades from 2019 will make one request to the server evey 3 seconds. So, the bigger the request, the slower.

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

//account is a pointer to a struct that matches the account information
trades, err := client.GetAllTrades(
    args.Market("ETHARS"),
    args.Start("2019-05-10"),
    args.End("2019-12-10"),
)
if err != nil {
    fmt.Errorf("Error getting the trades: %s", err)
}

```

## API Calls Examples

here we include some API calls examples

### Public endpoints

**Listing available markets**

```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
)
client := conn.NewClient(apiKey, apiSecret)

// marketList is a list of enabled markets
marketList, err := client.Markets()
if err != nil {
    fmt.Errorf("Error getting the market list: %s", err)
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

//sameAccount is the same information of account, but stored in a map[string]string
sameAccount, err := client.GetAccountAsMap()
if err != nil {
    fmt.Errorf("Error getting the account: %s", err)
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
    args.Type("buy"),
)
if err != nil {
    fmt.Errorf("Error making an order: %s", err)
}
```

**Create Wallet**


```golang
import (
    "github.com/cryptomkt/cryptomkt-go/conn"
    "github.com/cryptomkt/cryptomkt-go/args"
)
client := conn.NewClient(apiKey, apiSecret)

err := client.CreateWallet(
    args.Id("P2023132"),
    args.Token("xToY232aheSt8F"),
    args.Wallet("ETH"),
)
if err != nil {
    fmt.Errorf("Error creating the Wallet: %s", err)
}
```

