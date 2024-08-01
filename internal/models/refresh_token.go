package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Key         uuid.UUID
	UserID      int64
	Fingerprint string
	ExpiresAt   time.Time
}
