package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/yanzay/tbot"
)

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("assets/help.txt") // just pass the file name
	if err != nil {
		log.Println(err)
	}

	bot.Handle("/update", "update")
	bot.Handle("/help", string(b))
	bot.Handle("/start", "start")
	bot.Handle("/infoluno", "infoluno")
	bot.Handle("/fee", "fee")
	bot.Handle("/convert", "convert")

	if err := bot.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
