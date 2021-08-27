package finkgoes

import (
	"log"
)

type AggregateType string

func (at AggregateType) StreamNameFor(aggregateId string) string {
	return string(at) + "_" + aggregateId
}

type AggregateRoot interface {
	AggregateID() string
	CurrentVersion() int
	OriginalVersion() int
	When(event Event)
	Apply(aggregate AggregateRoot, event interface{})
	Fold(aggregate AggregateRoot, events []Event)
	GetChanges() []Event
	ClearChanges()
}

type AggregateBase struct {
	id      string
	currentVersion int
	originalVersion int
	changes  []Event
}

func NewAggregateBase(id string) *AggregateBase {
	return &AggregateBase{
		id:      id,
		changes:  []Event{},
		currentVersion: -1,
		originalVersion: -1,
	}
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

func (a *AggregateBase) Apply(aggregate AggregateRoot, event interface{}) {
	eventMessage := NewEventMessage(event)
	aggregate.When(eventMessage)
	a.changes = append(a.changes, eventMessage)
	log.Println(typeOf(event))
	a.currentVersion++
}

func (a *AggregateBase) Fold(aggregate AggregateRoot, events []Event) {
	for _, e := range events {
		a.currentVersion++
		a.originalVersion++
		aggregate.When(e)
	}
}

func (a *AggregateBase) GetChanges() []Event {
	return a.changes
}

func (a *AggregateBase) ClearChanges() {
	a.changes = []Event{}
}
