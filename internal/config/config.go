package config

import "os"

type Config struct {
	Path string
}

func NewConfig() *Config {
	cfg := new(Config)

	cfg.Path = os.Getenv("PORT")
	if cfg.Path == "" {
		cfg.Path = "8888"

	}
	return cfg
}
