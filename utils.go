package finkgoes

import (
	"github.com/satori/go.uuid"
	"reflect"

)

func typeOf(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}

func NewUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}