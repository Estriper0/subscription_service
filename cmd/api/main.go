package main

import (
	"github.com/Estriper0/subscription_service/internal/app"
	"github.com/Estriper0/subscription_service/internal/config"
	"github.com/Estriper0/subscription_service/internal/logger"
)

const configPath = "configs/config.yaml"

func main() {
	config := config.New(configPath)
	logger := logger.GetLogger(config.App.Env)

	app := app.New(logger, config)
	app.Run()
}
