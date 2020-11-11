package account

import (
	"github.com/chatzikalymnios/gocqrs"
	"github.com/rs/xid"
)

const (
	CreateAccount = gocqrs.CommandType("CreateAccount")
	CloseAccount  = gocqrs.CommandType("CloseAccount")
	CreditAccount = gocqrs.CommandType("CreditAccount")
	DebitAccount  = gocqrs.CommandType("DebitAccount")
)

type CreateAccountCommand struct {
	*gocqrs.BaseCommand
	Name            string
	StartingBalance int
}

func NewCreateAccountCommand(id xid.ID, name string, startingBalance int) *CreateAccountCommand {
	return &CreateAccountCommand{
		BaseCommand:     gocqrs.NewBaseCommand(AggregateType, id, CreateAccount),
		Name:            name,
		StartingBalance: startingBalance,
	}
}

type CloseAccountCommand struct {
	*gocqrs.BaseCommand
}

type CreditAccountCommand struct {
	*gocqrs.BaseCommand
	CreditAmount int
}

func NewCreditAccountCommand(id xid.ID, creditAmount int) *CreditAccountCommand {
	return &CreditAccountCommand{
		BaseCommand:  gocqrs.NewBaseCommand(AggregateType, id, CreditAccount),
		CreditAmount: creditAmount,
	}
}

type DebitAccountCommand struct {
	*gocqrs.BaseCommand
	DebitAmount int
}
