package data

import (
	"strings"

	"github.com/thomascastle/tarsk/internal/validator"
)

type Sort struct {
	Sort         string
	SortSafelist []string
}

func (s Sort) sortColumn() string {
	for _, safeValue := range s.SortSafelist {
		if s.Sort == safeValue {
			return strings.TrimPrefix(s.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + s.Sort)
}

func (s Sort) sortDesc() bool {
	return strings.HasPrefix(s.Sort, "-")
}

func (s Sort) Validate(v *validator.Validator) {
	if !validator.In(s.Sort, s.SortSafelist...) {
		v.AddError("sort", "invalid sort value")
	}
}
