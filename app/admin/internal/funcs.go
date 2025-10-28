package internal

import (
	"context"
	"database/sql"
	"time"
)

func GetProducts(ch chan ChanProducts, ctx context.Context, timeout int, db *sql.DB) {
	newcontext, c := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer c()
	var resp ChanProducts
	rows, err := db.QueryContext(newcontext, `select * from products`)
	if err != nil {
		resp.Err = err
		ch <- resp
		return
	}
	for rows.Next() {
		var prod Product
		err := rows.Scan(prod)
		if err != nil {
			resp.Err = err
			ch <- resp
		}
		resp.Products = append(resp.Products, prod)
	}
	resp.Err = nil
	ch <- resp

}
func GetCats(ch chan ChanCats, ctx context.Context, timeout int, db *sql.DB) {
	newcontext, c := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer c()
	var resp ChanCats
	rows, err := db.QueryContext(newcontext, `select * from categories`)
	if err != nil {
		resp.Err = err
		ch <- resp
		return
	}
	for rows.Next() {
		var prod string
		err := rows.Scan(prod)
		if err != nil {
			resp.Err = err
			ch <- resp
		}
		resp.Categories = append(resp.Categories, prod)
	}
	resp.Err = nil
	ch <- resp

}
func GetUsers() {

}
func GetMainPage(ctx context.Context, db *sql.DB, timeout int) (MainPage, error) {
	prod_chan := make(chan ChanProducts)
	cats_chan := make(chan ChanCats)
	var chanprods ChanProducts
	var chancats ChanCats
	var MainPage MainPage
	go GetProducts(prod_chan, ctx, timeout, db)
	go GetCats(cats_chan, ctx, timeout, db)
	chancats = <-cats_chan
	chanprods = <-prod_chan
	if chancats.Err != nil {
		return MainPage, chancats.Err
	}
	if chanprods.Err != nil {
		return MainPage, chanprods.Err
	}
	MainPage.Categories = chancats.Categories
	MainPage.Products = chanprods.Products
	return MainPage, nil
}
