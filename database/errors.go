package database

import (
	"github.com/mattn/go-sqlite3"
)

// IsConstraintViolationError checks if err is a constraint violation error
func IsConstraintViolationError(err error) bool {
	if e, ok := err.(sqlite3.Error); ok {
		return e.Code == sqlite3.ErrConstraint
	}
	return false
}
