package main

import "testing"

type err struct {
	str string
}

func (e err) Error() string {
	return e.str
}

type empData struct {
	id   string
	name string
	add  string
}
type testcase struct {
	desc    string
	Input   empData
	expResp empData
	err     error
}

func TestEmployee(t *testing.T) {

	tcase := []testcase{{
		desc: "post request check",
		Input: empData{
			id:   "1",
			name: "Rajan",
			add:  "Patna",
		},
		expResp: empData{
			id:   "1",
			name: "Rajan",
			add:  "Patna",
		},
		err: nil},
		{
			desc: "post request check without id",
			Input: empData{
				name: "Rajan",
				add:  "Patna",
			},
			expResp: empData{
				id:   "",
				name: "Rajan",
				add:  "Patna",
			},
			err: nil,
		},
		{
			desc: "post request check without name",
			Input: empData{
				id:   "1",
				name: "",
				add:  "Patna",
			},
			expResp: empData{
				id:  "1",
				add: "Patna",
			},
			err: nil,
		},
		{
			desc: "post request check without add",
			Input: empData{
				id:   "1",
				name: "Rajan",
			},
			expResp: empData{
				id:   "1",
				name: "Rajan",
				add:  "",
			},
			err: nil,
		},
	}
}
