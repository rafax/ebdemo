package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	cmap "github.com/streamrail/concurrent-map"

	"fmt"
)

type Store interface {
	Store(uid string, update PlayheadUpdate)
	Get(uid, mgid string) *PlayheadUpdate
	Ping() error
}

type Hits struct {
	factorial string
	hits      int64
}

type pgStore struct {
	db *sql.DB
}

type memStore struct {
	phs cmap.ConcurrentMap
}

func (pg pgStore) Store(uid string, update PlayheadUpdate) {
	res, err := pg.db.Exec("INSERT INTO ebdemo.playhead as p (uid,mgid,playhead) VALUES ($1,$2,$3) ON CONFLICT (uid,mgid) DO UPDATE SET playhead = $3 WHERE p.uid = $1 AND p.mgid = $2", uid, update.Mgid, update.Playhead)
	if err != nil {
		fmt.Printf("Could not store %v", err)
	} else {
		fmt.Printf("Stored, got %v", res)
	}
}

func (pg pgStore) Get(uid, mgid string) *PlayheadUpdate {
	var h PlayheadUpdate
	h.Mgid = mgid
	err := pg.db.QueryRow("SELECT p.playhead FROM ebdemo.playhead as p WHERE p.uid = $1 AND p.mgid = $2", uid, mgid).Scan(&h.Playhead)
	if err != nil {
		fmt.Printf("Error %v", err)
		return nil
	} else {
		return &h
	}
}

func (pg pgStore) Ping() error {
	return pg.db.Ping()
}

func key(uid, mgid string) string {
	return uid + ":" + mgid
}

func (m memStore) Store(uid string, update PlayheadUpdate) {
	m.phs.Set(key(uid, update.Mgid), update.Playhead)
}

func (m memStore) Get(uid, mgid string) *PlayheadUpdate {
	playhead, ok := m.phs.Get(key(uid, mgid))
	if ok {
		return &PlayheadUpdate{Mgid: mgid, Playhead: playhead.(string)}
	} else {
		return nil
	}
}

func (m memStore) Ping() error {
	return nil
}

type DbConfig struct {
	Url      string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Db       string `required:"true"`
}

func NewPgStore(dbc DbConfig) Store {
	db, _ := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbc.User, dbc.Password, dbc.Db, dbc.Url))
	db.SetMaxOpenConns(100)
	return pgStore{db: db}
}

func NewMemStore() memStore {
	return memStore{phs: cmap.New()}
}
