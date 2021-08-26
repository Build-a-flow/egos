package finkgoes

import (
	"log"
	"time"
)

var TimeNow = time.Now

type AggregateType string

func (at AggregateType) StreamNameFor(aggregateId string) string {
	return string(at) + "_" + aggregateId
}

type AggregateRoot interface {
	AggregateID() string
	Version() int
	IncrementVersion()
	SetVersion(version int)
	ApplyEvent(event Event) error
	AppendEvent(aggregate AggregateRoot, event interface{})
	GetEvents() []Event
	ClearEvents()
}

type AggregateBase struct {
	id      string
	version int
	events  []Event
}

func NewAggregateBase(id string) *AggregateBase {
	return &AggregateBase{
		id:      id,
		events:  []Event{},
		version: -1,
	}
}

func (a *AggregateBase) AggregateID() string {
	return a.id
}

func (a *AggregateBase) Version() int {
	return a.version
}

func (a *AggregateBase) IncrementVersion()  {
	a.version++
}

func (a *AggregateBase) SetVersion(version int) {
	a.version = version
}

func (a *AggregateBase) AppendEvent(aggregate AggregateRoot, event interface{}) {
	a.IncrementVersion()
	eventMessage := NewEventMessage(event, a.version, TimeNow())
	err := aggregate.ApplyEvent(eventMessage)
	if err != nil {
		log.Fatalf(err.Error())
	}
	a.events = append(a.events, eventMessage)
}

func (a *AggregateBase) GetEvents() []Event {
	return a.events
}

func (a *AggregateBase) ClearEvents() {
	a.events = []Event{}
}