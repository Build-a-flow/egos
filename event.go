package finkgoes

import (
	"encoding/json"
	"log"
)

type Event interface {
	GetData() interface{}
	EventType() string
	Serialize() (serializedData []byte, serializedHeaders []byte)
}

type EventDescriptor struct {
	Data   interface{}
	Headers map[string]interface{}
}

func NewEventMessage(data interface{}) *EventDescriptor {
	return &EventDescriptor{
		Data:   data,
		Headers: make(map[string]interface{}),
	}
}

func (e EventDescriptor) GetData() interface{} {
	return e.Data
}

func (e EventDescriptor) EventType() string {
	return typeOf(e.Data)
}

func (e EventDescriptor) Serialize() (serializedData []byte, serializedHeaders []byte) {
	serializedData, err := json.Marshal(e.Data)
	if err != nil {
		log.Fatalf("error serializing event data: %s", err)
	}
	serializedHeaders, err = json.Marshal(e.Headers)
	if err != nil {
		log.Fatalf("error serializing event headers: %s", err)
	}
	return
}
