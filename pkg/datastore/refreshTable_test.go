package datastore

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

const ycqlKeyspaceTest = "test"

type MockTesting struct {
	TotalErrors int
}

func (mt *MockTesting) Error(args ...interface{}) {
	mt.TotalErrors++
}

func (mt *MockTesting) Errorf(format string, args ...interface{}) {
	mt.TotalErrors++
}

func createYCQLKeyspace(ycql *YCQL) error {
	err := ycql.Session.Query(
		"CREATE KEYSPACE IF NOT EXISTS test  WITH REPLICATION = " +
			"{'class': 'SimpleStrategy', 'replication_factor': '1'} AND DURABLE_WRITES = true; ").Exec()

	return err
}

func initializeDB(t *testing.T, c *config.MockConfig) *Seeder {
	dc := DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	}

	orm, err := NewORM(&dc)
	if err != nil {
		t.Error(err)
	}

	store := new(DataStore)
	store.gorm = orm

	err = createTestTable(store)
	if err != nil {
		t.Error(err)
	}

	path, err := os.Getwd()
	if err != nil {
		t.Log(err)
	}

	return NewSeeder(store, path)
}

// createTestTable The function creates a table which will be utilized for testing the refresh table functions
func createTestTable(d *DataStore) error {
	_, err := d.DB().Exec("DROP TABLE IF EXISTS store")
	if err != nil {
		return err
	}

	var createTableQuery string

	switch d.GORM().Dialect().GetName() {
	case msSQL:
		createTableMSSQL(d)
	case pgSQL:
		createTableQuery = "CREATE TABLE store(id SERIAL PRIMARY KEY, name varchar(20))"
	default:
		createTableQuery = "CREATE TABLE store(id int NOT NULL PRIMARY KEY AUTO_INCREMENT, name varchar(20))"
	}

	_, err = d.DB().Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func createCassandraTestTable(d *DataStore) error {
	err := d.Cassandra.Session.Query("DROP TABLE IF EXISTS store").Exec()
	if err != nil {
		return err
	}

	return d.Cassandra.Session.Query(`CREATE TABLE store(id int, "name" varchar, PRIMARY KEY (id))`).Exec()
}

func createTableMSSQL(d *DataStore) {
	// since IDENTITY_INSERT cannot be done on master db, creating another db for tests.
	_, _ = d.DB().Exec("CREATE DATABASE tests")
	_, _ = d.DB().Exec("USE tests")
	_, _ = d.DB().Exec("DROP TABLE IF EXISTS student")
	// table with UUID
	_, _ = d.DB().Exec("CREATE TABLE student(id uniqueidentifier not null , name varchar(20))")
	// table with Primary Key
	_, _ = d.DB().Exec("CREATE TABLE store(id int not null PRIMARY KEY, name varchar(20))")
	_, _ = d.DB().Exec("DROP TABLE IF EXISTS employee")
	// table with Identity Column
	_, _ = d.DB().Exec("CREATE TABLE employee(id int PRIMARY KEY IDENTITY, name varchar(20))")
}

func Test_RefreshTablesAndVersionCheck(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")
	msSQLConf := DBConfig{c.Get("MSSQL_HOST"), c.Get("MSSQL_USER"),
		c.Get("MSSQL_PASSWORD"), c.Get("MSSQL_DB_NAME"),
		c.Get("MSSQL_PORT"), "mssql", "disable", "", "", "", 30}

	testcases := []struct {
		tableName string
		version   string
		config    DBConfig
		rows      int
	}{
		{"store", "wrong",
			DBConfig{os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
				"mysql", "", "", "", "", 30}, 5,
		},
		{"store", "incorrect", msSQLConf, 5},
		{"student", "incorrect", msSQLConf, 2},
		{"employee", "incorrect", msSQLConf, 3},
		{"store", "incorrect",
			DBConfig{os.Getenv("PGSQL_HOST"), os.Getenv("PGSQL_USER"),
				os.Getenv("PGSQL_PASSWORD"), os.Getenv("PGSQL_DB_NAME"), os.Getenv("PGSQL_PORT"),
				"postgres", "", "", "", "", 30}, 5,
		},
	}
	for _, testcase := range testcases {
		mock1 := &MockTesting{}
		mock2 := &MockTesting{}
		cgf := config.MockConfig{Data: map[string]string{
			"DB_HOST":     testcase.config.HostName,
			"DB_USER":     testcase.config.Username,
			"DB_PASSWORD": testcase.config.Password,
			"DB_NAME":     testcase.config.Database,
			"DB_PORT":     testcase.config.Port,
			"DB_DIALECT":  testcase.config.Dialect,
		}}

		db := initializeDB(t, &cgf)
		db.RefreshTables(mock1, testcase.tableName)

		testError := ""

		if mock1.TotalErrors != 0 {
			testError = "Refresh table failed."
		}

		// testing incorrect version hence errors should be greater than zero
		db.AssertVersion(mock2, testcase.version)

		if mock2.TotalErrors == 0 {
			testError += "Assert version failed"
		}

		db.AssertRowCount(t, testcase.tableName, testcase.rows)

		if testError != "" {
			t.Error(testError)
		}
	}
}

func Test_FileOpen(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../configs")
	cfg := config.MockConfig{Data: map[string]string{
		"DB_HOST":     c.Get("DB_HOST"),
		"DB_USER":     c.Get("DB_USER"),
		"DB_PASSWORD": c.Get("DB_PASSWORD"),
		"DB_NAME":     c.Get("DB_NAME"),
		"DB_PORT":     c.Get("DB_PORT"),
		"DB_DIALECT":  c.Get("DB_DIALECT"),
	},
	}
	db := initializeDB(t, &cfg)
	db.path = "wrong/path"
	mock := &MockTesting{}
	tableName := "store"
	db.RefreshTables(mock, tableName)

	if mock.TotalErrors != 1 {
		t.Error("An error should have been reported while opening the file location")
	}
}

func Test_CSVReader(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../configs")
	cfg := config.MockConfig{Data: map[string]string{
		"DB_HOST":     c.Get("DB_HOST"),
		"DB_USER":     c.Get("DB_USER"),
		"DB_PASSWORD": c.Get("DB_PASSWORD"),
		"DB_NAME":     c.Get("DB_NAME"),
		"DB_PORT":     c.Get("DB_PORT"),
		"DB_DIALECT":  c.Get("DB_DIALECT"),
	},
	}
	db := initializeDB(t, &cfg)
	goPath := c.Get("GOPATH")
	db.path = goPath + "/src/developer.zopsmart.com/go/gofr/pkg/datastore"
	mock := &MockTesting{}
	tableName := "random"
	db.RefreshTables(mock, tableName)

	if mock.TotalErrors != 2 {
		t.Error("errors should have been reported while reading the csv file")
	}
}

func Test_ClearTable(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../configs")
	cfg := config.MockConfig{Data: map[string]string{
		"DB_HOST":     c.Get("DB_HOST"),
		"DB_USER":     c.Get("DB_USER"),
		"DB_PASSWORD": c.Get("DB_PASSWORD"),
		"DB_NAME":     c.Get("DB_NAME"),
		"DB_PORT":     c.Get("DB_PORT"),
		"DB_DIALECT":  c.Get("DB_DIALECT"),
	},
	}
	db := initializeDB(t, &cfg)
	mock := &MockTesting{}
	tableName := "?SELECT * FROM store"
	db.ClearTable(mock, tableName)

	if mock.TotalErrors != 1 {
		t.Error("An error should have been reported while clearing the table")
	}
}

func Test_AssertRowCount(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../configs")
	cfg := config.MockConfig{Data: map[string]string{
		"DB_HOST":     c.Get("DB_HOST"),
		"DB_USER":     c.Get("DB_USER"),
		"DB_PASSWORD": c.Get("DB_PASSWORD"),
		"DB_NAME":     c.Get("DB_NAME"),
		"DB_PORT":     c.Get("DB_PORT"),
		"DB_DIALECT":  c.Get("DB_DIALECT"),
	},
	}
	db := initializeDB(t, &cfg)
	mock := &MockTesting{}
	db.AssertRowCount(mock, "store", 4)

	// due to incorrect number of rows.
	if mock.TotalErrors != 1 {
		t.Errorf("Error should have been reported while counting rows")
	}

	db.AssertRowCount(mock, "random", 4)

	if mock.TotalErrors == 1 {
		t.Errorf("Error should have been reported while counting rows")
	}
}

func TestFailedQuery(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../configs")
	cfg := DBConfig{HostName: c.Get("DB_HOST"), Port: c.Get("DB_PORT"), Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"), Database: c.Get("DB_NAME"), Dialect: c.Get("DB_DIALECT")}

	db, _ := NewORM(&cfg)
	d := DataStore{ORM: db}
	path, _ := os.Getwd()
	s := NewSeeder(&d, path)

	records, _ := s.getRecords("store")
	mock := &MockTesting{}
	s.populateTable(mock, "unknown", records)

	if mock.TotalErrors != 1 {
		t.Errorf("Got %v error(s)\tExpected 1 error", mock.TotalErrors)
	}
}

// TestForeignKey relies on the table that was created previously, store and employee.
// This checks if test data seeding can support tables with foreign keys
func TestForeignKey(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../configs")
	msSQLConf := DBConfig{c.Get("MSSQL_HOST"), c.Get("MSSQL_USER"),
		c.Get("MSSQL_PASSWORD"), c.Get("MSSQL_DB_NAME"),
		c.Get("MSSQL_PORT"), "mssql", "", "", "", "", 30}

	db, _ := NewORM(&msSQLConf)

	_ = db.Exec("DROP TABLE IF EXISTS store_employee")
	_ = db.Exec("CREATE table store_employee (store_id int, employee_id int, CONSTRAINT store_employee_1 FOREIGN KEY " +
		"(store_id) REFERENCES store	(id), CONSTRAINT store_employee_2 FOREIGN KEY (employee_id) REFERENCES employee	(id) )")

	defer func() {
		_ = db.Exec("DROP TABLE IF EXISTS store_employee")
	}()

	d := DataStore{ORM: db}

	path, _ := os.Getwd()
	s := NewSeeder(&d, path)

	tester := &MockTesting{}
	s.RefreshTables(tester, "store_employee", "store", "employee")

	s.AssertRowCount(tester, "store_employee", 3)

	// expecting 0 errors
	if tester.TotalErrors != 0 {
		t.Error("Test Data seeding failed")
	}
}

func TestSeeder_RefreshMongoCollections(t *testing.T) {
	mongo, err := GetMongoDBFromEnv(log.NewLogger())
	if err != nil {
		t.Error(err)
	}

	d := DataStore{MongoDB: mongo}
	path, _ := os.Getwd()
	s := NewSeeder(&d, path+"/test_data")

	tester := &MockTesting{}
	s.RefreshMongoCollections(tester, "customers")
	// expecting 0 errors
	if tester.TotalErrors != 0 {
		t.Errorf("Test Data seeding failed Expecting 0 errors but got: %v errors", tester.TotalErrors)
	}
}

func TestSeeder_RefreshMongoCollections_Errors(t *testing.T) {
	mongo, err := GetMongoDBFromEnv(log.NewLogger())
	if err != nil {
		t.Error(err)
	}

	d := DataStore{MongoDB: mongo}
	s := NewSeeder(&d, "some/invalid/path")

	tester := &MockTesting{}
	s.RefreshMongoCollections(tester, "customers")

	path, _ := os.Getwd()
	s = NewSeeder(&d, path+"/test_data")
	s.RefreshMongoCollections(tester, "unknown")

	// checking invalid json data
	s = NewSeeder(&d, path+"/test_data")
	s.RefreshMongoCollections(tester, "customers1")

	// expecting 3 errors
	expectedErrors := 3
	if tester.TotalErrors != expectedErrors {
		t.Errorf("Test Data seeding failed Expecting %v errors but got: %v errors", expectedErrors, tester.TotalErrors)
	}
}

func TestSeeder_RefreshCassandra(t *testing.T) {
	logger := log.NewLogger()
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	cassandraPort, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandraCfg := CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     cassandraPort,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "test",
	}

	cassDB, err := GetNewCassandra(logger, &cassandraCfg)
	if err != nil {
		t.Error(err)
	}

	d := DataStore{Cassandra: cassDB}

	err = createCassandraTestTable(&d)
	if err != nil {
		t.Error(err)
	}

	tester := &MockTesting{}
	path, _ := os.Getwd()
	s := NewSeeder(&d, path)
	s.RefreshCassandra(tester, "store")
	// expecting 0 errors
	expectedErrors := 0
	if tester.TotalErrors != expectedErrors {
		t.Errorf("Test Data seeding failed Expecting %v errors but got: %v errors", expectedErrors, tester.TotalErrors)
	}
}

func intialiseYCQL() *CassandraCfg {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../configs")

	port, err := strconv.Atoi(c.Get("YCQL_DB_PORT"))
	if err != nil || port == 0 {
		port = 9042
	}

	ycqlCfg := CassandraCfg{
		Hosts:       c.Get("CASS_DB_HOST"),
		Port:        port,
		Consistency: LocalQuorum,
		Username:    c.Get("YCQL_DB_USER"),
		Password:    c.Get("YCQL_DB_PASS"),
		Keyspace:    c.Get("CASS_DB_KEYSPACE"),
		Timeout:     600,
	}

	return &ycqlCfg
}

func TestSeeder_RefreshYCQL(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	ycqlCfg := intialiseYCQL()

	ycql, err := GetNewYCQL(logger, ycqlCfg)
	if err != nil {
		t.Error(err)
	}

	err = createYCQLKeyspace(&ycql)
	if err != nil {
		t.Error(err)
	}

	// We want to connect with keyspace  test instead of system
	ycqlCfg.Keyspace = ycqlKeyspaceTest

	ycql, err = GetNewYCQL(logger, ycqlCfg)
	if err != nil {
		t.Error(err)
	}

	tester := &MockTesting{}
	path, _ := os.Getwd()
	s := NewSeeder(&DataStore{YCQL: ycql}, path)

	expectedErrors := 0
	testCases := []struct {
		query string
	}{
		{
			query: `CREATE TABLE store (id varchar  PRIMARY KEY, name varchar ) WITH transactions = { 'enabled' : true }  ;`,
		},
		{
			query: `CREATE TABLE store (id int     PRIMARY KEY, name varchar ) WITH transactions = { 'enabled' : true }  ;`,
		},
		{
			query: `CREATE TABLE store (id float      PRIMARY KEY, name varchar ) WITH transactions = { 'enabled' : true }  ;`,
		},
	}

	for i, tc := range testCases {
		drop := `DROP TABLE IF EXISTS store`

		// Dropping the table if exist, different field type can use at same time.
		err = ycql.Session.Query(drop).Exec()
		if err != nil {
			t.Errorf("[FAILED]%v , Test Data seeding failed in droping table store, error : %v ", i, err)
		}

		err = ycql.Session.Query(tc.query).Exec()
		if err != nil {
			t.Errorf("[FAILED]%v , Test Data seeding failed due to store table creation, error:%v", i, err)
		}

		s.RefreshYCQL(tester, "store")

		if tester.TotalErrors != 0 {
			t.Errorf("[FAILED]%v , Test Data seeding failed Expecting %v errors but got: %v errors", i, expectedErrors, tester.TotalErrors)
		}
	}
}

func TestSeeder_RefreshYCQL_Table_Error(t *testing.T) {
	ycqlCfg := intialiseYCQL()
	logger := log.NewLogger()

	// We want to connect with keyspace with test instead of system
	ycqlCfg.Keyspace = ycqlKeyspaceTest

	ycql, err := GetNewYCQL(logger, ycqlCfg)
	if err != nil {
		t.Error(err)
	}

	tester := &MockTesting{}
	path, _ := os.Getwd()
	s := NewSeeder(&DataStore{YCQL: ycql}, path)
	s.RefreshYCQL(tester, "random")

	s = NewSeeder(&DataStore{YCQL: ycql}, "useless") // checking for invalid path
	s.RefreshYCQL(tester, "store")
	// expecting 2 errors one due to different store and one with different directoryPath
	expectedErrors := 2
	if tester.TotalErrors != expectedErrors {
		t.Errorf("Test Data seeding failed Expecting %v errors but got: %v errors", expectedErrors, tester.TotalErrors)
	}
}

func TestSeeder_RefreshCassandra_Error(t *testing.T) {
	logger := log.NewMockLogger(io.Discard)
	c := config.NewGoDotEnvProvider(logger, "../../configs")
	cassandraPort, _ := strconv.Atoi(c.Get("CASS_DB_PORT"))
	cassandraCfg := CassandraCfg{
		Hosts:    c.Get("CASS_DB_HOST"),
		Port:     cassandraPort,
		Username: c.Get("CASS_DB_USER"),
		Password: c.Get("CASS_DB_PASS"),
		Keyspace: "test",
	}

	cassDB, err := GetNewCassandra(logger, &cassandraCfg)
	if err != nil {
		t.Error(err)
	}

	d := DataStore{Cassandra: cassDB}

	err = createCassandraTestTable(&d)
	if err != nil {
		t.Error(err)
	}

	tester := &MockTesting{}
	s := NewSeeder(&d, "/some/invalid/path")
	s.RefreshCassandra(tester, "store")

	path, _ := os.Getwd()
	s = NewSeeder(&d, path)
	s.RefreshCassandra(tester, "unknown")
	// expecting 2 errors
	expectedErrors := 2
	if tester.TotalErrors != expectedErrors {
		t.Errorf("Test Data seeding failed Expecting %v errors but got: %v errors", expectedErrors, tester.TotalErrors)
	}
}

func TestSeeder_RefreshRedis(t *testing.T) {
	redisDB, err := NewRedisFromEnv(nil)
	if err != nil {
		t.Error(err)
	}

	d := DataStore{Redis: redisDB}
	path, _ := os.Getwd()
	s := NewSeeder(&d, path+"/test_data")
	tester := &MockTesting{}
	s.RefreshRedis(tester, "storeKeyVal", "customers")

	expectedErrors := 0
	if tester.TotalErrors != expectedErrors {
		t.Errorf("Test Data seeding failed Expecting %v errors but got: %v errors", expectedErrors, tester.TotalErrors)
	}
}

func TestSeeder_RefreshRedis_Error(t *testing.T) {
	redisDB, err := NewRedisFromEnv(nil)
	if err != nil {
		t.Error(err)
	}

	d := DataStore{Redis: redisDB}
	// case where invalid path is provided
	s := NewSeeder(&d, "/some/invalid/path")
	tester := &MockTesting{}
	s.RefreshRedis(tester, "storeKeyVal1")

	// case where invalid data is provided in csv
	path, _ := os.Getwd()
	s = NewSeeder(&d, path+"/test_data")
	s.RefreshRedis(tester, "storeKeyVal1")

	expectedErrors := 2
	if tester.TotalErrors != expectedErrors {
		t.Errorf("Test Data seeding failed Expecting %v errors but got: %v errors", expectedErrors, tester.TotalErrors)
	}
}

func TestSeeder_resetIdentitySequence(t *testing.T) {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(new(bytes.Buffer)), "../../configs")
	testcases := []struct {
		dbConfig          DBConfig
		tableCreateQuery  string
		beforeTransaction bool
	}{
		{DBConfig{
			HostName: c.Get("PGSQL_HOST"),
			Username: c.Get("PGSQL_USER"),
			Password: c.Get("PGSQL_PASSWORD"),
			Database: c.Get("PGSQL_DB_NAME"),
			Port:     c.Get("PGSQL_PORT"),
			Dialect:  pgSQL,
		}, `
 	   CREATE TABLE IF NOT EXISTS dummy (
	   id serial primary key,
	   name varchar (50))
	`, false},
		{DBConfig{
			HostName: c.Get("MSSQL_HOST"),
			Username: c.Get("MSSQL_USER"),
			Password: c.Get("MSSQL_PASSWORD"),
			Database: c.Get("MSSQL_DB_NAME"),
			Port:     c.Get("MSSQL_PORT"),
			Dialect:  msSQL,
		}, `IF NOT EXISTS
	(  SELECT [name]
		FROM sys.tables
      WHERE [name] = 'dummy'
   ) CREATE TABLE dummy (id int primary key identity(1,1),
	   name varchar (50))
	`, true},
		{
			DBConfig{
				HostName: c.Get("DB_HOST"),
				Username: c.Get("DB_USER"),
				Password: c.Get("DB_PASSWORD"),
				Database: c.Get("DB_NAME"),
				Port:     c.Get("DB_PORT"),
				Dialect:  mySQL,
			}, `
 	   CREATE TABLE IF NOT EXISTS dummy (
	   id int auto_increment,
	   name varchar (50),
	   PRIMARY KEY (id))
	`, true},
	}

	for _, tc := range testcases {
		orm, err := NewORM(&tc.dbConfig)
		if err != nil {
			t.Error(err)
		}

		ds := DataStore{ORM: orm}
		s := NewSeeder(&ds, "/dummy")
		s.ResetCounter = true

		tester := &MockTesting{}

		if _, err := ds.DB().Exec(tc.tableCreateQuery); err != nil {
			t.Errorf("got error sourcing the schema: %v", err)
		}

		if _, err := ds.DB().Exec(`INSERT INTO dummy(name) VALUES('Alice')`); err != nil {
			t.Errorf("unable to insert dummy data: %v", err)
		}

		s.resetIdentitySequence(tester, "dummy", tc.beforeTransaction)

		if tester.TotalErrors != 0 {
			t.Errorf("reset identity sequence failed. Expecting 0 errors but got %v errors", tester.TotalErrors)
		}
	}
}
