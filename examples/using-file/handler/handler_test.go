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
		desc       string
		expErr     error
		resp       interface{}
		fileName   string
		fileAccess file.Mode
		ctx        *gofr.Context
	}{
		{"Successful Read", nil, "Welcome to Zopsmart!", "../temp/test.txt", "rw", ctx},
		{desc: "Invalid file", expErr: &fs.PathError{Op: "open", Path: "", Err: syscall.Errno(2)}, fileAccess: "rw", ctx: ctx},
		{"Failed to Read", errors.New("EOF"), nil, "../temp/abc.txt", "rw", ctx},
	}
	for i, tc := range testcases {
		l, err := file.NewWithConfig(app.Config, tc.fileName, tc.fileAccess)
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
		desc       string
		expErr     error
		resp       interface{}
		ctx        *gofr.Context
		fileName   string
		fileAccess file.Mode
	}{
		{"successful write", nil, "File written successfully!", ctx, "../temp/test.txt", "rw"},
		{"failed write - invalid permission", &fs.PathError{Op: "write", Path: "../temp/test.txt", Err: syscall.Errno(9)},
			nil, ctx, "../temp/test.txt", "0444"},
		{desc: "failed write -invalid path", expErr: &fs.PathError{Op: "open", Path: "", Err: syscall.Errno(2)}, fileAccess: "rw", ctx: ctx},
	}

	for i, tc := range testcases {
		l, err := file.NewWithConfig(app.Config, tc.fileName, tc.fileAccess)
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
