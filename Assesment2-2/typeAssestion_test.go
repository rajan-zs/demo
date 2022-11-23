package Assesment2_2

import "testing"

type testcase struct {
	desc    string
	input   interface{}
	expResp interface{}
}

func TestTypeAssert(t *testing.T) {
	tcase := []testcase{
		{"passing integer", 10, 100},
		{"passing Float", 10, 105},
		{"passing String", "rajan", "Hello rajan"},
		{"passing bool", false, true},
		{"passing special characters", "@", "unknown"},
	}
	for i, tc := range tcase {
		resp, err := typeAssert(tc.input)
		if tc.expResp != resp && err != nil {
			t.Errorf("Test case number %v is faild Beacuse of %v", i, err)
		}
	}
}
