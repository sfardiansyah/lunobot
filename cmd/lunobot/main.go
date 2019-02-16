package main

import (
	"html"
	"io/ioutil"
	"log"
	"os"

	"github.com/yanzay/tbot/model"

	"github.com/yanzay/tbot"
)

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"), tbot.WithWebhook("https://luno-bot.herokuapp.com/", ":"+os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/update", "update")
	bot.Handle("/help", fileReader("assets/help.txt"))
	bot.HandleFunc("/start", startHandler)
	bot.Handle("/infoluno", "infoluno")
	bot.HandleFunc("/fee", feeHandler)
	bot.Handle("/convert", fileReader("assets/convert.txt"))

	// bot.HandleDefault(defaultHandler)

	if err := bot.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func feeHandler(m *tbot.Message) {
	// buttons := []map[string]string{make(map[string]string)}
	// buttons := make([]map[string]string, 1, 1)
	// buttons[0]["fee"] = "https://www.luno.com/id/countries/ID"
	buttons := [][]string{
		{"Some", "Test", "Buttons"},
		{"Another", "Row"},
	}
	m.ReplyKeyboard("Buttons example", buttons)
	// m.ReplyInlineKeyboard(fileReader("assets/fee.txt"), buttons, )
}

// func defaultHandler(m *tbot.Message) {
// 	// if len(m.Vars["new_chat_members"]) > 0 {
// 	// 	m.Reply("Halo!")
// 	// }
// 	m.Reply(m.Data)
// }

func startHandler(m *tbot.Message) {
	if m.ChatType == model.ChatTypePrivate {
		m.Reply(fileReader("assets/start.txt"))
		m.Reply(fileReader("assets/help.txt"))
	}
}

func fileReader(dir string) string {
	b, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Println(err)
	}

	return html.UnescapeString(string(b))
}
