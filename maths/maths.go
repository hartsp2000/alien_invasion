package maths

import (
	"math/rand"
	"time"
)

func RandBetween(min int, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	randnum := rand.New(source)

	return randnum.Intn(max-min+1) + min
}

func Abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}
