// Package evaluation provides types for detailed evaluation reports
// with severity-based findings and recommendations. This is suited for
// LLM-as-Judge style reviews like PRD and ARB evaluations.
package evaluation

// Severity represents the severity level of a finding.
// Based on InfoSec severity classifications.
type Severity string

const (
	SeverityCritical Severity = "critical" // Blocks approval, must fix
	SeverityHigh     Severity = "high"     // Blocks approval, must fix
	SeverityMedium   Severity = "medium"   // Should fix before approval
	SeverityLow      Severity = "low"      // Nice to fix
	SeverityInfo     Severity = "info"     // Informational only
)

// Icon returns the emoji icon for the severity.
func (s Severity) Icon() string {
	switch s {
	case SeverityCritical:
		return "🔴"
	case SeverityHigh:
		return "🔴"
	case SeverityMedium:
		return "🟡"
	case SeverityLow:
		return "🟢"
	case SeverityInfo:
		return "ℹ️"
	default:
		return "⚪"
	}
}

// IsBlocking returns true if this severity blocks approval.
func (s Severity) IsBlocking() bool {
	return s == SeverityCritical || s == SeverityHigh
}

// Weight returns a numeric weight for sorting (higher = more severe).
func (s Severity) Weight() int {
	switch s {
	case SeverityCritical:
		return 5
	case SeverityHigh:
		return 4
	case SeverityMedium:
		return 3
	case SeverityLow:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 0
	}
}

// AllSeverities returns all severity levels in order of severity.
func AllSeverities() []Severity {
	return []Severity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		SeverityInfo,
	}
}
