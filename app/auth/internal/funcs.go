package internal

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, ctx context.Context, reg Register, timeout int) (string, error) {
	newcontext, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	var id string

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		return id, err
	}
	err = db.QueryRowContext(newcontext, `insert into users(username,firstname,pass,email) value ($1,$2,$3,$4) returning id`, reg.Username, reg.Firstname, hashedpassword, reg.Email).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
func GetUserId(db *sql.DB, ctx context.Context, username string, timeout int) (string, error) {
	newcontext, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer cancel()
	var id string
	err := db.QueryRowContext(newcontext, `select id from users where username = $1`, username).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
func GetPass(db *sql.DB, ctx context.Context, username string, timeout int) ([]byte, error) {
	newcontext, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer cancel()
	var id []byte
	err := db.QueryRowContext(newcontext, `select pass from users where username = $1`, username).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
