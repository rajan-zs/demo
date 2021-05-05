package gofr

import (
	"strconv"
	"time"

	"github.com/zopsmart/gofr/pkg/datastore"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub/avro"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub/eventhub"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub/kafka"
)

// kafkaRetry retries connecting to kafka
// once connection is successful, retrying is terminated
func kafkaRetry(c *kafka.Config, avroConfig *avro.Config, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying Kafka connection")

		var err error

		k.PubSub, err = kafka.New(c, k.Logger)
		if err == nil {
			k.Logger.Info("Kafka initialized successfully")

			initializeAvro(avroConfig, k)

			break
		}
	}
}

// eventhubRetry retries connecting to eventhub
// once connection is successful, retrying is terminated
// also while retrying to connect to eventhub, initializes avro as well if configs are set
func eventhubRetry(c *eventhub.Config, avroConfig *avro.Config, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying Eventhub connection")

		var err error

		k.PubSub, err = eventhub.New(c)
		if err == nil {
			k.Logger.Info("Eventhub initialized successfully, Namespace: %v, Eventhub: %v\\n", c.Namespace, c.EventhubName)

			initializeAvro(avroConfig, k)

			break
		}
	}
}

func mongoRetry(c *datastore.MongoConfig, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying MongoDB connection")

		var err error

		k.MongoDB, err = datastore.GetNewMongoDB(k.Logger, c)

		if err == nil {
			k.Logger.Info("MongoDB initialized successfully")

			break
		}
	}
}

func yclRetry(c *datastore.CassandraCfg, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying Ycql connection")

		var err error

		k.YCQL, err = datastore.GetNewYCQL(k.Logger, c)
		if err == nil {
			k.Logger.Info("Ycql initialized successfully")
			break
		}
	}
}

func cassandraRetry(c *datastore.CassandraCfg, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying Cassandra connection")

		var err error

		k.Cassandra, err = datastore.GetNewCassandra(k.Logger, c)
		if err == nil {
			k.Logger.Info("Cassandra initialized successfully")

			break
		}
	}
}

func ormRetry(c *datastore.DBConfig, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying ORM connection")

		db, err := datastore.NewORM(c)
		if err == nil {
			k.SetORM(db)
			k.Logger.Info("ORM initialized successfully")

			break
		}
	}
}

func sqlxRetry(c *datastore.DBConfig, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnRetryDuration) * time.Second)

		k.Logger.Debug("Retrying SQLX connection")

		db, err := datastore.NewSQLX(c)
		if err == nil {
			k.SetORM(db)
			k.Logger.Info("SQLX initialized successfully")

			break
		}
	}
}

func redisRetry(c *datastore.RedisConfig, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnectionRetryDuration) * time.Second)

		k.Logger.Debug("Retrying Redis connection")

		var err error

		k.Redis, err = datastore.NewRedis(k.Logger, *c)
		if err == nil {
			k.Logger.Info("Redis initialized successfully")

			break
		}
	}
}

func elasticSearchRetry(c *datastore.ElasticSearchCfg, k *Gofr) {
	for {
		time.Sleep(time.Duration(c.ConnectionRetryDuration) * time.Second)

		k.Logger.Debug("Retrying ElasticSearch connection")

		var err error

		k.Elasticsearch, err = datastore.NewElasticsearchClient(c)

		if err == nil {
			k.Logger.Info("ElasticSearch initialized successfully")

			break
		}
	}
}

func getRetryDuration(envDuration string) int {
	retryDuration, _ := strconv.Atoi(envDuration)
	if retryDuration == 0 {
		// default duration 30 seconds
		retryDuration = 30
	}

	return retryDuration
}