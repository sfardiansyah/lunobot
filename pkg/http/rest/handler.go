package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler ...
func Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/admin", handler())

	return r
}

func handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hohoho"))
	}
}
