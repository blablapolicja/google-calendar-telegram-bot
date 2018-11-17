package botmanager

import (
	"strconv"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/util"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// BotManager represents Bot Manager
type BotManager struct {
	config                 config.BotConf
	logger                 *log.Entry
	botAPI                 *tgbotapi.BotAPI
	messageParser          *messageParser
	messageComposer        *messageComposer
	calendarManager        *calendarmanager.CalendarManager
	authorizedUsersCache   *cache.Cache
	unauthorizedUsersCache *cache.Cache
}

// NewBotManager creates new BotManager
func NewBotManager(
	config config.BotConf,
	logger *log.Entry,
	botAPI *tgbotapi.BotAPI,
	messageParser *messageParser,
	messageComposer *messageComposer,
	calendarManager *calendarmanager.CalendarManager,
	authorizedUsersCache *cache.Cache,
	unauthorizedUsersCache *cache.Cache,
) *BotManager {
	return &BotManager{
		config:                 config,
		logger:                 logger,
		botAPI:                 botAPI,
		messageParser:          messageParser,
		messageComposer:        messageComposer,
		calendarManager:        calendarManager,
		authorizedUsersCache:   authorizedUsersCache,
		unauthorizedUsersCache: unauthorizedUsersCache,
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

	switch operation.operationType {
	case Authorise:
		b.startAuth(operation.userID)
	default:
		b.sendDefaultResponse(operation.userID)
	}
}

func (b *BotManager) startAuth(userID int64) {
	state := util.GenerateStateString()
	authURL := b.calendarManager.CreateAuthCodeURL(state)
	message := tgbotapi.NewMessage(userID, "Open this [link]("+authURL+") to connect your Google Calendar (link will expire in 1 minute)")
	message.ParseMode = tgbotapi.ModeMarkdown

	b.unauthorizedUsersCache.Set(state, userID, cache.DefaultExpiration)

	_, err := b.botAPI.Send(message)

	if err != nil {
		b.logger.Errorf("Error while sending authorization link to user %s", err.Error())
		return
	}

	b.logger.Infof("Authorization link was sent to user %d", userID)
}

// FinishAuth saves client oauth2 token
func (b *BotManager) FinishAuth(state string, code string) {
	value, found := b.unauthorizedUsersCache.Get(state)

	if !found {
		b.logger.Warnf("User with state %s and code %s was not found", state, code)
		return
	}

	userID := value.(int64)
	token, err := b.calendarManager.CreateToken(code)

	if err != nil {
		b.logger.Errorf("Error while getting Oauth token: %s", err.Error())
		return
	}

	b.unauthorizedUsersCache.Delete(state)
	b.authorizedUsersCache.Set(strconv.FormatInt(userID, 10), token, cache.DefaultExpiration)

	message := tgbotapi.NewMessage(userID, "You have been successfully authorized!")

	// TODO: add error handling
	b.botAPI.Send(message)
	b.logger.Infof("User %d has been authorized", userID)

	b.sendCalendarEvents(userID, token)
}

func (b *BotManager) sendCalendarEvents(userID int64, token *oauth2.Token) {
	calendarClient, err := b.calendarManager.CreateClient(token)

	if err != nil {
		b.logger.Errorf("Can't create Google Calendar client for user %d", userID)
		return
	}

	events, err := b.calendarManager.GetCalendarEvents(calendarClient)

	if err != nil {
		b.logger.Errorf("Error while getting calendar events for user %d", userID)
		return
	}

	message := b.messageComposer.CreateEventsList(userID, events.Items)

	// TODO: add error handling
	b.botAPI.Send(message)
}

func (b *BotManager) sendDefaultResponse(userID int64) {
	message := tgbotapi.NewMessage(userID, "I'm not sure I know what do you want")

	// TODO: add error handling
	b.botAPI.Send(message)
}
