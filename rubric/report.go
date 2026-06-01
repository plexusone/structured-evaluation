package rubric

import (
	"fmt"
	"time"
)

// Rubric is the detailed rubric-based evaluation report for LLM-as-Judge reviews.
type Rubric struct {
	// Schema is the JSON Schema URL.
	Schema string `json:"$schema,omitempty"`

	// Metadata contains report identification and audit info.
	Metadata ReportMetadata `json:"metadata"`

	// ReviewType identifies the type of review (prd, arb, security, article, etc.).
	ReviewType string `json:"reviewType"`

	// Judge contains metadata about the LLM judge.
	Judge *JudgeMetadata `json:"judge,omitempty"`

	// RubricID references the rubric used for scoring.
	RubricID string `json:"rubricId,omitempty"`

	// RubricVersion is the version of the rubric used.
	RubricVersion string `json:"rubricVersion,omitempty"`

	// Reference contains gold/expected data for comparison.
	Reference *ReferenceData `json:"reference,omitempty"`

	// Categories contains results for each evaluation dimension.
	Categories []CategoryResult `json:"categories"`

	// Findings are all issues discovered during evaluation.
	Findings []Finding `json:"findings"`

	// PassCriteria defines the requirements for approval.
	PassCriteria PassCriteria `json:"passCriteria"`

	// Decision is the evaluation outcome.
	Decision Decision `json:"decision"`

	// OverallDecision is a simplified pass/conditional/fail status.
	OverallDecision string `json:"overallDecision"`

	// NextSteps provides actionable guidance.
	NextSteps NextSteps `json:"nextSteps"`

	// Summary is the overall assessment.
	Summary string `json:"summary"`
}

// ReportMetadata contains report identification.
type ReportMetadata struct {
	// Document is the filename or path being evaluated.
	Document string `json:"document"`

	// DocumentID is the document identifier (e.g., PRD ID).
	DocumentID string `json:"documentId,omitempty"`

	// DocumentTitle is the document title.
	DocumentTitle string `json:"documentTitle,omitempty"`

	// DocumentVersion is the document version.
	DocumentVersion string `json:"documentVersion,omitempty"`

	// GeneratedAt is when the report was created.
	GeneratedAt time.Time `json:"generatedAt"`

	// GeneratedBy identifies what created this report.
	GeneratedBy string `json:"generatedBy,omitempty"`

	// ReviewerID identifies the reviewer (agent or human).
	ReviewerID string `json:"reviewerId,omitempty"`
}

// NextSteps provides actionable workflow guidance.
type NextSteps struct {
	// RerunCommand is the command to re-run evaluation.
	RerunCommand string `json:"rerunCommand,omitempty"`

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

// NewRubric creates a new rubric-based evaluation report.
func NewRubric(reviewType, document string) *Rubric {
	return &Rubric{
		Metadata: ReportMetadata{
			Document:    document,
			GeneratedAt: time.Now().UTC(),
			GeneratedBy: "structured-evaluation",
		},
		ReviewType:   reviewType,
		Categories:   []CategoryResult{},
		Findings:     []Finding{},
		PassCriteria: DefaultPassCriteria(),
	}
}

// AddCategoryResult adds a category result.
func (r *Rubric) AddCategoryResult(cr CategoryResult) {
	r.Categories = append(r.Categories, cr)
	// Also collect findings from the category
	r.Findings = append(r.Findings, cr.Findings...)
}

// AddFinding adds a finding.
func (r *Rubric) AddFinding(f Finding) {
	r.Findings = append(r.Findings, f)
}

// Evaluate computes the decision based on findings and category results.
func (r *Rubric) Evaluate(rubricSet *RubricSet) Decision {
	r.Decision = EvaluateResults(r.Categories, r.Findings, r.PassCriteria, rubricSet)
	r.OverallDecision = string(r.Decision.Status)
	return r.Decision
}

// GenerateNextSteps creates actionable next steps.
func (r *Rubric) GenerateNextSteps(rerunCommand string) {
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

	// Add actions for failed/partial categories
	for _, cat := range r.Categories {
		if cat.Score == ScoreFail {
			r.NextSteps.Immediate = append(r.NextSteps.Immediate, ActionItem{
				Action:   "Address failing category: " + cat.Category,
				Category: cat.Category,
			})
		} else if cat.Score == ScorePartial {
			r.NextSteps.Recommended = append(r.NextSteps.Recommended, ActionItem{
				Action:   "Improve partial category: " + cat.Category,
				Category: cat.Category,
			})
		}
	}
}

// GenerateSummary creates the summary text.
func (r *Rubric) GenerateSummary() string {
	counts := r.Decision.CategoryCounts
	findings := r.Decision.FindingCounts

	summary := fmt.Sprintf("Categories: %d pass, %d partial, %d fail. ",
		counts.Pass, counts.Partial, counts.Fail)

	if findings.Total == 0 {
		summary += "No findings."
	} else {
		findingParts := []string{}
		if findings.Critical > 0 {
			findingParts = append(findingParts, fmt.Sprintf("%d critical", findings.Critical))
		}
		if findings.High > 0 {
			findingParts = append(findingParts, fmt.Sprintf("%d high", findings.High))
		}
		if findings.Medium > 0 {
			findingParts = append(findingParts, fmt.Sprintf("%d medium", findings.Medium))
		}
		if findings.Low > 0 {
			findingParts = append(findingParts, fmt.Sprintf("%d low", findings.Low))
		}
		summary += fmt.Sprintf("%d findings", findings.Total)
		if len(findingParts) > 0 {
			summary += " ("
			for i, part := range findingParts {
				if i > 0 {
					summary += ", "
				}
				summary += part
			}
			summary += ")"
		}
		summary += "."
	}

	summary += " Decision: " + string(r.Decision.Status) + "."

	r.Summary = summary
	return summary
}

// Finalize computes all derived fields.
func (r *Rubric) Finalize(rubricSet *RubricSet, rerunCommand string) {
	r.Evaluate(rubricSet)
	r.GenerateNextSteps(rerunCommand)
	r.GenerateSummary()
}

// SetJudge sets the judge metadata.
func (r *Rubric) SetJudge(judge *JudgeMetadata) {
	r.Judge = judge
	if judge != nil && judge.RubricID != "" {
		r.RubricID = judge.RubricID
	}
	if judge != nil && judge.RubricVersion != "" {
		r.RubricVersion = judge.RubricVersion
	}
}

// SetReference sets the reference data for comparison.
func (r *Rubric) SetReference(ref *ReferenceData) {
	r.Reference = ref
}

// SetRubricInfo sets the rubric ID and version.
func (r *Rubric) SetRubricInfo(rubricID, rubricVersion string) {
	r.RubricID = rubricID
	r.RubricVersion = rubricVersion
}

// SetPassCriteria sets the pass criteria.
func (r *Rubric) SetPassCriteria(criteria PassCriteria) {
	r.PassCriteria = criteria
}

// GetCategoryResult returns a category result by ID, or nil if not found.
func (r *Rubric) GetCategoryResult(categoryID string) *CategoryResult {
	for i := range r.Categories {
		if r.Categories[i].Category == categoryID {
			return &r.Categories[i]
		}
	}
	return nil
}
