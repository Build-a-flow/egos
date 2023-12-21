package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/finktek/egos"
	"github.com/finktek/egos/inmem"
	"github.com/google/uuid"
)

func main() {
	// Create an in-memory event store
	store := inmem.NewInMemEventStore()

	// Create an in-memory aggregate store
	aggregateStore, err := egos.NewAggregateStore(store, &TodoList{})
	if err != nil {
		panic(err)
	}

	commandHandler := &TodoCommandHandler{AggregateStore: aggregateStore}

	todoListID := uuid.New().String()
	cmd := egos.NewCommand(&CreateTodoList{Id: todoListID, Title: "My Todo"})
	err = commandHandler.Handle(context.Background(), cmd)
	if err != nil {
		panic(err)
	}

	todoItemID := uuid.New().String()
	cmd2 := egos.NewCommand(&AddTodoItem{Id: todoListID, TodoItemID: todoItemID, Description: "Do something good"})
	err = commandHandler.Handle(context.Background(), cmd2)
	if err != nil {
		panic(err)
	}

	todoItemID2 := uuid.New().String()
	cmd3 := egos.NewCommand(&AddTodoItem{Id: todoListID, TodoItemID: todoItemID2, Description: "Do nothing for the rest of the day"})
	err = commandHandler.Handle(context.Background(), cmd3)
	if err != nil {
		panic(err)
	}

	todoItemID3 := uuid.New().String()
	cmd4 := egos.NewCommand(&AddTodoItem{Id: todoListID, TodoItemID: todoItemID3, Description: "Sleep"})
	err = commandHandler.Handle(context.Background(), cmd4)
	if err != nil {
		panic(err)
	}

	cmd5 := egos.NewCommand(&MarkItemAsDone{Id: todoListID, TodoItemID: todoItemID})
	err = commandHandler.Handle(context.Background(), cmd5)
	if err != nil {
		panic(err)
	}

	cmd6 := egos.NewCommand(&MarkItemAsDone{Id: todoListID, TodoItemID: todoItemID2})
	err = commandHandler.Handle(context.Background(), cmd6)
	if err != nil {
		panic(err)
	}

	todo := Init(todoListID)
	if err := aggregateStore.Load(context.Background(), todo, todoListID); err != nil {
		log.Println("error loading todo list: ", err)
	}

	todoData, _ := json.Marshal(&todo)

	log.Println("TODO ", string(todoData))

}
