package subscriptions

import finkgoes "github.com/finktek/eventum"

type EventHandler interface {
	Handle(event finkgoes.Event)
}