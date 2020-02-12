package client

type Argument func(*Request)

func Amount(amount string) Argument {
	return func(request *Request) {
		request.addArgument("amount", amount)
	}
}

func Market(market string) Argument {
	return func(request *Request) {
		request.addArgument("market", market)
	}
}

func Type(atype string) Argument {
	return func(request *Request) {
		request.addArgument("type", atype)
	}
}

func Page(page string) Argument {
	return func(request *Request) {
		request.addArgument("page", page)
	}
}

func Limit(limit string) Argument {
	return func(request *Request) {
		request.addArgument("limit", limit)
	}
}

func Price(price string) Argument {
	return func(request *Request) {
		request.addArgument("price", price)
	}
}

func Currency(currency string) Argument {
	return func(request *Request) {
		request.addArgument("currency", currency)
	}
}

func Id(id string) Argument {
	return func(request *Request) {
		request.addArgument("id", id)
	}
}

func BankAccount(bankAccount string) Argument {
	return func(request *Request) {
		request.addArgument("bank_account", bankAccount)
	}
}

func Address(address string) Argument {
	return func(request *Request) {
		request.addArgument("address", address)
	}
}

func Memo(memo string) Argument {
	return func(request *Request) {
		request.addArgument("memo", memo)
	}
}

func CallbackUrl(callbackUrl string) Argument {
	return func(request *Request) {
		request.addArgument("callback_url", callbackUrl)
	}
}

func ErrorUrl(errorUrl string) Argument {
	return func(request *Request) {
		request.addArgument("error_url", errorUrl)
	}
}

func ExternalId(externalId string) Argument {
	return func(request *Request) {
		request.addArgument("external_id", externalId)
	}
}

func PaymentReceiver(paymentReceiver string) Argument {
	return func(request *Request) {
		request.addArgument("payment_receiver", paymentReceiver)
	}
}
			
func SuccessUrl(successUrl string) Argument {
	return func(request *Request) {
		request.addArgument("success_url", successUrl)
	}
}

func ToReceive(toReceive string) Argument {
	return func(request *Request) {
		request.addArgument("to_receive", toReceive)
	}
}

func ToReceiveCurrency(toReceiveCurrency string) Argument {
	return func(request *Request) {
		request.addArgument("to_receive_currency", toReceiveCurrency)
	}
}

func Token(token string) Argument {
	return func(request *Request) {
		request.addArgument("token", token)
	}
}

func Wallet(wallet string) Argument {
	return func(request *Request) {
		request.addArgument("wallet", wallet)
	}
}

func StartDate(startDate string) Argument {
	return func(request *Request) {
		request.addArgument("start_date", startDate)
	}
}

func EndDate(endDate string) Argument {
	return func(request *Request) {
		request.addArgument("end_date", endDate)
	}
}

func RefundMail(refundMail string) Argument {
	return func(request *Request) {
		request.addArgument("refund_mail", refundMail)
	}
}