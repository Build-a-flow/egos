package egos

import (
	"encoding/json"
	"log"
)

type Event interface {
	GetData() interface{}
	EventType() string
	Serialize() (serializedData []byte, serializedMetadata []byte)
}

type EventDescriptor struct {
	Data     interface{}
	Metadata *Metadata
}

func NewEventMessage(data interface{}, metadata *Metadata) *EventDescriptor {
	return &EventDescriptor{
		Data:     data,
		Metadata: metadata,
	}
}

func (e EventDescriptor) GetData() interface{} {
	return e.Data
}

func (e EventDescriptor) EventType() string {
	return typeOf(e.Data)
}

func (e EventDescriptor) Serialize() (serializedData []byte, serializedMetadata []byte) {
	serializedData, err := json.Marshal(e.Data)
	if err != nil {
		log.Fatalf("error serializing event data: %s", err)
	}
	serializedMetadata, err = json.Marshal(e.Metadata.All())
	if err != nil {
		log.Fatalf("error serializing event headers: %s", err)
	}
	return
}
