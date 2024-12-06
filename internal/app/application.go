package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg"
	"github.com/gorilla/mux"
)

// Конфиг для приложения
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

// Конструктор приложения
func New() *App {
	return &App{
		config: ConfigWhitEnv(),
	}
}

// Функция для запуска приложения в консоле
func (a *App) Run() error {
	for {
		reader := bufio.NewReader(os.Stdin)

		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("invalid expression from in the console")
		}

		text = strings.TrimSpace(text)

		if text == "exit" {
			log.Println("app succsefully close")
			return nil
		}

		result, err := pkg.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}

	}

}

type JsonReq struct {
	Expression string
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	req := new(JsonReq)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	result, err := pkg.Calc(req.Expression)

	if err != nil {
		if errors.Is(err, pkg.ErrInvalidExpr) {
			fmt.Fprintf(w, "err: %s", err.Error())
		} else {
			fmt.Fprintf(w, "unknown err")
		}

	} else {
		fmt.Fprintf(w, "result: %f", result)
	}

}

// Функция для запуска приложения в вебе
func (a *App) RunServer() error {
	mux := mux.NewRouter()

	mux.HandleFunc("/", CalcHandler)

	err := http.ListenAndServe(":"+a.config.Path, mux)
	if err != nil {
		return err
	}
	return nil
}
