package main

import (
	"context"
	"log"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
	es "github.com/finktek/egos/esdb"
)

func main() {
	eventStoreDbConfig, err := esdb.ParseConnectionString("esdb://localhost:2113?tls=true")
	if err != nil {
		panic(err)
	}
	eventStoreDbConfig.SkipCertificateVerification = true

	client, err := esdb.NewClient(eventStoreDbConfig)
	if err != nil {
		panic(err)
	}

	projection := NewProjectionHandler()
	subscription := es.NewAllStreamSubscription(client)
	subscription.AddHandler(projection)
	subscription.Start(context.Background())
}

type TodoProjectionHandler struct {
}

func NewProjectionHandler() egos.SubscriptionHandler {
	return &TodoProjectionHandler{}
}

func (h *TodoProjectionHandler) Handle(ctx context.Context, event egos.Event) error {
	log.Println("Event: ", event)
	return nil
}
