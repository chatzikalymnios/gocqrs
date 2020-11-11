package account

import (
	"github.com/chatzikalymnios/gocqrs"
	"github.com/rs/xid"
)

const AggregateType = gocqrs.AggregateType("Account")

type Status string

const (
	Open   = Status("Open")
	Closed = Status("Closed")
)

func init() {
	gocqrs.TypeRegistry[AggregateType] = func(id xid.ID, version int) gocqrs.Aggregate {
		return &Account{
			BaseAggregate: gocqrs.NewBaseAggregate(AggregateType, id, version),
		}
	}
}

type Account struct {
	*gocqrs.BaseAggregate
	Name    string `json:"name"`
	Balance int    `json:"balance"`
	Status  Status `json:"status"`
}

func (a *Account) Process(cmd gocqrs.Command) []gocqrs.Event {
	switch cmd.CommandType() {
	case CreateAccount:
		createCmd, ok := cmd.(*CreateAccountCommand)
		if !ok {
			panic("invalid CreateAccountCommand")
		}
		return a.processCreateCommand(*createCmd)
	case CloseAccount:
		closeCmd, ok := cmd.(*CloseAccountCommand)
		if !ok {
			panic("invalid CloseAccountCommand")
		}
		return a.processCloseCommand(*closeCmd)
	case CreditAccount:
		creditCmd, ok := cmd.(*CreditAccountCommand)
		if !ok {
			panic("invalid CreditAccountCommand")
		}
		return a.processCreditCommand(*creditCmd)
	case DebitAccount:
		debitCmd, ok := cmd.(*DebitAccountCommand)
		if !ok {
			panic("invalid DebitAccountCommand")
		}
		return a.processDebitCommand(*debitCmd)
	default:
		panic("unknown command type")
	}
}

func (a *Account) Apply(event gocqrs.Event) {
	switch event.EventType {
	case AccountCreated:
		createdEventData, ok := event.EventData.(*AccountCreatedEventData)
		if !ok {
			panic("invalid data for AccountCreatedEvent")
		}
		a.applyCreatedEvent(*createdEventData)
		break
	case AccountClosed:
		closedEventData, ok := event.EventData.(*AccountClosedEventData)
		if !ok {
			panic("invalid data for AccountClosedEvent")
		}
		a.applyClosedEvent(*closedEventData)
		break
	case AccountCredited:
		creditedEventData, ok := event.EventData.(*AccountCreditedEventData)
		if !ok {
			panic("invalid data for AccountCreditedEvent")
		}
		a.applyCreditedEvent(*creditedEventData)
	case AccountDebited:
		debitedEventData, ok := event.EventData.(*AccountDebitedEventData)
		if !ok {
			panic("invalid data for AccountDebitedEvent")
		}
		a.applyDebitedEvent(*debitedEventData)
	default:
		panic("unknown event type")
	}
}

func (a *Account) processCreateCommand(cmd CreateAccountCommand) []gocqrs.Event {
	return []gocqrs.Event{
		a.NewEvent(AccountCreated, AccountCreatedEventData{
			Name:            cmd.Name,
			StartingBalance: cmd.StartingBalance,
		}),
	}
}

func (a *Account) applyCreatedEvent(e AccountCreatedEventData) {
	a.Name = e.Name
	a.Balance = e.StartingBalance
	a.Status = Open
}

func (a *Account) processCloseCommand(cmd CloseAccountCommand) []gocqrs.Event {
	if a.Balance != 0 {
		panic("Can't close non-empty account!")
	}
	return []gocqrs.Event{
		a.NewEvent(AccountClosed, nil),
	}
}

func (a *Account) applyClosedEvent(e AccountClosedEventData) {
	a.Status = Closed
}

func (a *Account) processCreditCommand(cmd CreditAccountCommand) []gocqrs.Event {
	if a.Status != Open {
		panic("Can't credit not open account!")
	}
	return []gocqrs.Event{
		a.NewEvent(AccountCredited, AccountCreditedEventData{
			CreditAmount: cmd.CreditAmount,
		}),
	}
}

func (a *Account) applyCreditedEvent(e AccountCreditedEventData) {
	a.Balance += e.CreditAmount
}

func (a *Account) processDebitCommand(cmd DebitAccountCommand) []gocqrs.Event {
	if a.Status != Open {
		panic("Can't debit not open account!")
	}
	return []gocqrs.Event{
		a.NewEvent(AccountDebited, AccountDebitedEventData{
			DebitAmount: cmd.DebitAmount,
		}),
	}
}

func (a *Account) applyDebitedEvent(e AccountDebitedEventData) {
	a.Balance -= e.DebitAmount
}
