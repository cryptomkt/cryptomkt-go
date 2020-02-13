package client

type Argument func(*Request)

func Amount(val string) Argument {
	return func(request *Request) {
		request.addArgument("amount", val)
	}
}

func Market(val string) Argument {
	return func(request *Request) {
		request.addArgument("market", val)
	}
}

func Type(val string) Argument {
	return func(request *Request) {
		request.addArgument("type", val)
	}
}

func Page(val string) Argument {
	return func(request *Request) {
		request.addArgument("page", val)
	}
}

func Limit(val string) Argument {
	return func(request *Request) {
		request.addArgument("limit", val)
	}
}

func Price(val string) Argument {
	return func(request *Request) {
		request.addArgument("price", val)
	}
}

func Currency(val string) Argument {
	return func(request *Request) {
		request.addArgument("currency", val)
	}
}

func Id(val string) Argument {
	return func(request *Request) {
		request.addArgument("id", val)
	}
}

func BankAccount(val string) Argument {
	return func(request *Request) {
		request.addArgument("bank_account", val)
	}
}

func Address(val string) Argument {
	return func(request *Request) {
		request.addArgument("address", val)
	}
}

func Memo(val string) Argument {
	return func(request *Request) {
		request.addArgument("memo", val)
	}
}

func CallbackUrl(val string) Argument {
	return func(request *Request) {
		request.addArgument("callback_url", val)
	}
}

func ErrorUrl(val string) Argument {
	return func(request *Request) {
		request.addArgument("error_url", val)
	}
}

func ExternalId(val string) Argument {
	return func(request *Request) {
		request.addArgument("external_id", val)
	}
}

func PaymentReceiver(val string) Argument {
	return func(request *Request) {
		request.addArgument("payment_receiver", val)
	}
}
			
func SuccessUrl(val string) Argument {
	return func(request *Request) {
		request.addArgument("success_url", val)
	}
}

func ToReceive(val string) Argument {
	return func(request *Request) {
		request.addArgument("to_receive", val)
	}
}

func ToReceiveCurrency(val string) Argument {
	return func(request *Request) {
		request.addArgument("to_receive_currency", val)
	}
}

func Language(val string) Argument {
	return func(request *Request) {
		request.addArgument("language", val)
	}
}

func Token(val string) Argument {
	return func(request *Request) {
		request.addArgument("token", val)
	}
}

func Wallet(val string) Argument {
	return func(request *Request) {
		request.addArgument("wallet", val)
	}
}

func StartDate(val string) Argument {
	return func(request *Request) {
		request.addArgument("start_date", val)
	}
}

func EndDate(val string) Argument {
	return func(request *Request) {
		request.addArgument("end_date", val)
	}
}

func RefundMail(val string) Argument {
	return func(request *Request) {
		request.addArgument("refund_email", val)
	}
}