package main

import "github.com/stefanmcshane/commentsold/db"

// dbHandler is used to pass the db connection between endpoints
type dbHandler struct {
	db *db.CSVDatabase // This should be changed to an interface for easier mocking
}

// NewDBHandler instatiates the dbhandler with the passed in database
func NewDBHandler(db *db.CSVDatabase) *dbHandler {
	return &dbHandler{db: db}
}
