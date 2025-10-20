package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	redis "github.com/redis/go-redis/v9"
)

func NewPostgresDb(cnf *Config, logger *log.Logger) (error, *sql.DB) {
	var err error
	var db *sql.DB
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cnf.Postgres.Addr,
		cnf.Postgres.Port,
		cnf.Postgres.User,
		cnf.Postgres.Pass,
		cnf.Postgres.Dbname,
		cnf.Postgres.Sslmode,

		cnf.Postgres.DialTimeout,
	)
	fmt.Println("Postgres config:")
	fmt.Println("Addr:", cnf.Postgres.Addr)
	fmt.Println("Port:", cnf.Postgres.Port)
	fmt.Println("User:", cnf.Postgres.User)
	fmt.Println("Pass:", cnf.Postgres.Pass)
	fmt.Println("Dbname:", cnf.Postgres.Dbname)
	fmt.Println("Sslmode:", cnf.Postgres.Sslmode)
	fmt.Println("DialTimeout:", cnf.Postgres.DialTimeout)
	db, err = sql.Open("postgres", connstr)
	if err != nil {
		return err, db
	}

	for i := 0; i < 10; i++ {
		logger.Print("pinging postgres")
		if err = db.Ping(); err != nil {
			logger.Print(err.Error())

		} else {
			break
		}

	}
	logger.Print("returning Storage")

	return err, db
}
func NewRedisDb(cnf *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cnf.Redis.Addr,
		DialTimeout:  time.Duration(cnf.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cnf.Redis.RwTimeout) * time.Second,
		WriteTimeout: time.Duration(cnf.Redis.RwTimeout) * time.Second,
		DB:           0,
	})
	fmt.Println("redis conf:")
	fmt.Println("redis addr:", cnf.Redis.Addr)
	fmt.Println("redis dial timeout:", cnf.Redis.DialTimeout)
	var err error
	for i := 0; i < 10; i++ {
		err = rdb.Ping(context.Background()).Err()
		if err == nil {
			break
		}
	}
	return rdb, err
}
