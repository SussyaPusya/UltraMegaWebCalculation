package app

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	Path string
}

func ConfigWhitEnv() *Config {
	cfg := new(Config)

	cfg.Path = os.Getenv("PORT")
	if cfg.Path == "" {
		cfg.Path = "8080"

	}
	return cfg
}

// Application struct
type App struct {
	config *Config
}

func New() *App {
	return &App{
		config: ConfigWhitEnv(),
	}
}

func (a *App) Run() error {

}

func (a *App) RunServer() error {
	mux := mux.NewRouter()

	http.ListenAndServe(":"+a.config.Path, mux)
}
