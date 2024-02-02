package utils

import (
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var ErrDuplicatedKey error = &DBError{Code: sqlite3.SQLITE_CONSTRAINT_UNIQUE}

type DBError struct {
	Code    int
	Message string
	Err     error
}

func (e *DBError) Error() string {
	return e.Message
}

func (e *DBError) Unwrap() error {
	return e.Err
}

func (e *DBError) Is(target error) bool {
	dbErr, ok := target.(*DBError)
	return ok && e.Code == dbErr.Code
}

func WrapDBErr(err error) error {
	if sqliteErr := new(sqlite.Error); errors.As(err, &sqliteErr) &&
		sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return &DBError{
			Code:    sqlite3.SQLITE_CONSTRAINT_UNIQUE,
			Message: "record already exists",
			Err:     err,
		}
	}

	return err
}
