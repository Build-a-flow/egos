package domain

import (
	"context"
	finkgoes "github.com/finktek/eventum"
	"log"
)

type TodoCommandHandler struct {
	AggregateStore finkgoes.AggregateStore
}

func (h *TodoCommandHandler) Handle(ctx context.Context, command finkgoes.Command) error {
	switch cmd := command.Command().(type) {
	case *CreateTodoList:
		todo := InitTodoList(cmd.Id.String())
		err := todo.CreateTodoList(cmd.Title)
		if err != nil {
			log.Println("error creating todo list: ", err)
			return err
		}
		return h.AggregateStore.Store(ctx, todo)
	case *AddTodoItem:
		todo := InitTodoList(cmd.Id.String())
		if err := h.AggregateStore.Load(ctx, todo, cmd.Id.String()); err != nil {
			log.Println("error loading todo list: ", err)
		}
		if err := todo.AddItem(cmd.TodoItemID, cmd.Description); err != nil {
			log.Println("error adding todo item into list: ", err)
		}
		return h.AggregateStore.Store(ctx, todo)
	case *MarkItemAsDone:
		todo := InitTodoList(cmd.Id.String())
		if err := h.AggregateStore.Load(ctx, todo, cmd.Id.String()); err != nil {
			log.Println("error loading todo list: ", err)
		}
		if err := todo.ItemDone(cmd.TodoItemID); err != nil {
			log.Println("error marking todo item as Done: ", err)
		}
		return h.AggregateStore.Store(ctx, todo)
	default:
		log.Fatalf("TodoCommandHandler has received a command that it is does not know how to handle, %#v", cmd)
	}

	return nil
}
