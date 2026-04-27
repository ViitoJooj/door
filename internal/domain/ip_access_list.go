package domain

import "time"

type IPAccessEntry struct {
	ID        int
	IP        string
	CreatedBy int
	UpdatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
}
