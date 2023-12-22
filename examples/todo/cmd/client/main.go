package main

import (
	"context"
	"database/sql"
	"examples/todo/domain"
	"log"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
	es "github.com/finktek/egos/esdb"
	"github.com/finktek/egos/postgres"
	_ "github.com/lib/pq"
)

func main() {
	eventStoreDbConfig, err := esdb.ParseConnectionString("esdb://localhost:2113?tls=false")
	if err != nil {
		panic(err)
	}
	eventStoreDbConfig.SkipCertificateVerification = true

	client, err := esdb.NewClient(eventStoreDbConfig)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	checkointStore := postgres.NewPostgresCheckpointStore(db)
	projection := NewProjectionHandler()
	subscription := es.NewAllStreamSubscription("my-todo-subscription", client, checkointStore)
	subscription.AddHandler(projection)
	subscription.Start(context.Background())
}

type TodoProjectionHandler struct {
}

func NewProjectionHandler() egos.SubscriptionHandler {
	return &TodoProjectionHandler{}
}

func (h *TodoProjectionHandler) Handle(ctx context.Context, event egos.Event) error {
	switch e := event.GetData().(type) {
	//	case *domain.TodoListCreated:
	//		log.Println("TodoListCreated", e.Title)
	//	case *domain.TodoItemAdded:
	//		log.Println("TodoItemAdded", e.Description)
	case *domain.TodoItemDone:
		log.Println("TodoItemDone", e.TodoItemID)
	}
	return nil
}
