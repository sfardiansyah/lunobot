package rest

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

// PageVariables ...
type PageVariables struct {
	Date   string
	Time   string
	Params url.Values
}

// Handler ...
func Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/admin", adminHandler())

	return r
}

func handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hohoho"))
	}
}

func adminHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		now := time.Now()              // find the time right now
		HomePageVars := PageVariables{ //store the date and time in a struct
			Date:   now.Format("02-01-2006"),
			Time:   now.Format("15:04:05"),
			Params: params,
		}

		t, err := template.ParseFiles("web/homepage.html") //parse the html file homepage.html
		if err != nil {                                    // if there is an error
			log.Print("template parsing error: ", err) // log it
		}
		err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
		if err != nil {                  // if there is an error
			log.Print("template executing error: ", err) //log it
		}
	}
}
