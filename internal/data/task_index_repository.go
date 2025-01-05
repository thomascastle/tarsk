package data

import (
	"gorm.io/gorm"
)

type TaskIndexRepository struct {
	db *gorm.DB
}

func NewTaskIndexRepository(db *gorm.DB) TaskIndexRepository {
	return TaskIndexRepository{
		db: db,
	}
}

func (r TaskIndexRepository) Select(filters *Filters, search string) ([]*Task, error) {
	var tasks []*Task

	r.db.
		Select([]string{"description", "done", "due_at", "id", "priority", "started_at"}).
		Where(
			r.db.Where("to_tsvector('simple', description) @@ plainto_tsquery('simple', ?)", search).Or("?=''", search),
		).
		Where(map[string]interface{}(*filters)).
		Find(&tasks)

	return tasks, nil
}
