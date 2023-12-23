package domain

type CreateTodoList struct {
	UserID string
	Id     string
	Title  string
}

type AddTodoItem struct {
	Id          string
	TodoItemID  string
	Description string
}

type MarkItemAsDone struct {
	Id         string
	TodoItemID string
}
