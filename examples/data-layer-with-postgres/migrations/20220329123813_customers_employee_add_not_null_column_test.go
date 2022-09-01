package migrations

import (
	"errors"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

// nolint:dupl //Cannot use same testCase for migrateUP and migrateDOWN
func TestK20220329123813_Up(t *testing.T) {
	mock, db := initializeTests(t)
	k := K20220329123813{}

	// register mock calls for success case
	mock.ExpectExec(AddNotNullColumn).WillReturnResult(sqlmock.NewResult(1, 0))
	// register mock calls for failure case
	mock.ExpectExec(invalidQuery).WillReturnError(errors.New("invalid migration"))

	testCases := []struct {
		desc string
		err  error
	}{
		{"success", nil},
		{"failure", errors.New("invalid migration")},
	}

	for i, tc := range testCases {
		err := k.Up(&db, log.NewMockLogger(io.Discard))

		assert.IsType(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

// nolint:dupl //Cannot use same testCase for migrateUP and migrateDOWN
func TestK20220329123813_Down(t *testing.T) {
	mock, db := initializeTests(t)
	k := K20220329123813{}

	// register mock calls for success case
	mock.ExpectExec(DeleteNotNullColumn).WillReturnResult(sqlmock.NewResult(1, 0))
	// register mock calls for failure case
	mock.ExpectExec(invalidQuery).WillReturnError(errors.New("invalid migration"))

	testCases := []struct {
		desc string
		err  error
	}{
		{"success", nil},
		{"failure", errors.New("invalid migration")},
	}

	for i, tc := range testCases {
		err := k.Down(&db, log.NewMockLogger(io.Discard))

		assert.IsType(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
