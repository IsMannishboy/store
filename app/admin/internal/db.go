package internal

import (
	"database/sql"
	"fmt"

	"context"
	"time"

	_ "github.com/lib/pq"
	redis "github.com/redis/go-redis/v9"
)

func GetPostgresConn(cnf Config) (*sql.DB, error) {
	connstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cnf.Postgres.Host,
		cnf.Postgres.Port,
		cnf.Postgres.User,
		cnf.Postgres.Password,
		cnf.Postgres.Db,
		cnf.Postgres.Sslmode,
		cnf.Postgres.DialTimeout,
	)
	fmt.Println("Sslmode from config:", cnf.Postgres.Sslmode)
	fmt.Println("Connstr:", connstr)

	db, err := sql.Open("postgres", connstr)

	if err != nil {
		fmt.Println("err while oppening postgres conn:", err)
		return db, err
	}
	var PostgresPingError error
	for i := 0; i < 10; i++ {
		PostgresPingError = db.Ping()
		if PostgresPingError == nil {
			break
		}
		time.Sleep(2000)
	}
	if PostgresPingError != nil {
		fmt.Println("PostgresPingError:", PostgresPingError)
		return db, PostgresPingError
	}
	fmt.Println("pinged successfuly")
	return db, nil

}
func GetRedisConn(cnf Config) (*redis.Client, error) {
	redis_db := redis.NewClient(&redis.Options{
		Addr:         cnf.Redis.Addr,
		DialTimeout:  time.Duration(cnf.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cnf.Redis.RwTimeout) * time.Second,
		WriteTimeout: time.Duration(cnf.Redis.RwTimeout) * time.Second,
		DB:           cnf.Redis.DbIndex,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cnf.Redis.RwTimeout)*time.Second)
	defer cancel()
	var PingError error
	for i := 0; i < 10; i++ {
		PingError = redis_db.Ping(ctx).Err()
		if PingError == nil {
			break
		}
		time.Sleep(2000)
	}
	if PingError != nil {
		return redis_db, PingError
	}
	return redis_db, nil

}
