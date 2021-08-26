package main

import (
	"fmt"
	"github.com/finktek/eventum"
)

const TodoAggregateType finkgoes.AggregateType = "todo"

type Todo struct {
	*finkgoes.AggregateBase
	count     int
}

func NewTodo(id string) *Todo {
	return &Todo{
		AggregateBase: finkgoes.NewAggregateBase(id),
	}
}

func (a *Todo) Create() error {
	a.AggregateBase.AppendEvent(a, &Created{ID: a.AggregateID()})
	return nil
}

func (a *Todo) ApplyEvent(event finkgoes.Event) error {
	switch e := event.Event().(type) {
	case *Created:
		fmt.Println("EVENT ID %s", e.ID)
		a.count = 1
	}

	return nil
}
func (a *Todo) Count() int {
	return a.count
}

type Created struct {
	ID   string
}
