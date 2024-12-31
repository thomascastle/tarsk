package data

import (
	"database/sql"
	"errors"
)

var (
	ErrorEditConflict   = errors.New("edit conflict")
	ErrorRecordNotFound = errors.New("record was not found")
)

type Models struct {
	Tasks TaskModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tasks: TaskModel{DB: db},
	}
}
