package domain

import (
	"errors"
	"fmt"

	"github.com/build-a-flow/egos"
	"github.com/google/uuid"
)

type TodoList struct {
	*egos.AggregateBase
	Title string      `json:"title"`
	Items []*TodoItem `json:"items"`
}

type TodoItem struct {
	TodoItemID  uuid.UUID `json:"item_id"`
	Description string    `json:"description"`
	Done        bool      `json:"is_done"`
}

func InitTodoList(id string) *TodoList {
	return &TodoList{
		AggregateBase: egos.NewAggregateBase(id),
	}
}

func (t *TodoList) CreateTodoList(title string) error {
	t.AggregateBase.Apply(t, &TodoListCreated{ID: t.AggregateID(), Title: title})
	return nil
}

func (t *TodoList) AddItem(todoItemID uuid.UUID, todoDescription string) error {
	if todoDescription == "" {
		return errors.New("todo Description is required")
	}
	t.AggregateBase.Apply(t, &TodoItemAdded{ID: t.AggregateID(), TodoItemID: todoItemID.String(), Description: todoDescription})
	return nil
}

func (t *TodoList) ItemDone(todoItemID uuid.UUID) error {
	if item := t.findItem(todoItemID); item != nil && item.Done == false {
		t.AggregateBase.Apply(t, &TodoItemDone{ID: t.AggregateID(), TodoItemID: item.TodoItemID.String()})
		return nil
	}
	return fmt.Errorf("could not mark item as Done")
}

func (t *TodoList) When(event egos.Event) {
	switch e := event.GetData().(type) {
	case *TodoListCreated:
		t.Title = e.Title
	case *TodoItemAdded:
		t.Items = append(t.Items, &TodoItem{TodoItemID: uuid.Must(uuid.Parse(e.TodoItemID)), Description:e.Description, Done: false})
	case *TodoItemDone:
		if item := t.findItem(uuid.Must(uuid.Parse(e.TodoItemID))); item != nil {
			item.Done = true
		}
	}
}

func (t *TodoList) findItem(todoItemID uuid.UUID) *TodoItem {
	for _, i := range t.Items {
		if i.TodoItemID == todoItemID {
			return i
		}
	}
	return nil
}