package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zopsmart/gofr/pkg/datastore/pubsub"
	"github.com/zopsmart/gofr/pkg/errors"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/gofr/types"
)

type mockPubSub struct {
	id string
}

func (m *mockPubSub) CommitOffset(offsets pubsub.TopicPartition) {
}

func (m *mockPubSub) PublishEventWithOptions(key string, val interface{}, headers map[string]string, options *pubsub.PublishOptions) error {
	if m.id == "1" {
		return errors.EntityNotFound{ID: "1"}
	}

	return nil
}

func (m *mockPubSub) PublishEvent(key string, val interface{}, headers map[string]string) error {
	if m.id == "1" {
		return errors.EntityNotFound{ID: "1"}
	}

	return nil
}

func (m *mockPubSub) Subscribe() (*pubsub.Message, error) {
	return &pubsub.Message{}, nil
}

func (m *mockPubSub) SubscribeWithCommit(commitFunc pubsub.CommitFunc) (*pubsub.Message, error) {
	return &pubsub.Message{}, nil
}

func (m *mockPubSub) Bind(v []byte, target interface{}) error {
	return nil
}

func (m *mockPubSub) Ping() error {
	return nil
}

func (m *mockPubSub) HealthCheck() types.Health {
	return types.Health{}
}

func (m *mockPubSub) IsSet() bool {
	return true
}

func TestProducerHandler(t *testing.T) {
	k := gofr.New()
	m := mockPubSub{}
	k.PubSub = &m

	tests := []struct {
		name    string
		id      string
		want    interface{}
		wantErr bool
	}{
		{"error from publisher", "1", nil, true},
		{"success", "123", nil, false},
	}

	req := httptest.NewRequest("GET", "http://dummy", nil)
	context := gofr.NewContext(nil, request.NewHTTPRequest(req), k)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			context.SetPathParams(map[string]string{
				"id": tt.id,
			})

			m.id = tt.id
			got, err := Producer(context)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConsumerHandler(t *testing.T) {
	k := gofr.New()
	k.PubSub = &mockPubSub{}

	ctx := gofr.NewContext(nil, nil, k)
	_, err := Consumer(ctx)

	assert.Equal(t, nil, err)
}