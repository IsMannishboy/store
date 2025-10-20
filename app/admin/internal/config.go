package internal

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetConfig(logger *log.Logger) (Config, error) {
	err := godotenv.Load(os.Getenv("ENVPATH"))
	var cnf Config
	if err != nil {
		logger.Print(err)
		return cnf, err
	}
	cnf.Secret = os.Getenv("secret")
	cnf.HTMLPath = os.Getenv("htmlpath")
	var postgres Postgres
	postgres.Sslmode = os.Getenv("postgres_ssl")
	postgres.Host = os.Getenv("postgres_host")
	postgres.Port, err = strconv.Atoi(os.Getenv("postgres_port"))
	if err != nil {
		logger.Print(err)

		return cnf, err
	}
	postgres.Password = os.Getenv("postgres_pass")
	postgres.Db = os.Getenv("dbname")
	postgres.DialTimeout, err = strconv.Atoi(os.Getenv("postgres_dialtimeout"))
	postgres.RwTimeout, err = strconv.Atoi(os.Getenv("postgres_rwtimeout"))
	postgres.User = os.Getenv("postgres_user")
	var redis Redis
	redis.Addr = os.Getenv("redis_addr")
	redis.DialTimeout, err = strconv.Atoi(os.Getenv("redis_dialtimeout"))
	if err != nil {
		logger.Print(err)

		return cnf, err
	}
	redis.RwTimeout, err = strconv.Atoi(os.Getenv("redis_rwtimeout"))
	if err != nil {
		logger.Print(err)

		return cnf, err
	}
	redis.DbIndex, err = strconv.Atoi(os.Getenv("redis_db"))
	if err != nil {
		logger.Print(err)

		return cnf, err
	}
	cnf.Server.Port, err = strconv.Atoi(os.Getenv("port"))
	if err != nil {
		logger.Print(err)

		return cnf, err
	}
	cnf.Postgres = postgres
	cnf.Redis = redis
	return cnf, nil
}
