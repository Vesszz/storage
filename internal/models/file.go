package models

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	ID          int
	UserID      int64
	Key         uuid.UUID
	Path        string
	TimeCreated time.Time
	Name        string
	Description string
	TimesViewed int64
}
