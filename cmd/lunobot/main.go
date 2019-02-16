package main

import (
	"log"
	"os"

	"github.com/yanzay/tbot"
)

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Handle("/update", "update")
	bot.Handle("/help", "help")
	bot.Handle("/start", "start")
	bot.Handle("/infoluno", "infoluno")
	bot.Handle("/fee", "fee")
	bot.Handle("/convert", "convert")

	if err := bot.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	// http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
