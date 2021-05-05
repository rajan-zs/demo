package gofr

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zopsmart/gofr/pkg/log"
)

func TestRedirectHttps(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	k := New()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(k.Server.redirectHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

	if !strings.Contains(rr.Header().Get("Location"), "https://") {
		t.Errorf("handler returned a non https request")
	}

	if rr.Header().Get("Strict-Transport-Security") != "max-age=63072000; includeSubDomains" {
		t.Errorf("handler returned wrong header : ")
	}
}

// TestGofrHttpsStart tests if https server can be started while the port is already being used by another server
func TestGofrHttpsStart(t *testing.T) {
	httpsServer := &HTTPS{
		Port:            443,
		TLSConfig:       nil,
		CertificateFile: "../../examples/sample-https/configs/server.crt",
		KeyFile:         "../../examples/sample-https/configs/server.key.test",
	}

	// starting an https server on the same port.
	// nolint:errcheck // Error return value of http.ListenAndServeTLS not checked
	go http.ListenAndServeTLS(":443", httpsServer.CertificateFile, httpsServer.KeyFile, nil)

	time.Sleep(3 * time.Second)

	buf := new(bytes.Buffer)
	httpsServer.StartServer(log.NewMockLogger(buf), nil)

	if !strings.Contains(buf.String(), "unable to start HTTPS Server") {
		t.Errorf("was able to start https server on port while server was already running")
	}
}

func TestHTTPSFail(t *testing.T) {
	k := New()

	httpsServer := &HTTPS{
		Port:            9011,
		TLSConfig:       nil,
		CertificateFile: "../../examples/sample-https/configs/server.crt",
		KeyFile:         "failtestkey.pem",
	}
	httpsServer.StartServer(k.Logger, k.Server.Router)

	req, _ := http.NewRequest("GET", "https://localhost:9011/", nil)
	client := http.Client{}

	//nolint:bodyclose // no response body to close
	_, err := client.Do(req)
	if err == nil {
		t.Errorf("server started even certificate are wrong")
	}
}
