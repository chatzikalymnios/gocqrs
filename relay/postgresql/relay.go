package postgresql

import (
	"context"
	"encoding/json"
	"github.com/chatzikalymnios/gocqrs"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
	"time"
)

type Relay struct {
	db *pgx.Conn
	p  gocqrs.Publisher
}

func NewRelay(db *pgx.Conn, p gocqrs.Publisher) *Relay {
	return &Relay{
		db: db,
		p:  p,
	}
}

func (r *Relay) Relay() {
	for true {
		// poll the database for new events
		events, err := r.loadEvents(context.Background())
		if err != nil {
			panic("could not load events to relay!")
		}

		for _, event := range events {
			r.p.Events() <- &event

			// mark event as published
			// in reality this should only be done if we get confirmation form the publisher that this specific event
			// was actually published, but for now assume that delivery is always successful
			r.markEvent(context.Background(), event)
		}

		// wait for new events to come in
		time.Sleep(5 * time.Second)
	}
}

func (r *Relay) loadEvents(ctx context.Context) ([]gocqrs.Event, error) {
	rows, err := r.db.Query(
		ctx,
		"select event_id, event_type, event_data, entity_type, entity_id from events where published is false order by event_id",
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
		var entityType gocqrs.AggregateType
		var entityId xid.ID
		err = rows.Scan(&eventId, &eventType, &eventDataMarshalled, &entityType, &entityId)
		if err != nil {
			return nil, err
		}

		eventData := gocqrs.EventRegistry[eventType]()
		err = json.Unmarshal(eventDataMarshalled, &eventData)
		if err != nil {
			return nil, err
		}

		events = append(events, *gocqrs.NewEvent(eventId, eventType, eventData, entityType, entityId))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

func (r *Relay) markEvent(ctx context.Context, event gocqrs.Event) error {
	_, err := r.db.Exec(
		ctx,
		"update events set published = true where event_id = $1 and event_type = $2",
		event.EventId,
		event.EventType,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Relay) Publisher() gocqrs.Publisher {
	return r.p
}
