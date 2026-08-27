//go:build integration

package kafka

import (
	"context"
	"testing"
)

func TestNewProducer_Integration(t *testing.T) {
	// Requires a running Kafka container.
	// Run with: go test ./... -tags=integration
	brokers := []string{"localhost:29092"}

	producer, err := NewProducer(brokers)
	if err != nil {
		t.Fatalf("NewProducer failed: %v", err)
	}
	defer producer.Close()

	// Produce a test message to infra-health-check
	err = producer.Produce(context.Background(), "infra-health-check", "test-key", []byte(`{"phase":"3","test":true}`))
	if err != nil {
		t.Fatalf("Produce failed: %v", err)
	}
}
