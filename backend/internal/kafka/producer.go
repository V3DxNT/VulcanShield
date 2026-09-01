package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)



type Producer struct {
	client *kgo.Client
}




func NewProducer(brokers []string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ProducerBatchMaxBytes(1<<20), 
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating kafka producer: %w", err)
	}

	
	if err := client.Ping(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("kafka ping failed: %w", err)
	}

	return &Producer{client: client}, nil
}




func (p *Producer) Produce(ctx context.Context, topic, key string, value []byte) error {
	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}

	
	results := p.client.ProduceSync(ctx, record)
	if err := results.FirstErr(); err != nil {
		return fmt.Errorf("producing to topic %s: %w", topic, err)
	}
	return nil
}


func (p *Producer) Close() {
	p.client.Close()
}
