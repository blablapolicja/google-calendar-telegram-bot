package controller

import (
	"net/http"

	"github.com/blablapolicja/google-calendar-telegram-bot/internal/botmanager"

	log "github.com/sirupsen/logrus"
)

// Controller represents API controller
type Controller struct {
	logger     *log.Entry
	botManager *botmanager.BotManager
}

// NewController create new API controller
func NewController(logger *log.Entry, botManager *botmanager.BotManager) *Controller {
	return &Controller{logger, botManager}
}

// HandleOauthCallback handles Google Oauth2 callback call
func (c *Controller) HandleOauthCallback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()

		c.logger.Printf("Oauth callback is called with code %s and state %s", query["code"], query["state"])
	default:
		c.logger.Warnf("HandleOauthCallback was called with wrong method: %s", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
