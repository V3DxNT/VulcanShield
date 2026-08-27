package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Producer wraps a franz-go client to provide a simple Produce helper.
// Phase 3 establishes the producer foundation; consumers are introduced in Phase 4.
type Producer struct {
	client *kgo.Client
}

// NewProducer creates a Kafka producer and verifies broker connectivity
// via a metadata fetch. Returns an error if brokers are unreachable.
// Kafka connectivity is non-fatal: callers should log a warning and continue.
func NewProducer(brokers []string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ProducerBatchMaxBytes(1<<20), // 1 MiB
	)
	if err != nil {
		return nil, fmt.Errorf("creating kafka producer: %w", err)
	}

	// Verify connectivity: fetch cluster metadata
	if err := client.Ping(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("kafka ping failed: %w", err)
	}

	return &Producer{client: client}, nil
}

// Produce sends a single record to the given topic.
// key is used for partition routing (e.g. transaction_id or user_id).
// value is expected to be JSON-encoded bytes.
func (p *Producer) Produce(ctx context.Context, topic, key string, value []byte) error {
	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}

	// ProduceSync blocks until the record is acknowledged or an error occurs.
	results := p.client.ProduceSync(ctx, record)
	if err := results.FirstErr(); err != nil {
		return fmt.Errorf("producing to topic %s: %w", topic, err)
	}
	return nil
}

// Close flushes pending records and closes the underlying client.
func (p *Producer) Close() {
	p.client.Close()
}
