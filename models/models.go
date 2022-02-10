package models

import "github.com/gen1us2k/cloudnative_todo_list/grpc/v1/todolist"

type (
	// Todo model represents Todo
	Todo struct {
		ID      int64  `json:"id,omitempty"`
		Title   string `json:"title"`
		Status  string `json:"status"`
		OwnerID string `json:"owner_id"`
	}
	// User model
	User struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

// NewTodoFromPB converts data from protocolbuffers
// to Todo model
func NewTodoFromPB(todo *todolist.Todo) Todo {
	return Todo{
		ID:      todo.Id,
		Title:   todo.Title,
		Status:  todo.Status,
		OwnerID: todo.Owner.Id,
	}
}

// ToProto converts Todo to todolist.Todo
func (t Todo) ToProto() *todolist.Todo {
	return &todolist.Todo{
		Id:     t.ID,
		Title:  t.Title,
		Status: t.Status,
		Owner: &todolist.User{
			Id: t.OwnerID,
		},
	}
}
