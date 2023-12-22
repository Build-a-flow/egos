package egos

import "context"

type Subscription interface {
	Start(ctx context.Context) error
	Stop() error
	AddHandler(handler SubscriptionHandler) error
	Handle(event Event) error
}

type BaseSubscription struct {
	handlers []SubscriptionHandler
}

func (s *BaseSubscription) AddHandler(handler SubscriptionHandler) error {
	s.handlers = append(s.handlers, handler)
	return nil
}

func (s *BaseSubscription) Handle(event Event) error {
	for _, h := range s.handlers {
		err := h.Handle(context.Background(), event)
		if err != nil {
			return err
		}
	}
	return nil
}
