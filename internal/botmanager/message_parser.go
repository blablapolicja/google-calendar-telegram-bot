package botmanager

import "github.com/go-telegram-bot-api/telegram-bot-api"

// Operation represents command that was generated from user's message
type Operation struct {
	operationType int
	userID        int64
}

func newOperation(operationType int, userID int64) *Operation {
	return &Operation{operationType, userID}
}

const (
	authorise = 0
	unknown = 666
)

type messageParser struct{}

// NewMessageParser creates new Message Parser
func NewMessageParser() *messageParser {
	return &messageParser{}
}

// commands available in Bot
const (
	start = "start"
)

// ParseMessage parses message from user
func (messageParser) ParseMessage(m *tgbotapi.Message) *Operation {
	if m.IsCommand() {
		switch m.Command() {
		case start:
			return newOperation(authorise, m.Chat.ID)
		default:
			return newOperation(unknown, m.Chat.ID)
		}
	}

	return newOperation(unknown, m.Chat.ID)
}
