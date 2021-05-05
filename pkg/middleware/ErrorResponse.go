package middleware

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zopsmart/gofr/pkg/errors"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, logger logger, err errors.MultipleErrors) {
	// setting default content type to be application/json
	contentType := "application/json"
	reqContentType := r.Header.Get("Content-Type")
	errByte, _ := json.Marshal(err)
	errOutput := string(errByte)

	if logger != nil {
		logger.AddData(string(CorrelationIDKey), r.Context().Value(CorrelationIDKey))
		logger.Errorf("%v", err)
		// pushing error type to prometheus for 500s only
		if err.StatusCode == http.StatusInternalServerError {
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			// remove the trailing slash
			path = strings.TrimSuffix(path, "/")

			ErrorTypesStats.With(prometheus.Labels{"type": "UnknownError", "path": path, "method": r.Method}).Inc()
		}
	}

	switch reqContentType {
	case "text/xml", "application/xml":
		contentType = reqContentType
		errByte, _ := xml.Marshal(err)
		errOutput = string(errByte)
	case "text/plain":
		contentType = reqContentType
		errOutput = err.Error()
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(err.StatusCode)
	fmt.Fprintln(w, errOutput)
}

func FetchErrResponseWithCode(statusCode int, reason, code string) *errors.MultipleErrors {
	zone, _ := time.Now().Zone()

	return &errors.MultipleErrors{
		StatusCode: statusCode,
		Errors: []error{&errors.Response{
			Code:   code,
			Reason: reason,
			DateTime: errors.DateTime{
				Value:    time.Now().Format(time.RFC3339),
				TimeZone: zone,
			},
		},
		},
	}
}
