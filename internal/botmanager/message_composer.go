package botmanager

import (
	"bytes"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/calendar/v3"
)

type messageComposer struct{}

// NewMessageComposer creates new MessageComposer
func NewMessageComposer() *messageComposer {
	return &messageComposer{}
}

// CreateEventsList creates message with Google Calendar events
func (messageComposer) CreateEventsList(userID int64, events []*calendar.Event) tgbotapi.MessageConfig {
	var buffer bytes.Buffer

	buffer.WriteString("Your nearest events:\n\n")

	for _, event := range events {
		date := event.Start.DateTime

		if date == "" {
			date = event.Start.Date
		}

		buffer.WriteString(date)
		buffer.WriteString(" - ")
		buffer.WriteString(event.Summary)
		buffer.WriteString("\n")
	}

	return tgbotapi.NewMessage(userID, buffer.String())
}
