package corev

import (
	"math/rand"
)

const (
	runes1 = "abcdefghijklmnopqrstuvwxyz"
	runes2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	runes3 = "0123456789"
	runes4 = runes3 + runes1 + runes2
)

func RandStringLower(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = runes1[rand.Intn(len(runes1))]
	}
	return string(b)
}

func RandStringUpper(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = runes2[rand.Intn(len(runes2))]
	}
	return string(b)
}

func RandStringNumber(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = runes3[rand.Intn(len(runes3))]
	}
	return string(b)
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = runes4[rand.Intn(len(runes4))]
	}
	return string(b)
}
