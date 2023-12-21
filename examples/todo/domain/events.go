package domain

import "github.com/finktek/egos"

var _ = egos.RegisterEvent(TodoListCreated{})
var _ = egos.RegisterEvent(TodoItemAdded{})
var _ = egos.RegisterEvent(TodoItemDone{})

type TodoListCreated struct {
	Title string
}

type TodoItemAdded struct {
	TodoItemID  string
	Description string
}

type TodoItemDone struct {
	TodoItemID string
}
