package subscriptions

import "github.com/build-a-flow/egos"

type EventHandler interface {
	Handle(event egos.Event)
}
