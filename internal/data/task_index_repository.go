package data

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskIndexRepository struct {
	db *gorm.DB
}

func NewTaskIndexRepository(db *gorm.DB) TaskIndexRepository {
	return TaskIndexRepository{
		db: db,
	}
}

func (r TaskIndexRepository) Select(search string, filters *Filters, sort Sort, paginator Paginator) ([]*Task, Pagination, error) {
	var tasks []*Task
	r.db.
		Select([]string{"description", "done", "due_at", "id", "priority", "started_at"}).
		Where(
			r.db.Where("to_tsvector('simple', description) @@ plainto_tsquery('simple', ?)", search).Or("?=''", search),
		).
		Where(map[string]interface{}(*filters)).
		Order(clause.OrderByColumn{Column: clause.Column{Name: sort.sortColumn()}, Desc: sort.sortDesc()}).
		Offset(paginator.offset()).Limit(paginator.limit()).
		Find(&tasks)

	var total int64
	r.db.
		Model(&Task{}).
		Where(
			r.db.Where("to_tsvector('simple', description) @@ plainto_tsquery('simple', ?)", search).Or("?=''", search),
		).
		Where(map[string]interface{}(*filters)).
		Count(&total)

	pagination := buildPagination(paginator.Page, paginator.Limit, int(total))

	return tasks, pagination, nil
}
