package account

import (
	"github.com/chatzikalymnios/gocqrs"
)

const (
	AccountCreated  = gocqrs.EventType("AccountCreated")
	AccountClosed   = gocqrs.EventType("AccountClosed")
	AccountCredited = gocqrs.EventType("AccountCredited")
	AccountDebited  = gocqrs.EventType("AccountDebited")
)

func init() {
	gocqrs.EventRegistry[AccountCreated] = func() gocqrs.EventData {
		return &AccountCreatedEventData{}
	}

	gocqrs.EventRegistry[AccountClosed] = func() gocqrs.EventData {
		return &AccountClosedEventData{}
	}

	gocqrs.EventRegistry[AccountCredited] = func() gocqrs.EventData {
		return &AccountCreditedEventData{}
	}

	gocqrs.EventRegistry[AccountDebited] = func() gocqrs.EventData {
		return &AccountDebitedEventData{}
	}
}

type AccountCreatedEventData struct {
	Name            string `json:"name"`
	StartingBalance int    `json:"startingBalance"`
}

type AccountClosedEventData struct{}

type AccountCreditedEventData struct {
	CreditAmount int `json:"creditAmount"`
}

type AccountDebitedEventData struct {
	DebitAmount int `json:"debitAmount"`
}
