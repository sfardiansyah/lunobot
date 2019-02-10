package main

import (
	"log"
	"net/http"
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
	bot.ListenAndServe()
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
