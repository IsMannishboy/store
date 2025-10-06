package funcs

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	s "gin/internal/structs"
)

func GetProducts(db *sql.DB, ch chan s.Products_chan_struct, ctx context.Context, timeout int) {
	newcontext, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	var data s.Products_chan_struct
	data.Err = nil
	rows, err := db.QueryContext(newcontext, `select * from products`)
	if err != nil {
		data.Err = err
		ch <- data
	}

	for rows.Next() {
		var product s.Product
		err := rows.Scan(product.ID, product.Name, product.Description, product, product.Price, product.Stock, product.Category, product.Date_of_create)
		if err != nil {
			data.Err = err
			ch <- data
		}
		data.Products = append(data.Products, product)
	}
	ch <- data
}
func GetCategories(db *sql.DB, ch chan s.Cat_chan_struct, cxt context.Context, timeout int) {
	newcontext, cancel := context.WithTimeout(cxt, time.Duration(timeout)*time.Second)
	defer cancel()
	var cats s.Cat_chan_struct
	cats.Err = nil
	rows, err := db.QueryContext(newcontext, `select * from categories`)
	if err != nil {
		cats.Err = err
		ch <- cats
	}
	var cat s.Categorie
	for rows.Next() {
		err := rows.Scan(cat.ID, cat.Name)
		if err != nil {
			cats.Err = err
			ch <- cats
		}
	}
	ch <- cats
}
func GetMainPageData(db *sql.DB, ctx context.Context, timeout int) (error, s.MainpageData) {
	var MainData s.MainpageData
	products_chan := make(chan s.Products_chan_struct)
	cat_chan := make(chan s.Cat_chan_struct)
	go GetProducts(db, products_chan, ctx, timeout)
	go GetCategories(db, cat_chan, ctx, timeout)
	var cs s.Cat_chan_struct
	var ps s.Products_chan_struct
	cs = <-cat_chan
	if cs.Err != nil {
		return cs.Err, MainData
	}
	MainData.Categories = cs.Categories
	ps = <-products_chan
	fmt.Println(MainData.Categories)
	if ps.Err != nil {
		return ps.Err, MainData
	}
	MainData.Products = ps.Products
	fmt.Println(MainData.Products)
	MainData.Title = "main page"
	return nil, MainData
}
