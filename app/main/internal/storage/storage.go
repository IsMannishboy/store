package storage

import (
	"database/sql"
	"fmt"
	c "gin/internal/config"
)

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB, cnf *c.Config) (*Storage, error) {
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		cnf.Postgres.Addr,
		cnf.Postgres.Port,
		cnf.Postgres.User,
		cnf.Postgres.Pass,
		cnf.Postgres.Dbname,
		cnf.Postgres.Sslmode,

		cnf.Postgres.DialTimeout,
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}
