package calculter

import "testing"

type testCase struct {
	inputA, inputB float64
	oprt           string
	expResult      float64
	description    string
}

func TestCalculator(t *testing.T) {

	tests := []testCase{
		{3, 4, "+", 7, "Addition"},
		{7, 2, "-", 5, "Substraction"},
		{7, 2, "*", 14, "multiplication"},
		{4, 2, "/", 2, "Divison"},
		{4, 2, "%", 0, "modules"},
		{2, 4, "!", -1, "invalid operator"},
	}
	for i, tc := range tests {
		res, err := calculator(tc.inputA, tc.inputB, tc.oprt)
		if res == tc.expResult && err == nil {
			t.Log("Test case is passed  ", res)
		} else {
			t.Errorf("Test no : %v   is faild beacsue of %v", i, err)
		}

	}

}
func TestAdd(t *testing.T) {

	tests := []testCase{
		{-3, 4, "+", 1, "Addition of negative numbers"},
		{-1, -1, "+", -2, "Addition of negative numbers"},
		{1.0, 2, "+", 3, "Addition of float"},
		{1.0, 2.2, "+", 3.2, "Addition of float"},
		{2, 3, "+", 5, "Addition"},
	}
	for i, tc := range tests {
		res, err := add(tc.inputA, tc.inputB)
		if res == tc.expResult && err == nil {
			t.Log("Test case is passed  ", res)
		} else {
			t.Errorf("Test no : %v   is faild beacsue of %v", i, err)
		}

	}

}
func TestSub(t *testing.T) {

	tests := []testCase{
		{7, 2, "-", 5, "Substraction"},
		{7.0, 2, "-", 5, "Substraction of float"},
		{2, 7, "-", -5, "Substracting from smaller number"},
		{-2, -7, "-", 5, "Substraction of both -ve number"},
		{-7, -2, "-", -5, "Substraction of both -ve number"},
	}
	for i, tc := range tests {
		res, err := sub(tc.inputA, tc.inputB)
		if res == tc.expResult && err == nil {
			t.Log("Test case is passed  ", i, res)
		} else {
			t.Errorf("Test no : %v   is faild beacsue of %v", i, err)
		}

	}

}
func TestMulti(t *testing.T) {

	tests := []testCase{
		{0, 2, "*", 0, "multiplication with 0"},
		{2, 0, "*", 0, "multiplication with 0"},
		{-1, -2, "*", 2, "multiplication of both -ve number"},
		{2, -1, "*", -2, "multiplication with -ve number"},
		{-1, 2, "*", -2, "multiplication with -ve number"},
		{2.1, 1.2, "*", 2.52, "multiplication of float"},
	}
	for i, tc := range tests {
		res, err := multi(tc.inputA, tc.inputB)
		if res == tc.expResult && err == nil {
			t.Log("Test case is passed  ", i, res)
		} else {
			t.Errorf("Test no : %v   is faild beacsue of %v", i, err)
		}

	}

}
func TestDiv(t *testing.T) {

	tests := []testCase{
		{1, 2, "/", 0.5, "Divison with smaller number"},
		{6, 0, "/", -1, "Dividing with 0"},
		{0, 6, "/", 0, "Dividing to 0"},
		{4, 2, "/", 2, "Dividing"},
		{-4, -2, "/", 2, "Divison of -ve number"},
		{-4, 2, "/", -2, "Divison of -ve number"},
		{4, -2, "/", -2, "Divison of -ve number"},
		{4.2, 2.1, "/", 2, "Divison of float"},
	}
	for i, tc := range tests {
		res, err := div(tc.inputA, tc.inputB)
		if res == tc.expResult && err == nil {
			t.Log("Test case is passed  ", i, res)
		} else {
			t.Errorf("Test no : %v   is faild beacsue of %v", i, err)
		}

	}

}
func TestModule(t *testing.T) {

	tests := []testCase{
		{4, 2, "%", 0, "modules"},
		{0, 2, "%", 0, "modules"},
		{6, 0, "%", -1, "Dividing with 0"},
		{2, 4, "%", 2, "modules"},
	}
	for i, tc := range tests {
		res, err := mod(tc.inputA, tc.inputB)
		if res == tc.expResult && err == nil {
			t.Log("Test case is passed  ", i, res)
		} else {
			t.Errorf("Test no : %v   is faild beacsue of %v", i, err)
		}

	}

}
