package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sfardiansyah/lunobot/pkg/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	app, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	h := bot.NewHandler(app)

	app.Debug = true

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
		log.Printf("%+v\n", update.Message)
		h.Handle(update)
	}
}
