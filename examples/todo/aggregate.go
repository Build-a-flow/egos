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
	a.AggregateBase.Apply(a, &Created{ID: a.AggregateID()})
	return nil
}

func (a *Todo) SuperChange(count int) error {
	a.AggregateBase.Apply(a, &SuperChange{ID: a.AggregateID(), Count: count})
	return nil
}

func (a *Todo) When(event finkgoes.Event) {
	switch e := event.GetData().(type) {
	case *Created:
		fmt.Println("EVENT ID %s", e.ID)
		a.count = 1
	case *SuperChange:
		fmt.Println("EVENT ID 2 %s", e.ID)
		a.count = e.Count
	}
}

func (a *Todo) Count() int {
	return a.count
}

type Created struct {
	ID   string `json:"id"`
}


type SuperChange struct {
	ID   string `json:"id"`
	Count int  `json:"count"`
}
