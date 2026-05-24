package evaluation

// ScoreValue represents a categorical score value.
type ScoreValue string

const (
	ScorePass    ScoreValue = "pass"
	ScorePartial ScoreValue = "partial"
	ScoreFail    ScoreValue = "fail"
)

// IsPassing returns true if this score is considered passing.
func (s ScoreValue) IsPassing() bool {
	return s == ScorePass
}

// IsPartial returns true if this score is partial.
func (s ScoreValue) IsPartial() bool {
	return s == ScorePartial
}

// IsFailing returns true if this score is failing.
func (s ScoreValue) IsFailing() bool {
	return s == ScoreFail
}

// Icon returns the emoji icon for the score.
func (s ScoreValue) Icon() string {
	switch s {
	case ScorePass:
		return "🟢"
	case ScorePartial:
		return "🟡"
	case ScoreFail:
		return "🔴"
	default:
		return "⚪"
	}
}

// CategoryResult is the evaluation result for a single category.
type CategoryResult struct {
	// Category is the category ID.
	Category string `json:"category"`

	// Score is the assigned score (pass, partial, fail).
	// This is the authoritative score for decision-making.
	Score ScoreValue `json:"score"`

	// NumericScore is an optional numeric score (e.g., 1-5 Likert).
	// Used for human comparison, inter-rater reliability, and calibration.
	// The categorical Score takes precedence for pass/fail decisions.
	NumericScore *float64 `json:"numericScore,omitempty"`

	// Reasoning explains the score (chain-of-thought).
	Reasoning string `json:"reasoning"`

	// Evidence are specific quotes or observations.
	Evidence []string `json:"evidence,omitempty"`

	// Findings are issues discovered in this category.
	Findings []Finding `json:"findings,omitempty"`

	// ChecklistResults tracks checklist items (for checklist scales).
	ChecklistResults *ChecklistResults `json:"checklistResults,omitempty"`
}

// ChecklistResults tracks which items were found for checklist scales.
type ChecklistResults struct {
	// RequiredPresent are required items that were found.
	RequiredPresent []string `json:"requiredPresent,omitempty"`

	// RequiredMissing are required items that were not found.
	RequiredMissing []string `json:"requiredMissing,omitempty"`

	// OptionalPresent are optional items that were found.
	OptionalPresent []string `json:"optionalPresent,omitempty"`

	// OptionalMissing are optional items that were not found.
	OptionalMissing []string `json:"optionalMissing,omitempty"`
}

// NewCategoryResult creates a category result with the given score.
func NewCategoryResult(category string, score ScoreValue, reasoning string) *CategoryResult {
	return &CategoryResult{
		Category:  category,
		Score:     score,
		Reasoning: reasoning,
		Evidence:  []string{},
		Findings:  []Finding{},
	}
}

// NewCategoryResultWithNumeric creates a category result with both categorical and numeric scores.
// The numeric score is used for human comparison; categorical score is authoritative for decisions.
func NewCategoryResultWithNumeric(category string, score ScoreValue, numericScore float64, reasoning string) *CategoryResult {
	return &CategoryResult{
		Category:     category,
		Score:        score,
		NumericScore: &numericScore,
		Reasoning:    reasoning,
		Evidence:     []string{},
		Findings:     []Finding{},
	}
}

// NewCategoryResultFromLikert creates a category result from a Likert score.
// The categorical score is derived from the numeric score using the config thresholds.
func NewCategoryResultFromLikert(category string, likertScore int, config *LikertConfig, reasoning string) *CategoryResult {
	categoricalScore := LikertToCategorical(likertScore, config)
	numericScore := float64(likertScore)
	return &CategoryResult{
		Category:     category,
		Score:        categoricalScore,
		NumericScore: &numericScore,
		Reasoning:    reasoning,
		Evidence:     []string{},
		Findings:     []Finding{},
	}
}

// SetNumericScore sets the numeric score.
func (cr *CategoryResult) SetNumericScore(score float64) *CategoryResult {
	cr.NumericScore = &score
	return cr
}

// HasNumericScore returns true if a numeric score is set.
func (cr *CategoryResult) HasNumericScore() bool {
	return cr.NumericScore != nil
}

// GetNumericScore returns the numeric score, or 0 if not set.
func (cr *CategoryResult) GetNumericScore() float64 {
	if cr.NumericScore == nil {
		return 0
	}
	return *cr.NumericScore
}

// AddEvidence adds evidence to the result.
func (cr *CategoryResult) AddEvidence(evidence ...string) *CategoryResult {
	cr.Evidence = append(cr.Evidence, evidence...)
	return cr
}

// AddFinding adds a finding to the result.
func (cr *CategoryResult) AddFinding(f Finding) *CategoryResult {
	cr.Findings = append(cr.Findings, f)
	return cr
}

// SetChecklistResults sets the checklist results.
func (cr *CategoryResult) SetChecklistResults(results *ChecklistResults) *CategoryResult {
	cr.ChecklistResults = results
	return cr
}

// IsPassing returns true if this category passed.
func (cr *CategoryResult) IsPassing() bool {
	return cr.Score.IsPassing()
}

// CountCategoryResults counts results by score value.
type CategoryResultCounts struct {
	Pass    int `json:"pass"`
	Partial int `json:"partial"`
	Fail    int `json:"fail"`
	Total   int `json:"total"`
}

// CountResults counts category results by score.
func CountResults(results []CategoryResult) CategoryResultCounts {
	counts := CategoryResultCounts{}
	for _, r := range results {
		counts.Total++
		switch r.Score {
		case ScorePass:
			counts.Pass++
		case ScorePartial:
			counts.Partial++
		case ScoreFail:
			counts.Fail++
		}
	}
	return counts
}

// AllPassing returns true if all results are passing.
func (c CategoryResultCounts) AllPassing() bool {
	return c.Fail == 0 && c.Partial == 0
}

// AllRequiredPassing checks if all required categories passed.
func AllRequiredPassing(results []CategoryResult, rubric *RubricSet) bool {
	for _, result := range results {
		cat := rubric.GetCategory(result.Category)
		if cat != nil && cat.Required && !result.IsPassing() {
			return false
		}
	}
	return true
}
