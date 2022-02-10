package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/gen1us2k/cloudnative_todo_list/config"
	"github.com/gen1us2k/cloudnative_todo_list/database"
	"github.com/gen1us2k/cloudnative_todo_list/database/supabase"
	"github.com/gen1us2k/cloudnative_todo_list/grpc/v1/todolist"
	"github.com/gen1us2k/cloudnative_todo_list/middleware"
	"github.com/gen1us2k/cloudnative_todo_list/models"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Server implements gRPC and http APIs
type Server struct {
	db       database.Database
	config   *config.AppConfig
	errGroup *errgroup.Group
}

// NewServer configures server
func NewServer(c *config.AppConfig) (*Server, error) {
	db, err := supabase.NewSupabaseDatabase(c)
	if err != nil {
		return nil, err
	}
	return &Server{db: db, config: c, errGroup: new(errgroup.Group)}, nil
}

// CreateTodo API
func (s *Server) CreateTodo(ctx context.Context, todo *todolist.Todo) (*todolist.Todo, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	todo.Owner = &todolist.User{Id: userID}
	t := models.NewTodoFromPB(todo)
	t, err = s.db.CreateTodo(t)
	if err != nil {
		return nil, err
	}
	return t.ToProto(), nil
}

// ListTodos returns todos created by authenticated user
func (s *Server) ListTodos(ctx context.Context, e *emptypb.Empty) (*todolist.TodoListResponse, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	todos, err := s.db.ListTodos(userID)
	if err != nil {
		return nil, err
	}
	res := &todolist.TodoListResponse{}
	for _, todo := range todos {
		res.Todos = append(res.Todos, todo.ToProto())
	}
	return res, nil
}

// UpdateTodo updatesTodo
func (s *Server) UpdateTodo(ctx context.Context, todo *todolist.Todo) (*todolist.Todo, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	todo.Owner = &todolist.User{Id: userID}
	t := models.NewTodoFromPB(todo)
	t, err = s.db.UpdateTodo(t)
	if err != nil {
		return nil, err
	}
	return t.ToProto(), nil
}

// DeleteTodo deletes todo
func (s *Server) DeleteTodo(ctx context.Context, todo *todolist.Todo) (*todolist.DeleteResponse, error) {
	userID, err := s.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	todo.Owner = &todolist.User{Id: userID}
	if err := s.db.DeleteTodo(models.NewTodoFromPB(todo)); err != nil {
		return nil, err
	}

	return &todolist.DeleteResponse{Status: "success"}, nil
}
func (s *Server) startHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	gwmux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(middleware.GatewayResponseModifier),
		runtime.WithMetadata(middleware.GatewayMetadataAnnotator),
	)
	kratos := middleware.KratosMiddleware{
		APIURL: s.config.KratosAPIURL,
		UIURL:  s.config.KratosUIURL,
		Client: &http.Client{},
	}

	r := mux.NewRouter()
	r.Use(kratos.Middleware)
	r.PathPrefix("/").Handler(gwmux)
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("localhost:%d", s.config.GRPCPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	err = todolist.RegisterTodolistAPIServiceHandler(ctx, gwmux, conn)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.HTTPPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	return srv.ListenAndServe()
}
func (s *Server) startGRPC() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.config.GRPCPort))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	todolist.RegisterTodolistAPIServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
}

// Start starts both REST and gRPC services
func (s *Server) Start() {
	s.errGroup.Go(s.startGRPC)
	s.errGroup.Go(s.startHTTP)
}

// Wait for it. Just wait
func (s *Server) Wait() error {
	return s.errGroup.Wait()
}
func (s *Server) getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no session found in metadata")
	}
	data := md.Get("user_id")
	if len(data) == 0 {
		return "", errors.New("no user_id found in context metadata")
	}
	return data[0], nil
}
