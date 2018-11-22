package botmanager

import (
	"bytes"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/calendar/v3"
)

type messageComposer struct{}

// NewMessageComposer creates new MessageComposer
func NewMessageComposer() *messageComposer {
	return &messageComposer{}
}

func (messageComposer) CreateAuthMessage(userID int64, authURL string) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(userID, "Open this [link]("+authURL+") to connect your Google Calendar (link will expire in 1 minute)")
	message.ParseMode = tgbotapi.ModeMarkdown

	return message
}

// CreateEventsList creates message with Google Calendar events
func (messageComposer) CreateEventsList(userID int64, events []*calendar.Event) tgbotapi.MessageConfig {
	var buffer bytes.Buffer

	buffer.WriteString("Your events:\n\n")

	for _, event := range events {
		var date string

		if event.Start.DateTime != "" {
			date = event.Start.DateTime
			dateParsed, err := time.Parse(time.RFC3339, date)

			if err != nil {
				//TODO: add error handling
				continue
			}

			buffer.WriteString(dateParsed.Format(time.RubyDate)[0:16])
		} else {
			date = event.Start.Date

			buffer.WriteString(date)
		}

		buffer.WriteString(" - ")
		buffer.WriteString(event.Summary)
		buffer.WriteString("\n")
	}

	return tgbotapi.NewMessage(userID, buffer.String())
}
