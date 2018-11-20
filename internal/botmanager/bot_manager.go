package botmanager

import (
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/repositories"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/util"

	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// BotManager represents Bot Manager
type BotManager struct {
	config           config.BotConf
	logger           *log.Entry
	botAPI           *tgbotapi.BotAPI
	messageParser    *messageParser
	messageComposer  *messageComposer
	calendarManager  *calendarmanager.CalendarManager
	userIDRepository *repositories.UserIDRepository
	tokenRepository  *repositories.TokenRepository
}

// NewBotManager creates new BotManager
func NewBotManager(
	config config.BotConf,
	logger *log.Entry,
	botAPI *tgbotapi.BotAPI,
	messageParser *messageParser,
	messageComposer *messageComposer,
	calendarManager *calendarmanager.CalendarManager,
	userIDRepository *repositories.UserIDRepository,
	tokenRepository *repositories.TokenRepository,
) *BotManager {
	return &BotManager{
		config,
		logger,
		botAPI,
		messageParser,
		messageComposer,
		calendarManager,
		userIDRepository,
		tokenRepository,
	}
}

// Start starts Bot Manager
func (b *BotManager) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botAPI.GetUpdatesChan(u)

	if err != nil {
		return err
	}

	b.logger.Infof("Bot authorized on account %s", b.botAPI.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}

	return nil
}

func (b *BotManager) handleMessage(m *tgbotapi.Message) {
	operation := b.messageParser.ParseMessage(m)

	if !operation.needAuth {
		switch operation.operationType {
		case operationAuthorise:
			b.startAuth(operation.userID)
		case operationUnknown:
		default:
			b.sendDefaultResponse(operation.userID)
		}
		return
	}

	token, err := b.tokenRepository.Get(operation.userID)

	if err == redis.Nil {
		b.startAuth(operation.userID)
		return
	}

	if err != nil {
		b.logger.Errorf("Error while getting user ID: %s", err.Error())
		return
	}

	switch operation.operationType {
	case operationEvents:
		b.sendCalendarEvents(operation.userID, token, operation.params)
	default:
		b.sendDefaultResponse(operation.userID)
	}
}

func (b *BotManager) startAuth(userID int64) {
	state := util.GenerateStateString()
	authURL := b.calendarManager.CreateAuthCodeURL(state)
	message := tgbotapi.NewMessage(userID, "Open this [link]("+authURL+") to connect your Google Calendar (link will expire in 1 minute)")
	message.ParseMode = tgbotapi.ModeMarkdown

	if err := b.userIDRepository.Save(state, userID); err != nil {
		b.logger.Errorf("Error while saving user ID %s", err.Error())
		return
	}

	if _, err := b.botAPI.Send(message); err != nil {
		b.logger.Errorf("Error while sending authorization link to user %s", err.Error())
		return
	}

	b.logger.Infof("Authorization link was sent to user %d", userID)
}

// FinishAuth saves client oauth2 token
func (b *BotManager) FinishAuth(state string, code string) {
	userID, err := b.userIDRepository.Get(state)

	if err == redis.Nil {
		b.logger.Warnf("User with state %s was not found", state)
		return
	}

	if err != nil {
		b.logger.Errorf("Error while getting user ID: %s", err.Error())
		return
	}

	if err := b.userIDRepository.Delete(state); err != nil {
		b.logger.Errorf("Error while deleting user ID %s", err.Error())
	}

	token, err := b.calendarManager.CreateToken(code)

	if err != nil {
		b.logger.Errorf("Error while creating Oauth token: %s", err.Error())
		return
	}

	if err := b.tokenRepository.Save(userID, token); err != nil {
		b.logger.Errorf("Error while saving Oauth token: %s", err.Error())
		return
	}

	message := tgbotapi.NewMessage(userID, "You have been successfully authorized!")

	if _, err := b.botAPI.Send(message); err != nil {
		b.logger.Errorf("Error while sending authorization message to user %s", err.Error())
		return
	}

	b.logger.Infof("User %d has been authorized", userID)
	b.sendCalendarEvents(userID, token, nil)
}

func (b *BotManager) sendCalendarEvents(userID int64, token *oauth2.Token, params interface{}) {
	calendarClient, err := b.calendarManager.CreateClient(token)

	if err != nil {
		b.logger.Errorf("Can't create Google Calendar client for user %s", err.Error())
		return
	}

	events, err := b.calendarManager.GetCalendarEvents(calendarClient, params)

	if err != nil {
		b.logger.Errorf("Error while getting calendar events for user %s", err.Error())
		return
	}

	message := b.messageComposer.CreateEventsList(userID, events.Items)

	if _, err := b.botAPI.Send(message); err != nil {
		b.logger.Errorf("Error while sending calendar events to user %s", err.Error())
	}
}

func (b *BotManager) sendDefaultResponse(userID int64) {
	message := tgbotapi.NewMessage(userID, "I'm not sure I know what do you want")

	if _, err := b.botAPI.Send(message); err != nil {
		b.logger.Errorf("Error while sending default response to user %s", err.Error())
	}
}
