package gocqrs

import (
	"github.com/rs/xid"
)

type Entity interface {
	EntityId() xid.ID
}
