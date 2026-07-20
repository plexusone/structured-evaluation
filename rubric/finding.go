package rubric

// Finding represents an issue discovered during evaluation.
type Finding struct {
	// ID is the unique identifier for this finding.
	ID string `json:"id"`

	// Category is the evaluation category this relates to.
	Category string `json:"category"`

	// Code is the standardized reason code for this finding.
	// Enables automated repair workflows.
	Code ReasonCode `json:"code,omitempty"`

	// Severity indicates the impact level.
	Severity Severity `json:"severity"`

	// Title is a brief summary of the finding.
	Title string `json:"title"`

	// Description provides detailed explanation.
	Description string `json:"description"`

	// Recommendation explains how to fix the issue.
	Recommendation string `json:"recommendation"`

	// Location is a reference to where the issue was found (e.g., "REQ-12", "Section 3.2").
	Location string `json:"location,omitempty"`

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

// NewFinding creates a new finding with the required fields.
func NewFinding(id, category string, severity Severity, title, description string) *Finding {
	return &Finding{
		ID:          id,
		Category:    category,
		Severity:    severity,
		Title:       title,
		Description: description,
	}
}

// NewFindingWithCode creates a new finding with a reason code.
func NewFindingWithCode(id, category string, code ReasonCode, title, description string) *Finding {
	info := GetReasonCodeInfo(code)
	severity := SeverityMedium
	if info != nil {
		severity = info.DefaultSeverity
	}
	return &Finding{
		ID:          id,
		Category:    category,
		Code:        code,
		Severity:    severity,
		Title:       title,
		Description: description,
	}
}

// SetCode sets the reason code on the finding.
func (f *Finding) SetCode(code ReasonCode) *Finding {
	f.Code = code
	return f
}

// SetLocation sets the location reference on the finding.
func (f *Finding) SetLocation(location string) *Finding {
	f.Location = location
	return f
}

// SetRecommendation sets the recommendation on the finding.
func (f *Finding) SetRecommendation(recommendation string) *Finding {
	f.Recommendation = recommendation
	return f
}

// SetEvidence sets the evidence on the finding.
func (f *Finding) SetEvidence(evidence string) *Finding {
	f.Evidence = evidence
	return f
}

// SetOwner sets the owner on the finding.
func (f *Finding) SetOwner(owner string) *Finding {
	f.Owner = owner
	return f
}

// SetEffort sets the effort estimate on the finding.
func (f *Finding) SetEffort(effort string) *Finding {
	f.Effort = effort
	return f
}

// GetCodeInfo returns the reason code info for this finding's code.
func (f *Finding) GetCodeInfo() *ReasonCodeInfo {
	return GetReasonCodeInfo(f.Code)
}

// GetRepairHint returns the repair hint from the reason code registry.
// Deprecated: Use GetRepairPrompt instead.
func (f *Finding) GetRepairHint() string {
	return f.GetRepairPrompt()
}

// GetRepairPrompt returns the AI repair prompt from the reason code registry.
func (f *Finding) GetRepairPrompt() string {
	if info := f.GetCodeInfo(); info != nil {
		return info.RepairPrompt
	}
	return f.Recommendation
}

// CountFindingsByCode counts findings by reason code.
func CountFindingsByCode(findings []Finding) map[ReasonCode]int {
	counts := make(map[ReasonCode]int)
	for _, f := range findings {
		if f.Code != "" {
			counts[f.Code]++
		}
	}
	return counts
}

// GetBlockingCodes returns the reason codes from blocking findings.
func GetBlockingCodes(findings []Finding) []ReasonCode {
	var codes []ReasonCode
	seen := make(map[ReasonCode]bool)
	for _, f := range findings {
		if f.IsBlocking() && f.Code != "" && !seen[f.Code] {
			codes = append(codes, f.Code)
			seen[f.Code] = true
		}
	}
	return codes
}

// WorstSeverity returns the highest-weight severity among findings, or the
// zero value if findings is empty. Used to roll a category's findings up
// into a single severity for prioritization (e.g. CategoryResult.Severity).
func WorstSeverity(findings []Finding) Severity {
	var worst Severity
	for _, f := range findings {
		if worst == "" || f.Severity.Weight() > worst.Weight() {
			worst = f.Severity
		}
	}
	return worst
}
