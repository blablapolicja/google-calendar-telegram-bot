package calendarmanager

import "golang.org/x/oauth2"

// CalendarManager represents Calendar Manager
type CalendarManager struct {
	config *oauth2.Config
}

// NewCalendarManager creates new CalendarManager
func NewCalendarManager(config *oauth2.Config) *CalendarManager {
	return &CalendarManager{config}
}

// InitAuth - authorize user
func (m *CalendarManager) InitAuth() string {
	return m.config.AuthCodeURL("random-state-string")
}
