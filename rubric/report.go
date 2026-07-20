package rubric

import (
	"fmt"
	"time"
)

// Rubric is the detailed rubric-based evaluation report for LLM-as-Judge reviews.
type Rubric struct {
	// Schema is the JSON Schema URL.
	Schema string `json:"$schema,omitempty"`

	// SchemaVersion is the evaluation schema version (e.g., "v2").
	// Used for backwards compatibility.
	SchemaVersion string `json:"schemaVersion,omitempty"`

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

	// IntScore is the overall 1-5 integer score.
	// Preferred for LLM judges as they are unreliable at finer granularity.
	IntScore IntegerScore `json:"intScore,omitempty"`

	// Confidence is the overall confidence in the evaluation (0.0-1.0).
	// Low confidence evaluations may be routed to human review.
	Confidence float64 `json:"confidence,omitempty"`

	// Pass is an explicit pass/fail gate, orthogonal to score.
	// A spec can have a high score but still fail due to blocking issues.
	Pass bool `json:"pass"`

	// Blocking contains reason codes that caused failure.
	// Empty if Pass is true.
	Blocking []ReasonCode `json:"blocking,omitempty"`

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

	// Extensions contains domain-specific metadata.
	// Use this to store custom data without modifying the core schema.
	// Example: {"coverage": {...}, "metrics": {...}}
	Extensions map[string]any `json:"extensions,omitempty"`
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

// SchemaVersionV2 is the current schema version.
const SchemaVersionV2 = "v2"

// NewRubric creates a new rubric-based evaluation report.
func NewRubric(reviewType, document string) *Rubric {
	return &Rubric{
		SchemaVersion: SchemaVersionV2,
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

// SetIntScore sets the overall integer score.
func (r *Rubric) SetIntScore(score IntegerScore) *Rubric {
	r.IntScore = score
	return r
}

// SetConfidence sets the overall confidence value.
func (r *Rubric) SetConfidence(confidence float64) *Rubric {
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}
	r.Confidence = confidence
	return r
}

// SetPass sets the pass/fail status.
func (r *Rubric) SetPass(pass bool) *Rubric {
	r.Pass = pass
	return r
}

// AddBlocking adds a blocking reason code.
func (r *Rubric) AddBlocking(code ReasonCode) *Rubric {
	r.Blocking = append(r.Blocking, code)
	r.Pass = false
	return r
}

// SetBlocking sets the blocking reason codes.
func (r *Rubric) SetBlocking(codes []ReasonCode) *Rubric {
	r.Blocking = codes
	if len(codes) > 0 {
		r.Pass = false
	}
	return r
}

// HasLowConfidence returns true if confidence is below the threshold (default 0.7).
func (r *Rubric) HasLowConfidence(threshold ...float64) bool {
	t := 0.7
	if len(threshold) > 0 {
		t = threshold[0]
	}
	return r.Confidence > 0 && r.Confidence < t
}

// NeedsHumanReview returns true if this evaluation should be reviewed by a human.
func (r *Rubric) NeedsHumanReview(confidenceThreshold ...float64) bool {
	// Low confidence overall
	if r.HasLowConfidence(confidenceThreshold...) {
		return true
	}
	// Any category with low confidence
	for _, cat := range r.Categories {
		if cat.HasLowConfidence(confidenceThreshold...) {
			return true
		}
	}
	// Decision status is human_review
	if r.Decision.Status == DecisionHumanReview {
		return true
	}
	return false
}

// ComputeOverallIntScore calculates the overall integer score from category scores.
// Uses weighted average of category IntScores.
func (r *Rubric) ComputeOverallIntScore(rubricSet *RubricSet) IntegerScore {
	if len(r.Categories) == 0 {
		return ScoreAcceptable // Default to middle
	}

	var totalWeight float64
	var weightedSum float64

	for _, cat := range r.Categories {
		if cat.IntScore == 0 {
			continue // Skip categories without integer scores
		}
		weight := 1.0
		if rubricSet != nil {
			if catDef := rubricSet.GetCategory(cat.Category); catDef != nil && catDef.Weight > 0 {
				weight = catDef.Weight
			}
		}
		totalWeight += weight
		weightedSum += float64(cat.IntScore) * weight
	}

	if totalWeight == 0 {
		return ScoreAcceptable
	}

	avg := weightedSum / totalWeight
	return ParseIntegerScore(int(avg + 0.5)) // Round to nearest
}

// ComputeOverallConfidence calculates the overall confidence from category confidences.
// Uses minimum confidence across all categories (weakest link).
func (r *Rubric) ComputeOverallConfidence() float64 {
	if len(r.Categories) == 0 {
		return 0
	}

	minConfidence := 1.0
	hasConfidence := false

	for _, cat := range r.Categories {
		if cat.Confidence > 0 {
			hasConfidence = true
			if cat.Confidence < minConfidence {
				minConfidence = cat.Confidence
			}
		}
	}

	if !hasConfidence {
		return 0
	}
	return minConfidence
}

// CollectBlockingCodes gathers all blocking reason codes from findings.
func (r *Rubric) CollectBlockingCodes() []ReasonCode {
	return GetBlockingCodes(r.Findings)
}

// IsV2 returns true if this is a v2 schema evaluation.
func (r *Rubric) IsV2() bool {
	return r.SchemaVersion == SchemaVersionV2
}

// AddCategoryResult adds a category result.
func (r *Rubric) AddCategoryResult(cr CategoryResult) {
	if cr.Severity == "" {
		cr.Severity = WorstSeverity(cr.Findings)
	}
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
	// Compute category severities if not already set — a safety net for
	// categories appended directly to r.Categories, bypassing
	// AddCategoryResult/AddFinding.
	for i := range r.Categories {
		if r.Categories[i].Severity == "" {
			r.Categories[i].Severity = WorstSeverity(r.Categories[i].Findings)
		}
	}

	r.Decision = EvaluateResults(r.Categories, r.Findings, r.PassCriteria, rubricSet)
	r.OverallDecision = string(r.Decision.Status)

	// Compute v2 fields
	r.Pass = r.Decision.Passed
	r.Blocking = r.CollectBlockingCodes()

	// Compute overall integer score if not set
	if r.IntScore == 0 {
		r.IntScore = r.ComputeOverallIntScore(rubricSet)
	}

	// Compute overall confidence if not set
	if r.Confidence == 0 {
		r.Confidence = r.ComputeOverallConfidence()
	}

	// Check MinIntScore criteria (must be checked after IntScore is computed)
	if r.PassCriteria.MinIntScore > 0 && r.IntScore < r.PassCriteria.MinIntScore {
		r.Decision.Status = DecisionFail
		r.Decision.Passed = false
		r.Decision.Rationale = fmt.Sprintf("Score %d is below minimum required %d", r.IntScore, r.PassCriteria.MinIntScore)
		r.OverallDecision = string(r.Decision.Status)
		r.Pass = false
	}

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

// SetExtension sets a single extension value.
func (r *Rubric) SetExtension(key string, value any) {
	if r.Extensions == nil {
		r.Extensions = make(map[string]any)
	}
	r.Extensions[key] = value
}

// GetExtension retrieves an extension value by key.
// Returns nil if not found.
func (r *Rubric) GetExtension(key string) any {
	if r.Extensions == nil {
		return nil
	}
	return r.Extensions[key]
}

// HasExtension checks if an extension exists.
func (r *Rubric) HasExtension(key string) bool {
	if r.Extensions == nil {
		return false
	}
	_, ok := r.Extensions[key]
	return ok
}
