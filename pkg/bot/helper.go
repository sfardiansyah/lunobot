package bot

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func parse(m string) (string, []string) {
	vars := parseVariables(m)
	pattern := fmt.Sprintf("^%s$", replaceVariables(m))
	return pattern, vars
}

func parseVariables(pattern string) []string {
	var vars []string
	re := regexp.MustCompile("{([A-Za-z0-9_]*)}")
	matches := re.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		if len(match) > 0 {
			vars = append(vars, match[1])
		}
	}
	return vars
}

func replaceVariables(pattern string) string {
	re := regexp.MustCompile("{[A-Za-z0-9_]*}")
	return re.ReplaceAllString(pattern, "(.*)")
}

func fileReader(dir string) string {
	b, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Println(err)
	}

	return html.UnescapeString(string(b))
}

func inlineURLButtonsFromStrings(strs []map[string]string) [][]tgbotapi.InlineKeyboardButton {
	btns := make([][]tgbotapi.InlineKeyboardButton, len(strs))
	for i, buttonRow := range strs {
		btnsRow := []tgbotapi.InlineKeyboardButton{}
		for buttonText, buttonURL := range buttonRow {
			btnsRow = append(btnsRow, tgbotapi.NewInlineKeyboardButtonURL(buttonText, buttonURL))
		}
		btns[i] = tgbotapi.NewInlineKeyboardRow(btnsRow...)
	}
	return btns
}
