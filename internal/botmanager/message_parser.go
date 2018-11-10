package botmanager

import "github.com/go-telegram-bot-api/telegram-bot-api"

// CommandType - type of command
type CommandType int

// Command represents command that was generated from user's message
type Command struct {
	commandType CommandType
	replyTo int64
}

func newCommand(commandType CommandType, replyTo int64) *Command {
	return &Command{commandType, replyTo}
}

const (
	// Authorise - command for authorizing client
	Authorise CommandType = 0
)

type messageParser struct{}

// NewMessageParser creates new Message Parser
func NewMessageParser() *messageParser {
	return &messageParser{}
}

// ParseMessage parses message from client
func (mp *messageParser) ParseMessage(m *tgbotapi.Message) *Command {
	return newCommand(Authorise, m.Chat.ID)
}
