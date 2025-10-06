package structs

import "time"

type Product struct {
	ID             int
	Name           string
	Price          float64
	Description    string
	Category       string
	Stock          int
	Date_of_create time.Time
}
type Categorie struct {
	ID   int
	Name string
}
type MainpageData struct {
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
