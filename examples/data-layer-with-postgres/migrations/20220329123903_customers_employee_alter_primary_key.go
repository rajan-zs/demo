package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20220329123903 struct {
}

func (k K20220329123903) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(AlterPrimaryKey)
	if err != nil {
		return err
	}

	return nil
}

func (k K20220329123903) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(ResetPrimaryKey)
	if err != nil {
		return err
	}

	return nil
}
