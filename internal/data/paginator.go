package data

import (
	"math"

	"github.com/thomascastle/tarsk/internal/validator"
)

type Paginator struct {
	Page  int
	Limit int
}

func (p Paginator) limit() int {
	return p.Limit
}

func (p Paginator) offset() int {
	return (p.Page - 1) * p.Limit
}

func (p Paginator) Validate(v *validator.Validator) {
	if p.Limit < 1 {
		v.AddError("limit", "must be greater than zero")
	}
	if p.Limit > 100 {
		v.AddError("limit", "must be a maximum of 100")
	}
	if p.Page < 1 {
		v.AddError("page", "must be greater than zero")
	}
	if p.Page > 10_000_000 {
		v.AddError("page", "must be a maximum of 10 million")
	}
}

type Pagination struct {
	CurrentPage int `json:"current_page,omitempty"`
	Limit       int `json:"limit,omitempty"`
	FirstPage   int `json:"first_page,omitempty"`
	LastPage    int `json:"last_page,omitempty"`
	Total       int `json:"total,omitempty"`
}

func buildPagination(page, limit, total int) Pagination {
	if total == 0 {
		return Pagination{}
	}

	return Pagination{
		CurrentPage: page,
		Limit:       limit,
		FirstPage:   1,
		LastPage:    int(math.Ceil(float64(total) / float64(limit))),
		Total:       total,
	}
}
