package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CNY = "CNY"
)

// IsSupportedCurrency checks if the input currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CNY:
		return true
	}
	return false
}
