package evaluation

import (
	"fmt"
	"time"
)

// EvaluationReport is the detailed evaluation report for LLM-as-Judge reviews.
type EvaluationReport struct {
	// Schema is the JSON Schema URL.
	Schema string `json:"$schema,omitempty"`

	// Metadata contains report identification and audit info.
	Metadata ReportMetadata `json:"metadata"`

	// ReviewType identifies the type of review (prd, arb, security, etc.).
	ReviewType string `json:"review_type"`

	// Judge contains metadata about the LLM judge (v0.2.0).
	Judge *JudgeMetadata `json:"judge,omitempty"`

	// RubricID references the rubric used for scoring (v0.2.0).
	RubricID string `json:"rubric_id,omitempty"`

	// Reference contains gold/expected data for comparison (v0.2.0).
	Reference *ReferenceData `json:"reference,omitempty"`

	// Categories contains scores for each evaluation dimension.
	Categories []CategoryScore `json:"categories"`

	// Findings are all issues discovered during evaluation.
	Findings []Finding `json:"findings"`

	// WeightedScore is the overall weighted score.
	WeightedScore float64 `json:"weighted_score"`

	// PassCriteria defines the requirements for approval.
	PassCriteria PassCriteria `json:"pass_criteria"`

	// Decision is the evaluation outcome.
	Decision Decision `json:"decision"`

	// NextSteps provides actionable guidance.
	NextSteps NextSteps `json:"next_steps"`

	// Summary is the overall assessment.
	Summary string `json:"summary"`
}

// ReportMetadata contains report identification.
type ReportMetadata struct {
	// Document is the filename or path being evaluated.
	Document string `json:"document"`

	// DocumentID is the document identifier (e.g., PRD ID).
	DocumentID string `json:"document_id,omitempty"`

	// DocumentTitle is the document title.
	DocumentTitle string `json:"document_title,omitempty"`

	// DocumentVersion is the document version.
	DocumentVersion string `json:"document_version,omitempty"`

	// GeneratedAt is when the report was created.
	GeneratedAt time.Time `json:"generated_at"`

	// GeneratedBy identifies what created this report.
	GeneratedBy string `json:"generated_by,omitempty"`

	// ReviewerID identifies the reviewer (agent or human).
	ReviewerID string `json:"reviewer_id,omitempty"`
}

// NextSteps provides actionable workflow guidance.
type NextSteps struct {
	// RerunCommand is the command to re-run evaluation.
	RerunCommand string `json:"rerun_command"`

	// Immediate are blocking actions that must be completed.
	Immediate []ActionItem `json:"immediate,omitempty"`

	// Recommended are suggested improvements.
	Recommended []ActionItem `json:"recommended,omitempty"`
}

// ActionItem is a specific action to take.
type ActionItem struct {
	// Action describes what needs to be done.
	Action string `json:"action"`

	// Category is the related evaluation category.
	Category string `json:"category,omitempty"`

	// Severity is the related finding severity.
	Severity Severity `json:"severity,omitempty"`

	// Owner suggests who should do this.
	Owner string `json:"owner,omitempty"`

	// Effort estimates work required.
	Effort string `json:"effort,omitempty"`
}

// NewEvaluationReport creates a new evaluation report.
func NewEvaluationReport(reviewType, document string) *EvaluationReport {
	return &EvaluationReport{
		Metadata: ReportMetadata{
			Document:    document,
			GeneratedAt: time.Now().UTC(),
			GeneratedBy: "structured-evaluation",
		},
		ReviewType:   reviewType,
		Categories:   []CategoryScore{},
		Findings:     []Finding{},
		PassCriteria: DefaultPassCriteria(),
	}
}

// AddCategory adds a category score.
func (r *EvaluationReport) AddCategory(cs CategoryScore) {
	cs.ComputeStatus()
	r.Categories = append(r.Categories, cs)
}

// AddFinding adds a finding.
func (r *EvaluationReport) AddFinding(f Finding) {
	r.Findings = append(r.Findings, f)
}

// ComputeWeightedScore calculates the overall weighted score.
func (r *EvaluationReport) ComputeWeightedScore() float64 {
	var totalWeight float64
	var totalScore float64

	for _, c := range r.Categories {
		totalWeight += c.Weight
		totalScore += c.ComputeWeightedScore()
	}

	if totalWeight > 0 {
		r.WeightedScore = totalScore / totalWeight
	}
	return r.WeightedScore
}

// Evaluate computes the decision based on findings and score.
func (r *EvaluationReport) Evaluate() Decision {
	r.ComputeWeightedScore()
	r.Decision = Evaluate(r.Findings, r.WeightedScore, r.PassCriteria)
	return r.Decision
}

// GenerateNextSteps creates actionable next steps.
func (r *EvaluationReport) GenerateNextSteps(rerunCommand string) {
	r.NextSteps = NextSteps{
		RerunCommand: rerunCommand,
		Immediate:    []ActionItem{},
		Recommended:  []ActionItem{},
	}

	// Add immediate actions for blocking findings
	for _, f := range r.Findings {
		if f.IsBlocking() {
			r.NextSteps.Immediate = append(r.NextSteps.Immediate, ActionItem{
				Action:   f.Recommendation,
				Category: f.Category,
				Severity: f.Severity,
				Owner:    f.Owner,
				Effort:   f.Effort,
			})
		} else if f.Severity == SeverityMedium {
			r.NextSteps.Recommended = append(r.NextSteps.Recommended, ActionItem{
				Action:   f.Recommendation,
				Category: f.Category,
				Severity: f.Severity,
				Owner:    f.Owner,
				Effort:   f.Effort,
			})
		}
	}
}

// GenerateSummary creates the summary text.
func (r *EvaluationReport) GenerateSummary() string {
	counts := r.Decision.FindingCounts

	summary := fmt.Sprintf("Score: %.1f/10. ", r.WeightedScore)

	if counts.Total == 0 {
		summary += "No findings."
	} else {
		if counts.Critical > 0 {
			summary += fmt.Sprintf("%d critical, ", counts.Critical)
		}
		if counts.High > 0 {
			summary += fmt.Sprintf("%d high, ", counts.High)
		}
		if counts.Medium > 0 {
			summary += fmt.Sprintf("%d medium, ", counts.Medium)
		}
		if counts.Low > 0 {
			summary += fmt.Sprintf("%d low", counts.Low)
		}
		summary += " findings."
	}

	summary += " Decision: " + string(r.Decision.Status) + "."

	r.Summary = summary
	return summary
}

// Finalize computes all derived fields.
func (r *EvaluationReport) Finalize(rerunCommand string) {
	r.Evaluate()
	r.GenerateNextSteps(rerunCommand)
	r.GenerateSummary()
}

// SetJudge sets the judge metadata.
func (r *EvaluationReport) SetJudge(judge *JudgeMetadata) {
	r.Judge = judge
	if judge != nil && judge.RubricID != "" {
		r.RubricID = judge.RubricID
	}
}

// SetReference sets the reference data for comparison.
func (r *EvaluationReport) SetReference(ref *ReferenceData) {
	r.Reference = ref
}

// SetRubric sets the rubric ID.
func (r *EvaluationReport) SetRubric(rubricID string) {
	r.RubricID = rubricID
}
