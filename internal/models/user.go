package models

import "time"

type User struct {
	ID          int64
	Name        string
	Password    string
	TimeCreated time.Time
}
