package main

import (
	"database/sql"
	"fmt"
	c "gin/internal/config"
	handler "gin/internal/handlers"
	m "gin/internal/migrations"
	d "gin/internal/storage"

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

var db *sql.DB

// todo: script service
// todo: auth service
func main() {
	// Initialize configuration
	cnf := c.GetConfig()
	fmt.Println(cnf)
	// Setup logger based on environment
	logger := setupLogger(cnf.Env)
	logger.Info("Logger initialized", slog.String("env", cnf.Env))
	logger.Debug("Debugging information")
	storage, err := d.New(db, cnf)
	if err != nil {
		logger.Error("Failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	migrator := m.MustGetNewMigrator()
	if err := migrator.ApplyMigrations(storage.DB); err != nil {
		logger.Error("Failed to apply migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
	router := gin.Default()
	router.LoadHTMLGlob("D:/gin/app/main/internal/templates/html/*")

	router.GET("/main", handler.Mainhendler(logger, storage.DB))
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
