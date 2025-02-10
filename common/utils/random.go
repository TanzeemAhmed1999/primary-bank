package commonutils

import (
	"math/rand"
	"strings"
	"time"
)

const aplhabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(aplhabet)

	for i := 0; i < n; i++ {
		c := aplhabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandonOwner generates random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates random currency name
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "RUP"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
