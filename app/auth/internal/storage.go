package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func NewPostgresDb(db *sql.DB, cnf *Config) error {
	var err error
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cnf.Postgres.Addr,
		cnf.Postgres.Port,
		cnf.Postgres.User,
		cnf.Postgres.Pass,
		cnf.Postgres.Dbname,
		cnf.Postgres.Sslmode,

		cnf.Postgres.DialTimeout,
	)
	db, err = sql.Open("postgres", connstr)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		if err := db.Ping(); err != nil {
			return err
		}

	}

	return nil
}
func NewRedisDb(cnf *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cnf.Redis.Addr,
		DialTimeout:  time.Duration(cnf.Redis.DialTimeout),
		ReadTimeout:  time.Duration(cnf.Redis.RwTimeout),
		WriteTimeout: time.Duration(cnf.Redis.RwTimeout),
		DB:           0,
	})
	var err error
	for i := 0; i < 10; i++ {
		err = rdb.Ping(context.Background()).Err()
		if err == nil {
			break
		}
	}
	return rdb, err
}
