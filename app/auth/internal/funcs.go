package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, ctx context.Context, reg Register, timeout int) (int, error) {
	newcontext, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	var id int

	err := db.QueryRow(`select id from users where username = $1`, reg.Username).Scan(&id)
	switch err {
	case sql.ErrNoRows:
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
		if err != nil {
			return id, err
		}
		row := db.QueryRowContext(newcontext, `
			insert into users(username, firstname, pass, gmail)
			values ($1, $2, $3, $4)
			returning id`,
			reg.Username, reg.Firstname, hashedpassword, reg.Email,
		)

		err = row.Scan(&id)
		if err != nil {
			return id, err
		}
		return id, nil
	case nil:
		return id, fmt.Errorf("this username already exist")
	default:
		return id, err
	}

}
func GetUserIdAndPassAndRole(db *sql.DB, ctx context.Context, username string, timeout int) (int, []byte, string, error) {
	newcontext, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	defer cancel()
	var id int
	var pass []byte
	var role string
	var err error
	fmt.Println("username:", username)
	err = db.QueryRowContext(newcontext, `select id,pass,rolee from users where username = $1`, username).Scan(&id, &pass, &role)

	return id, pass, role, err
}
func GetRole(ref string) string {
	if ref == "http://localhost/admin/login" {
		return "admin"
	} else if ref == "http://localhost/login" {
		return "user"

	}
	return ""
}
