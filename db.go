package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"fmt"
)

type Store interface {
	Hit(n int, fact string)
	Get(n int) *Hits
	Ping() error
}

type Hits struct {
	factorial string
	hits      int64
}

type pgStore struct {
	db *sql.DB
}

func (pg pgStore) Hit(n int, factorial string) {
	pg.db.Exec("INSERT INTO ebdemo.fact as f (n,factorial,hits) VALUES ($1,$2,0) ON CONFLICT (n) DO UPDATE SET hits = f.hits +1 WHERE f.n = $1;", n, factorial)
}

func (pg pgStore) Get(n int) *Hits {
	var h Hits
	err := pg.db.QueryRow("SELECT factorial, hits FROM ebdemo.fact WHERE n = $1", n).Scan(&h.factorial, &h.hits)
	if err != nil {
		log.Printf("Error when getting for %d: %s", n, err.Error())
		return nil
	}
	return &h
}

func (pg pgStore) Ping() error {
	return pg.db.Ping()
}

type DbConfig struct {
	Url      string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Db       string `required:"true"`
}

func NewStore(dbc DbConfig) Store {
	db, _ := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbc.User, dbc.Password, dbc.Db, dbc.Url))
	return pgStore{db: db}
}
