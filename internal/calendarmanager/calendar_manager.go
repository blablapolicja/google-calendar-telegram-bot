package calendarmanager

import (
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

const (
	intervalDay   = "day"
	intervalWeek  = "week"
	intervalMonth = "month"
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
	return m.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
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

// GetCalendarEvents returns events from Calendar
func (m *CalendarManager) GetCalendarEvents(calendarClient *calendar.Service, params interface{}) (*calendar.Events, error) {
	switch params.(type) {
	case string:
		return m.getCalendarEventsByInterval(calendarClient, params.(string))
	}

	return m.getNearestEvents(calendarClient)
}

func (m *CalendarManager) getNearestEvents(calendarClient *calendar.Service) (*calendar.Events, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)

	return calendarClient.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(startOfDay).MaxResults(10).OrderBy("startTime").Do()
}

func (m *CalendarManager) getCalendarEventsByInterval(calendarClient *calendar.Service, interval string) (*calendar.Events, error) {
	var timeMax string
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch interval {
	case intervalDay:
		timeMax = startOfDay.AddDate(0, 0, 1).Format(time.RFC3339)
	case intervalWeek:
		timeMax = startOfDay.AddDate(0, 0, 7).Format(time.RFC3339)
	case intervalMonth:
		timeMax = startOfDay.AddDate(0, 1, 0).Format(time.RFC3339)
	default:
		return m.getNearestEvents(calendarClient)
	}

	return calendarClient.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(startOfDay.Format(time.RFC3339)).TimeMax(timeMax).OrderBy("startTime").Do()
}
