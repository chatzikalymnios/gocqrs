package gocqrs

import (
	"context"
	"github.com/rs/xid"
)

type AggregateStore interface {
	Save(ctx context.Context, aggregate Aggregate, events []Event) error
	Load(ctx context.Context, aggregateType AggregateType, aggregateId xid.ID) (Aggregate, error)
}
