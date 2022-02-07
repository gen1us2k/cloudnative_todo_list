package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	Env         string `envconfig:"ENV" default:"development"`
	GRPCPort    int    `envconfig:"GRPC_PORT" default:"8080"`
	HTTPPort    int    `envconfig:"HTTP_PORT" default:"8081"`
	SupabaseURL string `envconfig:"SUPABASE_URL"`
	SupabaseKey string `envconfig:"SUPABASE_KEY"`
}

func Parse() (*AppConfig, error) {
	var c AppConfig
	err := envconfig.Process("", &c)
	return &c, err
}
