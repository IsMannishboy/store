package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetConfig() *Config {
	err := godotenv.Load(os.Getenv("ENVPATH"))
	fmt.Println("ENVPATH =", os.Getenv("ENVPATH"))
	if err != nil {
		panic(err)
	}

	var server Server
	server.Port = os.Getenv("port")
	server.Addr = os.Getenv("addr")
	var postgres Postgres
	postgres.Addr = os.Getenv("postgres_addr")
	postgres.Port = os.Getenv("postgres_port")
	var postgres_dial_timeout string
	var postgres_rw_timeout string
	postgres_dial_timeout = os.Getenv("postgres_dial_timeout")
	dt, err := strconv.Atoi(postgres_dial_timeout)
	if err != nil {
		panic(err)
	}
	postgres_rw_timeout = os.Getenv("postgres_rw_timeout")
	rwt, err := strconv.Atoi(postgres_rw_timeout)
	if err != nil {
		panic(err)
	}
	postgres.DialTimeout = dt
	postgres.RwTimeout = rwt
	postgres.User = os.Getenv("postgres_user")
	postgres.Pass = os.Getenv("postgres_pass")
	postgres.Sslmode = os.Getenv("postgres_sslmode")
	postgres.Dbname = os.Getenv("postgres_db_name")
	var redis Redis
	redis.Addr = os.Getenv("redis_addr")
	var redis_dial_timeout string
	var redis_rw_timeout string
	redis_dial_timeout = os.Getenv("redis_dial_timeout")
	dt, err = strconv.Atoi(redis_dial_timeout)
	if err != nil {
		panic(err)
	}
	redis_rw_timeout = os.Getenv("redis_rw_timeout")
	rwt, err = strconv.Atoi(redis_rw_timeout)
	if err != nil {
		panic(err)
	}
	redis.DialTimeout = dt
	redis.RwTimeout = rwt
	var cnf Config
	cnf.Secret = os.Getenv("secret")
	cnf.Server = server
	cnf.Postgres = postgres
	cnf.Redis = redis
	cnf.TemplatePath = os.Getenv("template_path")
	return &cnf
}
