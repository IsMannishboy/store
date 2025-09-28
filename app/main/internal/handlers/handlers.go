package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Mainhendler(logger *slog.Logger, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Main handler called")
		c.HTML(200, "main.html", gin.H{
			"title": "Main website",
		})
	}
}
