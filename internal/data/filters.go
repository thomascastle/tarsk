package data

import (
	"net/url"
	"strconv"

	"github.com/thomascastle/tarsk/internal/validator"
)

type Filters map[string]interface{}

func (f Filters) Validate(v *validator.Validator) {
	if value, present := f["done"]; present {
		_, ok := value.(bool)
		if !ok {
			v.AddError("done", "invalid value")
		}
	}

	if value, present := f["priority"]; present {
		priority, ok := value.(Priority)
		if !ok {
			v.AddError("priority", "invalid value")
		}
		if !priority.Valid() {
			v.AddError("priority", "invalid value")
		}
	}
}

func ParseFilters(values url.Values) Filters {
	filters := make(map[string]interface{})
	done_string := values.Get("done")
	if done_string != "" {
		if done, e := strconv.ParseBool(done_string); e == nil {
			filters["done"] = done
		}
	}

	priority := values.Get("priority")
	if priority != "" {
		filters["priority"] = Priority(priority)
	}

	return filters
}
