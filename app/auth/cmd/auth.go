package main

import (
	i "auth/internal"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

var db *sql.DB
var redis_db *redis.Client
var cnf *i.Config
var logger *log.Logger

func LoginHandler(db *sql.DB, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var login i.Login
		if err := c.ShouldBindBodyWithJSON(&login); err != nil {
			c.Error(fmt.Errorf("internal server error"))
			return
		}

		MainContext := c.Request.Context()
		cookie_id, _ := c.Cookie("session_id")
		id, err := i.GetUserId(db, MainContext, login.Username, cnf.Postgres.RwTimeout)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		hash, err := i.GetPass(db, MainContext, login.Username, cnf.Postgres.RwTimeout)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		err = i.CheckPassword(login.Password, hash)
		if err != nil {
			c.String(403, "wrong password")
			return
		}
		stored_session, err := i.GetSession(redis_db, cookie_id, MainContext, cnf.Redis.RwTimeout)
		var session i.Session
		if err != nil {
			if err == redis.Nil {
				session, err = i.CreateSession(id, redis_db)
				if err != nil {
					logger.Print(err.Error())
					c.String(500, err.Error())
					return
				}
				err = i.AddSessionToCash(MainContext, session, redis_db, cnf.Redis.RwTimeout)
				if err != nil {
					logger.Print(err.Error())
					c.String(500, err.Error())
					return
				}
			} else {
				logger.Print(err.Error())
				c.String(500, "interna server error")
				return
			}

		} else {

			session.Id = cookie_id
			session.Exp = stored_session.Exp
			session.UserId = stored_session.UserId
			err = i.AddSessionToCash(MainContext, session, redis_db, cnf.Redis.RwTimeout)
			if err != nil {
				logger.Print(err.Error())
				c.String(500, "internal server error")
				return
			}
		}
		token, err := i.CreateCSRF(cnf.Secret, session.UserId)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		c.SetCookie("session_id", session.Id, 1800, "/", "localhost", false, true)
		c.JSON(200, gin.H{
			"csrf": token,
		})
	}
}
func RegisterHandler(db *sql.DB, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reg i.Register
		if err := c.ShouldBindBodyWithJSON(&reg); err != nil {
			logger.Print(err.Error())

			c.String(500, err.Error())
			return
		}
		MainContext := c.Request.Context()
		id, err := i.CreateUser(db, MainContext, reg, cnf.Postgres.RwTimeout)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		session, err := i.CreateSession(id, redis_db)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		err = i.AddSessionToCash(MainContext, session, redis_db, cnf.Redis.RwTimeout)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		c.SetCookie("session_id", session.Id, 1800, "/", "localhost", false, true)
		csrf, err := i.CreateCSRF(cnf.Secret, id)
		if err != nil {
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		c.JSON(200, gin.H{
			"csrf": csrf,
		})
	}
}
func main() {
	logger = log.Default()
	cnf = i.GetConfig()
	fmt.Println(cnf)
	err := i.NewPostgresDb(db, cnf)
	if err != nil {
		logger.Print(err.Error())
	}
	redis_db, err = i.NewRedisDb(cnf)
	if err != nil {
		logger.Print(err.Error())
	}
	r := gin.Default()
	r.POST("/login", LoginHandler(db, logger))
	r.POST("/register", RegisterHandler(db, logger))

}
