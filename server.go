package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	r *render.Render
)

func main() {
	var dbc DbConfig
	err := envconfig.Process("ebdemo", &dbc)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(dbc)

	store := NewStore(dbc)
	r = render.New()
	mux := mux.NewRouter().StrictSlash(false)

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	mux.HandleFunc("/{n}", calc)
	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		err := store.Ping()
		if err == nil {
			r.JSON(w, http.StatusOK, map[string]string{"status": "UP"})
		} else {
			r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "DOWN", "err": err.Error()})
		}
	})

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}

func calc(w http.ResponseWriter, req *http.Request) {
	n := mux.Vars(req)["n"]
	in, err := strconv.Atoi(n)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot parse to int: " + n, "err": err.Error()})
	}
	res := factorial(int64(in))
	r.JSON(w, http.StatusOK, map[string]int64{"n": int64(in), "factorial": res, "hits": 0})
}

func factorial(n int64) int64 {
	res := int64(1)
	for i := int64(2); i <= int64(n); i++ {
		res *= i
	}
	return res
}
