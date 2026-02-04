package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	_ "github.com/Estriper0/subscription_service/docs"
	"github.com/Estriper0/subscription_service/internal/config"
	"github.com/Estriper0/subscription_service/internal/handlers"
	"github.com/Estriper0/subscription_service/internal/repository/db"
	"github.com/Estriper0/subscription_service/internal/server"
	"github.com/Estriper0/subscription_service/internal/service"
	"github.com/Estriper0/subscription_service/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	logger *slog.Logger
	config *config.Config
	db     *pgxpool.Pool
	server *server.Server
}

func New(logger *slog.Logger, config *config.Config) *App {
	//Removing Gin logs in production
	if config.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	dbPool, err := postgres.New(config.DB.Url(), config.DB.PoolSize)
	if err != nil {
		panic(err)
	}

	validate := validator.New()
	err = registerCustomValidations(validate)
	if err != nil {
		panic(err)
	}

	subscriptionRepo := db.NewSubscriptionRepo(dbPool)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, logger)
	subscriptionGroup := router.Group("/subscription")

	handlers.NewSubscriptionHandler(subscriptionGroup, subscriptionService, validate)

	server := server.New(router, config)

	return &App{
		logger: logger,
		config: config,
		db:     dbPool,
		server: server,
	}
}

func (a *App) Run() {
	//Closing the connection to the database
	defer a.db.Close()

	a.logger.Info("Start application")

	a.logger.Info(fmt.Sprintf("Starting server on :%d", a.config.Server.Port))
	go a.server.Run()

	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//Catch the application signal or error
	select {
	case q := <-quit:
		a.logger.Info(fmt.Sprintf("Received signal: %s", q.String()))
	case err := <-a.server.Err():
		a.logger.Error(fmt.Sprintf("Server error: %s", err.Error()))
	}
	a.logger.Info("Initiating graceful shutdown...")

	//Graceful shutdown
	err := a.server.Stop()
	if err != nil {
		a.logger.Error("Incorrect server shutdown", slog.String("error", err.Error()))
	} else {
		a.logger.Info("Server shutdown gracefully")
	}
	a.logger.Info("Stop application")
}

func registerCustomValidations(v *validator.Validate) error {
	err := v.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		date := fl.Field().String()
		pattern := `^(0[1-9]|1[0-2])-(19\d{2}|20\d{2})$`
		matched, err := regexp.MatchString(pattern, date)
		return err == nil && matched
	})

	return err
}
