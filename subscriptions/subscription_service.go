package subscriptions

import "context"

type SubscriptionService interface {
	Start(ctx context.Context)
	Stop()
	AddHandler(handler EventHandler)
}