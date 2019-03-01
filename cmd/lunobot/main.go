package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sfardiansyah/lunobot/pkg/bot"
)

func main() {
	app, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	h := bot.NewHandler(app)
	// r := rest.Handler()

	app.Debug = false

	info, err := app.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsSet() {
		_, err = app.SetWebhook(tgbotapi.NewWebhook("https://luno-bot.herokuapp.com:443"))
		if err != nil {
			log.Fatal(err)
		}
	}

	updates := app.ListenForWebhook("/")
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		h.Handle(update)
	}
}
