package finkgoes

type EventStore interface {
	AppendEvents(streamName string, expectedVersion int, events []Event) error
	ReadEvents(streamName string, start int, limit int) []Event
	ReadEventsBackwards(streamName string, limit int) []Event
	ReadStream(streamName string, start int, callback func())
	TruncateStream(streamName string, position int, expectedVersion int)
	DeleteStream(streamName string, expectedVersion int)
}
