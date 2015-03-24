package api

import (
	"crypto/rand"
	"log"
	"math/big"
	mrand "math/rand"
	"time"
)

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var max = big.NewInt(int64(len(chars)))

//fallbackRand uses less random math/rand in case of failure
func fallbackRand(err error) int {
	log.Println("Error using crypto/rand:", err)

	mrand.Seed(time.Now().UTC().UnixNano())
	return mrand.Int() % len(chars)
}

//randString returns a random string of given length using crypto/rand
func randString(length int) string {
	str := make([]byte, length)
	for i := range str {
		k, err := rand.Int(rand.Reader, max)
		if err != nil {
			j := fallbackRand(err)
			str[i] = chars[j]
		} else {
			str[i] = chars[k.Int64()]
		}
	}
	return string(str)
}
