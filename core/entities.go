package core

import "fmt"

type JobState int

const (
	Discarded JobState = -101
	Cancelled JobState = -100
	Available JobState = 0
	Running   JobState = 1
	Completed JobState = 100
	Retryable JobState = 101
)

func (js JobState) String() string {
	switch js {
	case Discarded:
		return "Discarded"
	case Cancelled:
		return "Cancelled"
	case Available:
		return "Available"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Retryable:
		return "Retryable"
	default:
		return fmt.Sprintf("Unknown JobState (%d)", js)
	}
}
