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
		// pattern, variables := parse(u.Message.Text)
	}
}

func (h *handler) handleJoin(m *tgbotapi.Message) {
	var arr []string
	for _, member := range *m.NewChatMembers {
		arr = append(arr, member.FirstName)
	}

	members := strings.Join(arr, ", ")
	str := strings.Replace(fileReader("assets/join.txt"), "[Name]", members, 1)

	if err := h.replyText(m.Chat.ID, str); err != nil {
		log.Fatal(err)
	}
}

func (h *handler) replyText(cID int64, text string) error {
	msg := tgbotapi.NewMessage(cID, text)
	if _, err := h.a.Send(msg); err != nil {
		return err
	}
	return nil
}

func (h *handler) replyWithFile(cID int64, dir string) error {
	msg := tgbotapi.NewMessage(cID, fileReader(dir))
	if _, err := h.a.Send(msg); err != nil {
		return err
	}
	return nil
}
