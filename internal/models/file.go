package models

import "time"

type File struct {
	ID          int
	UserID      int64
	TimeCreated time.Time
	Name        string
	TimesViewed int64
}
