package domain

import egos "github.com/finktek/egos/core"

var (
	_ = egos.RegisterEvent(TodoListCreated{})
	_ = egos.RegisterEvent(TodoItemAdded{})
	_ = egos.RegisterEvent(TodoItemDone{})
)

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
