package gocqrs

import "github.com/rs/xid"

var TypeRegistry = make(map[AggregateType]func(id xid.ID, version int) Aggregate)

var EventRegistry = make(map[EventType]func() EventData)
