package rest

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handler ...
func Handler(b *tgbotapi.BotAPI) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", handler(b))

	updates := b.ListenForWebhook("/")

	for update := range updates {
		log.Printf("%+v\n", update.Message)
	}

	return r
}

func handler(b *tgbotapi.BotAPI) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hohoho")
	}
}
