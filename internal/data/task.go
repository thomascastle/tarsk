package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/thomascastle/tarsk/internal/validator"
)

type Task struct {
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	DueAt       time.Time `json:"due_at"`
	ID          string    `json:"id"`
	Priority    string    `json:"priority"`
	StartedAt   time.Time `json:"started_at"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m TaskModel) Insert(task *Task) error {
	query := `
		INSERT INTO tasks (description, due_at, priority, started_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{task.Description, task.DueAt, task.Priority, task.StartedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID)
}

func (m TaskModel) Select() ([]*Task, error) {
	query := `
		SELECT description, done, due_at, id, priority, started_at
		FROM tasks`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, e := m.DB.QueryContext(ctx, query)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		var task Task
		e := rows.Scan(
			&task.Description,
			&task.Done,
			&task.DueAt,
			&task.ID,
			&task.Priority,
			&task.StartedAt,
		)
		if e != nil {
			return nil, e
		}

		tasks = append(tasks, &task)
	}

	if e := rows.Err(); e != nil {
		return nil, e
	}

	return tasks, nil
}

func (m TaskModel) SelectOne(id string) (*Task, error) {
	query := `
		SELECT description, done, due_at, id, priority, started_at
		FROM tasks
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var task Task
	e := m.DB.QueryRowContext(ctx, query, id).Scan(
		&task.Description,
		&task.Done,
		&task.DueAt,
		&task.ID,
		&task.Priority,
		&task.StartedAt,
	)
	if e != nil {
		switch {
		case errors.Is(e, sql.ErrNoRows):
			return nil, ErrorRecordNotFound
		default:
			return nil, e
		}
	}

	return &task, nil
}

func (m TaskModel) Update(task *Task) error {
	query := `
		UPDATE tasks
		SET description=$1, done=$2, due_at=$3, priority=$4, started_at=$5
		WHERE id=$6`

	args := []interface{}{
		task.Description,
		task.Done,
		task.DueAt,
		task.Priority,
		task.StartedAt,
		task.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, e := m.DB.ExecContext(ctx, query, args...)
	if e != nil {
		switch {
		case errors.Is(e, sql.ErrNoRows):
			return ErrorEditConflict
		default:
			return e
		}
	}

	return nil
}

func (m TaskModel) Delete(id string) error {
	query := `
		DELETE FROM tasks
		WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, e := m.DB.ExecContext(ctx, query, id)
	if e != nil {
		return e
	}

	rowsAffected, e := result.RowsAffected()
	if e != nil {
		return e
	}

	if rowsAffected == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Description != "", "description", "is required")
	v.Check(len(task.Description) <= 512, "description", "must not be more than 512 bytes long")

	v.Check(!task.DueAt.IsZero(), "due_at", "is required")

	supportedPriorities := []string{"none", "low", "medium", "high"}
	v.Check(validator.In(task.Priority, supportedPriorities...), "priority", "invalid priority value")

	v.Check(!task.StartedAt.IsZero(), "started_at", "is required")
	v.Check(!task.StartedAt.After(task.DueAt), "started_at", "date started must not be after due date")
}
