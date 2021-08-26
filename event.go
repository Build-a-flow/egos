package finkgoes

import "time"

type Event interface {
	GetHeaders() map[string]interface{}
	SetHeader(string, interface{})
	Event() interface{}
	Version() int
	EventType() string
}

type EventDescriptor struct {
	event   interface{}
	headers map[string]interface{}
	version int
	timestamp time.Time
}

func NewEventMessage(event interface{}, version int, timestamp time.Time) *EventDescriptor {
	return &EventDescriptor{
		event:   event,
		headers: make(map[string]interface{}),
		version: version,
		timestamp: timestamp,
	}
}

func (c *EventDescriptor) EventType() string {
	return typeOf(c.event)
}

func (c *EventDescriptor) GetHeaders() map[string]interface{} {
	return c.headers
}

func (c *EventDescriptor) SetHeader(key string, value interface{}) {
	c.headers[key] = value
}

func (c *EventDescriptor) Event() interface{} {
	return c.event
}

func (c *EventDescriptor) Timestamp() time.Time {
	return c.timestamp
}

func (c *EventDescriptor) Version() int {
	return c.version
}