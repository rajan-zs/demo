package Triangle

import "testing"

func TestTraingle(t *testing.T) {
	type tcase struct {
		s1  int
		s2  int
		s3  int
		exp string
	}
	var testcase = []tcase{
		{2, 3, 4, "scalar"},
		{2, 4, 3, "scalar"},
		{4, 3, 2, "scalar"},
		{2, 2, 2, "ISO"},
		{3, 3, 5, "Eque"},
		{3, 5, 3, "Eque"},
		{5, 3, 3, "Eque"},
		{1, 1, 3, "not valid"},
		{0, 1, 3, "not valid"},
		{1, 0, 3, "not valid"},
		{3, 0, 1, "not valid"},
		{0, 0, 0, "not valid"},
		{-2, 3, 5, "not valid"},
		{-2, 3, 5, "not valid"},
		{2, -3, 5, "not valid"},
		{2, 3, -5, "not valid"},
	}
	for i := range testcase {
		act := traingle(testcase[i].s1, testcase[i].s2, testcase[i].s3)
		if act == testcase[i].exp {
			t.Log("Test case is passed and Traingle is " + act)
		} else {
			t.Errorf("Test cases is faild beacsue of sum of two side is not grater that 3rd side")
		}

	}

}
