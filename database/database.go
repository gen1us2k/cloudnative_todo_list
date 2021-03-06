package database

import (
	"strconv"

	"github.com/gen1us2k/cloudnative_todo_list/models"
)

const base = 10

// Database interface
// To add a new cloud database solution you can simply implement this interface and use it
//
type Database interface {
	CreateTodo(models.Todo) (models.Todo, error)
	ListTodos(string) ([]models.Todo, error)
	UpdateTodo(models.Todo) (models.Todo, error)
	DeleteTodo(models.Todo) error
}

// FormatInt64 formats int to string
func FormatInt64(id int64) string {
	return strconv.FormatInt(id, base)
}
