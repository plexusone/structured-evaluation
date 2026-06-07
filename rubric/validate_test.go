package rubric

import (
	"testing"
	"time"
)

func TestValidateReport_ValidReport(t *testing.T) {
	report := &Rubric{
		Metadata: ReportMetadata{
			Document:    "test.md",
			GeneratedAt: time.Now(),
		},
		ReviewType: "prd",
		Categories: []CategoryResult{
			{
				Category:  "problem_definition",
				Score:     ScorePass,
				Reasoning: "Good problem definition",
			},
			{
				Category:  "user_stories",
				Score:     ScorePartial,
				Reasoning: "Some stories missing",
			},
		},
		Findings: []Finding{
			{
				ID:       "F1",
				Category: "user_stories",
				Severity: SeverityMedium,
				Title:    "Missing stories",
			},
		},
		Decision: Decision{
			Status:    DecisionConditional,
			Passed:    true,
			Rationale: "Passed with partial scores",
			CategoryCounts: CategoryResultCounts{
				Pass:    1,
				Partial: 1,
				Fail:    0,
				Total:   2,
			},
			FindingCounts: FindingCounts{
				Medium: 1,
				Total:  1,
			},
		},
		OverallDecision: "conditional",
	}

	result := ValidateReport(report)

	if !result.Valid {
		t.Errorf("expected valid report, got invalid: %s", result.String())
	}
	if result.ErrorCount != 0 {
		t.Errorf("expected 0 errors, got %d", result.ErrorCount)
	}
}

func TestValidateReport_InvalidScore(t *testing.T) {
	report := &Rubric{
		Metadata: ReportMetadata{
			Document: "test.md",
		},
		ReviewType: "prd",
		Categories: []CategoryResult{
			{
				Category:  "test",
				Score:     ScoreValue("passed"), // Invalid: should be "pass"
				Reasoning: "test",
			},
		},
		Decision: Decision{
			Status: DecisionPass,
			Passed: true,
		},
	}

	result := ValidateReport(report)

	if result.Valid {
		t.Error("expected invalid report for bad score value")
	}
	if result.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", result.ErrorCount)
	}

	// Check that the error has the right details
	found := false
	for _, issue := range result.Issues {
		if issue.Code == "INVALID_SCORE" {
			found = true
			if issue.ActualValue != "passed" {
				t.Errorf("expected actual value 'passed', got %q", issue.ActualValue)
			}
			if len(issue.AllowedValues) != 3 {
				t.Errorf("expected 3 allowed values, got %d", len(issue.AllowedValues))
			}
		}
	}
	if !found {
		t.Error("expected INVALID_SCORE error not found")
	}
}

func TestValidateReport_InvalidSeverity(t *testing.T) {
	report := &Rubric{
		Metadata: ReportMetadata{
			Document: "test.md",
		},
		ReviewType: "prd",
		Categories: []CategoryResult{},
		Findings: []Finding{
			{
				ID:       "F1",
				Category: "test",
				Severity: Severity("blocker"), // Invalid: should be "critical"
				Title:    "Test finding",
			},
		},
		Decision: Decision{
			Status: DecisionPass,
			Passed: true,
		},
	}

	result := ValidateReport(report)

	if result.Valid {
		t.Error("expected invalid report for bad severity value")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Code == "INVALID_SEVERITY" {
			found = true
			if issue.ActualValue != "blocker" {
				t.Errorf("expected actual value 'blocker', got %q", issue.ActualValue)
			}
		}
	}
	if !found {
		t.Error("expected INVALID_SEVERITY error not found")
	}
}

func TestValidateReport_InvalidDecisionStatus(t *testing.T) {
	report := &Rubric{
		Metadata: ReportMetadata{
			Document: "test.md",
		},
		ReviewType: "prd",
		Categories: []CategoryResult{},
		Decision: Decision{
			Status: DecisionStatus("approved"), // Invalid
			Passed: true,
		},
	}

	result := ValidateReport(report)

	if result.Valid {
		t.Error("expected invalid report for bad decision status")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Code == "INVALID_DECISION_STATUS" {
			found = true
		}
	}
	if !found {
		t.Error("expected INVALID_DECISION_STATUS error not found")
	}
}

func TestValidateReport_InconsistentDecision(t *testing.T) {
	// Decision says "pass" but there's a critical finding
	report := &Rubric{
		Metadata: ReportMetadata{
			Document: "test.md",
		},
		ReviewType: "prd",
		Categories: []CategoryResult{
			{Category: "test", Score: ScorePass, Reasoning: "good"},
		},
		Findings: []Finding{
			{
				ID:       "F1",
				Category: "test",
				Severity: SeverityCritical,
				Title:    "Critical issue",
			},
		},
		Decision: Decision{
			Status: DecisionPass, // Should not be pass with critical finding
			Passed: true,
		},
	}

	result := ValidateReport(report)

	// Should be valid (no errors) but have warnings
	if !result.Valid {
		t.Errorf("expected valid (warnings only), got errors: %s", result.String())
	}
	if result.WarningCount == 0 {
		t.Error("expected at least one warning for inconsistent decision")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Code == "INCONSISTENT_DECISION" {
			found = true
		}
	}
	if !found {
		t.Error("expected INCONSISTENT_DECISION warning not found")
	}
}

func TestValidateReport_IncorrectCounts(t *testing.T) {
	report := &Rubric{
		Metadata: ReportMetadata{
			Document: "test.md",
		},
		ReviewType: "prd",
		Categories: []CategoryResult{
			{Category: "a", Score: ScorePass, Reasoning: "good"},
			{Category: "b", Score: ScorePass, Reasoning: "good"},
		},
		Decision: Decision{
			Status: DecisionPass,
			Passed: true,
			CategoryCounts: CategoryResultCounts{
				Pass:  5, // Wrong: should be 2
				Total: 5, // Wrong: should be 2
			},
		},
	}

	result := ValidateReport(report)

	if result.WarningCount < 2 {
		t.Errorf("expected at least 2 warnings for incorrect counts, got %d", result.WarningCount)
	}

	countWarnings := 0
	for _, issue := range result.Issues {
		if issue.Code == "INCORRECT_COUNT" {
			countWarnings++
		}
	}
	if countWarnings < 2 {
		t.Errorf("expected at least 2 INCORRECT_COUNT warnings, got %d", countWarnings)
	}
}

func TestValidateReport_MissingRequiredFields(t *testing.T) {
	report := &Rubric{
		// Missing metadata.document and reviewType
		Metadata:   ReportMetadata{},
		ReviewType: "",
		Categories: []CategoryResult{},
		Decision: Decision{
			Status: DecisionPass,
		},
	}

	result := ValidateReport(report)

	if result.Valid {
		t.Error("expected invalid report for missing required fields")
	}
	if result.ErrorCount < 2 {
		t.Errorf("expected at least 2 errors for missing fields, got %d", result.ErrorCount)
	}
}

func TestValidScoreValues(t *testing.T) {
	values := ValidScoreValues()
	if len(values) != 3 {
		t.Errorf("expected 3 score values, got %d", len(values))
	}

	expected := map[string]bool{"pass": true, "partial": true, "fail": true}
	for _, v := range values {
		if !expected[v] {
			t.Errorf("unexpected score value: %s", v)
		}
	}
}

func TestValidSeverityValues(t *testing.T) {
	values := ValidSeverityValues()
	if len(values) != 5 {
		t.Errorf("expected 5 severity values, got %d", len(values))
	}

	expected := map[string]bool{
		"critical": true,
		"high":     true,
		"medium":   true,
		"low":      true,
		"info":     true,
	}
	for _, v := range values {
		if !expected[v] {
			t.Errorf("unexpected severity value: %s", v)
		}
	}
}

func TestValidDecisionStatusValues(t *testing.T) {
	values := ValidDecisionStatusValues()
	if len(values) != 4 {
		t.Errorf("expected 4 decision status values, got %d", len(values))
	}

	expected := map[string]bool{
		"pass":         true,
		"conditional":  true,
		"fail":         true,
		"human_review": true,
	}
	for _, v := range values {
		if !expected[v] {
			t.Errorf("unexpected decision status value: %s", v)
		}
	}
}

func TestValidationResult_String(t *testing.T) {
	result := &ValidationResult{
		Valid:        false,
		ErrorCount:   1,
		WarningCount: 1,
		Issues: []ValidationIssue{
			{
				Path:          "categories[0].score",
				Code:          "INVALID_SCORE",
				Message:       "invalid score value",
				Severity:      ValidationError,
				ActualValue:   "passed",
				AllowedValues: []string{"pass", "partial", "fail"},
			},
			{
				Path:     "decision.status",
				Code:     "INCONSISTENT_DECISION",
				Message:  "decision inconsistent with findings",
				Severity: ValidationWarning,
			},
		},
	}

	str := result.String()
	if str == "" {
		t.Error("expected non-empty string")
	}
	if !strContains(str, "ERROR") {
		t.Error("expected ERROR in output")
	}
	if !strContains(str, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strContains(str, "passed") {
		t.Error("expected actual value in output")
	}
}

func strContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
