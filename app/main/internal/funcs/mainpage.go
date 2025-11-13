package funcs

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	s "gin/internal/structs"

	"github.com/lib/pq"
)

func GetProducts(db *sql.DB, ch chan s.Products_chan_struct, ctx context.Context, timeout int) {
	fmt.Println("GetProducts query")
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
	close(ch)
}

func GetCart(id int, db *sql.DB, ctx context.Context, timeout int, ch chan s.Cart_chan_struct) {
	fmt.Println("GetCart query")
	fmt.Printf("id type: %T, value: %v\n", id, id)

	newcontext, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	var cart s.Cart_chan_struct
	cart.Err = nil
	rows, err := db.QueryContext(newcontext, `select product_id,quantity from cart where user_id = $1`, id)
	if err != nil {
		fmt.Println("cart err:", err.Error())
		cart.Err = err
		ch <- cart
		return
	}
	var product_id int64
	var quantity int64
	var ids []int64
	for rows.Next() {
		err := rows.Scan(&product_id, &quantity)
		if err != nil {
			cart.Err = err
			ch <- cart
			return
		}
		ids = append(ids, product_id)
		var pro s.Product
		pro.ID = product_id
		pro.Stock = quantity
		cart.CartProducts = append(cart.CartProducts, pro)
	}
	prod_rows, err := db.QueryContext(newcontext, `select * from products where id = any($1)`, pq.Array(ids))
	if err != nil {
		cart.Err = err
		ch <- cart
		return
	}
	var product s.Product
	i := 0
	for prod_rows.Next() {
		err := rows.Scan(&product)
		if err != nil {
			cart.Err = err
			ch <- cart
		}
		cart.CartProducts[i] = product
	}
	ch <- cart
	close(ch)
}
func GetMainPageData(id int, db *sql.DB, ctx context.Context, timeout int) (error, s.MainpageData) {
	var MainData s.MainpageData
	products_chan := make(chan s.Products_chan_struct)
	cart_chan := make(chan s.Cart_chan_struct)

	go GetProducts(db, products_chan, ctx, timeout)
	go GetCart(id, db, ctx, timeout, cart_chan)
	var ps s.Products_chan_struct
	var cart s.Cart_chan_struct
	var err error

	ps = <-products_chan
	fmt.Println(MainData.Categories)
	if ps.Err != nil {
		fmt.Println("err while products query:", ps.Err.Error())
		err = ps.Err
	}
	MainData.Products = ps.Products
	fmt.Println(MainData.Products)
	cart = <-cart_chan
	if cart.Err != nil {
		fmt.Println("err while cart query:", cart.Err.Error())
		err = cart.Err
	}
	MainData.Cart = cart.CartProducts
	MainData.Title = "main page"
	return err, MainData
}
