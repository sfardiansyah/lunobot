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
	bot.Handle("/update", "42")
	bot.Handle("/help", "42")
	bot.Handle("/start", "42")
	bot.Handle("/infoluno", "42")
	bot.Handle("/fee", "42")
	bot.Handle("/convert", "42")
	bot.ListenAndServe()
}
