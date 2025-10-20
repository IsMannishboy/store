package main

import (
	i "admin/internal"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func MainHandler(logger *log.Logger, redis_db *redis.Client, db *sql.DB, cnf *i.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		csrf := c.Param("csrf")
		user_id, new_csrf, err := i.CheckCSRF(csrf, []byte(cnf.Secret))
		if err != nil {
			logger.Print("CheckCSRF err:", err)
			c.String(403, err.Error())
			return
		}
		var session_id string
		for _, c := range c.Request.Cookies() {
			if len(c.Name) > len("session_id_") {
				id := c.Name[len("session_id_"):]
				int_id, err := strconv.Atoi(id)
				if err != nil {
					continue
				}
				if int_id == user_id {
					session_id = c.Value
					break
				}
			}
		}
		MainContext := c.Request.Context()
		err, _ = i.CheckSession(MainContext, session_id, redis_db, time.Duration(cnf.Redis.RwTimeout))
		if err != nil {
			logger.Print("CheckSession err:", err)
			c.String(403, err.Error())
			return
		}
		c.HTML(200, "main.html", gin.H{
			"csrf": new_csrf,
		})
	}
}

func main() {
	logger := log.Default()
	cnf, err := i.GetConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Print("config:", cnf)
	db, err := i.GetPostgresConn(cnf)
	if err != nil {
		logger.Fatal(err)
	}
	redis_db, err := i.GetRedisConn(cnf)
	if err != nil {
		logger.Fatal(err)
	}
	router := gin.Default()
	router.LoadHTMLGlob(cnf.HTMLPath + "/main.html")
	router.GET("/main/:csrf", MainHandler(logger, redis_db, db, &cnf))
	http.ListenAndServe(":80", router)
}
