package evaluation

// ScoreStatus represents the pass/warn/fail status for a category score.
type ScoreStatus string

const (
	ScoreStatusPass          ScoreStatus = "pass"              // Score >= 7.0
	ScoreStatusWarn          ScoreStatus = "warn"              // Score >= 5.0 && < 7.0
	ScoreStatusFail          ScoreStatus = "fail"              // Score < 5.0
	CategoryPending          ScoreStatus = "pending"           // Not yet evaluated
	CategoryNeedsImprovement ScoreStatus = "needs_improvement" // Requires attention
)

// Icon returns the emoji icon for the score status.
func (s ScoreStatus) Icon() string {
	switch s {
	case ScoreStatusPass:
		return "🟢"
	case ScoreStatusWarn:
		return "🟡"
	case ScoreStatusFail:
		return "🔴"
	default:
		return "⚪"
	}
}

// CategoryScore represents a score for a single evaluation category.
type CategoryScore struct {
	// Category is the name/ID of the category.
	Category string `json:"category"`

	// Weight is the category weight (0.0-1.0, should sum to 1.0).
	Weight float64 `json:"weight"`

	// Score is the category score (0.0-10.0).
	Score float64 `json:"score"`

	// MaxScore is the maximum possible score (default 10.0).
	MaxScore float64 `json:"max_score"`

	// Status is the derived status (pass/warn/fail).
	Status ScoreStatus `json:"status"`

	// Justification explains why this score was given.
	Justification string `json:"justification"`

	// Evidence provides specific supporting evidence.
	Evidence string `json:"evidence,omitempty"`

	// Findings are issues found in this category.
	Findings []Finding `json:"findings,omitempty"`
}

// ComputeStatus calculates the status from the score.
func (c *CategoryScore) ComputeStatus() ScoreStatus {
	switch {
	case c.Score >= 7.0:
		c.Status = ScoreStatusPass
	case c.Score >= 5.0:
		c.Status = ScoreStatusWarn
	default:
		c.Status = ScoreStatusFail
	}
	return c.Status
}

// ComputeWeightedScore calculates the weighted contribution of this category.
func (c *CategoryScore) ComputeWeightedScore() float64 {
	return c.Score * c.Weight
}

// NewCategoryScore creates a category score with computed status.
func NewCategoryScore(category string, weight, score float64, justification string) CategoryScore {
	cs := CategoryScore{
		Category:      category,
		Weight:        weight,
		Score:         score,
		MaxScore:      10.0,
		Justification: justification,
	}
	cs.ComputeStatus()
	return cs
}
