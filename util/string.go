package util

import (
	"math/rand"
	"time"
)

type randString struct {
	rand     *rand.Rand
	alphabet string
}

// Generates a pseudo-random string of a given length from the provided alphabet
func (n *randString) String(length int) string {
	bytes := make([]byte, length)
	for idx := range bytes {
		bytes[idx] = n.alphabet[n.rand.Intn(len(n.alphabet))]
	}

	return string(bytes)
}

var alphanumeric = randString{
	rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	alphabet: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
}

// UID generates pseudo-random 6 char long String uid
func UID() string {
	return alphanumeric.String(6)
}
