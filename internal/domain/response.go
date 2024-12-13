package domain

import "time"

// Response response when a user enters a prompt
type Response struct {
	Text         string
	StartTime    time.Time
	EndTime      time.Time
	ResponseTime time.Duration
}
