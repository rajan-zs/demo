package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/zopsmart/gofr/examples/universal-example/avro/handlers"
	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/log"
)

func TestMain(m *testing.M) {
	k := gofr.New()

	cassandraTableInitialization(k)

	postgresTableInitialization(k)

	// avro schema registry test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := map[string]interface{}{
			"subject": "employee-value",
			"version": 3,
			"id":      303,
			"schema": "{\"type\":\"record\",\"name\":\"employee\"," +
				"\"fields\":[{\"name\":\"Id\",\"type\":\"string\"}," +
				"{\"name\":\"Name\",\"type\":\"string\"}," +
				"{\"name\":\"Phone\",\"type\":\"string\"}," +
				"{\"name\":\"Email\",\"type\":\"string\"}," +
				"{\"name\":\"City\",\"type\":\"string\"}]}",
		}

		reBytes, _ := json.Marshal(re)
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write(reBytes)
	}))

	schemaURL := os.Getenv("AVRO_SCHEMA_URL")
	os.Setenv("AVRO_SCHEMA_URL", ts.URL)

	topic := os.Getenv("KAFKA_TOPIC")
	os.Setenv("KAFKA_TOPIC", "avro-pubsub")

	defer func() {
		os.Setenv("AVRO_SCHEMA_URL", schemaURL)
		os.Setenv("KAFKA_TOPIC", topic)
	}()

	os.Exit(m.Run())
}

func TestUniversalIntegration(t *testing.T) {
	// call the main function
	go main()
	// sleep, so that every data stores get initialized properly
	time.Sleep(5 * time.Second)

	testDataStores(t)

	testKafkaDataStore(t)

	testEventhub(t)
}

func testDataStores(t *testing.T) {
	testcases := []struct {
		testID             int
		method             string
		endpoint           string
		expectedStatusCode int
		body               []byte
	}{
		// Cassandra
		{1, "GET", "/cassandra/employee?name=Aman", 200, nil},
		{2, "POST", "/cassandra/employee", 201,
			[]byte(`{"id": 5, "name": "Sukanya", "phone": "01477", "email":"sukanya@zopsmart.com", "city":"Guwahati"}`)},
		{3, "GET", "/cassandra/unknown", 404, nil},
		// Redis
		{4, "GET", "/redis/config/key123", 500, nil},
		{5, "POST", "/redis/config", 201, []byte(`{}`)},
		// Postgres
		{6, "GET", "/pgsql/employee", 200, nil},
		{7, "POST", "/pgsql/employee", 201,
			[]byte(`{"id": 5, "name": "Sukanya", "phone": "01477", "email":"sukanya@zopsmart.com", "city":"Guwahati"}`)},
	}
	for _, tc := range testcases {
		req, _ := request.NewMock(tc.method, "http://localhost:9095"+tc.endpoint, bytes.NewBuffer(tc.body))
		cl := http.Client{}
		resp, err := cl.Do(req)

		if err != nil {
			t.Errorf("TestCase[%v] \t FAILED \nGot Error: %v", tc.testID, err)
			return
		}

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("Testcase[%v] Failed.\tExpected %v\tGot %v\n", tc.testID, tc.expectedStatusCode, resp.StatusCode)
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

// nolint:gocognit // can't break the function because of retry logic
func testKafkaDataStore(t *testing.T) {
	tcs := []struct {
		testID             int
		method             string
		endpoint           string
		expectedResponse   string
		expectedStatusCode int
	}{
		{8, "GET", "http://localhost:9095/avro/pub?id=1", "", 200},
		{9, "GET", "http://localhost:9095/avro/sub", "1", 200},
	}

	for _, tc := range tcs {
		req, _ := request.NewMock(tc.method, tc.endpoint, nil)
		c := http.Client{}

		for i := 0; i < 5; i++ {
			resp, _ := c.Do(req)

			if resp != nil && resp.StatusCode != tc.expectedStatusCode {
				// retry is required since, creation of topic takes time
				if checkRetry(resp.Body) {
					time.Sleep(3 * time.Second)
					continue
				}

				t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.testID, tc.expectedStatusCode, resp.StatusCode)

				return
			}

			// checks whether bind avro.Unmarshal functionality works fine
			if tc.expectedResponse != "" && resp.Body != nil {
				body, _ := io.ReadAll(resp.Body)

				m := struct {
					Data handlers.Employee `json:"data"`
				}{}
				_ = json.Unmarshal(body, &m)

				if m.Data.ID != tc.expectedResponse {
					t.Errorf("Expected: %v, Got: %v", tc.expectedResponse, m.Data.ID)
				}
			}

			if resp != nil {
				resp.Body.Close()
			}

			break
		}
	}
}

//nolint:gocognit // braking down the function will reduce the readability
func testEventhub(t *testing.T) {
	testcase := []struct {
		testID             int
		method             string
		endpoint           string
		expectedResponse   string
		expectedStatusCode int
	}{
		{10, "GET", "http://localhost:9095/eventhub/pub?id=1", "", 200},
		{11, "GET", "http://localhost:9095/eventhub/sub", "1", 200},
	}

	for _, tc := range testcase {
		req, _ := request.NewMock(tc.method, tc.endpoint, nil)
		c := http.Client{}
		resp, _ := c.Do(req)

		if resp != nil && resp.StatusCode != tc.expectedStatusCode {
			// required because eventhub is shared and there can be messages with avro or without avro
			// messages without avro would return 200 as we do json.Marshal to a map
			// messages with avro would return 206 as it would have to go through avro.Marshal
			// we can't use any avro schema as any schema can be used
			if resp.StatusCode != 206 {
				t.Errorf("Test %v: Failed.\tExpected %v\tGot %v\n", tc.testID, tc.expectedStatusCode, resp.StatusCode)
			}
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

// Cassandra Table initialization, Remove table if already exists
func cassandraTableInitialization(k *gofr.Gofr) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "configs")

	// Keyspace Creation for cassandra
	cassandraPort, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandraCfg := datastore.CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     cassandraPort,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "system",
	}

	cassDB, err := datastore.GetNewCassandra(logger, &cassandraCfg)
	if err == nil {
		err = cassDB.Session.Query("CREATE KEYSPACE test WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};").Exec()
	}

	for i := 0; err != nil && i < 10; i++ {
		time.Sleep(5 * time.Second)

		cassDB, err = datastore.GetNewCassandra(logger, &cassandraCfg)
		if err != nil {
			continue
		}

		err = cassDB.Session.Query(
			"CREATE KEYSPACE IF NOT EXISTS test WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};").Exec()
	}

	queryStr := "DROP TABLE IF EXISTS employees"
	if e := k.Cassandra.Session.Query(queryStr).Exec(); e != nil {
		k.Logger.Errorf("Got error while dropping the existing table employees: ", e)
	}

	queryStr = "CREATE TABLE IF NOT EXISTS employees (id int, name text, phone text, email text, city text, PRIMARY KEY (id) )"

	err = k.Cassandra.Session.Query(queryStr).Exec()
	if err != nil {
		k.Logger.Errorf("Failed creation of Table employees :%v", err)
	} else {
		k.Logger.Info("Table employees created Successfully")
	}
}

// Postgres Table initialization, Remove table if already exists
func postgresTableInitialization(k *gofr.Gofr) {
	if k.DB() == nil {
		return
	}

	query := `DROP TABLE IF EXISTS employees`
	if _, err := k.DB().Exec(query); err != nil {
		k.Logger.Errorf("Got error while dropping the existing table employees: ", err)
	}

	queryTable := `
 	   CREATE TABLE IF NOT EXISTS employees (
	   id         int primary key,
	   name       varchar (50),
 	   phone      varchar(50),
 	   email      varchar(50) ,
 	   city       varchar(50))
	`

	if _, err := k.DB().Exec(queryTable); err != nil {
		k.Logger.Errorf("Got error while sourcing the schema: ", err)
	}
}

func checkRetry(respBody io.Reader) bool {
	body, _ := io.ReadAll(respBody)

	errResp := struct {
		Errors []errors.Response `json:"errors"`
	}{}

	if len(errResp.Errors) == 0 {
		return false
	}

	_ = json.Unmarshal(body, &errResp)

	return strings.Contains(errResp.Errors[0].Reason, "Leader Not Available: the cluster is in the middle of a leadership election")
}
