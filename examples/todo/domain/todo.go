package domain

import (
	"context"
	"errors"

	egos "github.com/finktek/egos/core"
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

func Init(id string) TodoList {
	aggregate := TodoList{}
	aggregate.AggregateBase = egos.NewAggregateBase(aggregate.When)
	aggregate.AggregateBase.SetAggregateID(id)
	return aggregate
}

func (t *TodoList) CreateTodoList(ctx context.Context, title string) error {
	t.AggregateBase.Apply(ctx, &TodoListCreated{Title: title})
	return nil
}

func (t *TodoList) AddItem(ctx context.Context, todoItemID string, todoDescription string) error {
	if todoItemID == "" {
		return errors.New("todo ID is required")
	}
	if todoDescription == "" {
		return errors.New("todo Description is required")
	}
	t.AggregateBase.Apply(ctx, &TodoItemAdded{TodoItemID: todoItemID, Description: todoDescription})
	return nil
}

func (t *TodoList) ItemDone(ctx context.Context, todoItemID string) error {
	if item := t.findItem(todoItemID); item != nil && item.Done == false {
		t.AggregateBase.Apply(ctx, &TodoItemDone{TodoItemID: item.TodoItemID})
		return nil
	}
	return errors.New("could not mark item as Done")
}

func (t *TodoList) findItem(todoItemID string) *TodoItem {
	for _, i := range t.Items {
		if i.TodoItemID == todoItemID {
			return i
		}
	}
	return nil
}
