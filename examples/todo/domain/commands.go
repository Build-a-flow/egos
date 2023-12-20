package domain

import "github.com/google/uuid"

type CreateTodoList struct {
	Id    uuid.UUID
	Title string
}

type AddTodoItem struct {
	Id          uuid.UUID
	TodoItemID  uuid.UUID
	Description string
}

type MarkItemAsDone struct {
	Id         uuid.UUID
	TodoItemID uuid.UUID
}
