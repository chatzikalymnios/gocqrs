package kafka

import (
	"encoding/json"
	"github.com/chatzikalymnios/gocqrs"
	"github.com/rs/xid"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Publisher struct {
	producer *kafka.Producer
	events   chan *gocqrs.Event
	ackChan  chan xid.ID
	errChan  chan error
}

func NewPublisher(producer *kafka.Producer, events chan *gocqrs.Event, errChan chan error) *Publisher {
	return &Publisher{
		producer: producer,
		events:   events,
		errChan:  errChan,
	}
}

func (p *Publisher) Publish() {
	for event := range p.events {
		eventData, err := json.Marshal(event)
		if err != nil {
			p.errChan <- err
		}

		topic := string(event.EntityType)

		p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          eventData,
			Key:            event.EntityId.Bytes(),
		}, nil)
	}
}

func (p *Publisher) Events() chan *gocqrs.Event {
	return p.events
}

func (p *Publisher) Err() chan error {
	return p.errChan
}
