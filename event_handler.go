package finkgoes

type EventHandler interface {
	Handle(Event)
}