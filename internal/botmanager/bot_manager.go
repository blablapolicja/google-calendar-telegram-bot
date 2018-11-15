package botmanager

import (
	"strconv"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/util"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// BotManager represents Bot Manager
type BotManager struct {
	config                 config.BotConf
	logger                 *log.Entry
	botAPI                 *tgbotapi.BotAPI
	messageParser          *messageParser
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
	calendarManager *calendarmanager.CalendarManager,
	authorizedUsersCache *cache.Cache,
	unauthorizedUsersCache *cache.Cache,
) *BotManager {
	return &BotManager{
		config:                 config,
		logger:                 logger,
		botAPI:                 botAPI,
		messageParser:          messageParser,
		calendarManager:        calendarManager,
		authorizedUsersCache:   authorizedUsersCache,
		unauthorizedUsersCache: unauthorizedUsersCache,
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
		b.startAuth(command.replyTo)
	}
}

func (b *BotManager) startAuth(userID int64) {
	state := util.GetRandomStateString()
	authURL := b.calendarManager.GetAuthCodeURL(state)
	message := tgbotapi.NewMessage(userID, "Open this link to connect your Google Calendar:\n\n"+authURL+"\n\n(link will expire in 1 minute)")

	b.unauthorizedUsersCache.Set(state, userID, cache.DefaultExpiration)
	b.botAPI.Send(message)
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
	token, err := b.calendarManager.GetToken(code)

	if err != nil {
		b.logger.Errorf("Error while getting Oauth token: %s", err.Error())
		return
	}

	b.unauthorizedUsersCache.Delete(state)
	b.authorizedUsersCache.Set(strconv.FormatInt(userID, 10), token, cache.DefaultExpiration)

	message := tgbotapi.NewMessage(userID, "You have been successfully authorized!")
	b.botAPI.Send(message)
	b.logger.Infof("User %d has been authorized", userID)
}
