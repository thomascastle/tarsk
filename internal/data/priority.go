package data

type Priority string

const (
	PriorityNone   Priority = "none"
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityNone, PriorityLow, PriorityMedium, PriorityHigh:
		return true
	}

	return false
}
