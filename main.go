package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/botmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/controller"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/database"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
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

	dbConnection, err := database.NewDatabaseConn(config.DatabaseConfig)

	if err != nil {
		log.Fatalf("Can't connect to database: %s", err.Error())
	}

	log.WithField("logger", "init").Infof("Connected to database: %s:%d", config.DatabaseConfig.Host, config.DatabaseConfig.Port)

	if err := database.CreateTables(dbConnection); err != nil {
		log.Fatalf("Can't create tables in database: %s", err.Error())
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
	authorizedUsersCache := cache.New(0, 0)
	unauthorizedUsersCache := cache.New(time.Minute, 5*time.Minute)
	botManager := botmanager.NewBotManager(
		config.BotConfig,
		botManagerLogger,
		botAPI,
		messageParser,
		messageComposer,
		calendarManager,
		authorizedUsersCache,
		unauthorizedUsersCache,
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
