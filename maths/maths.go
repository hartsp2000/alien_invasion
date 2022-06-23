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

func MinMax(list []int) (min int, max int) {
	min = list[0]
	max = list[0]
	for _, value := range list {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func Offset(val int) int {
	if val == 0 || val > 0 {
		return 0
	}
	offset := 0
	for i := val; i < 0; i++ {
		offset++
	}
	return offset
}

func Count(min, max int) (elements int) {
	for i := min; i < max; i++ {
		elements++
	}
	return elements
}

func AbsoluteValue(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}
