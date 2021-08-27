package main

import (
	"context"
	"fmt"
	"github.com/finktek/eventum"
	"github.com/finktek/eventum/eventstore"
		"github.com/gofrs/uuid"
	"log"
)

func main()  {
	eventStore, _ := eventstore.NewEventStoreDbClient("esdb://localhost:2113?tls=false")
	todoAggregateStore, _:= finkgoes.NewAggregateStore(eventStore, TodoAggregateType)

	commandHandler := &TestCommandHandler{
		AggregateStore: todoAggregateStore,
	}
	cmdId, _ := uuid.FromString("b80abf6a-d2a0-4303-b84d-3ebce1edab23")
	cmd := finkgoes.NewCommand(&TestCommand{Id: cmdId})
	err := commandHandler.Handle(context.Background(), cmd)
	if err != nil {
		fmt.Println(err)
	}
}

type TestCommandHandler struct {
	AggregateStore finkgoes.AggregateStore
}

func (h *TestCommandHandler) Handle(ctx context.Context, command finkgoes.Command) error {
	switch cmd := command.Command().(type) {
	case *CreateCommand:
		todo := NewTodo(cmd.Id.String())
		todo.Create()
		todo.SuperChange(22)
		return h.AggregateStore.Store(ctx, todo)
	case *TestCommand:
		aa, _ := h.AggregateStore.Load(ctx, cmd.Id.String())
		fmt.Println("I AM TEST COMMMAND")
		log.Println(aa.AggregateID())
	default:
		log.Fatalf("TestCommandHandlers has received a command that it is does not know how to handle, %#v", cmd)
	}

	return nil
}


type TestCommand struct {
	Id uuid.UUID
}

type CreateCommand struct {
	Id uuid.UUID
}
