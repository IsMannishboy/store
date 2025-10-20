package storage

import (
	"context"
	"database/sql"
	"fmt"
	c "gin/internal/config"
	"log/slog"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Storage struct {
	DB *sql.DB
}
type Cash struct {
	Redis_db *redis.Client
}

func NewPostgresDb(cnf *c.Config, logger *slog.Logger) (*Storage, error) {
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
	fmt.Println("len addr:", len(cnf.Postgres.Addr))
	fmt.Println("Port:", cnf.Postgres.Port)
	fmt.Println("len port:", len(cnf.Postgres.Port))
	fmt.Println("User:", cnf.Postgres.User)
	fmt.Println("Pass:", cnf.Postgres.Pass)
	fmt.Println("Dbname:", cnf.Postgres.Dbname)
	fmt.Println("Sslmode:", cnf.Postgres.Sslmode)
	fmt.Println("DialTimeout:", cnf.Postgres.DialTimeout)
	fmt.Println("redis addr:", cnf.Redis.Addr)
	logger.Debug("openning db connection")
	db, err = sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 10; i++ {
		logger.Debug("pinging postgres")
		if err = db.Ping(); err != nil {
			logger.Debug(err.Error())
			continue
		}
		break

	}
	logger.Debug("returning Storage")
	return &Storage{DB: db}, err
}
func NewRedisDb(cnf *c.Config) (*Cash, error) {
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
		fmt.Println("pinging redis")
		err = rdb.Ping(context.Background()).Err()
		if err != nil {
			fmt.Println(err)
			continue
		}

		break

	}

	cash := Cash{
		Redis_db: rdb,
	}

	return &cash, err
}
