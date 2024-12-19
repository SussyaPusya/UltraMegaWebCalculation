package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg"
	"github.com/gorilla/mux"
)

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

type App struct {
	config *Config
}

func New() *App {
	return &App{
		config: ConfigWhitEnv(),
	}
}

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

// Мидлвaря няшка
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
			w.WriteHeader(http.StatusUnprocessableEntity)
			var vid = map[string]string{"error": err.Error()}

			_ = json.NewEncoder(w).Encode(vid)

		} else {
			log.Println("error: Internal server error", err)

			w.WriteHeader(http.StatusInternalServerError)
			var vid = map[string]string{"error": "Internal server error"}

			_ = json.NewEncoder(w).Encode(vid)

		}

	} else {
		w.WriteHeader(http.StatusOK)
		var vid = map[string]string{"result": strconv.Itoa(int(result))}

		_ = json.NewEncoder(w).Encode(vid)

		log.Printf("result: %f", result)
	}

}

func (a *App) RunServer() error {
	mux := mux.NewRouter()

	// jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	defaultLogger := slog.New(slog.Default().Handler())

	mux.Use(LoggingMiddleware(defaultLogger))

	mux.HandleFunc("/api/v1/calculate", CalcHandler)
	defaultLogger.Info("Server has started")
	err := http.ListenAndServe(":"+a.config.Path, mux)
	if err != nil {
		defaultLogger.Error("server has crashed")
		return err
	}
	return nil
}
