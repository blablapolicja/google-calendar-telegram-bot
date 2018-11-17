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
		state := query["state"]
		code := query["code"]

		if len(state) != 1 || len(code) != 1 {
			c.logger.Warnf("Oauth callback was called with code %s and state %s", code, state)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		c.botManager.FinishAuth(query["state"][0], query["code"][0])
	default:
		c.logger.Warnf("Oauth callback was called with wrong method: %s", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
