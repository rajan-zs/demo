package datastore

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub"
	"github.com/zopsmart/gofr/pkg/gofr/types"
	"github.com/zopsmart/gofr/pkg/log"
)

type DataStore struct {
	rdb  SQLClient
	gorm GORMClient
	sqlx SQLXClient

	Logger        log.Logger
	MongoDB       MongoDB
	Redis         Redis
	ORM           interface{}
	Cassandra     Cassandra
	YCQL          YCQL
	PubSub        pubsub.PublisherSubscriber
	Solr          Client
	Elasticsearch Elasticsearch
}

type QueryLogger struct {
	Hosts     string     `json:"host,omitempty"`
	Query     []string   `json:"query"`
	Duration  int64      `json:"duration"`
	StartTime time.Time  `json:"-"`
	Logger    log.Logger `json:"-"`
	DataStore string     `json:"datastore"`
}

func (ds *DataStore) GORM() *gorm.DB {
	if ds.gorm.DB != nil {
		return ds.gorm.DB
	}

	if db, ok := ds.ORM.(GORMClient); ok {
		ds.gorm = db
		if db.DB != nil {
			ds.rdb.DB = db.DB.DB()
		}

		return db.DB
	}

	if gormDB, ok := ds.ORM.(*gorm.DB); ok {
		return gormDB
	}

	return nil
}

func (ds *DataStore) SQLX() *sqlx.DB {
	if ds.sqlx.DB != nil {
		return ds.sqlx.DB
	}

	if db, ok := ds.ORM.(SQLXClient); ok {
		ds.sqlx = db
		if db.DB != nil {
			ds.rdb.DB = db.DB.DB
		}

		return db.DB
	}

	if sqlxDB, ok := ds.ORM.(*sqlx.DB); ok {
		return sqlxDB
	}

	return nil
}

func (ds *DataStore) DB() *SQLClient {
	if ds.rdb.DB != nil {
		return &ds.rdb
	}

	if db := ds.GORM(); db != nil {
		return &SQLClient{DB: ds.GORM().DB(), config: ds.gorm.config, logger: ds.Logger}
	}

	if db := ds.SQLX(); db != nil {
		return &SQLClient{DB: ds.SQLX().DB, config: ds.sqlx.config, logger: ds.Logger}
	}

	if db, ok := ds.ORM.(*sql.DB); ok {
		ds.rdb.DB = db
		return &SQLClient{DB: db, config: ds.rdb.config, logger: ds.Logger}
	}

	return nil
}

func (ds *DataStore) SetORM(client interface{}) {
	// making sure that either gorm or sqlx is set and not both
	if ds.ORM != nil {
		return
	}

	switch v := client.(type) {
	case GORMClient:
		ds.gorm = v

		if v.DB != nil {
			ds.rdb.DB, ds.rdb.config, ds.rdb.logger = v.DB.DB(), v.config, ds.Logger
			ds.ORM = v.DB
		}
	case SQLXClient:
		if v.DB != nil {
			ds.ORM = v.DB
		}

		ds.sqlx = v
	}
}

func (ds *DataStore) SQLHealthCheck() types.Health {
	return ds.gorm.HealthCheck()
}

func (ds *DataStore) SQLXHealthCheck() types.Health {
	return ds.sqlx.HealthCheck()
}

func (ds *DataStore) CQLHealthCheck() types.Health {
	return ds.Cassandra.HealthCheck()
}

func (ds *DataStore) YCQLHealthCheck() types.Health {
	return ds.YCQL.HealthCheck()
}

func (ds *DataStore) ElasticsearchHealthCheck() types.Health {
	return ds.Elasticsearch.HealthCheck()
}

func (ds *DataStore) MongoHealthCheck() types.Health {
	return ds.MongoDB.HealthCheck()
}

func (ds *DataStore) RedisHealthCheck() types.Health {
	return ds.Redis.HealthCheck()
}

func (ds *DataStore) PubSubHealthCheck() types.Health {
	return ds.PubSub.HealthCheck()
}