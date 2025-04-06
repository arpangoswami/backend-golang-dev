package util

import (
	"database/sql"
	"math"
	"math/rand"
	"strings"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.New(rand.NewSource(rand.Int63()))
}

// RandomInt returns a random integer in the range min -> max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString returns a random string of size "length"
func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < length; i++ {
		ch := alphabet[rand.Intn(k)]
		sb.WriteByte(ch)
	}
	return sb.String()
}

// RandomOwner returns a random string of length 6
func RandomOwner() string {
	return RandomString(6)
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// RandomMoney returns a random amount of money
func RandomMoney(maxAmount int64) float64 {
	return ToFixed(float64(RandomInt(0, maxAmount))+rand.Float64(), 2)
}

// RandomCurrencyCodes returns a random Currency code and its corresponding country code from given list
type CurrencyCountryCode struct {
	CurrencyCode string
	CountryCode  sql.NullInt32
}

func RandomCurrencyCodeCountryCode() CurrencyCountryCode {
	currencyCodes := []CurrencyCountryCode{
		{
			CurrencyCode: "USD",
			CountryCode:  sql.NullInt32{Int32: 0, Valid: true},
		},
		{
			CurrencyCode: "INR",
			CountryCode:  sql.NullInt32{Int32: 1, Valid: true},
		},
		{
			CurrencyCode: "EUR",
			CountryCode:  sql.NullInt32{Int32: 2, Valid: true},
		},
		{
			CurrencyCode: "GBP",
			CountryCode:  sql.NullInt32{Int32: 3, Valid: true},
		},
		{
			CurrencyCode: "JPY",
			CountryCode:  sql.NullInt32{Int32: 4, Valid: true},
		},
	}
	n := len(currencyCodes)
	return currencyCodes[rand.Intn(n)]
}
