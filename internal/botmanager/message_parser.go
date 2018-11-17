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
	// Authorise - command for authorizing client
	Authorise = 0
	// Unknown - command to identify unknown Bot command from user
	Unknown = 666
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
func (mp *messageParser) ParseMessage(m *tgbotapi.Message) *Operation {
	if m.IsCommand() {
		switch m.Command() {
		case start:
			return newOperation(Authorise, m.Chat.ID)
		default:
			return newOperation(Unknown, m.Chat.ID)
		}
	}

	return newOperation(Unknown, m.Chat.ID)
}
