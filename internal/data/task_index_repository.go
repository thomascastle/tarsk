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

func (r TaskIndexRepository) Select(filters *Filters) ([]*Task, error) {
	var tasks []*Task

	r.db.Where(map[string]interface{}(*filters)).
		Select([]string{"description", "done", "due_at", "id", "priority", "started_at"}).
		Find(&tasks)

	return tasks, nil
}
