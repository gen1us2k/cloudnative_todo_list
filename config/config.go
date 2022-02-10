package config

import "github.com/kelseyhightower/envconfig"

const (
	// KratosSessionKey contant stores cookie identifier for Ory Kratos
	KratosSessionKey = "ory_kratos_session"
	// KratosTraitsKey used to pass traits from gRPC gateway to gRPC server
	KratosTraitsKey = CtxKey("kratos_traits")
	// EnvProduction indicates production environment
	EnvProduction = "production"
	// EnvDevelopment indicates development environment
	EnvDevelopment = "development"
)

type (
	// CtxKey is a type that overrides collisions using context package
	CtxKey string
	// AppConfig is a configuration struct for the whole application
	AppConfig struct {
		Env          string `envconfig:"ENV" default:"development"`
		GRPCPort     int    `envconfig:"GRPC_PORT" default:"8080"`
		HTTPPort     int    `envconfig:"HTTP_PORT" default:"8081"`
		SupabaseURL  string `envconfig:"SUPABASE_URL"`
		SupabaseKey  string `envconfig:"SUPABASE_KEY"`
		KratosAPIURL string `envconfig:"KRATOS_URL" default:"http://127.0.0.1:4433"`
		KratosUIURL  string `envconfig:"KRATOS_UI_URL" default:"http://127.0.0.1:4455"`
	}
)

// Parse parses configuration set by environment variables
func Parse() (*AppConfig, error) {
	var c AppConfig
	err := envconfig.Process("", &c)
	return &c, err
}
