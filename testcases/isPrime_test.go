package testcases

import (
	"testing"
)

func TestIsPrime(t *testing.T) {
	type tcase struct {
		number         int
		expected_ouput bool
	}
	var testcase = []tcase{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{-4, false},
		{51, true},
	}
	for i := range testcase {
		actual_output := isPrime(testcase[i].number)
		if actual_output == testcase[i].expected_ouput {
			t.Errorf("value %d actual result %t is eual to %t", testcase[i].n, actual_output, testcase[i].expected_ouput)
		}

	}

}
