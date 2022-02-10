package server

import (
	"context"
	"testing"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/grpc/v1/todolist"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func TestSmokeTest(t *testing.T) {
	// Smoke test to test integration with database
	ownerID := "6e615410-8fe0-403b-8e5a-56674ac5aa92"
	md := metadata.Pairs("user_id", ownerID)
	todo := &todolist.Todo{
		Title:  "A simple test todo item",
		Status: "active",
	}
	reqCtx := metadata.NewIncomingContext(context.TODO(), md)

	c, err := config.Parse()
	assert.NoError(t, err)
	s, err := NewServer(c)
	assert.NoError(t, err)
	todo, err = s.CreateTodo(reqCtx, todo)
	assert.NoError(t, err)
	assert.Equal(t, ownerID, todo.Owner.Id)
	assert.NotEqual(t, 0, todo.Id)

	todo.Status = "closed"
	updated, err := s.UpdateTodo(reqCtx, todo)
	assert.NoError(t, err)
	assert.Equal(t, "closed", updated.Status)

	todos, err := s.ListTodos(reqCtx, &emptypb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, updated, todos.Todos[0])

	_, err = s.DeleteTodo(reqCtx, updated)
	assert.NoError(t, err)
}

func TestWithoutMetadata(t *testing.T) {
	// Test that API endpoints return error
	// once we have no userId
	reqCtx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs())
	c, err := config.Parse()
	assert.NoError(t, err)
	s, err := NewServer(c)
	assert.NoError(t, err)
	todo := &todolist.Todo{
		Title:  "A simple test todo item",
		Status: "active",
	}

	_, err = s.CreateTodo(reqCtx, todo)
	assert.Error(t, err)
	_, err = s.ListTodos(reqCtx, &emptypb.Empty{})
	assert.Error(t, err)
	_, err = s.UpdateTodo(reqCtx, todo)
	assert.Error(t, err)
	_, err = s.DeleteTodo(reqCtx, todo)
	assert.Error(t, err)
}
