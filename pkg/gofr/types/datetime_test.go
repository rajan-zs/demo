package types

import (
	"reflect"
	"testing"

	"github.com/zopsmart/gofr/pkg/errors"
)

func TestDatetime_Check(t *testing.T) {
	tests := []struct {
		name     string
		dateTime Datetime
		err      error
	}{
		{"empty datetime struct", Datetime{}, errors.InvalidParam{Param: []string{"datetime"}}},
		{"empty timezone struct", Datetime{Value: "2018-07-14T05:00:00Z", Timezone: ".."}, errors.InvalidParam{Param: []string{"timezone"}}},
		{"correct datetime struct", Datetime{Value: "2018-07-14T05:00:00Z", Timezone: "America/New_York"}, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.dateTime)

			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("%v, Failed. Got :%v\tExpected: %v", tt.name, err, tt.err)
			}
		})
	}
}

func Test_DatetimeJson(t *testing.T) {
	tests := []struct {
		name      string
		addstruct interface{}
		err       error
	}{
		{"datetime value marshal fail", make(chan int), errors.InvalidParam{Param: []string{"datetime"}}},
		{"datetime value unmarshal fail ", struct {
			Value []string
		}{Value: []string{"hello"}}, errors.InvalidParam{Param: []string{"datetime"}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(Datetime{})

			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("%v, Failed. Got :%v\tExpected: %v", tt.name, err, tt.err)
			}
		})
	}
}
