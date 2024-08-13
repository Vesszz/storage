package models

type Collection struct {
	ID     int
	UserID int
	Name   string
}

type CollectionFile struct {
	CollectionID int
	FileID       int
}
