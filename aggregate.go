package gocqrs

import "github.com/rs/xid"

type AggregateType string

type Aggregate interface {
	AggregateType() AggregateType
	AggregateId() xid.ID
	Version() int
	Process(cmd Command) []Event
	Apply(event Event)
}

type BaseAggregate struct {
	t  AggregateType
	Id xid.ID `json:"id"`
	V  int    `json:"version"`
}

func NewBaseAggregate(t AggregateType, id xid.ID, version int) *BaseAggregate {
	return &BaseAggregate{
		t:  t,
		Id: id,
		V:  version,
	}
}

func (b *BaseAggregate) AggregateType() AggregateType {
	return b.t
}

func (b *BaseAggregate) AggregateId() xid.ID {
	return b.Id
}

func (b *BaseAggregate) Version() int {
	return b.V
}

func (b *BaseAggregate) Apply(event Event) {
	panic("implement me")
}

func (b *BaseAggregate) NewEvent(eventType EventType, eventData interface{}) Event {
	return Event{
		EventId:    xid.New(),
		EventType:  eventType,
		EventData:  eventData,
		EntityType: b.t,
		EntityId:   b.Id,
	}
}
