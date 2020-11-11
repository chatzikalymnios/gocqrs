package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chatzikalymnios/gocqrs"
	"github.com/chatzikalymnios/gocqrs/aggregatestore/postgresql"
	"github.com/chatzikalymnios/gocqrs/example/accounts/domain/account"
	kafka2 "github.com/chatzikalymnios/gocqrs/publisher/kafka"
	postgresql2 "github.com/chatzikalymnios/gocqrs/relay/postgresql"
	"github.com/jackc/pgx/v4"
	"github.com/rs/xid"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"net/http"
	"os"
	"strings"
)

var aggregateStore gocqrs.AggregateStore

type accountAdapter struct {
	Name            string `json:"name"`
	StartingBalance int    `json:"startingBalance"`
}

type creditAdapter struct {
	CreditAmount int `json:"creditAmount"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi!")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	cmd := account.NewCreateAccountCommand(xid.New(), "Jane Doe", 10)
	account := gocqrs.TypeRegistry[account.AggregateType](cmd.AggregateId(), 0)
	events := account.Process(cmd)
	aggregateStore.Save(context.Background(), account, events)

	js, _ := json.Marshal(account)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		d := json.NewDecoder(r.Body)
		a := &accountAdapter{}
		err := d.Decode(a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// create the account
		cmd := account.NewCreateAccountCommand(xid.New(), a.Name, a.StartingBalance)
		aggregate := gocqrs.TypeRegistry[account.AggregateType](cmd.AggregateId(), 0)
		events := aggregate.Process(cmd)
		aggregateStore.Save(context.Background(), aggregate, events)

		// return the created account
		actAggregate, _ := aggregateStore.Load(context.Background(), account.AggregateType, aggregate.AggregateId())
		actualAccount, _ := actAggregate.(*account.Account)
		js, _ := json.Marshal(actualAccount)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break
	case "GET":
		accountId, _ := xid.FromString(strings.TrimPrefix(r.URL.Path, "/accounts/"))
		actAggregate, _ := aggregateStore.Load(context.Background(), account.AggregateType, accountId)
		actualAccount, _ := actAggregate.(*account.Account)
		js, _ := json.Marshal(actualAccount)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break
	case "PUT":
		accountId, _ := xid.FromString(strings.TrimPrefix(r.URL.Path, "/accounts/"))

		d := json.NewDecoder(r.Body)
		c := &creditAdapter{}
		err := d.Decode(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// credit the account
		cmd := account.NewCreditAccountCommand(accountId, c.CreditAmount)
		aggregate, _ := aggregateStore.Load(context.Background(), account.AggregateType, accountId)
		theAccount, _ := aggregate.(*account.Account)
		events := theAccount.Process(cmd)
		aggregateStore.Save(context.Background(), theAccount, events)

		// return the created account
		actAggregate, _ := aggregateStore.Load(context.Background(), account.AggregateType, aggregate.AggregateId())
		actualAccount, _ := actAggregate.(*account.Account)
		js, _ := json.Marshal(actualAccount)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		break
	}
}

func main() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())
	aggregateStore = postgresql.NewAggregateStore(db)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS")})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to kafka: %v\n", err)
		os.Exit(1)
	}
	defer p.Close()

	eventsChannel := make(chan *gocqrs.Event)
	publishErrChannel := make(chan error)
	publisher := kafka2.NewPublisher(p, eventsChannel, publishErrChannel)

	// start the kafka producer
	go publisher.Publish()

	relay := postgresql2.NewRelay(db, publisher)

	// start relaying events to the publisher
	go relay.Relay()

	http.HandleFunc("/", handler)
	http.HandleFunc("/accounts/", accountHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
