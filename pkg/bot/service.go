package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handler ...
type Handler interface {
	Handle(tgbotapi.Update)
}

type handler struct {
	a *tgbotapi.BotAPI
}

// NewHandler ...
func NewHandler(a *tgbotapi.BotAPI) Handler {
	return &handler{a}
}

func (h *handler) Handle(u tgbotapi.Update) {
	if u.Message != nil {
		if u.Message.NewChatMembers != nil {
			h.handleJoin(u.Message)
			return
		}
		pattern, _ := parse(h.trimBotName(u.Message.Text))
		log.Println(pattern)
		cID := u.Message.Chat.ID

		switch pattern {
		case "/start":
			h.handlerStart(u.Message)
		case "/fee":
			h.replyWithInline(cID, fileReader("assets/fee.txt"), "Kunjungi Rincian Biaya Luno")
		case "/convert":
			h.replyWithInline(cID, fileReader("assets/convert.txt"), "Luno Price Chart")
		case "/help":
			h.replyText(cID, fileReader("assets/help.txt"))
		}
	}
}

func (h *handler) handlerStart(m *tgbotapi.Message) {
	if m.Chat.IsPrivate() {
		h.replyText(m.Chat.ID, fileReader("assets/start.txt"))
		h.replyText(m.Chat.ID, fileReader("assets/start.txt"))
	}
}

func (h *handler) handleJoin(m *tgbotapi.Message) {
	var arr []string
	for _, member := range *m.NewChatMembers {
		arr = append(arr, member.FirstName)
	}

	members := strings.Join(arr, ", ")
	str := strings.Replace(fileReader("assets/join.txt"), "[Name]", members, 1)

	if err := h.replyWithInline(m.Chat.ID, str, "Daftar Luno"); err != nil {
		log.Fatal(err)
	}
}

func (h *handler) replyWithInline(cID int64, text, helper string) error {
	str := strings.Split(text, "||")
	btn := []map[string]string{map[string]string{helper: str[1]}}

	if err := h.replyInlineKeyboard(cID, str[0], btn); err != nil {
		return err
	}
	return nil
}

func (h *handler) replyInlineKeyboard(cID int64, text string, buttons []map[string]string) error {
	msg := tgbotapi.NewMessage(cID, text)
	btns := inlineURLButtonsFromStrings(buttons)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(btns...)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := h.a.Send(msg); err != nil {
		return err
	}
	return nil
}

func (h *handler) replyText(cID int64, text string) {
	msg := tgbotapi.NewMessage(cID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := h.a.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func (h *handler) replyWithFile(cID int64, dir string) error {
	msg := tgbotapi.NewMessage(cID, fileReader(dir))
	if _, err := h.a.Send(msg); err != nil {
		return err
	}
	return nil
}

func (h *handler) trimBotName(message string) string {
	parts := strings.SplitN(message, " ", 2)
	command := parts[0]
	command = strings.TrimSuffix(command, "@"+h.a.Self.UserName)
	command = strings.TrimSuffix(command, "@"+h.a.Self.FirstName)
	parts[0] = command
	return strings.Join(parts, " ")
}
