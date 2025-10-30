package internal

import (
	"context"
	"database/sql"
	"fmt"
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
		err := rows.Scan(&prod.Id, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.Category)
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
		var cat Category
		err := rows.Scan(&cat.Id, &cat.Name)
		if err != nil {
			resp.Err = err
			ch <- resp
		}
		resp.Categories = append(resp.Categories, cat)
	}
	resp.Err = nil
	ch <- resp

}
func GetUsers(ch chan ChanUsers, ctx context.Context, db *sql.DB, timeout int) {
	var users ChanUsers
	users.Err = nil
	newcontext, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	rows, err := db.QueryContext(newcontext, `select * from users`)
	if err != nil {
		users.Err = err
		ch <- users
		return
	}
	for rows.Next() {
		var id int
		var username string
		var firstname string
		var email string
		var Password string
		var role string
		err = rows.Scan(&id, &username, &firstname, &email, &Password, &role)
		if err != nil {
			users.Err = err
			ch <- users
		}
		var user = User{Id: id, Username: username, Firstname: firstname, Email: email, Password: Password, Role: role}
		users.Users = append(users.Users, user)
	}
	ch <- users
}
func GetMainPage(ctx context.Context, db *sql.DB, timeout int) (MainPage, error) {
	prod_chan := make(chan ChanProducts)
	cats_chan := make(chan ChanCats)
	users_chan := make(chan ChanUsers)
	var chanprods ChanProducts
	var chancats ChanCats
	var MainPage MainPage
	go GetProducts(prod_chan, ctx, timeout, db)
	go GetCats(cats_chan, ctx, timeout, db)
	go GetUsers(users_chan, ctx, db, timeout)
	chancats = <-cats_chan
	chanprods = <-prod_chan
	chanusers := <-users_chan
	if chancats.Err != nil {
		fmt.Println("GetCats err:", chancats.Err.Error())
		return MainPage, chancats.Err
	}
	if chanprods.Err != nil {
		fmt.Println("GetProducts err:", chanprods.Err.Error())
		return MainPage, chanprods.Err
	}
	if chanusers.Err != nil {
		fmt.Println("GetUsers err:", chanusers.Err.Error())
		return MainPage, chanusers.Err
	}
	MainPage.Categories = chancats.Categories
	MainPage.Products = chanprods.Products
	MainPage.Users = chanusers.Users
	return MainPage, nil
}
