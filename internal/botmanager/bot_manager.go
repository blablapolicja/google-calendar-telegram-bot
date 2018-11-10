package botmanager

import (
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// BotManager represents Bot Manager
type BotManager struct {
	config          config.BotConf
	logger          *log.Entry
	botAPI          *tgbotapi.BotAPI
	messageParser   *messageParser
	calendarManager *calendarmanager.CalendarManager
}

// NewBotManager creates new BotManager
func NewBotManager(
	config config.BotConf,
	logger *log.Entry,
	botAPI *tgbotapi.BotAPI,
	messageParser *messageParser,
	calendarManager *calendarmanager.CalendarManager,
) *BotManager {
	return &BotManager{
		config:          config,
		logger:          logger,
		botAPI:          botAPI,
		messageParser:   messageParser,
		calendarManager: calendarManager,
	}
}

// Start starts Bot Manager
func (b *BotManager) Start() error {
	b.botAPI.Debug = config.BotConfig.Debug

	b.logger.Infof("Authorized on account %s", b.botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botAPI.GetUpdatesChan(u)

	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}

	return nil
}

func (b *BotManager) handleMessage(m *tgbotapi.Message) {
	b.logger.Info(m.From.UserName, m.Text)

	command := b.messageParser.ParseMessage(m)

	if command.commandType == Authorise {
		authURL := b.calendarManager.InitAuth()
		message := tgbotapi.NewMessage(command.replyTo, authURL)

		b.botAPI.Send(message)
	}
}

// FinishAuth saves client oauth2 token
func (b *BotManager) FinishAuth() {

}
