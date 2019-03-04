package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler ...
func Handler() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/admin").Subrouter()

	s.HandleFunc("/", handler())
	s.HandleFunc("/admin", adminHandler())

	return s
}

func handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hohoho"))
	}
}

func adminHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hehehe"))
	}
}
