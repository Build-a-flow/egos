package egos

type AggregateRoot interface {
	AggregateID() string
	SetAggregateID(id string)
	CurrentVersion() int
	OriginalVersion() int
	GetChanges() []Event
	Fold(aggregate AggregateRoot, events []Event)
	ClearChanges()
	EmptyMetadata() map[string]interface{}
	Apply(event interface{})
	ApplyWithMetadata(event interface{}, metadata map[string]interface{})
	When
}

type AggregateBase struct {
	id              string
	currentVersion  int
	originalVersion int
	changes         []Event
	when            when
}

type When interface {
	When(event Event) error
}

type when func(event Event) error

func NewAggregateBase(when when) *AggregateBase {
	return &AggregateBase{
		changes:         []Event{},
		currentVersion:  -1,
		originalVersion: -1,
		when:            when,
	}
}

func (a *AggregateBase) SetAggregateID(id string) {
	a.id = id
}

func (a *AggregateBase) AggregateID() string {
	return a.id
}

func (a *AggregateBase) CurrentVersion() int {
	return a.currentVersion
}

func (a *AggregateBase) OriginalVersion() int {
	return a.originalVersion
}

func (a *AggregateBase) EmptyMetadata() map[string]interface{} {
	return make(map[string]interface{})
}

func (a *AggregateBase) Apply(event interface{}) {
	a.ApplyWithMetadata(event, a.EmptyMetadata())
}

func (a *AggregateBase) ApplyWithMetadata(event interface{}, metadata map[string]interface{}) {
	eventMessage := NewEventMessage(event, metadata)
	a.when(eventMessage)
	a.changes = append(a.changes, eventMessage)
	a.currentVersion++
}

func (a *AggregateBase) GetChanges() []Event {
	return a.changes
}

func (a *AggregateBase) ClearChanges() {
	a.changes = []Event{}
}

func (a *AggregateBase) Fold(aggregate AggregateRoot, events []Event) {
	for _, e := range events {
		a.currentVersion++
		a.originalVersion++
		aggregate.When(e)
	}
}

func StreamNameFor(aggregate AggregateRoot, aggregateId string) string {
	return typeOf(aggregate) + "_" + aggregateId
}
