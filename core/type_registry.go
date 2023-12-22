package egos

import (
	"reflect"
)

var eventRegistry = make(map[string]reflect.Type)

func RegisterEvent(event interface{}) error {
	t := reflect.TypeOf(event)
	eventRegistry[t.Name()] = t
	return nil
}

func GetEventInstance(eventName string) interface{} {
	if eventRegistry[eventName] == nil {
		return nil
	}
	return reflect.New(eventRegistry[eventName]).Interface()
}
