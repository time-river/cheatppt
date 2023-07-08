package sql

const currencyPower = 100.0

func RMB2Coins(amount float32) int64 {
	return int64(currencyPower * amount)
}

func Coins2RMB(amount int64) float32 {
	return float32(amount) / currencyPower
}
