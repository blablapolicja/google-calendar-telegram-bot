package botmanager

import "github.com/go-telegram-bot-api/telegram-bot-api"

// Operation represents command that was generated from user's message
type Operation struct {
	operationType int
	userID        int64
	needAuth      bool
	params        interface{}
}

func newOperation(operationType int, userID int64, needAuth bool, params interface{}) *Operation {
	return &Operation{operationType, userID, needAuth, params}
}

const (
	// operation types
	operationAuthorise = 0
	operationEvents    = 1
	operationUnknown   = 666

	// commands available for users
	commandStart = "start"
	commandList  = "list"

	// arguments for commands
	argumentDay   = "day"
	argumentWeek  = "week"
	argumentMonth = "month"
)

type messageParser struct{}

// NewMessageParser creates new Message Parser
func NewMessageParser() *messageParser {
	return &messageParser{}
}

// ParseMessage parses message from user
func (p *messageParser) ParseMessage(m *tgbotapi.Message) *Operation {
	switch command := m.CommandWithAt(); command {
	case commandStart:
		return newOperation(operationAuthorise, m.Chat.ID, false, nil)
	case commandList:
		return newOperation(operationEvents, m.Chat.ID, true, m.CommandArguments())
	default:
		return newOperation(operationUnknown, m.Chat.ID, false, nil)
	}
}
