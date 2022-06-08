package eventbridge

import (
	"encoding/json"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/prometheus/client_golang/prometheus"
)

type Client struct {
	client *eventbridge.EventBridge
	cfg    *Config
}

type Config struct {
	ConnRetryDuration int
	EventBus          string
	EventSource       string
	Region            string
	AccessKeyID       string
	SecretAccessKey   string
}

type customProvider struct {
	keyID     string
	secretKey string
}

//nolint // The declared global variable can be accessed across multiple functions
var (
	publishSuccessCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_publish_success_count",
		Help: "Counter for the number of events successfully published",
	}, []string{"topic", "consumerGroup"})

	publishFailureCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_publish_failure_count",
		Help: "Counter for the number of failed publish operations",
	}, []string{"topic", "consumerGroup"})

	publishTotalCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "zs_pubsub_publish_total_count",
		Help: "Counter for the total number of publish operations",
	}, []string{"topic", "consumerGroup"})
)

func (cp customProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{AccessKeyID: cp.keyID, SecretAccessKey: cp.secretKey}, nil
}

func (cp customProvider) IsExpired() bool {
	return false
}

// New returns new client
func New(cfg *Config) (*Client, error) {
	_ = prometheus.Register(publishFailureCount)
	_ = prometheus.Register(publishSuccessCount)
	_ = prometheus.Register(publishTotalCount)

	awsCfg := aws.NewConfig().WithRegion(cfg.Region)
	awsCfg.Credentials = credentials.NewCredentials(customProvider{cfg.AccessKeyID, cfg.SecretAccessKey})

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, err
	}

	client := eventbridge.New(sess, awsCfg)

	return &Client{
		client: client,
		cfg:    cfg,
	}, nil
}

// PublishEvent publishes the event to eventbridge
func (c *Client) PublishEvent(detailType string, detail interface{}, headers map[string]string) error {
	publishTotalCount.WithLabelValues(c.cfg.EventBus, "").Inc()

	payload, err := json.Marshal(detail)
	if err != nil {
		publishFailureCount.WithLabelValues(c.cfg.EventBus, "").Inc()
		return err
	}

	input := &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				Detail:       aws.String(string(payload)),
				DetailType:   aws.String(detailType),
				EventBusName: aws.String(c.cfg.EventBus),
				Source:       aws.String(c.cfg.EventSource),
			},
		},
	}

	_, err = c.client.PutEvents(input)
	if err != nil {
		publishFailureCount.WithLabelValues(c.cfg.EventBus, "").Inc()
		return err
	}

	publishSuccessCount.WithLabelValues(c.cfg.EventBus, "").Inc()

	return nil
}

// PublishEventWithOptions not implemented for Eventbridge
func (c *Client) PublishEventWithOptions(key string, value interface{}, headers map[string]string,
	options *pubsub.PublishOptions) (err error) {
	return nil
}

// Subscribe not implemented for Eventbridge
func (c *Client) Subscribe() (*pubsub.Message, error) {
	return nil, nil
}

// SubscribeWithCommit not implemented for Eventbridge
func (c *Client) SubscribeWithCommit(f pubsub.CommitFunc) (*pubsub.Message, error) {
	return nil, nil
}

// Bind not implemented for Eventbridge
func (c *Client) Bind(message []byte, target interface{}) error {
	return json.Unmarshal(message, &target)
}

// Ping not implemented for Eventbridge
func (c *Client) Ping() error {
	return nil
}

// CommitOffset not implemented for Eventbridge
func (c *Client) CommitOffset(offsets pubsub.TopicPartition) {

}

// HealthCheck checks eventbridge health.
func (c *Client) HealthCheck() types.Health {
	if c == nil {
		return types.Health{
			Name:   pkg.EventBridge,
			Status: pkg.StatusDown,
		}
	}

	resp := types.Health{
		Name:     pkg.EventBridge,
		Status:   pkg.StatusDown,
		Host:     c.cfg.Region,
		Database: c.cfg.EventBus,
	}

	if c.client == nil {
		return resp
	}

	resp.Status = pkg.StatusUp

	return resp
}

// IsSet checks whether eventbridge is initialized or not
func (c *Client) IsSet() bool {
	if c == nil {
		return false
	}

	if c.client == nil {
		return false
	}

	return true
}
