package data

import (
	"database/sql"
	"errors"
)

var (
	ErrorEditConflict   = errors.New("edit conflict")
	ErrorRecordNotFound = errors.New("record was not found")
)

type Repositories struct {
	Tasks TaskRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Tasks: TaskRepository{DB: db},
	}
}
