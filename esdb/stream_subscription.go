package esdb

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
)

type StreamSubscription struct {
	*egos.BaseSubscription
	client *esdb.Client
}

func NewStreamSubscription(client *esdb.Client) egos.Subscription {
	return &StreamSubscription{
		BaseSubscription: &egos.BaseSubscription{},
		client:           client,
	}
}

func (s *StreamSubscription) Start(ctx context.Context) error {
	return nil
}

func (s *StreamSubscription) Stop() error {
	return nil
}
