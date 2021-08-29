package finkgoes

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
	return reflect.New(eventRegistry[eventName]).Interface()
}
