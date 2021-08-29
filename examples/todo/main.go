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
	todoAggregateStore, _:= finkgoes.NewAggregateStore(eventStore, &Todo{})

	commandHandler := &TestCommandHandler{
		AggregateStore: todoAggregateStore,
	}
	cmdId := uuid.Must(uuid.NewV4())
	cmd := finkgoes.NewCommand(&CreateCommand{Id: cmdId})

	err := commandHandler.Handle(context.Background(), cmd)
	if err != nil {
		fmt.Println(err)
	}
	cmd2 := finkgoes.NewCommand(&TestCommand{Id: cmdId})
	err = commandHandler.Handle(context.Background(), cmd2)
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
		todo := InitTodo(cmd.Id.String())
		todo.Create()
		todo.SuperChange(22)
		log.Println(todo.GetCount())
		return h.AggregateStore.Store(ctx, todo)
	case *TestCommand:
		todo := InitTodo(cmd.Id.String())
		h.AggregateStore.Load(ctx, todo, cmd.Id.String())
		log.Println(todo.GetCount())
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

var _ = finkgoes.RegisterEvent(Created{})
var _ = finkgoes.RegisterEvent(SuperChanged{})
