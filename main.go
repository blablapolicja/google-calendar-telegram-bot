package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/botmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/calendarmanager"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/config"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/controller"
	"github.com/blablapolicja/google-calendar-telegram-bot/internal/database"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
		Scopes:       []string{"https://www.googleapis.com/auth/calendar"},
	}
	botManagerLogger := log.WithField("logger", "bot_manager")
	messageParser := botmanager.NewMessageParser()
	calendarManager := calendarmanager.NewCalendarManager(googleOauthConfig)
	botManager := botmanager.NewBotManager(
		config.BotConfig,
		botManagerLogger,
		botAPI,
		messageParser,
		calendarManager,
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
