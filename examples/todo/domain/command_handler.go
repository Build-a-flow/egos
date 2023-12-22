package domain

import (
	"context"
	"errors"

	"github.com/finktek/egos"
	"github.com/gofrs/uuid"
)

type TodoCommandHandler struct {
	AggregateStore egos.AggregateStore
}

func (h *TodoCommandHandler) Handle(ctx context.Context, command egos.Command) error {
	switch cmd := command.Command().(type) {
	case *CreateTodoList:
		todo := Init(cmd.Id)
		err := todo.CreateTodoList(cmd.Title)
		if err != nil {
			return err
		}
		return h.AggregateStore.Store(ctx, todo)
	case *AddTodoItem:
		todo := Init(cmd.Id)
		if err := h.AggregateStore.Load(ctx, todo, cmd.Id); err != nil {
			return err
		}
		if err := todo.AddItem(cmd.TodoItemID, cmd.Description); err != nil {
			return err
		}
		if err := todo.AddItem(uuid.Must(uuid.NewV4()).String(), cmd.Description); err != nil {
			return err
		}
		return h.AggregateStore.Store(ctx, todo)
	case *MarkItemAsDone:
		todo := Init(cmd.Id)
		if err := h.AggregateStore.Load(ctx, todo, cmd.Id); err != nil {
			return err
		}
		if err := todo.ItemDone(cmd.TodoItemID); err != nil {
			return err
		}
		return h.AggregateStore.Store(ctx, todo)
	default:
		return errors.New("TodoCommandHandler has received a command that it is does not know how to handle")
	}
}
