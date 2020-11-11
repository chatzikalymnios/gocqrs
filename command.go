package gocqrs

import "github.com/rs/xid"

type CommandType string

type Command interface {
	AggregateType() AggregateType
	AggregateId() xid.ID
	CommandType() CommandType
}

type BaseCommand struct {
	aggregateType AggregateType
	aggregateId   xid.ID
	commandType   CommandType
}

func NewBaseCommand(aggregateType AggregateType, aggregateId xid.ID, commandType CommandType) *BaseCommand {
	return &BaseCommand{
		aggregateType: aggregateType,
		aggregateId:   aggregateId,
		commandType:   commandType,
	}
}

func (b *BaseCommand) AggregateType() AggregateType {
	return b.aggregateType
}

func (b *BaseCommand) AggregateId() xid.ID {
	return b.aggregateId
}

func (b *BaseCommand) CommandType() CommandType {
	return b.commandType
}
