package main

import (
	"errors"

	"github.com/finktek/egos"
)

type TodoList struct {
	*egos.AggregateBase
	Title string
	Items []*TodoItem
}

type TodoItem struct {
	TodoItemID  string
	Description string
	Done        bool
}

func Init(id string) *TodoList {
	return &TodoList{
		AggregateBase: egos.NewAggregateBase(id),
	}
}

func (t *TodoList) CreateTodoList(title string) error {
	t.AggregateBase.Apply(t, &TodoListCreated{Title: title})
	return nil
}

func (t *TodoList) AddItem(todoItemID string, todoDescription string) error {
	if todoDescription == "" {
		return errors.New("todo Description is required")
	}
	t.AggregateBase.Apply(t, &TodoItemAdded{TodoItemID: todoItemID, Description: todoDescription})
	return nil
}

func (t *TodoList) ItemDone(todoItemID string) error {
	if item := t.findItem(todoItemID); item != nil && item.Done == false {
		t.AggregateBase.Apply(t, &TodoItemDone{TodoItemID: item.TodoItemID})
		return nil
	}
	return errors.New("could not mark item as Done")
}

func (t *TodoList) When(event egos.Event) {
	switch e := event.GetData().(type) {
	case *TodoListCreated:
		t.Title = e.Title
	case *TodoItemAdded:
		t.Items = append(t.Items, &TodoItem{TodoItemID: e.TodoItemID, Description: e.Description, Done: false})
	case *TodoItemDone:
		if item := t.findItem(e.TodoItemID); item != nil {
			item.Done = true
		}
	}
}

func (t *TodoList) findItem(todoItemID string) *TodoItem {
	for _, i := range t.Items {
		if i.TodoItemID == todoItemID {
			return i
		}
	}
	return nil
}
