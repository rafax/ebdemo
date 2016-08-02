package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
	"gopkg.in/unrolled/render.v1"
)

func main() {
	r := render.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, map[string]string{"status": "UP"})
	})

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
