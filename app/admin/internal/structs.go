package internal

import (
	"time"
)

type Config struct {
	HTMLPath string
	Secret   string
	Server   Server
	Postgres Postgres
	Redis    Redis
}
type Postgres struct {
	Host        string
	Port        int
	Db          string
	RwTimeout   int
	DialTimeout int
	Password    string
	User        string
	Sslmode     string
}
type Redis struct {
	Addr        string
	DialTimeout int
	RwTimeout   int
	DbIndex     int
}
type Server struct {
	Port int
}
type Session struct {
	Id     string
	UserId int
	Exp    time.Time
	Role   string
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
type Category struct {
	Id   int
	Name string
}
type Product struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	DateOfCreate time.Time `json:"created_at"`
}
type ChanProducts struct {
	Id       int
	Products []Product
	Err      error
}
type ChanCats struct {
	Categories []Category
	Err        error
}
type User struct {
	Id        int
	Username  string
	Firstname string
	Email     string
	Password  string
	Role      string
}
type ChanUsers struct {
	Users []User
	Err   error
}
type MainPage struct {
	Products   []Product
	Categories []Category
	Users      []User
}
