package supabase

import (
	"testing"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	c, err := config.Parse()
	assert.NoError(t, err)
	s, err := NewSupabaseDatabase(c)
	assert.NoError(t, err)
	todo, err := s.CreateTodo(models.Todo{
		Title:   "Test",
		Status:  "active",
		OwnerID: "uiod",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Test", todo.Title)
	assert.NotEqual(t, 0, todo.ID)
}
