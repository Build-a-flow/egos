package finkgoes

type Event interface {
	AggregateID() string
	GetHeaders() map[string]interface{}
	SetHeader(string, interface{})
	Event() interface{}
	EventType() string
	Version() *int
}

type EventDescriptor struct {
	id      string
	event   interface{}
	headers map[string]interface{}
	version *int
}

func NewEventMessage(aggregateID string, event interface{}, version *int) *EventDescriptor {
	return &EventDescriptor{
		id:      aggregateID,
		event:   event,
		headers: make(map[string]interface{}),
		version: version,
	}
}

func (c *EventDescriptor) EventType() string {
	return typeOf(c.event)
}

func (c *EventDescriptor) AggregateID() string {
	return c.id
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

func (c *EventDescriptor) Version() *int {
	return c.version
}