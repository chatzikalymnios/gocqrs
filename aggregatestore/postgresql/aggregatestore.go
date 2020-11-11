package postgresql

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/chatzikalymnios/gocqrs"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
)

type AggregateStore struct {
	db *pgx.Conn
}

func NewAggregateStore(db *pgx.Conn) *AggregateStore {
	return &AggregateStore{db: db}
}

func (a *AggregateStore) Save(ctx context.Context, aggregate gocqrs.Aggregate, events []gocqrs.Event) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	for _, event := range events {
		err := a.insertEvent(ctx, event)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	if aggregate.Version() == 0 {
		err = a.insertEntity(ctx, aggregate, len(events))
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	} else {
		err = a.updateEntity(ctx, aggregate, aggregate.Version()+len(events))
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return nil
}

func (a *AggregateStore) insertEvent(ctx context.Context, event gocqrs.Event) error {
	data, err := json.Marshal(event.EventData)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(ctx,
		"insert into events(event_id, event_type, event_data, entity_type, entity_id) values($1, $2, $3, $4, $5)",
		event.EventId,
		event.EventType,
		data,
		event.EntityType,
		event.EntityId,
	)
	return err
}

func (a *AggregateStore) insertEntity(ctx context.Context, aggregate gocqrs.Aggregate, newVersion int) error {
	_, err := a.db.Exec(
		ctx,
		"insert into entities(entity_type, entity_id, entity_version) values($1, $2, $3)",
		aggregate.AggregateType(),
		aggregate.AggregateId(),
		newVersion,
	)
	return err
}

func (a *AggregateStore) updateEntity(ctx context.Context, aggregate gocqrs.Aggregate, newVersion int) error {
	cmd, err := a.db.Exec(
		ctx,
		"update entities set entity_version = $1 where entity_type = $2 and entity_id = $3 and entity_version = $4",
		newVersion,
		aggregate.AggregateType(),
		aggregate.AggregateId(),
		aggregate.Version(),
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		// concurrent update error
		return errors.New("concurrent update conflict")
	}

	return nil
}

func (a *AggregateStore) Load(ctx context.Context, aggregateType gocqrs.AggregateType, aggregateId xid.ID) (gocqrs.Aggregate, error) {
	aggregate, err := a.loadEntity(ctx, aggregateType, aggregateId)
	if err != nil {
		return nil, err
	}

	events, err := a.loadEvents(ctx, aggregateType, aggregateId)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		aggregate.Apply(event)
	}

	return aggregate, nil
}

func (a *AggregateStore) loadEntity(ctx context.Context, aggregateType gocqrs.AggregateType, aggregateId xid.ID) (gocqrs.Aggregate, error) {
	var entityType gocqrs.AggregateType
	var entityId xid.ID
	var entityVersion int
	err := a.db.QueryRow(
		ctx,
		"select entity_type, entity_id, entity_version from entities where entity_type = $1 and entity_id = $2",
		aggregateType,
		aggregateId,
	).Scan(&entityType, &entityId, &entityVersion)
	if err != nil {
		return nil, err
	}

	aggregate := gocqrs.TypeRegistry[entityType](entityId, entityVersion)
	return aggregate, nil
}

func (a *AggregateStore) loadEvents(ctx context.Context, aggregateType gocqrs.AggregateType, aggregateId xid.ID) ([]gocqrs.Event, error) {
	rows, err := a.db.Query(
		ctx,
		"select event_id, event_type, event_data from events where entity_type = $1 and entity_id = $2 order by event_id",
		aggregateType,
		aggregateId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []gocqrs.Event

	for rows.Next() {
		var eventId xid.ID
		var eventType gocqrs.EventType
		var eventDataMarshalled []byte
		err = rows.Scan(&eventId, &eventType, &eventDataMarshalled)
		if err != nil {
			return nil, err
		}

		eventData := gocqrs.EventRegistry[eventType]()
		err = json.Unmarshal(eventDataMarshalled, &eventData)
		if err != nil {
			return nil, err
		}

		events = append(events, *gocqrs.NewEvent(eventId, eventType, eventData, aggregateType, aggregateId))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

//func (a *AggregateStore) Load(ctx context.Context, aggregateType gocqrs.AggregateType, aggregateId xid.ID) (gocqrs.Aggregate, error) {
//	var entityType gocqrs.AggregateType
//	var entityId xid.ID
//	var entityVersion int
//	err := a.db.QueryRow(
//		ctx,
//		"select entity_type, entity_id, entity_version from entities where entity_type = $1 and entity_id = $2",
//		aggregateType,
//		aggregateId,
//	).Scan(&entityType, &entityId, &entityVersion)
//	if err != nil {
//		return nil, fmt.Errorf("load %v<%v>: %v", aggregateType, aggregateId, err)
//	}
//
//	aggregate := gocqrs.TypeRegistry[entityType](entityId, entityVersion)
//
//	rows, err := a.db.Query(
//		ctx,
//		"select event_id, event_type, event_data from events where entity_type = $1 and entity_id = $2 order by event_id",
//		aggregateType,
//		aggregateId,
//	)
//	if err != nil {
//		return nil, fmt.Errorf("load events for %v<%v>: %v", aggregateType, aggregateId, err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var eventId xid.ID
//		var eventType gocqrs.EventType
//		var eventDataMarshalled []byte
//		err = rows.Scan(&eventId, &eventType, &eventDataMarshalled)
//		if err != nil {
//			return nil, fmt.Errorf("read event for %v<%v>: %v", aggregateType, aggregateId, err)
//		}
//
//		eventData := gocqrs.EventRegistry[eventType]()
//		err = json.Unmarshal(eventDataMarshalled, &eventData)
//		if err != nil {
//			return nil, fmt.Errorf("unmarshall event data for %v<%v>: %v", aggregateType, aggregateId, err)
//		}
//
//		aggregate.Apply(gocqrs.NewDomainEvent(eventId, eventType, eventData, aggregateType, aggregateId))
//	}
//
//	if rows.Err() != nil {
//		return nil, rows.Err()
//	}
//
//	return aggregate, nil
//}
