package terminal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/plexusone/structured-evaluation/evaluation"
)

func TestRenderer_Render(t *testing.T) {
	report := evaluation.NewEvaluationReport("article", "test-article.md")
	report.Metadata.DocumentTitle = "Test Article"

	report.AddCategoryResult(evaluation.CategoryResult{
		Category:  "technical_accuracy",
		Score:     evaluation.ScorePass,
		Reasoning: "All details correct",
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

	// Check for box drawing characters
	if !strings.Contains(output, "╔") {
		t.Error("Expected output to contain box drawing characters")
	}

	// Check for UTF8 icons
	if !strings.Contains(output, "🟢") || !strings.Contains(output, "🟡") {
		t.Error("Expected output to contain UTF8 score icons")
	}

	// Check for ANSI codes (colored output)
	if !strings.Contains(output, "\033[") {
		t.Error("Expected output to contain ANSI color codes")
	}

	// Check for content
	if !strings.Contains(output, "ARTICLE EVALUATION") {
		t.Error("Expected output to contain review type header")
	}
}

func TestRenderer_NoColor(t *testing.T) {
	report := evaluation.NewEvaluationReport("test", "doc.md")
	report.AddCategoryResult(evaluation.CategoryResult{
		Category: "quality",
		Score:    evaluation.ScorePass,
	})
	report.Finalize(nil, "")

	var buf bytes.Buffer
	renderer := NewNoColor(&buf)
	err := renderer.Render(report)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Should NOT contain ANSI codes
	if strings.Contains(output, "\033[") {
		t.Error("NoColor renderer should not contain ANSI codes")
	}

	// Should still contain UTF8 icons
	if !strings.Contains(output, "🟢") {
		t.Error("Expected output to contain UTF8 icons even without color")
	}
}

func TestRenderer_SetColor(t *testing.T) {
	report := evaluation.NewEvaluationReport("test", "doc.md")
	report.AddCategoryResult(evaluation.CategoryResult{
		Category: "quality",
		Score:    evaluation.ScorePass,
	})
	report.Finalize(nil, "")

	// Start with color
	var buf1 bytes.Buffer
	renderer := New(&buf1)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(buf1.String(), "\033[") {
		t.Error("Expected ANSI codes when color enabled")
	}

	// Disable color
	var buf2 bytes.Buffer
	renderer = New(&buf2)
	renderer.SetColor(false)
	if err := renderer.Render(report); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if strings.Contains(buf2.String(), "\033[") {
		t.Error("Expected no ANSI codes when color disabled")
	}
}

func TestRenderer_WithFindings(t *testing.T) {
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
		Description:    "Details don't match",
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

	// Check for findings section
	if !strings.Contains(output, "FINDINGS") {
		t.Error("Expected output to contain FINDINGS section")
	}
	if !strings.Contains(output, "CVE mismatch") {
		t.Error("Expected output to contain finding title")
	}
	if !strings.Contains(output, "→") {
		t.Error("Expected output to contain recommendation arrow")
	}
}

func TestRenderer_DecisionColors(t *testing.T) {
	tests := []struct {
		status   evaluation.DecisionStatus
		contains string
	}{
		{evaluation.DecisionPass, "PASS"},
		{evaluation.DecisionConditional, "CONDITIONAL"},
		{evaluation.DecisionFail, "BLOCKED"},
		{evaluation.DecisionHumanReview, "HUMAN REVIEW"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			report := evaluation.NewEvaluationReport("test", "doc.md")
			report.Decision.Status = tt.status
			report.Decision.CategoryCounts.Total = 1

			var buf bytes.Buffer
			renderer := New(&buf)
			if err := renderer.Render(report); err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if !strings.Contains(buf.String(), tt.contains) {
				t.Errorf("Expected output to contain %q for status %s", tt.contains, tt.status)
			}
		})
	}
}

func TestVisualLengthWithANSI(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"\033[31mhello\033[0m", 5},     // red "hello"
		{"\033[1m\033[32mhi\033[0m", 2}, // bold green "hi"
		{"🟢 pass", 7},                   // emoji + space + "pass"
		{"\033[31m🔴\033[0m", 2},         // colored emoji
	}

	for _, tt := range tests {
		result := visualLengthWithANSI(tt.input)
		if result != tt.expected {
			t.Errorf("visualLengthWithANSI(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}
