package main

import (
	i "auth/internal"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	redis "github.com/redis/go-redis/v9"
)

var cnf *i.Config
var logger *log.Logger

func LoginHandler(redis_db *redis.Client, db *sql.DB, logger *log.Logger) gin.HandlerFunc {
	logger.Print("login request")
	return func(c *gin.Context) {
		var login i.Login
		if err := c.ShouldBindBodyWithJSON(&login); err != nil {
			c.Error(fmt.Errorf("internal server error"))
			return
		}

		MainContext := c.Request.Context()
		/// logincon section
		id, hashed, err := i.GetUserIdAndPass(db, MainContext, login.Username, cnf.Postgres.RwTimeout)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Print("error while getting user id:", err.Error())
				c.String(404, err.Error())
				return
			}
			logger.Print("error while getting user id:", err.Error())
			c.String(500, err.Error())
			return
		}

		err = i.CheckPassword(login.Password, hashed)
		if err != nil {
			c.String(403, "wrong password")
			return
		}
		/// session section
		session, err := i.CreateSession(id)
		if err != nil {
			logger.Print("CreateSession err:", err)
			c.String(500, err.Error())
			return
		}
		err = i.AddSessionToCash(MainContext, session, redis_db, cnf.Redis.RwTimeout)
		if err != nil {
			logger.Print("CreateSession err:", err)
			c.String(500, err.Error())
			return
		}
		// csrf section
		token, err := i.CreateCSRF(cnf.Secret, session.UserId)
		if err != nil {
			logger.Print("CreateCSRF error:", err.Error())
			c.String(500, err.Error())
			return
		}
		fmt.Println("session_id to cookie:", session.Id)
		// response section
		c.SetCookie("session_id"+"_"+strconv.Itoa(id), session.Id, 1800, "/", "localhost", false, true)
		c.JSON(200, gin.H{
			"csrf": token,
		})

	}
}
func RegisterHandler(redis_db *redis.Client, db *sql.DB, logger *log.Logger) gin.HandlerFunc {
	logger.Print("register request")
	return func(c *gin.Context) {
		var reg i.Register
		if err := c.ShouldBindBodyWithJSON(&reg); err != nil {
			logger.Print("error while decoding json:", err.Error())

			c.String(500, err.Error())
			return
		}
		MainContext := c.Request.Context()
		logger.Print("reg data:", reg)
		// register section
		id, err := i.CreateUser(db, MainContext, reg, cnf.Postgres.RwTimeout)

		if err != nil {
			if err.Error() == "this username already exist" {
				logger.Print(err.Error())
				c.String(403, err.Error())
				return
			}
			logger.Print(err.Error())
			c.String(500, err.Error())
			return
		}
		session, err := i.CreateSession(id)
		if err != nil {
			logger.Print("CreateSession err:", err)
			c.String(500, err.Error())
			return
		}
		err = i.AddSessionToCash(MainContext, session, redis_db, cnf.Redis.RwTimeout)
		if err != nil {
			logger.Print("CreateSession err:", err)
			c.String(500, err.Error())
			return
		}
		c.SetCookie("session_id"+"_"+strconv.Itoa(id), session.Id, 1800, "/", "localhost", false, true)
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

	err, db := i.NewPostgresDb(cnf, logger)
	if err != nil {
		logger.Print("err while db connections:", err.Error())
	}
	redis_db, err := i.NewRedisDb(cnf)
	r := gin.Default()
	r.POST("/login", LoginHandler(redis_db, db, logger))
	r.POST("/register", RegisterHandler(redis_db, db, logger))
	http.ListenAndServe(":80", r)
}
