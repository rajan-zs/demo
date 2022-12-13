package handler

import (
	"errors"
	"io/fs"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/file"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
)

func Test_Read(t *testing.T) {
	c := config.MockConfig{Data: map[string]string{
		"FILE_STORE": "LOCAL",
	}}
	app := gofr.NewWithConfig(&c)
	ctx := gofr.NewContext(nil, nil, app)

	testcases := []struct {
		desc     string
		expErr   error
		resp     interface{}
		fileName string
		fileMode string
		ctx      *gofr.Context
	}{
		{"Successful Read", nil, "Welcome to Zopsmart!", "../test.txt", "rw", ctx},
		{desc: "Invalid file", expErr: &fs.PathError{Op: "open", Path: "", Err: syscall.Errno(2)}, fileMode: "rw", ctx: ctx},
		{"Failed to Read", errors.New("EOF"), nil, "../abc.txt", "rw", ctx},
	}

	for i, tc := range testcases {
		l, err := file.NewWithConfig(app.Config, tc.fileName, file.Mode(tc.fileMode))
		if err != nil {
			t.Error("Unable to initialize", err)
		}

		f := New(l)

		resp, err := f.Read(tc.ctx)
		assert.Equalf(t, tc.resp, resp, "Test case %v failed.\nExpected: %v, got: %v", i, tc.resp, resp)

		assert.Equalf(t, tc.expErr, err, "Test case %v failed.\nExpected: %v, got: %v", i, tc.expErr, err)
	}
}

func Test_Write(t *testing.T) {
	c := config.MockConfig{Data: map[string]string{
		"FILE_STORE": "LOCAL",
	}}
	app := gofr.NewWithConfig(&c)
	ctx := gofr.NewContext(nil, nil, app)

	testcases := []struct {
		expErr   error
		resp     interface{}
		ctx      *gofr.Context
		fileName string
		fileMode file.Mode
	}{
		{nil, "File written successfully!", ctx, "../test.txt", "rw"},
		{&fs.PathError{Op: "write", Path: "../test.txt", Err: syscall.Errno(9)}, nil, ctx, "../test.txt", "0444"},
		{expErr: &fs.PathError{Op: "open", Path: "", Err: syscall.Errno(2)}, fileMode: "rw", ctx: ctx},
	}

	for i, tc := range testcases {
		l, err := file.NewWithConfig(app.Config, tc.fileName, tc.fileMode)
		if err != nil {
			t.Error("Unable to initialize", err)
		}

		f := New(l)
		resp, err := f.Write(tc.ctx)
		assert.Equalf(t, tc.resp, resp, "Test case %v failed.\nExpected: %v, got: %v", i, tc.resp, resp)
		assert.Equalf(t, tc.expErr, err, "Test case %v failed.\nExpected: %v, got: %v", i, tc.expErr, err)
	}
}

func Test_List(t *testing.T) {
	c := &config.MockConfig{Data: map[string]string{
		"FILE_STORE": "LOCAL",
	}}
	app := gofr.NewWithConfig(c)
	ctx := gofr.NewContext(nil, nil, app)

	l, err := file.NewWithConfig(app.Config, "test.txt", "rw")
	if err != nil {
		t.Error("Unable to initialize", err)
	}

	f := New(l)

	expRes := []string{"handler.go", "handler_test.go"}
	res, err := f.List(ctx)
	assert.Equal(t, expRes, res, "Test case failed. Expected: %v, got: %v", expRes, res)
	assert.Equal(t, nil, err, "Test case failed. Expected: %v, got: %v", nil, err)
}
