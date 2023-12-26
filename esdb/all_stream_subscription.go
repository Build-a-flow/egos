package esdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
)

type AllStreamSubscription struct {
	*egos.BaseSubscription
	subscriptionName string
	client           *esdb.Client
	subscription     *esdb.Subscription
	checkpointStore  egos.CheckpointStore
}

func NewAllStreamSubscription(subscriptionName string, client *esdb.Client, checkpointStore egos.CheckpointStore) egos.Subscription {
	return &AllStreamSubscription{
		BaseSubscription: &egos.BaseSubscription{},
		subscriptionName: subscriptionName,
		client:           client,
		checkpointStore:  checkpointStore,
	}
}

func (s *AllStreamSubscription) Start(ctx context.Context) error {
	s.subscribe(ctx)
	return nil
}

func (s *AllStreamSubscription) Stop() error {
	return nil
}

func (s *AllStreamSubscription) subscribe(ctx context.Context) {
	checkoint := s.checkpointStore.GetLastCheckpoint(s.subscriptionName)
	options := esdb.SubscribeToAllOptions{
		From: esdb.Position{
			Commit:  checkoint.Position,
			Prepare: checkoint.Position,
		},
		Filter: esdb.ExcludeSystemEventsFilter(),
	}
	subscription, err := s.client.SubscribeToAll(ctx, options)

	if err != nil {
		panic(err)
	}
	s.subscription = subscription

	for {
		event := s.subscription.Recv()

		if event.SubscriptionDropped != nil {
			fmt.Println("Subscription dropped", event.SubscriptionDropped.Error.Error())
			options.From = event.EventAppeared.OriginalEvent().Position
			break
		}

		if event.EventAppeared != nil {
			data := resolvedEventToEvent(event.EventAppeared)
			s.Handle(data)
			checkoint.Position = event.EventAppeared.OriginalEvent().Position.Commit
			s.checkpointStore.StoreCheckpoint(checkoint)
		}
	}
}

func resolvedEventToEvent(resolvedEvent *esdb.ResolvedEvent) egos.Event {
	eventData := egos.GetEventInstance(resolvedEvent.Event.EventType)
	m := make(map[string]interface{})
	err := json.Unmarshal(resolvedEvent.Event.UserMetadata, &m)
	if err != nil {
		fmt.Println("error:", err)
	}

	if eventData != nil {
		return &egos.EventDescriptor{
			Data:     eventData,
			Metadata: egos.NewMetadata(&m),
		}
	}
	return nil
}
