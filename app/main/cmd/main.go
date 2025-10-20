package main

import (
	c "gin/internal/config"
	handler "gin/internal/handlers"
	m "gin/internal/migrations"
	d "gin/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"log/slog"
	"os"
)

const (
	EnvDev   = "dev"
	EnvProd  = "prod"
	EnvLocal = "local"
)

// todo: script service
// todo: auth service
func main() {
	// Initialize configuration
	cnf := c.GetConfig()

	// Setup logger based on environment
	logger := setupLogger(cnf.Env)
	logger.Info("Logger initialized", slog.String("env", cnf.Env))
	logger.Debug("Debugging information")
	//postgres connection
	storage, err := d.NewPostgresDb(cnf, logger)
	if err != nil {
		logger.Error("Failed to initialize storage", slog.String("error", err.Error()))

	}
	//cash inticialiation
	cash, err := d.NewRedisDb(cnf)
	if err != nil {
		logger.Debug("redis connection error:", err.Error())
	}
	//migrations
	migrator := m.MustGetNewMigrator()
	if err := migrator.ApplyMigrations(storage.DB, logger); err != nil {
		logger.Error("Failed to apply migrations", slog.String("error", err.Error()))
	}

	router := gin.Default()
	main_path := cnf.TemplatePath + "/main.html"
	router.LoadHTMLGlob(main_path)

	router.GET("/main/:csrf", handler.Mainhendler(logger, storage.DB, cash.Redis_db, cnf))
	http.ListenAndServe(":80", router)
}
func setupLogger(env string) *slog.Logger {

	var log *slog.Logger
	switch env {
	case EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
