package main

import (
	"html"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	bot.HandleFunc("/fee", func(m *tbot.Message) {
		feeHandler(m, fileReader("assets/fee.txt"), "Kunjungi Rincian Biaya LUNO")
	})
	bot.Handle("/convert", fileReader("assets/convert.txt"))

	// bot.HandleDefault(defaultHandler)

	if err := bot.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func feeHandler(m *tbot.Message, t, helper string) {
	str := strings.Split(t, "||")
	btn := []map[string]string{map[string]string{helper: str[1]}}

	m.ReplyInlineKeyboard(str[0], btn, tbot.WithURLInlineButtons)
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
