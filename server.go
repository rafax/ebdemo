package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	r     *render.Render
	store Store
)

func main() {
	var dbc DbConfig
	err := envconfig.Process("ebdemo", &dbc)
	if err != nil {
		log.Fatal(err.Error())
	}
	var sc ServerConfig
	err = envconfig.Process("ebdemo", &sc)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(sc)

	store = NewStore(dbc)
	r = render.New()
	mux := mux.NewRouter().StrictSlash(false)

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, sc.HelloMessage)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		err := store.Ping()
		if err == nil {
			r.JSON(w, http.StatusOK, map[string]string{"status": "UP"})
		} else {
			r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "DOWN", "err": err.Error()})
		}
	})
	mux.HandleFunc("/{n}", calc)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe("0.0.0.0:"+sc.Port, n)
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
		r.JSON(w, http.StatusOK, map[string]string{"n": n, "factorial": h.factorial, "hits": strconv.FormatInt(h.hits, 10)})
	} else {
		log.Println(fmt.Sprintf("Not found %s, calculating", n))
		res := factorial(in)
		store.Hit(in, res.String())
		r.JSON(w, http.StatusOK, map[string]string{"n": n, "factorial": res.String(), "hits": "0", "calculated": "true"})
	}

}

func factorial(n int) *big.Int {
	return big.NewInt(1).MulRange(1, int64(n))
}

type ServerConfig struct {
	Port         string `default:"3000"`
	HelloMessage string `default:"Hi there"`
}
