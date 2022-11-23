package assesment2_1

import (
	"reflect"
	"testing"
)

type testcase struct {
	desc    string
	input   []int
	expResp []int
}

func TestReverse(t *testing.T) {

	tcase := []testcase{
		{"Reverse of slice", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"raverse of empty slice", []int{}, []int{}},
		{"Reverse of slice", []int{1.0, 2, 3, 4}, []int{4, 3, 2, 1}},
	}

	for i, tc := range tcase {
		resp, err := reverse(tc.input)
		if reflect.DeepEqual(resp, tc.expResp) && err != 0 {
			t.Errorf("Test case %v is faild because of %v", i, err)
		}
	}
}
