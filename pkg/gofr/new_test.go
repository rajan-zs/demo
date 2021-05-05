package gofr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zopsmart/gofr/pkg/datastore/pubsub"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub/kafka"
	"github.com/zopsmart/gofr/pkg/gofr/config"
	"github.com/zopsmart/gofr/pkg/log"
)

type mockConfig struct {
	testCase string
}

// nolint
func (m mockConfig) Get(key string) string {
	switch m.testCase {
	case "redis error", "db error":
		return "mock"
	case "kafka error":
		if key == "KAFKA_HOSTS" {
			return ""
		}
	case "kafka":
		if key == "PUBSUB_BACKEND" {
			return "KAFKA"
		}
	case "avro":
		if key == "AVRO_SCHEMA_URL" {
			return "http://localhost:8081"
		}

		return ""
	case "avroerr":
		return ""
	default:
		c := &config.GoDotEnvProvider{}
		return c.Get(key)
	}

	return "mock"
}

func (m mockConfig) GetOrDefault(key, d string) string {
	return d
}

func Test_initializeRedis(t *testing.T) {
	tcs := []struct {
		c      Config
		expStr string // expected in the logs, logged by k.Logger
	}{
		{mockConfig{testCase: "redis error"}, "could not connect to Redis"},
		{mockConfig{}, "Redis connected"},
	}

	for _, tc := range tcs {
		k := New()
		b := new(bytes.Buffer)

		k.Logger = log.NewMockLogger(b)
		initializeRedis(tc.c, k)

		if !strings.Contains(b.String(), tc.expStr) {
			t.Errorf("FAILED, expected: `%v` in the logs, got: %v", tc.expStr, b.String())
		}
	}
}

func Test_initializeDB(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	c := config.NewGoDotEnvProvider(logger, "../../configs")

	hostName := c.Get("DB_HOST")
	port := c.Get("DB_PORT")

	testcases := []struct {
		host        string
		port        string
		ORM         string
		expectedLog string
	}{
		{"", "", "", ""},
		{"incorrect-url", "7", "", "could not connect to DB"},
		{hostName, port, "", "DB connected, HostName: " + hostName + ", Port: " + port},
		{"incorrect-url", "7", "SQLX", "could not connect to DB"},
		{hostName, port, "SQLX", "DB connected, HostName: " + hostName + ", Port: " + port},
	}

	for i, tc := range testcases {
		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)

		mockConfig := config.MockConfig{
			Data: map[string]string{"DB_HOST": tc.host, "DB_USER": c.Get("DB_USER"), "DB_PASSWORD": c.Get("DB_PASSWORD"),
				"DB_NAME": c.Get("DB_NAME"), "DB_PORT": tc.port, "DB_DIALECT": c.Get("DB_DIALECT"), "DB_ORM": tc.ORM},
		}

		k := NewWithConfig(&mockConfig)
		k.Logger = logger

		initializeDB(&mockConfig, k)

		if !strings.Contains(b.String(), tc.expectedLog) {
			t.Errorf("[TESTCASE %d] Failed. Got: %v\tExpected: %v\n", i+1, b.String(), tc.expectedLog)
		}
	}
}

func Test_InitializeElasticsearch(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(new(bytes.Buffer)), "../../configs")
	host := c.Get("ELASTIC_SEARCH_HOST")
	port := c.Get("ELASTIC_SEARCH_PORT")
	user := c.Get("ELASTIC_SEARCH_USER")
	pass := c.Get("ELASTIC_SEARCH_PASS")

	testcases := []struct {
		host        string
		port        string
		user        string
		pass        string
		expectedLog string
	}{
		{"", "", "", "", ""},
		{"incorrect-url", "7", "", "", "could not connect to Elasticsearch"},
		{host, port, user, pass, "connected to Elasticsearch, HOST: " + host + ", PORT: " + port},
	}

	for i, v := range testcases {
		b := new(bytes.Buffer)
		logger := log.NewMockLogger(b)

		var conf = config.MockConfig{
			Data: map[string]string{"ELASTIC_SEARCH_HOST": v.host, "ELASTIC_SEARCH_PORT": v.port,
				"ELASTIC_SEARCH_USER": v.user, "ELASTIC_SEARCH_PASS": v.pass},
		}

		k := NewWithConfig(&conf)
		k.Logger = logger

		initializeElasticsearch(&conf, k)

		if !strings.Contains(b.String(), v.expectedLog) {
			t.Errorf("[TESTCASE%d]Failed. Got %v\tExpected %v\n", i+1, b.String(), v.expectedLog)
		}
	}
}

func Test_initializeMongoDB(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	config.NewGoDotEnvProvider(logger, "../../configs")

	tcs := []struct {
		configLoc   Config
		expectedStr string
	}{
		{mockConfig{}, "Mongo DB connected."},
		{configLoc: &config.MockConfig{Data: map[string]string{"MONGO_DB_HOST": "fakehost", "MONGO_DB_USER": "admin",
			"MONGO_DB_PASS": "admin123",
			"MONGO_DB_PORT": "27017"}}, expectedStr: "error while pinging to Mongo DB"},
	}

	for _, tc := range tcs {
		k := New()
		b := new(bytes.Buffer)

		k.Logger = log.NewMockLogger(b)
		initializeMongoDB(tc.configLoc, k)

		if !strings.Contains(b.String(), tc.expectedStr) {
			t.Errorf("FAILED, expected: `%v` in the logs, got: %v", tc.expectedStr, b.String())
		}
	}
}

func Test_initializeCassandra(t *testing.T) {
	// this is done to so that it doesnt affects the other tests related to cassandra
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	config.NewGoDotEnvProvider(logger, "../../configs")

	tcs := []struct {
		configLoc   Config
		expectedStr string
	}{
		{mockConfig{}, "Connected to cassandra"},
		{configLoc: &config.MockConfig{Data: map[string]string{"CASS_DB_HOST": "cassandra", "CASS_DB_PORT": "2003"}}},
	}

	for _, tc := range tcs {
		k := &Gofr{}
		k.Logger = logger

		initializeCassandra(mockConfig{}, k)

		if !strings.Contains(b.String(), tc.expectedStr) {
			t.Errorf("FAILED, expected: `%v` in the logs, got: %v", tc.expectedStr, b.String())
		}
	}
}

func Test_initializeCassandra_InvalidDialect(t *testing.T) {
	c := config.MockConfig{Data: map[string]string{"CASS_DB_DIALECT": "invalid", "CASS_DB_HOST": "localhost", "CASS_DB_PORT": "20112"}}

	expectedLog := "invalid dialect"
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)

	k := &Gofr{}
	k.Logger = logger

	initializeCassandra(&c, k)

	if !strings.Contains(b.String(), expectedLog) {
		t.Errorf("FAILED, expected: `%v` in the logs, got: %v", expectedLog, b.String())
	}
}

func Test_getYcqlConfigs(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	cfg := &config.MockConfig{
		Data: map[string]string{
			"CASS_DB_DIALECT":  "YCQL",
			"CASS_DB_PASS":     c.Get("YCQL_DB_PASS"),
			"CASS_DB_USER":     c.Get("YCQL_DB_USER"),
			"CASS_DB_PORT":     c.Get("YCQL_DB_PORT"),
			"CASS_DB_KEYSPACE": c.Get("CASS_DB_KEYSPACE"),
			"CASS_DB_TIMEOUT":  c.Get("CASS_DB_TIMEOUT"),
		},
	}

	testCases := []struct {
		host        string
		expectedStr string
	}{
		{"localhost", "Connected to YCQL"},
		{"invalidhost", "error while connecting to YCQL"},
	}

	for i, tc := range testCases {
		b.Reset()

		k := &Gofr{}
		k.Logger = logger

		cfg.Data["CASS_DB_HOST"] = tc.host

		initializeCassandra(cfg, k)

		if !strings.Contains(b.String(), tc.expectedStr) {
			t.Errorf("FAILED case`%v`, expected: `%v` in the logs, got: %v", i, tc.expectedStr, b.String())
		}
	}
}

func Test_PubSub(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	conf := config.NewGoDotEnvProvider(logger, "../../configs")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := map[string]interface{}{
			"subject": "gofr-value",
			"version": 3,
			"id":      303,
			"schema": "{\"type\":\"record\",\"name\":\"person\"," +
				"\"fields\":[{\"name\":\"Id\",\"type\":\"string\"}," +
				"{\"name\":\"Name\",\"type\":\"string\"}," +
				"{\"name\":\"Email\",\"type\":\"string\"}]}",
		}

		reBytes, _ := json.Marshal(re)
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write(reBytes)
	}))

	k := &Gofr{Logger: logger}

	testCases := []struct {
		configLoc   Config
		expectedStr string
	}{
		{mockConfig{}, "Kafka initialized"},
		{&config.MockConfig{Data: map[string]string{
			"EVENTHUB_NAMESPACE":  "zsmisc-dev",
			"EVENTHUB_NAME":       "healthcheck",
			"AZURE_CLIENT_ID":     conf.Get("AZURE_CLIENT_ID"),
			"AZURE_CLIENT_SECRET": conf.Get("AZURE_CLIENT_SECRET"),
			"AZURE_TENANT_ID":     conf.Get("AZURE_TENANT_ID"),
			"PUBSUB_BACKEND":      "EVENTHUB",
		}}, "Azure Eventhub initialized"},
		{&config.MockConfig{Data: map[string]string{
			"EVENTHUB_NAMESPACE":  "zsmisc-dev",
			"EVENTHUB_NAME":       "healthcheck",
			"AZURE_CLIENT_ID":     conf.Get("AZURE_CLIENT_ID"),
			"AZURE_CLIENT_SECRET": conf.Get("AZURE_CLIENT_SECRET"),
			"AZURE_TENANT_ID":     conf.Get("AZURE_TENANT_ID"),
			"PUBSUB_BACKEND":      "EVENTHUB",
			"AVRO_SCHEMA_URL":     ts.URL,
		}}, "Avro initialized"},
	}

	for i, tc := range testCases {
		b.Reset()
		initializePubSub(tc.configLoc, k)

		if !strings.Contains(b.String(), tc.expectedStr) {
			t.Errorf("[FAILED %v], expected: `%v` in the logs, got: %v", i, tc.expectedStr, b.String())
		}
	}
}

func Test_initializeAvro(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := map[string]interface{}{
			"subject": "gofr-value",
			"version": 3,
			"id":      303,
			"schema": "{\"type\":\"record\",\"name\":\"person\"," +
				"\"fields\":[{\"name\":\"Id\",\"type\":\"string\"}," +
				"{\"name\":\"Name\",\"type\":\"string\"}," +
				"{\"name\":\"Email\",\"type\":\"string\"}]}",
		}

		reBytes, _ := json.Marshal(re)
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write(reBytes)
	}))

	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	k := &Gofr{Logger: logger}
	topic := c.Get("KAFKA_TOPIC") // CSV string
	topics := strings.Split(topic, ",")
	kafkaCfg := &kafka.Config{
		Brokers: c.Get("KAFKA_HOSTS"),
		Topics:  topics,
	}
	kafkaObj, _ := kafka.New(kafkaCfg, logger)
	tests := []struct {
		c           Config
		ps          pubsub.PublisherSubscriber
		expectedStr string
	}{
		{&config.MockConfig{Data: map[string]string{"AVRO_SCHEMA_URL": ts.URL, "AVRO_SUBJECT": "gofr-value"}},
			kafkaObj, "Avro initialized!"},
		{&config.MockConfig{Data: map[string]string{"AVRO_SCHEMA_URL": ts.URL, "AVRO_SUBJECT": "gofr-value"}},
			nil, "Kafka/Eventhub not present, cannot use Avro"},
		{&config.MockConfig{Data: map[string]string{"AVRO_SCHEMA_URL": "", "AVRO_SUBJECT": "gofr-value"}},
			kafkaObj, "Schema registry URL is required for Avro"},
	}

	for _, tt := range tests {
		k.PubSub = tt.ps
		avroConfig := avroConfigFromEnv(tt.c)
		initializeAvro(avroConfig, k)

		if !strings.Contains(b.String(), tt.expectedStr) {
			t.Errorf("FAILED, expected: `%v` in the logs, got: %v", tt.expectedStr, b.String())
		}
	}
}

func Test_initializeSolr(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	testCases := []struct {
		configLoc   config.MockConfig
		expectedStr string
	}{
		{
			config.MockConfig{Data: map[string]string{
				"SOLR_HOST": c.Get("SOLR_HOST"),
				"SOLR_PORT": c.Get("SOLR_PORT"),
			}},
			"Solr connected",
		},
		{
			config.MockConfig{Data: map[string]string{
				"SOLR_HOST": "",
				"SOLR_PORT": "",
			}},
			"",
		},
	}

	k := &Gofr{Logger: logger}

	for _, tc := range testCases {
		initializeSolr(&tc.configLoc, k)

		if !strings.Contains(b.String(), tc.expectedStr) {
			t.Errorf("FAILED, expected: `%v` in the logs, got: %v", tc.expectedStr, b.String())
		}

		b = new(bytes.Buffer)
	}
}

func Test_GofrCMDConfig(t *testing.T) {
	k := NewCMD()
	if k.Redis == nil {
		t.Errorf("expected redis to be connected through configs")
	}
}

// This function verifies that metric server is starting successfully or not, when NewCMD is called.
func Test_GofrCMDMetricServer(t *testing.T) {
	b := new(bytes.Buffer)
	logger := log.NewMockLogger(b)
	k := NewCMD()
	k.cmd.metricSvr.port = 12221
	k.Logger = logger
	k.Start()

	if k.cmd.metricSvr.server == nil || strings.Contains(b.String(), "error in metrics server") {
		t.Errorf("failed to start the metric server")
	}
}