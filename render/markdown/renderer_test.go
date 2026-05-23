package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/plexusone/structured-evaluation/evaluation"
)

func TestRenderer_Render(t *testing.T) {
	report := evaluation.NewEvaluationReport("article", "test-article.md")
	report.Metadata.DocumentTitle = "Test Article"
	report.RubricID = "vulnerability-article-v1"
	report.RubricVersion = "1.0.0"

	report.AddCategoryResult(evaluation.CategoryResult{
		Category:  "technical_accuracy",
		Score:     evaluation.ScorePass,
		Reasoning: "All details correct",
		Evidence:  []string{"CVE-2026-12345 matches NVD"},
	})
	report.AddCategoryResult(evaluation.CategoryResult{
		Category:  "completeness",
		Score:     evaluation.ScorePartial,
		Reasoning: "Missing optional sections",
	})

	report.Finalize(nil, "sevaluation check test-article.md")

	var buf bytes.Buffer
	renderer := New(&buf)
	err := renderer.Render(report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check for expected content
	checks := []string{
		"## Evaluation Report",
		"### Summary",
		"**Overall Decision:",
		"### Category Results",
		"| Category |",
		"technical_accuracy",
		"completeness",
		"🟢 Pass",
		"🟡 Partial",
		"### Evidence",
		"CVE-2026-12345",
		"structured-evaluation",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected output to contain %q", check)
		}
	}
}

func TestRenderer_RenderWithFindings(t *testing.T) {
	report := evaluation.NewEvaluationReport("article", "test.md")
	report.AddCategoryResult(evaluation.CategoryResult{
		Category: "accuracy",
		Score:    evaluation.ScoreFail,
	})
	report.AddFinding(evaluation.Finding{
		ID:             "F1",
		Category:       "accuracy",
		Severity:       evaluation.SeverityHigh,
		Title:          "CVE mismatch",
		Description:    "CVE details don't match NVD",
		Recommendation: "Verify against NVD",
	})
	report.Finalize(nil, "")

	var buf bytes.Buffer
	renderer := New(&buf)
	err := renderer.Render(report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	checks := []string{
		"### Findings",
		"🔴 High",
		"CVE mismatch",
		"Recommendation",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected output to contain %q", check)
		}
	}
}

func TestRenderer_RenderWithRubric(t *testing.T) {
	rubric := evaluation.NewRubricSet("test-v1", "Test Rubric", "1.0.0")
	cat := evaluation.NewCategory("accuracy", "Accuracy", "Test accuracy").
		SetRequired(true).
		SetWeight(2.0)
	rubric.AddCategory(*cat)

	report := evaluation.NewEvaluationReport("test", "doc.md")
	report.AddCategoryResult(evaluation.CategoryResult{
		Category: "accuracy",
		Score:    evaluation.ScorePass,
	})
	report.Finalize(rubric, "")

	var buf bytes.Buffer
	renderer := New(&buf)
	err := renderer.RenderWithRubric(report, rubric)
	if err != nil {
		t.Fatalf("RenderWithRubric failed: %v", err)
	}

	output := buf.String()

	// Should include weight and required info
	if !strings.Contains(output, "2.0") {
		t.Error("Expected output to contain weight '2.0'")
	}
	if !strings.Contains(output, "✅") {
		t.Error("Expected output to contain required checkmark")
	}
}

func TestRenderer_NoFindings(t *testing.T) {
	report := evaluation.NewEvaluationReport("test", "doc.md")
	report.AddCategoryResult(evaluation.CategoryResult{
		Category: "quality",
		Score:    evaluation.ScorePass,
	})
	report.Finalize(nil, "")

	var buf bytes.Buffer
	renderer := New(&buf)
	err := renderer.Render(report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "**None** - No issues identified") {
		t.Error("Expected output to indicate no findings")
	}
}
