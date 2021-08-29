package main

import (
	"fmt"
	"github.com/finktek/eventum"
)

type Todo struct {
	*finkgoes.AggregateBase
	count     int
}

func InitTodo(id string) *Todo {
	return &Todo{
		AggregateBase: finkgoes.NewAggregateBase(id),
	}
}

func (a *Todo) Create() error {
	a.AggregateBase.Apply(a, &Created{ID: a.AggregateID()})
	return nil
}

func (a *Todo) SuperChange(count int) error {
	a.AggregateBase.Apply(a, &SuperChanged{ID: a.AggregateID(), Count: count})
	return nil
}

func (a *Todo) When(event finkgoes.Event) {
	switch e := event.GetData().(type) {
	case *Created:
		fmt.Println("EVENT ID %s", e.ID)
		a.count = 1
	case *SuperChanged:
		fmt.Println("EVENT ID 2 %s", e.ID)
		a.count = e.Count
	}
}

func (a *Todo) GetCount() int {
	return a.count
}

type Created struct {
	ID   string `json:"id"`
}


type SuperChanged struct {
	ID   string `json:"id"`
	Count int  `json:"count"`
}
