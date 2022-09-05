package eventbridge

import (
	"encoding/json"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg"
	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	pkgEventbridge "github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	cfg := &Config{
		Region:      "us-east-1",
		EventBus:    "Gofr",
		EventSource: "application",
	}

	_, err := New(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestEventBridge_PublishEvent(t *testing.T) {
	ch := make(chan int)
	tcs := []struct {
		region string
		detail interface{}
		err    error
	}{
		{"", "sample payload", awserr.New("MissingRegion", "could not find region configuration", nil)},
		{"us-east-1", ch, &json.UnsupportedTypeError{Type: reflect.TypeOf(ch)}},
		{"us-east-1", "sample payload", nil},
	}

	for i, tc := range tcs {
		var eBridge Client

		awscfg := aws.NewConfig().WithRegion(tc.region)
		awscfg.Credentials = credentials.NewStaticCredentials("AKID", "SECRET_KEY", "TOKEN")

		eBridge.client = pkgEventbridge.New(mock.Session, awscfg)
		eBridge.cfg = &Config{EventBus: "gofr", EventSource: "application"}
		eb := &eBridge
		err := eb.PublishEvent("myDetailType", tc.detail, map[string]string{})
		assert.Equal(t, tc.err, err, "Test case failed [%v]\n Expected: %v, got: %v", i, tc.err, err)
	}
}

func TestEventBridge_HealthCheck(t *testing.T) {
	var eBridge Client

	awscfg := aws.NewConfig().WithRegion("us-west-2")
	awscfg.Credentials = credentials.NewStaticCredentials("AKID", "SECRET_KEY", "TOKEN")

	eBridge.client = pkgEventbridge.New(mock.Session, awscfg)
	eBridge.cfg = &Config{EventBus: "gofr", EventSource: "application"}
	eb := &eBridge

	testcases := []struct {
		client  *Client
		expResp types.Health
	}{
		{client: nil, expResp: types.Health{Name: pkg.EventBridge, Status: pkg.StatusDown}},
		{client: &Client{client: nil, cfg: &Config{EventBus: "gofr", Region: "us-west-2"}},
			expResp: types.Health{Name: pkg.EventBridge, Status: pkg.StatusDown, Host: "us-west-2", Database: "gofr"}},
		{client: eb, expResp: types.Health{Name: pkg.EventBridge, Status: pkg.StatusUp, Host: "", Database: "gofr"}},
	}
	for i, tc := range testcases {
		resp := tc.client.HealthCheck()
		assert.Equalf(t, tc.expResp, resp, "Test case failed [%v]. Expected: %v, got: %v", i, tc.expResp, resp)
	}
}

func TestEventBridge_PublishEventWithOptions(t *testing.T) {
	c := Client{}

	err := c.PublishEventWithOptions("", "", map[string]string{}, &pubsub.PublishOptions{})
	if err != nil {
		t.Error("Test case failed.")
	}
}

func TestEventBridge_Subscribe(t *testing.T) {
	c := Client{}

	_, err := c.Subscribe()
	if err != nil {
		t.Error("Test case failed")
	}
}

func TestEventBridge_SubscribeWithCommit(t *testing.T) {
	c := Client{}
	f := func(message *pubsub.Message) (bool, bool) { return false, false }

	_, err := c.SubscribeWithCommit(f)
	if err != nil {
		t.Error("Test case failed")
	}
}

func TestEventBridge_Bind(t *testing.T) {
	var k string

	c := Client{}

	err := c.Bind([]byte(`{"test":"test"}`), k)
	if err != nil {
		t.Error("Test case failed.")
	}
}

func TestEventBridge_Ping(t *testing.T) {
	c := Client{}

	err := c.Ping()
	if err != nil {
		t.Error("Test case failed")
	}
}

func TestEventBridge_IsSet(t *testing.T) {
	var eBridge Client

	awscfg := aws.NewConfig().WithRegion("us-west-2")
	awscfg.Credentials = credentials.NewStaticCredentials("AKID", "SECRET_KEY", "TOKEN")

	eBridge.client = pkgEventbridge.New(mock.Session, awscfg)
	eBridge.cfg = &Config{EventBus: "gofr", EventSource: "application"}
	eb := &eBridge

	testcases := []struct {
		client  *Client
		expResp bool
	}{
		{client: nil, expResp: false},
		{client: &Client{client: nil, cfg: &Config{}}, expResp: false},
		{client: eb, expResp: true},
	}
	for i, tc := range testcases {
		resp := tc.client.IsSet()
		assert.Equalf(t, tc.expResp, resp, "Test case failed [%v]. \n Expected: %v, got %v", i, tc.expResp, resp)
	}
}
