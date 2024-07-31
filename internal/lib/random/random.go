package random

import (
	"math/rand"
	"time"
)

func NewRandomAlias(aliasLength int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	b := make([]byte, aliasLength)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
