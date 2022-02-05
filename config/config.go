package config

type AppConfig struct {
	Env string `envconfig:"ENV" default:"development"`
}

func Parse() (*AppConfig, error) {
	var c AppConfig
	err := envconfig.Parse("", &c)
	return &c, err
}
