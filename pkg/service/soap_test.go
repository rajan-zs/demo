package service

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/log"
)

// It is a test server that behaves like SOAP API. Based on the different SOAP actions, it returns the different desired responses.
func testSOAPServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var res string
		switch r.Header.Get("SOAPAction") {
		case "gfr":
			res = `<?xml version="1.0" encoding="utf-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
							<soap:Body>
        						<m:CompanyNameResponse xmlns:m="http://www.zopsmart.com">
           						 	<m:CompanyNameResult>Gofr</m:CompanyNameResult>
								</m:CompanyNameResponse>
   							 </soap:Body>
						</soap:Envelope>`
		case "zop":
			res = `<?xml version="1.0" encoding="utf-8"?>
						<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
							<soap:Body>
        						<m:CompanyNameResponse xmlns:m="http://www.zopsmart.com">
           						 	<m:CompanyNameResult>ZopSmart</m:CompanyNameResult>
								</m:CompanyNameResponse>
   							 </soap:Body>
						</soap:Envelope>`
		}

		resBytes, _ := xml.Marshal(res)
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write(resBytes)
	}))

	return ts
}

// TestCallWithHeaders_SOAP tests the SOAP client with a test soap server
func TestSOAPServer(t *testing.T) {
	ts := testSOAPServer()
	defer ts.Close()

	tests := []struct {
		name string
		// Input
		action string
		// Output
		out string
	}{
		{"action/gfr", "gfr", "Gofr"},
		{"action/zop", "zop", "ZopSmart"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ps := NewSOAPClient(ts.URL, log.NewMockLogger(io.Discard), "", "")
			res, err := ps.Call(context.Background(), tc.action, nil)

			if !strings.Contains(string(res.Body), tc.out) {
				t.Errorf("Unexpected Response : %v", string(res.Body))
			}

			if err != nil {
				t.Errorf("Test %v:\t  error = %v", tc.name, err)
			}
		})
	}
}

// TestCallWithHeaders_SOAP tests the SOAP client by passing custom headers
func TestCallWithHeaders_SOAP(t *testing.T) {
	ts := testSOAPServer()
	defer ts.Close()

	tests := []struct {
		action      string
		headers     map[string]string
		expectedLog string
	}{
		{"zop", map[string]string{"X-Trace-Id": "a123ru", "X-B-Trace-Id": "198d7sf3d"}, `"X-B-Trace-Id":"198d7sf3d","X-Trace-Id":"a123ru"`},
		{"gfr", map[string]string{"X-Zopsmart-Tenant": "zopsmart"}, `"X-Zopsmart-Tenant":"zopsmart"`},
		{"gfr", nil, ``},
	}

	for i, tc := range tests {
		b := new(bytes.Buffer)

		soapClient := NewSOAPClient(ts.URL, log.NewMockLogger(b), "basic-user", "password")

		_, err := soapClient.CallWithHeaders(context.Background(), tc.action, nil, tc.headers)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if !strings.Contains(b.String(), tc.expectedLog) {
			t.Errorf("test id  %d headers is not logged", i+1)
		}
	}
}

func Test_Bind(t *testing.T) {
	ts := testSOAPServer()

	test := struct {
		desc   string
		resp   []byte
		result interface{}
	}{
		desc: "SOAP Bind successfully",
		resp: []byte("<note></note>"),
	}
	b := new(bytes.Buffer)
	soapClient := NewSOAPClient(ts.URL, log.NewMockLogger(b), "basic-user", "password")

	err := soapClient.Bind(test.resp, test.result)
	if err != nil {
		t.Errorf("test case us failed because of %v", err)
	}
}
func Test_BindStrict(t *testing.T) {
	ts := testSOAPServer()

	test := struct {
		desc   string
		resp   []byte
		result interface{}
	}{
		desc: "SOAP BindStrict successfully",
		resp: []byte("<note></note>"),
	}
	b := new(bytes.Buffer)
	soapClient := NewSOAPClient(ts.URL, log.NewMockLogger(b), "basic-user", "password")

	err := soapClient.BindStrict(test.resp, test.result)
	if err != nil {
		t.Errorf("test case us failed because of %v", err)
	}
}
