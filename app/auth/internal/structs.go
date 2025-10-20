package internal

import (
	"time"
)

type Register struct {
	Firstname string `json:"firstname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Session struct {
	Id     string
	UserId int
	Exp    time.Time
}
type SessionValue struct {
	UserId int
	Exp    time.Time
	Role   string
}
type CSRFvalue struct {
	UserId int
	Exp    time.Time
}
type Server struct {
	Addr string
	Port string
}
type Postgres struct {
	Dbname      string
	User        string
	Pass        string
	Sslmode     string
	Addr        string
	Port        string
	DialTimeout int
	RwTimeout   int
}
type Redis struct {
	Addr        string
	DialTimeout int
	RwTimeout   int
}
type Config struct {
	Secret       string
	TemplatePath string
	Server       Server
	Postgres     Postgres
	Redis        Redis
}
