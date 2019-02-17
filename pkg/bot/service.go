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
	msg := tgbotapi.NewMessage(m.Chat.ID, "Halo")
	if _, err := h.a.Send(msg); err != nil {
		log.Fatal(err)
	}
}
