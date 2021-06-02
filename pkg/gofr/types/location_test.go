package types

import (
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

func TestLocation_Check(t *testing.T) {
	tests := []struct {
		name string
		lat  Latitude
		lng  Longitude
		err  error
	}{
		{"invalid latitude passed", 1234, 0, errors.InvalidParam{Param: []string{"lat"}}},
		{"invalid longitude passed", 34.00, 1234, errors.InvalidParam{Param: []string{"lng"}}},
		{"correct location struct", 34.00, 12, nil},
		{"invalid latitude and longitude", -100, 1110, errors.InvalidParam{Param: []string{"lat"}}},
	}
	for _, tt := range tests {
		tt := tt
		locstruct := Location{
			Latitude:  &tt.lat,
			Longitude: &tt.lng,
		}

		t.Run(tt.name, func(t *testing.T) {
			err := Validate(locstruct)

			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("%v, Failed. Got :%v\tExpected: %v", tt.name, err, tt.err)
			}
		})
	}
}

func TestLocation_CheckEmpty(t *testing.T) {
	var testLat Latitude

	var testLong Longitude

	testLat = 24.00
	testLong = 12.00
	tests := []struct {
		name string
		lat  *Latitude
		lng  *Longitude
		err  error
	}{
		{"empty latitude", nil, &testLong, errors.InvalidParam{Param: []string{"lat is nil"}}},
		{"empty longitude", &testLat, nil, errors.InvalidParam{Param: []string{"lng is nil"}}},
		{"empty latitude and longitude", nil, nil,
			errors.MultipleErrors{Errors: []error{errors.InvalidParam{Param: []string{"lat is nil"}},
				errors.InvalidParam{Param: []string{"lng is nil"}}}}},
	}

	for _, tt := range tests {
		tt := tt
		locstruct := Location{
			Latitude:  tt.lat,
			Longitude: tt.lng,
		}

		t.Run(tt.name, func(t *testing.T) {
			err := Validate(locstruct)

			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("%v, Failed. Got :%v\tExpected: %v", tt.name, err, tt.err)
			}
		})
	}
}
