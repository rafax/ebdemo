package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"fmt"
)

type Store interface {
	Hit(n int, fact int64)
	Get(n int) Hits
	Ping() error
}

type Hits struct {
	factorial int64
	hits      int64
}

type pgStore struct {
	db *sql.DB
}

func (pg pgStore) Hit(n int, factorial int64) {
	pg.db.Exec("INSERT INTO fact (n,factorial,hits) VALUES ($1,$2,0) ON CONFLICT UPDATE SET hits = HITS +1", n, factorial)
}

func (pg pgStore) Get(n int) Hits {
	var h Hits
	rows, err := pg.db.Query("SELECT factorial, hits FROM fact WHERE n = $1", n)
	if err != nil {
		log.Println("Error when getting for n" + err.Error())
	}
	rows.Scan(&h)
	if err != nil {
		log.Println("Error when parsing result" + err.Error())
	}
	return h
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
