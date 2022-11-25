package factory

import (
	"database/sql"
	"net/http"
)

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) handler {
	return handler{
		db: db,
	}
}
func (h handler) empHandlerGet(w http.ResponseWriter, e *http.Request) {
	return h.db.Query("")
}
