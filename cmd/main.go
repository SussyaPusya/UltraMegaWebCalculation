package main

import (
	"context"
	"log/slog"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"
	httpapi "github.com/SussyaPusya/UltraMegaWebCalculation/internal/transport/httpApi"
)

func main() {

	ctx := context.Background()

	logger := slog.New(slog.Default().Handler())

	ctx = context.WithValue(ctx, httpapi.Key, logger)

	cfg := config.NewConfig()

	handler := httpapi.NewHandler()

	router := httpapi.NewRouter(ctx, handler, cfg)

	router.Run()
}
