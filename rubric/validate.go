package rubric

import (
	"fmt"
	"strings"
)

// ValidationIssue represents a single validation problem.
type ValidationIssue struct {
	// Path is the JSON path to the problematic field (e.g., "categories[0].score").
	Path string `json:"path"`

	// Code is a machine-readable error code.
	Code string `json:"code"`

	// Message describes the issue.
	Message string `json:"message"`

	// Severity indicates how serious the issue is.
	Severity ValidationSeverity `json:"severity"`

	// ActualValue is the invalid value found.
	ActualValue string `json:"actualValue,omitempty"`

	// AllowedValues lists valid options (for enum violations).
	AllowedValues []string `json:"allowedValues,omitempty"`
}

// ValidationSeverity indicates the severity of a validation issue.
type ValidationSeverity string

const (
	// ValidationError is a fatal issue that must be fixed.
	ValidationError ValidationSeverity = "error"

	// ValidationWarning is a non-fatal issue that should be fixed.
	ValidationWarning ValidationSeverity = "warning"
)

// ValidationResult contains all validation issues found.
type ValidationResult struct {
	// Valid is true if no errors were found (warnings allowed).
	Valid bool `json:"valid"`

	// Issues contains all validation problems found.
	Issues []ValidationIssue `json:"issues"`

	// ErrorCount is the number of error-level issues.
	ErrorCount int `json:"errorCount"`

	// WarningCount is the number of warning-level issues.
	WarningCount int `json:"warningCount"`
}

// HasErrors returns true if there are any error-level issues.
func (r *ValidationResult) HasErrors() bool {
	return r.ErrorCount > 0
}

// HasWarnings returns true if there are any warning-level issues.
func (r *ValidationResult) HasWarnings() bool {
	return r.WarningCount > 0
}

// addError adds an error-level issue.
func (r *ValidationResult) addError(path, code, message string) {
	r.Issues = append(r.Issues, ValidationIssue{
		Path:     path,
		Code:     code,
		Message:  message,
		Severity: ValidationError,
	})
	r.ErrorCount++
	r.Valid = false
}

// addEnumError adds an error for an invalid enum value.
func (r *ValidationResult) addEnumError(path, code, message, actual string, allowed []string) {
	r.Issues = append(r.Issues, ValidationIssue{
		Path:          path,
		Code:          code,
		Message:       message,
		Severity:      ValidationError,
		ActualValue:   actual,
		AllowedValues: allowed,
	})
	r.ErrorCount++
	r.Valid = false
}

// addWarning adds a warning-level issue.
func (r *ValidationResult) addWarning(path, code, message string) {
	r.Issues = append(r.Issues, ValidationIssue{
		Path:     path,
		Code:     code,
		Message:  message,
		Severity: ValidationWarning,
	})
	r.WarningCount++
}

// String returns a human-readable summary.
func (r *ValidationResult) String() string {
	if r.Valid && r.WarningCount == 0 {
		return "Valid: no issues found"
	}

	var sb strings.Builder
	if r.Valid {
		sb.WriteString(fmt.Sprintf("Valid with %d warning(s)\n", r.WarningCount))
	} else {
		sb.WriteString(fmt.Sprintf("Invalid: %d error(s), %d warning(s)\n", r.ErrorCount, r.WarningCount))
	}

	for _, issue := range r.Issues {
		icon := "ERROR"
		if issue.Severity == ValidationWarning {
			icon = "WARN"
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n", icon, issue.Path, issue.Message))
		if issue.ActualValue != "" {
			sb.WriteString(fmt.Sprintf("         got: %q\n", issue.ActualValue))
		}
		if len(issue.AllowedValues) > 0 {
			sb.WriteString(fmt.Sprintf("         allowed: %v\n", issue.AllowedValues))
		}
	}

	return sb.String()
}

// ValidScoreValues returns all valid score values.
func ValidScoreValues() []string {
	return []string{string(ScorePass), string(ScorePartial), string(ScoreFail)}
}

// ValidSeverityValues returns all valid severity values.
func ValidSeverityValues() []string {
	return []string{
		string(SeverityCritical),
		string(SeverityHigh),
		string(SeverityMedium),
		string(SeverityLow),
		string(SeverityInfo),
	}
}

// ValidDecisionStatusValues returns all valid decision status values.
func ValidDecisionStatusValues() []string {
	return []string{
		string(DecisionPass),
		string(DecisionConditional),
		string(DecisionFail),
		string(DecisionHumanReview),
	}
}

// ValidEvaluationTypeValues returns all valid evaluation type values.
func ValidEvaluationTypeValues() []string {
	return []string{
		string(EvaluationTypeAnalytic),
		string(EvaluationTypeHolistic),
	}
}

// ValidScaleTypeValues returns all valid scale type values.
func ValidScaleTypeValues() []string {
	return []string{
		string(ScaleTypeCategorical),
		string(ScaleTypeChecklist),
		string(ScaleTypeBinary),
		string(ScaleTypeLikert),
	}
}

// isValidScore checks if a score value is valid.
func isValidScore(s ScoreValue) bool {
	switch s {
	case ScorePass, ScorePartial, ScoreFail:
		return true
	default:
		return false
	}
}

// isValidSeverity checks if a severity value is valid.
func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo:
		return true
	default:
		return false
	}
}

// isValidDecisionStatus checks if a decision status is valid.
func isValidDecisionStatus(s DecisionStatus) bool {
	switch s {
	case DecisionPass, DecisionConditional, DecisionFail, DecisionHumanReview:
		return true
	default:
		return false
	}
}

// ValidateReport validates a rubric evaluation report for correctness.
// It checks:
// - All enum values are valid (scores, severities, decision status)
// - Required fields are present
// - Counts are accurate
// - Decision is consistent with findings/categories
func ValidateReport(r *Rubric) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Validate required metadata
	if r.Metadata.Document == "" {
		result.addError("metadata.document", "REQUIRED_FIELD", "document is required")
	}
	if r.ReviewType == "" {
		result.addError("reviewType", "REQUIRED_FIELD", "reviewType is required")
	}

	// Validate categories
	for i, cat := range r.Categories {
		path := fmt.Sprintf("categories[%d]", i)

		// Validate score enum
		if !isValidScore(cat.Score) {
			result.addEnumError(
				path+".score",
				"INVALID_SCORE",
				fmt.Sprintf("invalid score value %q", cat.Score),
				string(cat.Score),
				ValidScoreValues(),
			)
		}

		// Validate category ID
		if cat.Category == "" {
			result.addError(path+".category", "REQUIRED_FIELD", "category ID is required")
		}

		// Validate findings within category
		for j, finding := range cat.Findings {
			findingPath := fmt.Sprintf("%s.findings[%d]", path, j)
			validateFinding(finding, findingPath, result)
		}
	}

	// Validate top-level findings
	for i, finding := range r.Findings {
		path := fmt.Sprintf("findings[%d]", i)
		validateFinding(finding, path, result)
	}

	// Validate decision
	if !isValidDecisionStatus(r.Decision.Status) {
		result.addEnumError(
			"decision.status",
			"INVALID_DECISION_STATUS",
			fmt.Sprintf("invalid decision status %q", r.Decision.Status),
			string(r.Decision.Status),
			ValidDecisionStatusValues(),
		)
	}

	// Validate overallDecision matches decision.status
	if r.OverallDecision != "" && r.OverallDecision != string(r.Decision.Status) {
		result.addWarning(
			"overallDecision",
			"INCONSISTENT_DECISION",
			fmt.Sprintf("overallDecision %q does not match decision.status %q", r.OverallDecision, r.Decision.Status),
		)
	}

	// Validate count accuracy
	validateCounts(r, result)

	// Validate semantic consistency
	validateSemanticConsistency(r, result)

	return result
}

// validateFinding validates a single finding.
func validateFinding(f Finding, path string, result *ValidationResult) {
	// Validate severity enum
	if !isValidSeverity(f.Severity) {
		result.addEnumError(
			path+".severity",
			"INVALID_SEVERITY",
			fmt.Sprintf("invalid severity value %q", f.Severity),
			string(f.Severity),
			ValidSeverityValues(),
		)
	}

	// Validate required fields
	if f.Title == "" {
		result.addError(path+".title", "REQUIRED_FIELD", "finding title is required")
	}
}

// validateCounts verifies that reported counts match actual data.
func validateCounts(r *Rubric, result *ValidationResult) {
	// Validate category counts
	actualCatCounts := CountResults(r.Categories)
	reportedCatCounts := r.Decision.CategoryCounts

	if reportedCatCounts.Total > 0 { // Only validate if counts are set
		if actualCatCounts.Pass != reportedCatCounts.Pass {
			result.addWarning(
				"decision.categoryCounts.pass",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d pass categories but found %d", reportedCatCounts.Pass, actualCatCounts.Pass),
			)
		}
		if actualCatCounts.Partial != reportedCatCounts.Partial {
			result.addWarning(
				"decision.categoryCounts.partial",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d partial categories but found %d", reportedCatCounts.Partial, actualCatCounts.Partial),
			)
		}
		if actualCatCounts.Fail != reportedCatCounts.Fail {
			result.addWarning(
				"decision.categoryCounts.fail",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d fail categories but found %d", reportedCatCounts.Fail, actualCatCounts.Fail),
			)
		}
		if actualCatCounts.Total != reportedCatCounts.Total {
			result.addWarning(
				"decision.categoryCounts.total",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d total categories but found %d", reportedCatCounts.Total, actualCatCounts.Total),
			)
		}
	}

	// Validate finding counts
	actualFindingCounts := CountFindings(r.Findings)
	reportedFindingCounts := r.Decision.FindingCounts

	if reportedFindingCounts.Total > 0 { // Only validate if counts are set
		if actualFindingCounts.Critical != reportedFindingCounts.Critical {
			result.addWarning(
				"decision.findingCounts.critical",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d critical findings but found %d", reportedFindingCounts.Critical, actualFindingCounts.Critical),
			)
		}
		if actualFindingCounts.High != reportedFindingCounts.High {
			result.addWarning(
				"decision.findingCounts.high",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d high findings but found %d", reportedFindingCounts.High, actualFindingCounts.High),
			)
		}
		if actualFindingCounts.Medium != reportedFindingCounts.Medium {
			result.addWarning(
				"decision.findingCounts.medium",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d medium findings but found %d", reportedFindingCounts.Medium, actualFindingCounts.Medium),
			)
		}
		if actualFindingCounts.Low != reportedFindingCounts.Low {
			result.addWarning(
				"decision.findingCounts.low",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d low findings but found %d", reportedFindingCounts.Low, actualFindingCounts.Low),
			)
		}
		if actualFindingCounts.Info != reportedFindingCounts.Info {
			result.addWarning(
				"decision.findingCounts.info",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d info findings but found %d", reportedFindingCounts.Info, actualFindingCounts.Info),
			)
		}
		if actualFindingCounts.Total != reportedFindingCounts.Total {
			result.addWarning(
				"decision.findingCounts.total",
				"INCORRECT_COUNT",
				fmt.Sprintf("reported %d total findings but found %d", reportedFindingCounts.Total, actualFindingCounts.Total),
			)
		}
	}
}

// validateSemanticConsistency checks that the decision makes sense given the data.
func validateSemanticConsistency(r *Rubric, result *ValidationResult) {
	findingCounts := CountFindings(r.Findings)
	catCounts := CountResults(r.Categories)

	// If there are blocking findings (critical/high), decision should not be pass
	if findingCounts.HasBlocking() && r.Decision.Status == DecisionPass {
		result.addWarning(
			"decision.status",
			"INCONSISTENT_DECISION",
			fmt.Sprintf("decision is 'pass' but there are %d blocking findings (critical: %d, high: %d)",
				findingCounts.BlockingCount(), findingCounts.Critical, findingCounts.High),
		)
	}

	// If there are failed categories, decision should not be pass
	if catCounts.Fail > 0 && r.Decision.Status == DecisionPass {
		result.addWarning(
			"decision.status",
			"INCONSISTENT_DECISION",
			fmt.Sprintf("decision is 'pass' but %d categories failed", catCounts.Fail),
		)
	}

	// If decision is fail but no blocking issues found, it's suspicious
	if r.Decision.Status == DecisionFail && !findingCounts.HasBlocking() && catCounts.Fail == 0 {
		result.addWarning(
			"decision.status",
			"INCONSISTENT_DECISION",
			"decision is 'fail' but no blocking findings or failed categories found",
		)
	}

	// Validate decision.passed matches decision.status
	expectedPassed := r.Decision.Status == DecisionPass || r.Decision.Status == DecisionConditional
	if r.Decision.Passed != expectedPassed {
		result.addWarning(
			"decision.passed",
			"INCONSISTENT_DECISION",
			fmt.Sprintf("decision.passed is %v but decision.status is %q (expected passed=%v)",
				r.Decision.Passed, r.Decision.Status, expectedPassed),
		)
	}
}

// ValidateRubricSet validates a rubric definition (not a report).
// Returns a ValidationResult instead of []string for consistency.
func ValidateRubricSetV2(rs *RubricSet) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Use existing validation
	issues := rs.Validate()
	for _, issue := range issues {
		result.addError("", "RUBRIC_VALIDATION", issue)
	}

	// Additional enum validation for scale types
	for i, cat := range rs.Categories {
		path := fmt.Sprintf("categories[%d]", i)

		switch cat.Scale.Type {
		case ScaleTypeCategorical, ScaleTypeChecklist, ScaleTypeBinary, ScaleTypeLikert:
			// Valid
		default:
			result.addEnumError(
				path+".scale.type",
				"INVALID_SCALE_TYPE",
				fmt.Sprintf("invalid scale type %q", cat.Scale.Type),
				string(cat.Scale.Type),
				ValidScaleTypeValues(),
			)
		}
	}

	// Validate evaluation type
	if rs.EvaluationType != "" {
		switch rs.EvaluationType {
		case EvaluationTypeAnalytic, EvaluationTypeHolistic:
			// Valid
		default:
			result.addEnumError(
				"evaluationType",
				"INVALID_EVALUATION_TYPE",
				fmt.Sprintf("invalid evaluation type %q", rs.EvaluationType),
				string(rs.EvaluationType),
				ValidEvaluationTypeValues(),
			)
		}
	}

	return result
}
