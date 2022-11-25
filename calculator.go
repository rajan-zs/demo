package calculter

import "errors"

func calculator(inputA, inputB float64, opt string) (float64, error) {
	switch {
	case opt == "+":
		res, err := add(inputA, inputB)
		return res, err

	case opt == "-":
		res, err := sub(inputA, inputB)
		return res, err

	case opt == "*":
		res, err := multi(inputA, inputB)
		return res, err

	case opt == "/":
		res, err := div(inputA, inputB)
		return res, err

	case opt == "%":
		res, err := mod(inputA, inputB)
		return res, err

	default:
		return -1, errors.New("invalid operator")
	}

}
func add(inputA, inputB float64) (float64, error) {
	return inputA + inputB, nil
}
func sub(inputA, inputB float64) (float64, error) {
	return inputA - inputB, nil

}
func multi(inputA, inputB float64) (float64, error) {
	return inputA * inputB, nil

}
func div(inputA, inputB float64) (float64, error) {
	if inputB == 0 {
		return -1, errors.New("Dividing with 0")
	} else {
		return inputA / inputB, nil
	}
}
func mod(inputA, inputB float64) (float64, error) {
	if inputB != 0 {
		return float64(int64(inputA) % int64(inputB)), nil
	} else {
		return -1, errors.New("Dividing with 0")
	}
}
