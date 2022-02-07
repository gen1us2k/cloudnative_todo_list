package database

import "github.com/gen1us2k/cloudnative_todo_list/models"

type Database interface {
	CreateTodo(models.Todo) (models.Todo, error)
	ListTodos() ([]models.Todo, error)
	UpdateTodo(models.Todo) (models.Todo, error)
	DeleteTodo(models.Todo) error
}
