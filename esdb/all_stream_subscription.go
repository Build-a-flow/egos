package esdb

import (
	"context"
	"fmt"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
)

type AllStreamSubscription struct {
	*egos.BaseSubscription
	client       *esdb.Client
	subscription *esdb.Subscription
}

func NewAllStreamSubscription(client *esdb.Client) egos.Subscription {
	return &AllStreamSubscription{
		BaseSubscription: &egos.BaseSubscription{},
		client:           client,
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
	subscription, err := s.client.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{})

	if err != nil {
		panic(err)
	}
	s.subscription = subscription

	for {
		event := s.subscription.Recv()

		if event.EventAppeared != nil {
			fmt.Println(fmt.Sprintf("%s:%s", event.EventAppeared.Event.StreamID, event.EventAppeared.Event.EventType))
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

}
