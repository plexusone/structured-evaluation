package evaluation

// Rubric defines the scoring criteria for an evaluation category.
// It provides explicit anchors for what each score level means,
// improving consistency and reproducibility of LLM-as-Judge evaluations.
type Rubric struct {
	// Category is the evaluation dimension this rubric applies to.
	Category string `json:"category"`

	// Description explains what this category measures.
	Description string `json:"description"`

	// Anchors define what each score level means.
	// Key is the score (e.g., "10", "7", "5", "3", "1") or range (e.g., "8-10").
	Anchors []ScoreAnchor `json:"anchors"`

	// Examples provide sample inputs/outputs for each score level.
	Examples []RubricExample `json:"examples,omitempty"`
}

// ScoreAnchor defines the criteria for a specific score or score range.
type ScoreAnchor struct {
	// Score is the numeric score this anchor represents.
	// Use -1 for range-based anchors where MinScore/MaxScore are set.
	Score float64 `json:"score,omitempty"`

	// MinScore is the minimum score for range-based anchors.
	MinScore float64 `json:"min_score,omitempty"`

	// MaxScore is the maximum score for range-based anchors.
	MaxScore float64 `json:"max_score,omitempty"`

	// Label is a short name for this level (e.g., "Excellent", "Good", "Poor").
	Label string `json:"label"`

	// Description explains what qualifies for this score.
	Description string `json:"description"`

	// Criteria are specific requirements for this score level.
	Criteria []string `json:"criteria,omitempty"`
}

// RubricExample provides a concrete example for a score level.
type RubricExample struct {
	// Score is the score this example demonstrates.
	Score float64 `json:"score"`

	// Input is the example input/prompt.
	Input string `json:"input,omitempty"`

	// Output is the example output being scored.
	Output string `json:"output"`

	// Explanation describes why this output receives this score.
	Explanation string `json:"explanation"`
}

// RubricSet is a collection of rubrics for a complete evaluation.
type RubricSet struct {
	// ID is the unique identifier for this rubric set.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Version tracks rubric iterations.
	Version string `json:"version"`

	// Description explains what this rubric set evaluates.
	Description string `json:"description,omitempty"`

	// Rubrics are the category-specific rubrics.
	Rubrics []Rubric `json:"rubrics"`
}

// NewRubric creates a new rubric for a category.
func NewRubric(category, description string) *Rubric {
	return &Rubric{
		Category:    category,
		Description: description,
		Anchors:     []ScoreAnchor{},
		Examples:    []RubricExample{},
	}
}

// AddAnchor adds a score anchor to the rubric.
func (r *Rubric) AddAnchor(score float64, label, description string, criteria ...string) *Rubric {
	r.Anchors = append(r.Anchors, ScoreAnchor{
		Score:       score,
		Label:       label,
		Description: description,
		Criteria:    criteria,
	})
	return r
}

// AddRangeAnchor adds a range-based score anchor.
func (r *Rubric) AddRangeAnchor(minScore, maxScore float64, label, description string, criteria ...string) *Rubric {
	r.Anchors = append(r.Anchors, ScoreAnchor{
		Score:       -1,
		MinScore:    minScore,
		MaxScore:    maxScore,
		Label:       label,
		Description: description,
		Criteria:    criteria,
	})
	return r
}

// AddExample adds an example to the rubric.
func (r *Rubric) AddExample(score float64, output, explanation string) *Rubric {
	r.Examples = append(r.Examples, RubricExample{
		Score:       score,
		Output:      output,
		Explanation: explanation,
	})
	return r
}

// GetAnchorForScore returns the anchor that applies to the given score.
func (r *Rubric) GetAnchorForScore(score float64) *ScoreAnchor {
	for i := range r.Anchors {
		anchor := &r.Anchors[i]
		if anchor.Score >= 0 && anchor.Score == score {
			return anchor
		}
		if anchor.Score < 0 && score >= anchor.MinScore && score <= anchor.MaxScore {
			return anchor
		}
	}
	return nil
}

// DefaultPRDRubricSet returns a standard rubric set for PRD evaluation.
func DefaultPRDRubricSet() *RubricSet {
	return &RubricSet{
		ID:          "prd-evaluation-v1",
		Name:        "PRD Evaluation Rubric",
		Version:     "1.0",
		Description: "Standard rubric for evaluating Product Requirements Documents",
		Rubrics: []Rubric{
			{
				Category:    "problem_definition",
				Description: "Clarity and completeness of problem statement",
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Problem is crystal clear with strong evidence of customer pain"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Problem is well-defined with supporting data"},
					{MinScore: 5, MaxScore: 6.9, Label: "Adequate", Description: "Problem stated but lacks depth or evidence"},
					{MinScore: 3, MaxScore: 4.9, Label: "Weak", Description: "Problem is vague or poorly articulated"},
					{MinScore: 0, MaxScore: 2.9, Label: "Missing", Description: "No clear problem statement"},
				},
			},
			{
				Category:    "solution_clarity",
				Description: "How well the proposed solution is articulated",
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "Solution is detailed, actionable, and addresses all aspects of the problem"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Solution is clear with most implementation details"},
					{MinScore: 5, MaxScore: 6.9, Label: "Adequate", Description: "Solution described but missing key details"},
					{MinScore: 3, MaxScore: 4.9, Label: "Weak", Description: "Solution is unclear or incomplete"},
					{MinScore: 0, MaxScore: 2.9, Label: "Missing", Description: "No clear solution proposed"},
				},
			},
			{
				Category:    "success_metrics",
				Description: "Quality and measurability of success criteria",
				Anchors: []ScoreAnchor{
					{MinScore: 9, MaxScore: 10, Label: "Excellent", Description: "SMART metrics with baselines and targets"},
					{MinScore: 7, MaxScore: 8.9, Label: "Good", Description: "Clear metrics that are measurable"},
					{MinScore: 5, MaxScore: 6.9, Label: "Adequate", Description: "Metrics defined but not fully measurable"},
					{MinScore: 3, MaxScore: 4.9, Label: "Weak", Description: "Vague or unmeasurable metrics"},
					{MinScore: 0, MaxScore: 2.9, Label: "Missing", Description: "No success metrics defined"},
				},
			},
		},
	}
}
