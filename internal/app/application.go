package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

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
		cfg.Path = "8888"

	}
	return cfg
}

// Application struct
type App struct {
	config *Config
}

// Stupid Constructor
func New() *App {
	return &App{
		config: ConfigWhitEnv(),
	}
}

// Run app in the console
func (a *App) Run() error {
	for {
		log.Println("Enter expression")
		reader := bufio.NewReader(os.Stdin)

		text, err := reader.ReadString('\n')
		if err != nil {

			log.Println(pkg.ErrInvalidExpr, err)
		}

		text = strings.TrimSpace(text)

		if text == "exit" {
			log.Println("app succsefully close")
			return nil
		}

		result, err := pkg.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed with error: ", err)
		} else {
			log.Println(text, "=", result)
		}

	}

}

// Мидлвaря няшка* Middllware for mux
func LoggingMiddleware(DefaLogger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			DefaLogger.Info("finished", slog.Group("request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String())),
				slog.Int("status", http.StatusOK),
				slog.Duration("duration", time.Second),
			)
			next.ServeHTTP(w, r)

		})
	}
}

type JsonReq struct {
	Expression string `json:"expression"`
}

// Handler for calc func
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	req := new(JsonReq)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		http.Error(w, pkg.ErrIvalJson.Error(), http.StatusBadRequest)
		log.Println(pkg.ErrIvalJson.Error()+":", req, "error:", err)
		return
	}
	log.Println("expression received:", req)
	result, err := pkg.Calc(req.Expression)

	if err != nil {
		if errors.Is(err, pkg.ErrInvalidExpr) {
			log.Println("Invalid expression on the calculating")
			fmt.Fprintf(w, "err: %s", err.Error())
			http.Error(w, "", http.StatusNotAcceptable)
		} else {
			log.Println("unknown err", err)
			fmt.Fprintf(w, "unknown err")
		}

	} else {
		fmt.Fprintf(w, "result: %f", result)
		log.Printf("result: %f", result)
	}

}

// Run app in web
func (a *App) RunServer() error {
	mux := mux.NewRouter()

	// jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	defaultLogger := slog.New(slog.Default().Handler())

	mux.Use(LoggingMiddleware(defaultLogger))

	mux.HandleFunc("/", CalcHandler)
	defaultLogger.Info("Server has started")
	err := http.ListenAndServe(":"+a.config.Path, mux)
	if err != nil {
		defaultLogger.Error("server has crashed")
		return err
	}
	return nil
}
