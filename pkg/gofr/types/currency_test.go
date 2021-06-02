package types

import (
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

func TestCurrency_Check(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		err      error
	}{
		{"empty string passed as currency", "", errors.InvalidParam{Param: []string{"currency"}}},
		{"wrong currency code passed", "usuk 22", errors.InvalidParam{Param: []string{"currencyCountryCode"}}},
		{"wrong currency value passed", "USD ABCD", errors.InvalidParam{Param: []string{"currencyValue"}}},
		{"correct currency format", "USD 10.15", nil},
		{"currency format with negative value", "USD -66.6", nil},
		{"wrong currency format passed", "USD ABCD efg", errors.InvalidParam{Param: []string{"currency"}}},
		{"symbolic representation of currency", "$ 123.00", errors.InvalidParam{Param: []string{"currencyCountryCode"}}},
	}
	for _, tt := range tests {
		tt := tt
		err := Validate(tt.currency)

		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("%v, Failed. Got :%v\tExpected: %v", tt.name, err, tt.err)
		}
	}
}
