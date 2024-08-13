package models

import "time"

type Comment struct {
	ID          int
	UserID      int
	FileID      int
	TimeCreated time.Time
	Text        string
}
