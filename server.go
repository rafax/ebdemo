package main

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	r        *render.Render
	store    Store
	sc       ServerConfig
	hostname string
)

func main() {
	initialize()

	mux := mux.NewRouter()

	mux.HandleFunc("/mem/{userId}", updateMemory).Methods("POST")
	mux.HandleFunc("/status", health)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe("0.0.0.0:"+sc.Port, n)
}

func updateMemory(w http.ResponseWriter, req *http.Request) {
	var pu PlayheadUpdate
	json.NewDecoder(req.Body).Decode(&pu)
	userId := mux.Vars(req)["userId"]
	store.Store(userId, pu)
	r.JSON(w, http.StatusOK, map[string]string{"userId": userId, "mgid": pu.Mgid, "playhead": pu.Playhead})
}

func initialize() {
	var dbc DbConfig
	err := envconfig.Process("ebdemo", &dbc)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = envconfig.Process("ebdemo", &sc)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(sc)

	store = NewPgStore(dbc)
	r = render.New()
	if h, err := os.Hostname(); err == nil {
		hostname = h
	} else {
		hostname = "unknown"
	}

}

func health(w http.ResponseWriter, req *http.Request) {
	err := store.Ping()
	if err == nil {
		r.JSON(w, http.StatusOK, map[string]string{"status": "UP", "hostname": hostname})
	} else {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "DOWN", "err": err.Error()})
	}
}

func factorial(n int) *big.Int {
	return big.NewInt(1).MulRange(1, int64(n))
}

type ServerConfig struct {
	Port         string `default:"3000"`
	HelloMessage string `default:"Hi there"`
}

type PlayheadUpdate struct {
	Mgid     string
	Playhead string
}
