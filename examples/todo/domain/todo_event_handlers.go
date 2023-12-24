package domain

import egos "github.com/finktek/egos/core"

func (t *TodoList) When(event egos.Event) error {
	switch e := event.GetData().(type) {
	case *TodoListCreated:
		t.whenTodoListCreated(e)
	case *TodoItemAdded:
		t.whenTodoItemAdded(e)
	case *TodoItemDone:
		t.whenTodoItemDone(e)
	}
	return nil
}

func (t *TodoList) whenTodoListCreated(e *TodoListCreated) error {
	t.Title = e.Title
	return nil
}

func (t *TodoList) whenTodoItemAdded(e *TodoItemAdded) error {
	t.Items = append(t.Items, &TodoItem{TodoItemID: e.TodoItemID, Description: e.Description, Done: false})
	return nil
}

func (t *TodoList) whenTodoItemDone(e *TodoItemDone) error {
	if item := t.findItem(e.TodoItemID); item != nil {
		item.Done = true
	}
	return nil
}
