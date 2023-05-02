package rdbms

import (
	"database/sql"
	"errors"
	"strings"
)

type RDBMS interface {
	MigrateUp(source string) error

	MigrateDown(source string) error

	Create(query string, args []interface{}) (uint64, error)

	Read(query string, args []interface{}, dest ...interface{}) error

	Update(query string, args []interface{}) error

	Delete(query string, args []interface{}) error
}

type rdbms struct {
	db *sql.DB
}

const (
	errPrepareStatement = "error when tying to prepare statement"

	errCreate    = "error when tying to create entry"
	errDuplicate = "entry exists"

	errRead         = "error when tying to read entry"
	errReadNotFound = "there is no entry with provided arguments"

	errUpdate = "error when tying to update entry"

	errDelete = "error when tying to delete entry"
)

func (db *rdbms) Create(query string, args []interface{}) (uint64, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return 0, errors.New(errPrepareStatement)
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(args)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return 0, errors.New(errDuplicate)
		}

		return 0, errors.New(errCreate)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return 0, errors.New(errCreate)
	}

	return uint64(id), nil
}

func (db *rdbms) Read(query string, args []interface{}, dest ...interface{}) error {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return errors.New(errPrepareStatement)
	}
	defer stmt.Close()

	result := stmt.QueryRow(args)
	err = result.Scan(dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(errReadNotFound)
		}

		return errors.New(errRead)
	}

	return nil
}

func (db *rdbms) Update(query string, args []interface{}) error {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return errors.New(errPrepareStatement)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(args); err != nil {
		return errors.New(errUpdate)
	}

	return nil
}

func (db *rdbms) Delete(query string, args []interface{}) error {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return errors.New(errPrepareStatement)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(args); err != nil {
		return errors.New(errDelete)
	}

	return nil
}
