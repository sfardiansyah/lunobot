package bot

import (
	"log"

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
		// pattern, variables := parse(u.Message.Text)
	}
}

func (h *handler) handleJoin(m *tgbotapi.Message) {
	if err := h.replyWithFile(m.Chat.ID, "assets/join.txt"); err != nil {
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
