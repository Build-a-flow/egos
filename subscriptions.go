package finkgoes

type SubscriptionService interface {
	Start()
	Stop()
	AddHandler(EventHandler, ...Event)
}