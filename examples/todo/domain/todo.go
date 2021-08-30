package domain

import (
	"errors"
	"fmt"
	"github.com/finktek/eventum"
	"github.com/gofrs/uuid"
)

type TodoList struct {
	*finkgoes.AggregateBase
	Title string      `json:"Title"`
	Items []*TodoItem `json:"Items"`
}

type TodoItem struct {
	TodoItemID  uuid.UUID `json:"item_id"`
	Description string    `json:"Description"`
	Done        bool      `json:"Done"`
}

func InitTodoList(id string) *TodoList {
	return &TodoList{
		AggregateBase: finkgoes.NewAggregateBase(id),
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
	for _, i := range t.Items {
		if i.TodoItemID == todoItemID && i.Done == false {
			t.AggregateBase.Apply(t, &TodoItemDone{ID: t.AggregateID(), TodoItemID: i.TodoItemID.String()})
			return nil
		}
	}
	return fmt.Errorf("could not mark item as Done")
}

func (t *TodoList) When(event finkgoes.Event) {
	switch e := event.GetData().(type) {
	case *TodoListCreated:
		t.Title = e.Title
	case *TodoItemAdded:
		t.Items = append(t.Items, &TodoItem{TodoItemID: uuid.Must(uuid.FromString(e.TodoItemID)), Description:e.Description, Done: false})
	case *TodoItemDone:
		for _, i := range t.Items {
			if i.TodoItemID == uuid.Must(uuid.FromString(e.TodoItemID)) {
				i.Done = true
			}
		}
	}
}