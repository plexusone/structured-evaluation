// Package summary provides types for summary-style evaluation reports
// with GO/WARN/NO-GO status per task. This is suited for deterministic
// checks like tests, linting, and vulnerability scans.
package summary

// Status represents the pass/fail status for a task or team.
// Based on NASA Go/No-Go terminology.
type Status string

const (
	StatusGo   Status = "GO"
	StatusWarn Status = "WARN"
	StatusNoGo Status = "NO-GO"
	StatusSkip Status = "SKIP"
)

// Icon returns the emoji icon for the status.
func (s Status) Icon() string {
	switch s {
	case StatusGo:
		return "🟢"
	case StatusWarn:
		return "🟡"
	case StatusNoGo:
		return "🔴"
	case StatusSkip:
		return "⚪"
	default:
		return "⚪"
	}
}

// IsPassing returns true if the status is GO or WARN (not blocking).
func (s Status) IsPassing() bool {
	return s == StatusGo || s == StatusWarn || s == StatusSkip
}

// IsBlocking returns true if the status is NO-GO.
func (s Status) IsBlocking() bool {
	return s == StatusNoGo
}

// ComputeStatus determines the overall status from multiple statuses.
// Priority: NO-GO > WARN > GO > SKIP
func ComputeStatus(statuses []Status) Status {
	hasNoGo := false
	hasWarn := false
	hasGo := false

	for _, s := range statuses {
		switch s {
		case StatusNoGo:
			hasNoGo = true
		case StatusWarn:
			hasWarn = true
		case StatusGo:
			hasGo = true
		}
	}

	if hasNoGo {
		return StatusNoGo
	}
	if hasWarn {
		return StatusWarn
	}
	if hasGo {
		return StatusGo
	}
	return StatusSkip
}
