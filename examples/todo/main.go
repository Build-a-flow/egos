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
	cmdId, _ := uuid.NewV4()
	cmd := finkgoes.NewCommand("string-id", &CreateCommand{Id: cmdId})
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
		fmt.Println(cmd.Id )
		todo := NewTodo(cmd.Id.String())
		fmt.Println(todo.Count())
		todo.Create()
		fmt.Println(todo.Count())
		return h.AggregateStore.Store(ctx, todo)
	case *TestCommand:
		fmt.Println("I AM TEST COMMMAND")
	default:
		log.Fatalf("TestCommandHandlers has received a command that it is does not know how to handle, %#v", cmd)
	}

	return nil
}


type TestCommand struct {
}

type CreateCommand struct {
	Id uuid.UUID
}
