package domain

import "time"

type RequestLog struct {
	ID             int
	Method         string
	Path           string
	QueryString    string
	StatusCode     int
	ResponseTimeMs int64
	IP             string
	Country        string
	UserAgent      string
	Referer        string
	RequestSize    int
	ResponseSize   int
	Internal       bool
	CreatedAt      time.Time
}
