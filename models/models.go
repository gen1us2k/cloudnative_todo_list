package models

import "github.com/gen1us2k/cloudnative_todo_list/grpc/v1/todolist"

type (
	Todo struct {
		ID      int64  `json:"id,omitempty"`
		Title   string `json:"title"`
		Status  string `json:"status"`
		OwnerID string `json:"owner_id"`
	}

	User struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func NewTodoFromPB(todo *todolist.Todo) Todo {
	return Todo{
		Title:   todo.Title,
		Status:  todo.Status,
		OwnerID: todo.Owner.Id,
	}
}

func (t Todo) ToProto() *todolist.Todo {
	return &todolist.Todo{
		Id:     t.ID,
		Title:  t.Title,
		Status: t.Status,
	}
}
