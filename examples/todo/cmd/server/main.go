package main

import (
	"context"
	"encoding/json"
	"examples/todo/domain"
	"log"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	egos "github.com/finktek/egos/core"
	es "github.com/finktek/egos/esdb"
	"github.com/google/uuid"
)

func main() {
	eventStoreDbConfig, err := esdb.ParseConnectionString("esdb+discover://localhost:2113?tls=false")
	if err != nil {
		panic(err)
	}
	eventStoreDbConfig.SkipCertificateVerification = true

	client, err := esdb.NewClient(eventStoreDbConfig)
	if err != nil {
		panic(err)
	}

	store := es.NewEsdbEventStore(client)

	// Create an in-memory aggregate store
	aggregateStore, err := egos.NewAggregateStore(store, &domain.TodoList{})
	if err != nil {
		panic(err)
	}

	commandHandler := &domain.TodoCommandHandler{AggregateStore: aggregateStore}

	todoListID := uuid.New().String()
	cmd := egos.NewCommand(&domain.CreateTodoList{UserID: "user-g12g3g3h1", Id: todoListID, Title: "My Todo"})
	err = commandHandler.Handle(context.Background(), cmd)
	if err != nil {
		panic(err)
	}

	todoItemID := uuid.New().String()
	cmd2 := egos.NewCommand(&domain.AddTodoItem{Id: todoListID, TodoItemID: todoItemID, Description: "Do something good"})
	err = commandHandler.Handle(context.Background(), cmd2)
	if err != nil {
		panic(err)
	}

	todoItemID2 := uuid.New().String()
	cmd3 := egos.NewCommand(&domain.AddTodoItem{Id: todoListID, TodoItemID: todoItemID2, Description: "Do nothing for the rest of the day"})
	err = commandHandler.Handle(context.Background(), cmd3)
	if err != nil {
		panic(err)
	}

	todoItemID3 := uuid.New().String()
	cmd4 := egos.NewCommand(&domain.AddTodoItem{Id: todoListID, TodoItemID: todoItemID3, Description: "Sleep"})
	err = commandHandler.Handle(context.Background(), cmd4)
	if err != nil {
		panic(err)
	}

	cmd5 := egos.NewCommand(&domain.MarkItemAsDone{Id: todoListID, TodoItemID: todoItemID})
	err = commandHandler.Handle(context.Background(), cmd5)
	if err != nil {
		panic(err)
	}

	cmd6 := egos.NewCommand(&domain.MarkItemAsDone{Id: todoListID, TodoItemID: todoItemID2})
	err = commandHandler.Handle(context.Background(), cmd6)
	if err != nil {
		panic(err)
	}

	todo := domain.Init(todoListID)
	if err := aggregateStore.Load(context.Background(), &todo, todoListID); err != nil {
		log.Println("error loading todo list: ", err)
	}

	todoData, _ := json.Marshal(&todo)

	log.Println("TODO ", string(todoData))
}