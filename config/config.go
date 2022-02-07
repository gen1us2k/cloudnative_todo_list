package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	Env         string `envconfig:"ENV" default:"development"`
	SupabaseURL string `envconfig:"SUPABASE_URL"`
	SupabaseKey string `envconfig:"SUPABASE_KEY"`
}

func Parse() (*AppConfig, error) {
	var c AppConfig
	err := envconfig.Process("", &c)
	return &c, err
}
