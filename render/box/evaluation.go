// Package box provides box-format terminal rendering for reports.
package box

import (
	"fmt"
	"io"
	"strings"

	"github.com/plexusone/structured-evaluation/rubric"
)

// EvaluationRenderer renders evaluation reports in ASCII box format.
type EvaluationRenderer struct {
	w io.Writer
}

// NewEvaluationRenderer creates a new evaluation box renderer.
func NewEvaluationRenderer(w io.Writer) *EvaluationRenderer {
	return &EvaluationRenderer{w: w}
}

// Render outputs the evaluation report in ASCII box format.
func (r *EvaluationRenderer) Render(report *rubric.Rubric) error {
	var b strings.Builder

	// Header
	b.WriteString(header())
	b.WriteString("\n")
	title := "EVALUATION REPORT"
	if report.Metadata.DocumentTitle != "" {
		title = strings.ToUpper(report.Metadata.DocumentTitle)
		if len(title) > boxWidth-4 {
			title = title[:boxWidth-7] + "..."
		}
	}
	b.WriteString(centerLine(title))
	b.WriteString("\n")
	b.WriteString(separator())
	b.WriteString("\n")

	// Decision summary
	decisionIcon := decisionIcon(report.Decision.Status)
	decisionText := strings.ToUpper(string(report.Decision.Status))
	b.WriteString(centerLine(fmt.Sprintf("%s  %s  %s", decisionIcon, decisionText, decisionIcon)))
	b.WriteString("\n")
	b.WriteString(separator())
	b.WriteString("\n")

	// Category table header
	b.WriteString(paddedLine("CATEGORY                   ST SCORE  DETAIL"))
	b.WriteString("\n")
	b.WriteString(thinSeparator())
	b.WriteString("\n")

	// Categories
	for _, cr := range report.Categories {
		b.WriteString(paddedLine(formatCategoryLine(cr)))
		b.WriteString("\n")
	}

	// Summary counts
	b.WriteString(separator())
	b.WriteString("\n")
	catCounts := report.Decision.CategoryCounts
	findCounts := report.Decision.FindingCounts

	b.WriteString(paddedLine(fmt.Sprintf("Categories: %d pass, %d partial, %d fail",
		catCounts.Pass, catCounts.Partial, catCounts.Fail)))
	b.WriteString("\n")
	b.WriteString(paddedLine(fmt.Sprintf("Findings:   %d critical, %d high, %d medium, %d low",
		findCounts.Critical, findCounts.High, findCounts.Medium, findCounts.Low)))
	b.WriteString("\n")

	// Findings section (if any high or critical)
	if findCounts.Critical > 0 || findCounts.High > 0 || findCounts.Medium > 0 {
		b.WriteString(separator())
		b.WriteString("\n")
		b.WriteString(paddedLine("FINDINGS"))
		b.WriteString("\n")
		b.WriteString(thinSeparator())
		b.WriteString("\n")

		for _, f := range report.Findings {
			if f.Severity == rubric.SeverityCritical || f.Severity == rubric.SeverityHigh || f.Severity == rubric.SeverityMedium {
				icon := severityIcon(f.Severity)
				text := f.Title
				if text == "" {
					text = f.Description
				}
				line := fmt.Sprintf("%s [%s] %s", icon, f.Category, truncateStr(text, 50))
				b.WriteString(paddedLine(line))
				b.WriteString("\n")
			}
		}
	}

	// Next steps (if any immediate)
	if len(report.NextSteps.Immediate) > 0 {
		b.WriteString(separator())
		b.WriteString("\n")
		b.WriteString(paddedLine("IMMEDIATE ACTIONS REQUIRED"))
		b.WriteString("\n")
		b.WriteString(thinSeparator())
		b.WriteString("\n")
		for _, action := range report.NextSteps.Immediate {
			line := fmt.Sprintf("[ ] %s", truncateStr(action.Action, 70))
			b.WriteString(paddedLine(line))
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString(footer())
	b.WriteString("\n")

	_, err := fmt.Fprint(r.w, b.String())
	return err
}

func formatCategoryLine(cr rubric.CategoryResult) string {
	name := cr.Category
	if len(name) > 26 {
		name = name[:23] + "..."
	}

	var scoreStr string
	var statusIcon string

	if cr.NumericScore != nil {
		score := int(*cr.NumericScore)
		scoreStr = fmt.Sprintf("%d/5", score)
		statusIcon = numericIcon(score)
	} else {
		scoreStr = string(cr.Score)
		statusIcon = categoricalIcon(cr.Score)
	}

	// Truncate reasoning to fit in remaining space
	reasoning := truncateStr(cr.Reasoning, 38)

	return fmt.Sprintf("%-26s %s %3s  %s", name, statusIcon, scoreStr, reasoning)
}

func numericIcon(score int) string {
	switch {
	case score >= 5:
		return "🟢"
	case score >= 3:
		return "🟡"
	default:
		return "🔴"
	}
}

func categoricalIcon(score rubric.ScoreValue) string {
	switch score {
	case rubric.ScorePass:
		return "🟢"
	case rubric.ScorePartial:
		return "🟡"
	case rubric.ScoreFail:
		return "🔴"
	default:
		return "⚪"
	}
}

func decisionIcon(status rubric.DecisionStatus) string {
	switch status {
	case rubric.DecisionPass:
		return "✅"
	case rubric.DecisionConditional:
		return "⚠️"
	case rubric.DecisionFail:
		return "❌"
	case rubric.DecisionHumanReview:
		return "👤"
	default:
		return "📋"
	}
}

func severityIcon(sev rubric.Severity) string {
	switch sev {
	case rubric.SeverityCritical:
		return "🔴"
	case rubric.SeverityHigh:
		return "🔴"
	case rubric.SeverityMedium:
		return "🟡"
	case rubric.SeverityLow:
		return "🟢"
	case rubric.SeverityInfo:
		return "ℹ️"
	default:
		return "⚪"
	}
}

func thinSeparator() string {
	return "╟" + strings.Repeat("─", boxWidth) + "╢"
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
