package domain

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

func Init() *TodoList {
	todoAggregate := &TodoList{}
	base := egos.NewAggregateBase(todoAggregate.When)
	todoAggregate.AggregateBase = base
	return todoAggregate
}

func (t *TodoList) CreateTodoList(id string, title string) error {
	t.AggregateBase.Apply(&TodoListCreated{Id: id, Title: title})
	return nil
}

func (t *TodoList) AddItem(todoItemID string, todoDescription string) error {
	if todoItemID == "" {
		return errors.New("todo ID is required")
	}
	if todoDescription == "" {
		return errors.New("todo Description is required")
	}
	t.AggregateBase.Apply(&TodoItemAdded{TodoItemID: todoItemID, Description: todoDescription})
	return nil
}

func (t *TodoList) ItemDone(todoItemID string) error {
	if item := t.findItem(todoItemID); item != nil && item.Done == false {
		t.AggregateBase.Apply(&TodoItemDone{TodoItemID: item.TodoItemID})
		return nil
	}
	return errors.New("could not mark item as Done")
}

func (t *TodoList) When(event egos.Event) error {
	switch e := event.GetData().(type) {
	case *TodoListCreated:
		t.SetAggregateID(e.Id)
		t.Title = e.Title
	case *TodoItemAdded:
		t.Items = append(t.Items, &TodoItem{TodoItemID: e.TodoItemID, Description: e.Description, Done: false})
	case *TodoItemDone:
		if item := t.findItem(e.TodoItemID); item != nil {
			item.Done = true
		}
	}
	return nil
}

func (t *TodoList) findItem(todoItemID string) *TodoItem {
	for _, i := range t.Items {
		if i.TodoItemID == todoItemID {
			return i
		}
	}
	return nil
}
