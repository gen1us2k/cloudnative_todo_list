package supabase

import (
	"fmt"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/database"
	"github.com/gen1us2k/cloudnative_todo_list/models"
	"github.com/supabase/postgrest-go"
)

// Supabase implements communication protocol with subabase.io/database
type Supabase struct {
	database.Database
	conn   *postgrest.Client
	config *config.AppConfig
}

// NewSupabaseDatabase creates Subabase database adapter by given configuration
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

// CreateTodo creates todo
func (s *Supabase) CreateTodo(todo models.Todo) (models.Todo, error) {
	var todos []models.Todo
	q := s.conn.From("Todos").Insert(todo, false, "do-nothing", "", "")
	_, err := q.ExecuteTo(&todos)
	return todos[0], err
}

// ListTodos returns todos for specifies user
func (s *Supabase) ListTodos(userID string) ([]models.Todo, error) {
	var todos []models.Todo
	q := s.conn.From("Todos").Select("*", "10", false).Match(map[string]string{"owner_id": userID})
	_, err := q.ExecuteTo(&todos)
	return todos, err
}

// UpdateTodo updates todo
func (s *Supabase) UpdateTodo(todo models.Todo) (models.Todo, error) {
	var todos []models.Todo
	q := s.conn.From("Todos").Update(todo, "", "").Match(map[string]string{"id": database.FormatInt64(todo.ID)})
	_, err := q.ExecuteTo(&todos)
	return todos[0], err
}

// DeleteTodo deletes Todo
func (s *Supabase) DeleteTodo(todo models.Todo) error {
	q := s.conn.From("Todos").Delete("", "").Match(map[string]string{"id": database.FormatInt64(todo.ID)})
	_, _, err := q.Execute()

	return err
}
