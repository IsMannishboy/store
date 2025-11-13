package main

import (
	i "admin/internal"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func MainHandler(logger *log.Logger, redis_db *redis.Client, db *sql.DB, cnf *i.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		csrf := c.Param("csrf")
		fmt.Println("csrf:", csrf)
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
		decodedValue, _ := url.QueryUnescape(session_id)

		MainContext := c.Request.Context()
		err = i.CheckSession(MainContext, decodedValue, redis_db, time.Duration(cnf.Redis.RwTimeout))
		if err != nil {
			logger.Print("CheckSession err:", err)
			c.String(403, err.Error())
			return
		}
		MainPage, err := i.GetMainPage(MainContext, db, cnf.Postgres.RwTimeout)
		if err != nil {
			logger.Print("MainPage:", err)
			c.String(500, err.Error())
			return
		}
		c.HTML(200, "main.html", gin.H{
			"csrf":       new_csrf,
			"products":   MainPage.Products,
			"categories": MainPage.Categories,
			"users":      MainPage.Users,
		})
	}
}
func DeleteHandler(logger *log.Logger, redis_db *redis.Client, db *sql.DB, cnf *i.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		MainContext := c.Request.Context()
		csrf := c.Param("csrf")
		id, new_csrf, err := i.CheckCSRF(csrf, []byte(cnf.Secret))
		if err != nil {
			logger.Print(err)
		}
		item_id := c.Param("id")
		int_item_id, err := strconv.Atoi(item_id)
		logger.Printf(item_id)
		var sessionId = ""
		for _, c := range c.Request.Cookies() {
			if len(c.Name) > len("session_id_") {
				user_id := c.Name[len("session_id_"):]
				int_id, err := strconv.Atoi(user_id)
				if err != nil {
					logger.Print(err)
					continue
				}
				if int_id == id {
					logger.Print("user id:", id)
					sessionId = c.Value
					break
				}
			}
		}
		if sessionId == "" {
			c.String(403, "")
			return
		}
		sessionId, _ = url.QueryUnescape(sessionId)

		err = i.CheckSession(MainContext, sessionId, redis_db, time.Duration(cnf.Redis.RwTimeout))
		if err != nil {
			if err != redis.Nil || err.Error() != "wrong role" || err.Error() != "session exp" {
				logger.Print("CheckSession err:", err)
				c.String(500, err.Error())
				return
			}
			c.String(403, err.Error())
		}
		_, err = db.Exec("delete from products where id = $1", int_item_id)
		if err != nil {
			logger.Print(err)
			c.String(500, err.Error())
			return
		}
		c.JSON(200, gin.H{
			"csrf": new_csrf,
		})
	}
}
func AddProduct(logger *log.Logger, db *sql.DB, redis_db *redis.Client, cnf *i.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		MainContext := c.Request.Context()
		csrf := c.Param("csrf")
		id, new_csrf, err := i.CheckCSRF(csrf, []byte(cnf.Secret))
		if err != nil {
			logger.Print(err)
		}

		var sessionId = ""
		for _, c := range c.Request.Cookies() {
			if len(c.Name) > len("session_id_") {
				user_id := c.Name[len("session_id_"):]
				int_id, err := strconv.Atoi(user_id)
				if err != nil {
					logger.Print(err)
					continue
				}
				if int_id == id {
					logger.Print("user id:", id)
					sessionId = c.Value
					break
				}
			}
		}
		if sessionId == "" {
			c.String(403, "")
			return
		}
		sessionId, _ = url.QueryUnescape(sessionId)
		fmt.Println("session id:", sessionId)
		err = i.CheckSession(MainContext, sessionId, redis_db, time.Duration(cnf.Redis.RwTimeout))
		if err != nil {
			if err != redis.Nil || err.Error() != "wrong role" || err.Error() != "session exp" {
				logger.Print("CheckSession err:", err)
				c.String(500, err.Error())
				return
			}
			c.String(403, err.Error())
		}
		var newProduct i.Product
		data, err := c.GetRawData()
		if err != nil {
			logger.Print(err)
			c.String(500, err.Error())
			return
		}
		err = json.Unmarshal(data, &newProduct)
		if err != nil {
			logger.Print(err)
			c.String(500, err.Error())
			return
		}
		fmt.Println("data:", newProduct)
		_, err = db.Exec("insert into products (prod_name, description, price, stock, category,created_at) values ($1,$2,$3,$4,$5,$6)", newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Stock, newProduct.Category, time.Now())
		if err != nil {
			logger.Print(err)
			c.String(500, err.Error())
			return
		}
		c.JSON(200, gin.H{
			"csrf": new_csrf,
			"Id":   newProduct.Id,
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
	router.LoadHTMLGlob(cnf.HTMLPath + "/*.html")

	router.GET("/main/:csrf", MainHandler(logger, redis_db, db, &cnf))
	router.POST("/delete/:id/:csrf", DeleteHandler(logger, redis_db, db, &cnf))
	router.POST("/add/:csrf", AddProduct(logger, db, redis_db, &cnf))
	fmt.Println("hello world")
	http.ListenAndServe(":80", router)
}
