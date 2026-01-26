package evaluation

// Finding represents an issue discovered during evaluation.
type Finding struct {
	// ID is the unique identifier for this finding.
	ID string `json:"id"`

	// Category is the evaluation category this relates to.
	Category string `json:"category"`

	// Severity indicates the impact level.
	Severity Severity `json:"severity"`

	// Title is a brief summary of the finding.
	Title string `json:"title"`

	// Description provides detailed explanation.
	Description string `json:"description"`

	// Recommendation explains how to fix the issue.
	Recommendation string `json:"recommendation"`

	// Evidence provides specific examples or references.
	Evidence string `json:"evidence,omitempty"`

	// Owner suggests who should address this finding.
	Owner string `json:"owner,omitempty"`

	// Effort estimates the work required (low, medium, high).
	Effort string `json:"effort,omitempty"`
}

// IsBlocking returns true if this finding blocks approval.
func (f *Finding) IsBlocking() bool {
	return f.Severity.IsBlocking()
}

// FindingCounts tracks the number of findings by severity.
type FindingCounts struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Info     int `json:"info"`
	Total    int `json:"total"`
}

// CountFindings counts findings by severity.
func CountFindings(findings []Finding) FindingCounts {
	counts := FindingCounts{}
	for _, f := range findings {
		counts.Total++
		switch f.Severity {
		case SeverityCritical:
			counts.Critical++
		case SeverityHigh:
			counts.High++
		case SeverityMedium:
			counts.Medium++
		case SeverityLow:
			counts.Low++
		case SeverityInfo:
			counts.Info++
		}
	}
	return counts
}

// BlockingCount returns the number of blocking findings.
func (c FindingCounts) BlockingCount() int {
	return c.Critical + c.High
}

// HasBlocking returns true if there are any blocking findings.
func (c FindingCounts) HasBlocking() bool {
	return c.BlockingCount() > 0
}
