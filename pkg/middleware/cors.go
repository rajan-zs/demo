package middleware

import (
	"net/http"
)

const (
	allowedHeaders = "Authorization, Content-Type, x-requested-with, origin, true-client-ip, x-correlation-id, x-zopsmart-tenant"
	allowedMethods = "PUT, POST, GET, DELETE, OPTIONS"
)

func CORS(envHeaders map[string]string) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			corsHeadersConfig := getValidCORSHeaders(envHeaders)
			for k, v := range corsHeadersConfig {
				w.Header().Set(k, v)
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			inner.ServeHTTP(w, r)
		})
	}
}

// getValidCORSHeaders returns a validated map of CORS headers.
// values specified in env are present in envHeaders
func getValidCORSHeaders(envHeaders map[string]string) map[string]string {
	validCORSHeadersAndValues := make(map[string]string)

	for _, header := range AllowedCORSHeader() {
		// If config is set, use that
		if val, ok := envHeaders[header]; ok && val != "" {
			validCORSHeadersAndValues[header] = val
			continue
		}

		// If config is not set - for the three headers, set default value.
		switch header {
		case "Access-Control-Allow-Origin":
			validCORSHeadersAndValues[header] = "*"
		case "Access-Control-Allow-Headers":
			validCORSHeadersAndValues[header] = allowedHeaders
		case "Access-Control-Allow-Methods":
			validCORSHeadersAndValues[header] = allowedMethods
		}
	}

	val := validCORSHeadersAndValues["Access-Control-Allow-Headers"]

	if val != allowedHeaders {
		validCORSHeadersAndValues["Access-Control-Allow-Headers"] = allowedHeaders + ", " + val
	}

	return validCORSHeadersAndValues
}

func AllowedCORSHeader() []string {
	return []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Credentials",
		"Access-Control-Expose-Headers",
		"Access-Control-Max-Age",
	}
}
