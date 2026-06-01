package summary

import (
	"encoding/json"
	"testing"
)

func TestSummaryReport_EmbedReports(t *testing.T) {
	report := NewSummaryReport("test-project", "1.0.0", "VALIDATION")

	// Create a mock evaluation report structure
	mockEval := map[string]any{
		"metadata": map[string]any{
			"document":    "test.md",
			"generatedAt": "2026-01-01T00:00:00Z",
		},
		"reviewType": "article",
		"decision": map[string]any{
			"status": "pass",
			"passed": true,
		},
	}

	// Create a mock claims report structure
	mockClaims := map[string]any{
		"metadata": map[string]any{
			"document":    "test.md",
			"generatedAt": "2026-01-01T00:00:00Z",
		},
		"claims": []any{},
		"decision": map[string]any{
			"status": "pass",
			"passed": true,
		},
	}

	// Embed reports
	if err := report.EmbedEvaluationReport("quality-review", mockEval); err != nil {
		t.Fatalf("Failed to embed evaluation report: %v", err)
	}

	if err := report.EmbedClaimsReport("source-validation", mockClaims); err != nil {
		t.Fatalf("Failed to embed claims report: %v", err)
	}

	// Verify embedded reports exist
	if !report.HasEmbeddedReports() {
		t.Error("Expected HasEmbeddedReports to return true")
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal report: %v", err)
	}

	// Unmarshal back
	var parsed SummaryReport
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal report: %v", err)
	}

	// Verify embedded reports survived round-trip
	if !parsed.HasEmbeddedReports() {
		t.Error("Expected HasEmbeddedReports to return true after round-trip")
	}

	if len(parsed.EmbeddedReports.Evaluations) != 1 {
		t.Errorf("Expected 1 embedded evaluation, got %d", len(parsed.EmbeddedReports.Evaluations))
	}

	if len(parsed.EmbeddedReports.Claims) != 1 {
		t.Errorf("Expected 1 embedded claims report, got %d", len(parsed.EmbeddedReports.Claims))
	}

	// Retrieve and verify evaluation
	var retrievedEval map[string]any
	if err := parsed.GetEmbeddedEvaluation("quality-review", &retrievedEval); err != nil {
		t.Fatalf("Failed to get embedded evaluation: %v", err)
	}

	if retrievedEval["reviewType"] != "article" {
		t.Errorf("Expected reviewType 'article', got %v", retrievedEval["reviewType"])
	}

	// Retrieve and verify claims
	var retrievedClaims map[string]any
	if err := parsed.GetEmbeddedClaims("source-validation", &retrievedClaims); err != nil {
		t.Fatalf("Failed to get embedded claims: %v", err)
	}

	decision, ok := retrievedClaims["decision"].(map[string]any)
	if !ok {
		t.Fatal("Failed to get decision from claims")
	}
	if decision["status"] != "pass" {
		t.Errorf("Expected claims decision status 'pass', got %v", decision["status"])
	}
}

func TestSummaryReport_NoEmbeddedReports(t *testing.T) {
	report := NewSummaryReport("test-project", "1.0.0", "VALIDATION")

	if report.HasEmbeddedReports() {
		t.Error("Expected HasEmbeddedReports to return false for new report")
	}

	// Retrieving from nil should not error
	var target map[string]any
	if err := report.GetEmbeddedEvaluation("nonexistent", &target); err != nil {
		t.Errorf("Expected no error for missing report, got: %v", err)
	}
}
