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

// GetAuthCodeURL - authorize user
func (m *CalendarManager) GetAuthCodeURL(state string) string {
	return m.config.AuthCodeURL(state)
}

// GetToken - get oauth2 token from code
func (m *CalendarManager) GetToken(code string) (*oauth2.Token, error) {
	return m.config.Exchange(oauth2.NoContext, code)
}
