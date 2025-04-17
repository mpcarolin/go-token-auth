package models

import "api/internal/db"

// Struct representing app's database queries and a function for closing it
type AppDB struct {
	Queries *db.Queries
	Close func()
} 