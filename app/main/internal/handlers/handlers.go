package handlers

import (
	"database/sql"
	"fmt"
	a "gin/internal/auth"
	c "gin/internal/config"
	f "gin/internal/funcs"
	s "gin/internal/structs"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Mainhendler(logger *slog.Logger, db *sql.DB, redis_db *redis.Client, cnf *c.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Main handler called")
		Maincontext := c.Request.Context()
		csrf := c.Param("csrf")
		user_id, new_csrf, err := a.VerifyCSRF(cnf.Secret, csrf)
		if err != nil {
			logger.Debug("VerifyCSRF:", err.Error())
			c.String(403, err.Error())
			return
		}
		var session_id = ""
		for _, ck := range c.Request.Cookies() {
			if len(ck.Name) > len("session_id_") {
				cookie_id := ck.Name[len("session_id_"):]
				fmt.Println("cookie_id:", cookie_id)
				int_id, err := strconv.Atoi(cookie_id)

				if err != nil {
					fmt.Println("err strconv:", err)
					continue
				}
				fmt.Println("int_cookie_id:", int_id)

				if int_id == user_id {
					session_id = ck.Value
					fmt.Println("id from cookie:", int_id, " and from csrf:", user_id, " matches")
					break
				}
			}
		}

		if session_id == "" {
			logger.Debug("err while getting session id")
			c.String(403, fmt.Errorf("err while getting session id").Error())
			return
		}
		decodedValue, _ := url.QueryUnescape(session_id)
		fmt.Println("session from cookie:", decodedValue)
		err, id := a.CheckSession(decodedValue, redis_db, cnf.Redis.RwTimeout, Maincontext)
		if err != nil {
			logger.Debug("CheckSession err:", err)
			c.String(500, err.Error())
			return
		}
		var data s.MainpageData
		err, data = f.GetMainPageData(id, db, Maincontext, cnf.Postgres.RwTimeout)
		if err != nil && err != sql.ErrNoRows {
			c.String(500, err.Error())
			return
		}
		c.HTML(200, "main.html", gin.H{
			"csrf":       new_csrf,
			"cart":       data.Cart,
			"title":      data.Title,
			"Products":   data.Products,
			"Categories": data.Categories,
		})
	}
}
