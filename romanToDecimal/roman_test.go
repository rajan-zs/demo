package romanToDecimal

import (
	"fmt"
	"testing"
)

type testcase struct {
	desc      string
	input     string
	expResult int
	err       error
}
type err struct {
	str string
}

func (e err) Error() string {
	return e.str
}
func TestRoman(t *testing.T) {
	tests := []testcase{
		{"Roman to decimal", "I", 1, nil},
		{"Roman to decimal", "III", 3, nil},
		{"Roman to decimal", "Iv", 4, nil},
		{"Roman to decimal", "V", 5, nil},
		{"Roman to decimal", "VI", 6, nil},
		{"Roman to decimal", "X", 10, nil},
		{"Roman to decimal", "XI", 11, nil},
		{"Roman to decimal", "L", 50, nil},
		{"Roman to decimal", "IL", 49, nil},
		{"Roman to decimal", "LIV", 54, nil},
		{"Roman to decimal", "C", 100, nil},
		{"Roman to decimal", "D", 500, nil},
		{"Roman to decimal", "M", 1000, nil},
		{"invlaid roman", "Mk", -1, err{"error - invalid roman"}},
		{"invlaid roman", "", -1, err{"error - empty input"}},
	}
	for i, tc := range tests {
		res, resErr := roman(tc.input)
		if res != tc.expResult && tc.err != resErr {
			fmt.Errorf("invalide input Test case : %d is faild because of %v", i, tc.err)
		}

	}

}
