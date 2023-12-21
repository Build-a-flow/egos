package esdb

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/finktek/egos"
)

type iterator struct {
	stream *esdb.ReadStream
}

// Close closes the stream
func (i *iterator) Close() {
	i.stream.Close()
}

// Next returns next event from the stream
func (i *iterator) Next() (egos.Event, error) {

	eventESDB, err := i.stream.Recv()
	if errors.Is(err, io.EOF) {
		return egos.EventDescriptor{}, errors.New("EOF")
	}
	if err, ok := esdb.FromError(err); !ok {
		if err.Code() == esdb.ErrorCodeResourceNotFound {
			return egos.EventDescriptor{}, errors.New("EOF")
		}
	}
	if err != nil {
		return egos.EventDescriptor{}, err
	}

	eventData := egos.GetEventInstance(eventESDB.Event.EventType)
	if eventData == nil {
		return egos.EventDescriptor{}, errors.New("event type not found")
	}
	json.Unmarshal(eventESDB.Event.Data, &eventData)
	event := &egos.EventDescriptor{
		Data:    eventData,
		Headers: make(map[string]interface{}),
	}
	return event, nil
}
