package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/botmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/controller"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/redisstorage"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/repositories"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	log.SetOutput(os.Stdout)

	if err := config.Init(); err != nil {
		log.Fatalf("Can't init config: %s", err.Error())
	}

	botAPI, err := tgbotapi.NewBotAPI(config.BotConfig.Token)

	if err != nil {
		log.Fatalf("Can't create new bot API: %s", err.Error())
	}

	botAPI.Debug = config.BotConfig.Debug

	if botAPI.Debug {
		botLogger := log.WithField("logger", "bot")

		if err := tgbotapi.SetLogger(botLogger); err != nil {
			log.Fatalf("Can't set bot logger: %s", err.Error())
		}
	}

	redisClient, err := redisstorage.NewRedisClient(config.RedisConfig)

	if err != nil {
		log.Fatalf("Can't create Redis client: %s", err.Error())
	}

	googleOauthConfig := &oauth2.Config{
		ClientID:     config.OauthConfig.ClientID,
		ClientSecret: config.OauthConfig.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  config.OauthConfig.RedirectURL,
		Scopes:       []string{calendar.CalendarScope, calendar.CalendarEventsScope},
	}
	botManagerLogger := log.WithField("logger", "bot_manager")
	messageParser := botmanager.NewMessageParser()
	messageComposer := botmanager.NewMessageComposer()
	calendarManager := calendarmanager.NewCalendarManager(googleOauthConfig)
	tokenRepository := repositories.NewTokenRepository(redisClient)
	userIDRepository := repositories.NewUserIDRepository(redisClient)
	botManager := botmanager.NewBotManager(
		config.BotConfig,
		botManagerLogger,
		botAPI,
		messageParser,
		messageComposer,
		calendarManager,
		userIDRepository,
		tokenRepository,
	)

	go func() {
		if err := botManager.Start(); err != nil {
			log.Fatalf("Can't start Bot Manager: %s", err.Error())
		}
	}()

	controllerLogger := log.WithField("logger", "controller")
	apiController := controller.NewController(controllerLogger, botManager)

	http.HandleFunc("/oauthCallback", apiController.HandleOauthCallback)

	if err := http.ListenAndServe(":"+strconv.Itoa(config.ServerConfig.Port), nil); err != nil {
		log.Fatalf("Can't start http server: %s", err.Error())
	}
}
