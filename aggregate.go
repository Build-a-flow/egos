package finkgoes

type AggregateRoot interface {
	GetAggregateID() string
	GetOriginalVersion() int
	GetCurrentVersion() int
	Apply(events Event, isNew bool)
	AddChange(Event)
	GetChanges() []Event
	ClearChanges()
}

type AggregateBase struct {
	id      string
	version int
	changes []Event
}

func NewAggregateBase(id string) *AggregateBase {
	return &AggregateBase{
		id:      id,
		changes: []Event{},
		version: -1,
	}
}

func (a *AggregateBase) GetAggregateID() string {
	return a.id
}

func (a *AggregateBase) GetOriginalVersion() int {
	return a.version
}

func (a *AggregateBase) GetCurrentVersion() int {
	return a.version + len(a.changes)
}

func (a *AggregateBase) AddChange(event Event) {
	a.changes = append(a.changes, event)
}

func (a *AggregateBase) GetChanges() []Event {
	return a.changes
}

func (a *AggregateBase) ClearChanges() {
	a.changes = []Event{}
}