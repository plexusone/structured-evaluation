package summary

// TaskResult represents the outcome of a single check or task.
type TaskResult struct {
	// ID is the unique identifier for this task.
	ID string `json:"id"`

	// Status is the task outcome (GO, WARN, NO-GO, SKIP).
	Status Status `json:"status"`

	// Detail provides a brief description of the result.
	Detail string `json:"detail,omitempty"`

	// DurationMs is the execution time in milliseconds.
	DurationMs int64 `json:"duration_ms,omitempty"`

	// Metadata contains additional task-specific data.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ComputeStatusFromTasks determines the overall status from task results.
func ComputeStatusFromTasks(tasks []TaskResult) Status {
	statuses := make([]Status, len(tasks))
	for i, t := range tasks {
		statuses[i] = t.Status
	}
	return ComputeStatus(statuses)
}
