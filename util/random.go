package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		char := alphabets[rand.Intn(k)]

		sb.WriteByte(char)
	}

	return sb.String()
}

func RandomOwner() string {
	return randomString(6)
}

func RandomAmount() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "KES", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", randomString(6))
}