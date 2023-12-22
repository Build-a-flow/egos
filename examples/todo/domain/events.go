package domain

import "github.com/finktek/egos"

var (
	_ = egos.RegisterEvent(TodoListCreated{})
	_ = egos.RegisterEvent(TodoItemAdded{})
	_ = egos.RegisterEvent(TodoItemDone{})
)

type TodoListCreated struct {
	Id    string
	Title string
}

type TodoItemAdded struct {
	TodoItemID  string
	Description string
}

type TodoItemDone struct {
	TodoItemID string
}
