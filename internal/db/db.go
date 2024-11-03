package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Db struct {
	Pg sqlx.DB
}

func New(connStr string) (*Db, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Db{Pg: *db}, nil
}
