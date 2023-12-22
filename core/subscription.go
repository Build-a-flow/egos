package egos

import "context"

type Subscription interface {
	Start(ctx context.Context) error
	Stop() error
	AddHandler(handler SubscriptionHandler) error
}

type BaseSubscription struct {
	handlers []SubscriptionHandler
}

func (s *BaseSubscription) AddHandler(handler SubscriptionHandler) error {
	s.handlers = append(s.handlers, handler)
	return nil
}
