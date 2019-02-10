package main

import (
	"log"
	"time"

	telebot "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := telebot.NewBot(telebot.Settings{
		Token: "776702570:AAHnERHvgIwVFnc5M5WTXwJxCjfCSBoC8kg",
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		URL:    "https://luno-bot.herokuapp.com",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *telebot.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()
}