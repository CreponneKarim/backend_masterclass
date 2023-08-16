package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min,max int64) int64 {
	return	min+rand.Int63n(max - min + 1)
}

func RandomString(n int) string{
	var sb strings.Builder

	k := len(alphabet)

	for i:=0;i<n; i++ {
		l := alphabet[rand.Intn(k)]
		sb.WriteByte(l)
	}
	return sb.String()
}

func RandomOwnerName() string {
	return	RandomString(8)
}

func RandomBalance() int64 {
	return	RandomInt(0,1000)
}

func RandomCurrency() string {
	currencies := []string {"EU", "USD", "DZ", "CAD"}
	n :=len(currencies)
	return currencies[rand.Intn(n)]
}