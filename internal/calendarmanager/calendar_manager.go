package calendarmanager

import (
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

// CalendarManager represents Calendar Manager
type CalendarManager struct {
	config *oauth2.Config
}

// NewCalendarManager creates new CalendarManager
func NewCalendarManager(config *oauth2.Config) *CalendarManager {
	return &CalendarManager{config}
}

// CreateAuthCodeURL - authorize user
func (m *CalendarManager) CreateAuthCodeURL(state string) string {
	return m.config.AuthCodeURL(state)
}

// CreateToken creates oauth2 token from code
func (m *CalendarManager) CreateToken(code string) (*oauth2.Token, error) {
	return m.config.Exchange(oauth2.NoContext, code)
}

// CreateClient creates new Google Calendar client
func (m *CalendarManager) CreateClient(token *oauth2.Token) (*calendar.Service, error) {
	httpClient := m.config.Client(oauth2.NoContext, token)

	return calendar.New(httpClient)
}

func (m *CalendarManager) GetCalendarEvents(calendarClient *calendar.Service) (*calendar.Events, error) {
	t := time.Now().Format(time.RFC3339)

	return calendarClient.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
}
