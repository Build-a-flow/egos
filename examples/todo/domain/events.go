package domain

import finkgoes "github.com/finktek/eventum"

var _ = finkgoes.RegisterEvent(TodoListCreated{})
var _ = finkgoes.RegisterEvent(TodoItemAdded{})
var _ = finkgoes.RegisterEvent(TodoItemDone{})

type TodoListCreated struct {
	ID   string `json:"id"`
	Title string `json:"Title"`
}

type TodoItemAdded struct {
	ID   		string 	`json:"id"`
	TodoItemID	string 	`json:"todo_item_id"`
	Description string  `json:"Description"`
}

type TodoItemDone struct {
	ID   		string 	`json:"id"`
	TodoItemID	string 	`json:"todo_item_id"`
}
