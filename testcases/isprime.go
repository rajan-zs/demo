package testcases

import (
	"math"
)

func isPrime(num1 int) bool {
	Prime := true
	for i := 2; i <= int(math.Sqrt(float64(num1))); i++ {
		if num1%i == 0 {
			Prime = false
			return Prime
		}
	}
	return Prime
}
