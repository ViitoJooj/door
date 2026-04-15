package dtos

import "time"

type RequestLogData struct {
	ID             int       `json:"id"`
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	QueryString    string    `json:"query_string"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	IP             string    `json:"ip"`
	Country        string    `json:"country"`
	UserAgent      string    `json:"user_agent"`
	Referer        string    `json:"referer"`
	RequestSize    int       `json:"request_size"`
	ResponseSize   int       `json:"response_size"`
	Internal       bool      `json:"internal"`
	CreatedAt      time.Time `json:"created_at"`
}

type RequestLogListOutput struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []RequestLogData `json:"data"`
}
