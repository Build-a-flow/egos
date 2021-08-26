package finkgoes

import (
	"reflect"
	"github.com/gofrs/uuid"

)

func typeOf(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}

func NewUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}