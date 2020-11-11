package gocqrs

import (
	"github.com/rs/xid"
)

type EventType string
type EventData interface{}

type Event struct {
	EventId    xid.ID        `json:"eventId"`
	EventType  EventType     `json:"eventType"`
	EventData  EventData     `json:"eventData"`
	EntityType AggregateType `json:"entityType"`
	EntityId   xid.ID        `json:"entityId"`
}

func NewEvent(id xid.ID, eventType EventType, data EventData, aggregateType AggregateType, entityId xid.ID) *Event {
	return &Event{
		EventId:    id,
		EventType:  eventType,
		EventData:  data,
		EntityType: aggregateType,
		EntityId:   entityId,
	}
}
