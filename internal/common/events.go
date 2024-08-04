package common

var (
	UserDataExchange string = "exchange.customer.user_data_" + "local"
	UserDataSent     string = UserDataExchange + "_login"

	CreditLimitExchange   string = "exchange.customer.credit_limit_" + "local"
	CreditLimitDataQueue  string = CreditLimitExchange + "_check"
	CreditLimitDataUpdate string = CreditLimitExchange + "_update"
)
