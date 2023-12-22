package egos

import (
	"context"
)

type SubscriptionHandler interface {
	Handle(ctx context.Context, event Event) error
}
