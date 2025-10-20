package structs

import "time"

type Product struct {
	ID             int64
	Name           string
	Price          float64
	Description    string
	Category       string
	Stock          int64
	Date_of_create time.Time
}
type Categorie struct {
	ID   int
	Name string
}
type MainpageData struct {
	Cart       []Product
	Categories []Categorie
	Products   []Product
	Title      string
}
type Products_chan_struct struct {
	Err      error
	Products []Product
}
type Cat_chan_struct struct {
	Err        error
	Categories []Categorie
}

type Cart_chan_struct struct {
	CartProducts []Product
	Err          error
}
type Session struct {
	Id      string
	User_id int
	Exp     time.Time
}
type SessionValue struct {
	UserId int
	Exp    time.Time
}
type CSRFvalue struct {
	UserId int
	Exp    time.Time
}
