package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

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

	mux.HandleFunc("/", hello)
	mux.HandleFunc("/status", health)
	mux.HandleFunc("/{n}", calc)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe("0.0.0.0:"+sc.Port, n)
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, sc.HelloMessage)
}

func calc(w http.ResponseWriter, req *http.Request) {
	n := mux.Vars(req)["n"]
	log.Println("BICZ " + req.URL.Query().Get("nocache"))
	nocache := req.URL.Query().Get("nocache") == "true"
	in, err := strconv.Atoi(n)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot parse to int: " + n, "err": err.Error()})
	}
	h := store.Get(in)
	if !nocache && h != nil {
		store.Hit(in, h.factorial)
		r.JSON(w, http.StatusOK, map[string]string{"n": n, "factorial": h.factorial, "hits": strconv.FormatInt(h.hits+1, 10)})
	} else {
		log.Println(fmt.Sprintf("Not found %s, calculating", n))
		res := factorial(in)
		store.Hit(in, res.String())
		r.JSON(w, http.StatusOK, map[string]string{"n": n, "factorial": res.String(), "hits": "0", "calculated": "true"})
	}

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

	store = NewStore(dbc)
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
