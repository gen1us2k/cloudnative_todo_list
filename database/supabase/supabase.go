package supabase

import (
	"fmt"
	"strconv"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/database"
	"github.com/gen1us2k/cloudnative_todo_list/models"
	"github.com/supabase/postgrest-go"
)

type Supabase struct {
	database.Database
	conn   *postgrest.Client
	config *config.AppConfig
}

func NewSupabaseDatabase(c *config.AppConfig) (*Supabase, error) {
	conn := postgrest.NewClient(c.SupabaseURL, "", map[string]string{
		"apikey":        c.SupabaseKey,
		"Authorization": fmt.Sprintf("Bearer %s", c.SupabaseKey),
	})
	if conn.ClientError != nil {
		return nil, conn.ClientError
	}
	return &Supabase{
		config: c,
		conn:   conn,
	}, nil

}

func (s *Supabase) CreateTodo(todo models.Todo) (models.Todo, error) {
	var todos []models.Todo
	q := s.conn.From("Todos").Insert(todo, false, "do-nothing", "", "")
	_, err := q.ExecuteTo(&todos)
	return todos[0], err
}
func (s *Supabase) ListTodos() ([]models.Todo, error) {
	var todos []models.Todo
	q := s.conn.From("Todos").Select("*", "10", false)
	_, err := q.ExecuteTo(&todos)
	return todos, err

}
func (s *Supabase) UpdateTodo(todo models.Todo) (models.Todo, error) {
	var todos []models.Todo
	q := s.conn.From("Todos").Update(todo, "", "")
	_, err := q.ExecuteTo(&todos)
	return todos[0], err

}
func (s *Supabase) DeleteTodo(todo models.Todo) error {
	q := s.conn.From("Todos").Delete("", "").Match(map[string]string{"id": strconv.FormatInt(todo.ID, 10)})
	_, _, err := q.Execute()

	return err
}
