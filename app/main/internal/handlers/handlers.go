package handlers

import (
	"database/sql"
	c "gin/internal/config"
	f "gin/internal/funcs"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Mainhendler(logger *slog.Logger, db *sql.DB, cnf *c.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Main handler called")
		ctx := c.Request.Context()
		err, data := f.GetMainPageData(db, ctx, cnf.Postgres.RwTimeout)
		if err != nil && err != sql.ErrNoRows {
			logger.Error("Internal Server Error", slog.String("error", err.Error()))
			c.String(500, "Internal Server Error")
			return
		}
		c.HTML(200, "main.html", gin.H{
			"title":      data.Title,
			"Products":   data.Products,
			"Categories": data.Categories,
		})
	}
}
