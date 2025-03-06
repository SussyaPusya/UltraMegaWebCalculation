package main

import (
	"context"
	"time"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/agent"
	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/api"
	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/config"
	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/queue"
	"github.com/SussyaPusya/UltraMegaWebCalculation/pkg/logger"
)

func main() {

	ctx := context.Background()

	ctx, _ = logger.New(ctx)

	cfg, _ := config.NewConfig()

	agent := agent.NewAgent(cfg)

	service := queue.NewService(cfg.EnvConfig)

	handler := api.NewOrchestrator(cfg.EnvConfig, service)

	router := api.NewRouter(ctx, handler, cfg)

	go router.Run()
	time.Sleep(1 * time.Second)
	go agent.Run(ctx)

	select {}

}
