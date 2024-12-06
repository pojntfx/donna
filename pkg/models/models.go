// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

import (
	"database/sql"
	"time"
)

type Activity struct {
	ID          int32
	Name        string
	Date        time.Time
	ContactID   int32
	Description string
}

type Contact struct {
	ID        int32
	FirstName string
	LastName  string
	Nickname  string
	Email     string
	Pronouns  string
	Namespace string
	Birthday  sql.NullTime
	Address   string
	Notes     string
}

type Debt struct {
	ID          int32
	Amount      float64
	Currency    string
	ContactID   int32
	Description string
}

type JournalEntry struct {
	ID        int32
	Title     string
	Date      time.Time
	Body      string
	Rating    int32
	Namespace string
}

type Todo struct {
	ID         int32
	Name       string
	Deadline   time.Time
	Importance int32
	Pending    bool
	Namespace  string
}
