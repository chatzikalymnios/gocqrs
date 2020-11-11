package service

import (
	"context"
	"github.com/chatzikalymnios/gocqrs"
	"github.com/rs/xid"
)

type AccountService struct {
	accountStore *gocqrs.AggregateStore
	publisher    *gocqrs.Publisher
}

func (s AccountService) CreateAccount(ctx context.Context, name string, startingBalance int) (xid.ID, error) {
	//accountId := xid.New()
	//createAccountCommand := account.NewCreateAccountCommand(accountId, name, startingBalance)
	//
	//account := account.NewAccount(accountId)
	//events := account.ProcessCreateCommand(*createAccountCommand)
	// publish events
	// commit

	return xid.New(), nil
}

func CloseAccount(ctx context.Context, accountId xid.ID) error {
	// start tx
	// load aggregate
	// create and process cmd (returns domain events)
	// emit domain events
	// commit tx

	// if at any point something fails, rollback transaction
	// return any errors
	return nil
}

func CreditAccount(ctx context.Context, accountId xid.ID, creditAmount int) error {
	return nil
}

func DebitAccount(ctx context.Context, accountId xid.ID, debitAmount int) error {
	return nil
}
