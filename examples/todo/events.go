package main

import "github.com/finktek/egos"

var _ = egos.RegisterEvent(TodoListCreated{})
var _ = egos.RegisterEvent(TodoItemAdded{})
var _ = egos.RegisterEvent(TodoItemDone{})

type TodoListCreated struct {
	ID    string `json:"id"`
	Title string `json:"Title"`
}

type TodoItemAdded struct {
	ID          string `json:"id"`
	TodoItemID  string `json:"todo_item_id"`
	Description string `json:"Description"`
}

type TodoItemDone struct {
	ID         string `json:"id"`
	TodoItemID string `json:"todo_item_id"`
}
