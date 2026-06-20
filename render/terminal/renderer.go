// Package terminal provides ANSI-colored terminal rendering for evaluation reports.
package terminal

import (
	"fmt"
	"io"
	"strings"

	"github.com/plexusone/structured-evaluation/rubric"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"

	// Foreground colors
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	// Bright foreground colors
	BrightRed    = "\033[91m"
	BrightGreen  = "\033[92m"
	BrightYellow = "\033[93m"
	BrightBlue   = "\033[94m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"

	// Background colors
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
)

const boxWidth = 78

// Renderer renders evaluation reports with ANSI colors.
type Renderer struct {
	w        io.Writer
	useColor bool
}

// New creates a new colored terminal renderer.
func New(w io.Writer) *Renderer {
	return &Renderer{w: w, useColor: true}
}

// NewNoColor creates a renderer without ANSI colors.
func NewNoColor(w io.Writer) *Renderer {
	return &Renderer{w: w, useColor: false}
}

// SetColor enables or disables ANSI colors.
func (r *Renderer) SetColor(enabled bool) {
	r.useColor = enabled
}

// color applies color if enabled (c is the ANSI color code).
//
//nolint:unparam // c is parameterized for future flexibility
func (r *Renderer) color(c, text string) string {
	if !r.useColor {
		return text
	}
	return c + text + Reset
}

// Render outputs the evaluation report with ANSI colors.
func (r *Renderer) Render(report *rubric.Rubric) error {
	var b strings.Builder

	// Header
	b.WriteString(r.color(Cyan, header()))
	b.WriteString("\n")
	title := strings.ToUpper(report.ReviewType) + " EVALUATION"
	b.WriteString(r.coloredCenterLine(title, BrightWhite+Bold))
	b.WriteString("\n")
	b.WriteString(r.color(Cyan, separator()))
	b.WriteString("\n")

	// Document info
	b.WriteString(r.paddedLine(fmt.Sprintf("%sDocument:%s %s",
		r.c(Gray), r.c(Reset), truncate(report.Metadata.Document, 60))))
	b.WriteString("\n")

	if report.Metadata.DocumentTitle != "" {
		b.WriteString(r.paddedLine(fmt.Sprintf("%sTitle:%s    %s",
			r.c(Gray), r.c(Reset), truncate(report.Metadata.DocumentTitle, 60))))
		b.WriteString("\n")
	}

	// Category summary
	catCounts := report.Decision.CategoryCounts
	summaryLine := fmt.Sprintf("%sResults:%s  %s%d pass%s, %s%d partial%s, %s%d fail%s",
		r.c(Gray), r.c(Reset),
		r.c(Green), catCounts.Pass, r.c(Reset),
		r.c(Yellow), catCounts.Partial, r.c(Reset),
		r.c(Red), catCounts.Fail, r.c(Reset))
	b.WriteString(r.paddedLineRaw(summaryLine))
	b.WriteString("\n")

	// Decision with color
	decisionColor := r.decisionColor(report.Decision.Status)
	findCounts := report.Decision.FindingCounts
	decisionLine := fmt.Sprintf("%sDecision:%s %s%s%s",
		r.c(Gray), r.c(Reset),
		r.c(decisionColor+Bold), strings.ToUpper(string(report.Decision.Status)), r.c(Reset))
	if findCounts.Total > 0 {
		decisionLine += fmt.Sprintf(" (%s%d Critical%s, %s%d High%s, %s%d Medium%s)",
			r.c(Red), findCounts.Critical, r.c(Reset),
			r.c(Red), findCounts.High, r.c(Reset),
			r.c(Yellow), findCounts.Medium, r.c(Reset))
	}
	b.WriteString(r.paddedLineRaw(decisionLine))
	b.WriteString("\n")

	// Category scores section
	b.WriteString(r.color(Cyan, separator()))
	b.WriteString("\n")
	b.WriteString(r.coloredPaddedLine("RESULTS BY CATEGORY", BrightWhite+Bold))
	b.WriteString("\n")
	b.WriteString(r.color(Cyan, separator()))
	b.WriteString("\n")

	for _, cr := range report.Categories {
		line := r.formatCategoryLine(cr)
		b.WriteString(r.paddedLineRaw(line))
		b.WriteString("\n")
	}

	// Findings by severity
	if len(report.Findings) > 0 {
		b.WriteString(r.color(Cyan, separator()))
		b.WriteString("\n")
		findingsTitle := fmt.Sprintf("FINDINGS (%s%d Critical%s, %s%d High%s, %s%d Medium%s)",
			r.c(Red), findCounts.Critical, r.c(Reset),
			r.c(Red), findCounts.High, r.c(Reset),
			r.c(Yellow), findCounts.Medium, r.c(Reset))
		b.WriteString(r.paddedLineRaw(findingsTitle))
		b.WriteString("\n")
		b.WriteString(r.color(Cyan, separator()))
		b.WriteString("\n")

		// Group by severity
		for _, sev := range rubric.AllSeverities() {
			for _, f := range report.Findings {
				if f.Severity == sev {
					sevColor := r.severityColor(sev)
					icon := sev.Icon()
					b.WriteString(r.paddedLineRaw(fmt.Sprintf("%s %s%-8s%s [%s%s%s]",
						icon, r.c(sevColor), strings.ToUpper(string(sev)), r.c(Reset),
						r.c(Cyan), f.Category, r.c(Reset))))
					b.WriteString("\n")
					b.WriteString(r.paddedLineRaw(fmt.Sprintf("          %s", truncate(f.Title, 60))))
					b.WriteString("\n")
					if f.Recommendation != "" {
						b.WriteString(r.paddedLineRaw(fmt.Sprintf("          %s→%s %s",
							r.c(Green), r.c(Reset), truncate(f.Recommendation, 58))))
						b.WriteString("\n")
					}
					b.WriteString(r.paddedLine(""))
					b.WriteString("\n")
				}
			}
		}
	}

	// Next steps
	if len(report.NextSteps.Immediate) > 0 || report.NextSteps.RerunCommand != "" {
		b.WriteString(r.color(Cyan, separator()))
		b.WriteString("\n")
		b.WriteString(r.coloredPaddedLine("NEXT STEPS", BrightWhite+Bold))
		b.WriteString("\n")
		b.WriteString(r.color(Cyan, separator()))
		b.WriteString("\n")

		for _, action := range report.NextSteps.Immediate {
			b.WriteString(r.paddedLineRaw(fmt.Sprintf("  %s🔴%s %s",
				r.c(Red), r.c(Reset), truncate(action.Action, 65))))
			b.WriteString("\n")
		}

		for _, action := range report.NextSteps.Recommended {
			b.WriteString(r.paddedLineRaw(fmt.Sprintf("  %s🟡%s %s",
				r.c(Yellow), r.c(Reset), truncate(action.Action, 65))))
			b.WriteString("\n")
		}

		if report.NextSteps.RerunCommand != "" {
			b.WriteString(r.paddedLine(""))
			b.WriteString("\n")
			b.WriteString(r.paddedLineRaw(fmt.Sprintf("%sRe-run:%s %s",
				r.c(Gray), r.c(Reset), report.NextSteps.RerunCommand)))
			b.WriteString("\n")
		}
	}

	// Final message
	b.WriteString(r.color(Cyan, separator()))
	b.WriteString("\n")
	finalMsg := r.finalMessage(report)
	b.WriteString(r.coloredCenterLine(finalMsg, r.decisionColor(report.Decision.Status)+Bold))
	b.WriteString("\n")
	b.WriteString(r.color(Cyan, footer()))
	b.WriteString("\n")

	_, err := fmt.Fprint(r.w, b.String())
	return err
}

// c returns color code if colors enabled, empty string otherwise.
func (r *Renderer) c(color string) string {
	if !r.useColor {
		return ""
	}
	return color
}

func (r *Renderer) formatCategoryLine(cr rubric.CategoryResult) string {
	name := cr.Category
	if len(name) > 24 {
		name = name[:21] + "..."
	}

	var scoreText string
	var scoreColor string
	var icon string

	if cr.NumericScore != nil {
		score := int(*cr.NumericScore)
		scoreText = fmt.Sprintf("%d/5", score)
		scoreColor = r.numericScoreColor(score)
		icon = numericScoreIcon(score)
	} else {
		scoreColor = r.scoreColor(cr.Score)
		icon = cr.Score.Icon()
		scoreText = strings.ToUpper(string(cr.Score))
	}

	reasoning := truncate(cr.Reasoning, 35)

	return fmt.Sprintf("  %-24s %s %s%-7s%s  %s%s%s",
		name, icon, r.c(scoreColor), scoreText, r.c(Reset),
		r.c(Gray), reasoning, r.c(Reset))
}

func (r *Renderer) numericScoreColor(score int) string {
	switch {
	case score >= 5:
		return Green
	case score >= 3:
		return Yellow
	default:
		return Red
	}
}

func numericScoreIcon(score int) string {
	switch {
	case score >= 5:
		return "🟢"
	case score >= 3:
		return "🟡"
	default:
		return "🔴"
	}
}

func (r *Renderer) finalMessage(report *rubric.Rubric) string {
	catCounts := report.Decision.CategoryCounts
	switch report.Decision.Status {
	case rubric.DecisionPass:
		return fmt.Sprintf("✅ %s PASSED (%d/%d categories)",
			strings.ToUpper(report.ReviewType), catCounts.Pass, catCounts.Total)
	case rubric.DecisionConditional:
		return fmt.Sprintf("⚠️ %s CONDITIONAL (%d pass, %d partial)",
			strings.ToUpper(report.ReviewType), catCounts.Pass, catCounts.Partial)
	case rubric.DecisionFail:
		return fmt.Sprintf("❌ %s BLOCKED - %d issues to resolve",
			strings.ToUpper(report.ReviewType), report.Decision.FindingCounts.BlockingCount())
	case rubric.DecisionHumanReview:
		return fmt.Sprintf("👤 %s NEEDS HUMAN REVIEW", strings.ToUpper(report.ReviewType))
	default:
		return fmt.Sprintf("📋 %s: %d/%d categories passed",
			strings.ToUpper(report.ReviewType), catCounts.Pass, catCounts.Total)
	}
}

func (r *Renderer) decisionColor(status rubric.DecisionStatus) string {
	switch status {
	case rubric.DecisionPass:
		return Green
	case rubric.DecisionConditional:
		return Yellow
	case rubric.DecisionFail:
		return Red
	case rubric.DecisionHumanReview:
		return Magenta
	default:
		return White
	}
}

func (r *Renderer) scoreColor(score rubric.ScoreValue) string {
	switch score {
	case rubric.ScorePass:
		return Green
	case rubric.ScorePartial:
		return Yellow
	case rubric.ScoreFail:
		return Red
	default:
		return White
	}
}

func (r *Renderer) severityColor(sev rubric.Severity) string {
	switch sev {
	case rubric.SeverityCritical, rubric.SeverityHigh:
		return Red
	case rubric.SeverityMedium:
		return Yellow
	case rubric.SeverityLow:
		return Green
	case rubric.SeverityInfo:
		return Cyan
	default:
		return White
	}
}

// Box drawing functions

func header() string {
	return "╔" + strings.Repeat("═", boxWidth) + "╗"
}

func separator() string {
	return "╠" + strings.Repeat("═", boxWidth) + "╣"
}

func footer() string {
	return "╚" + strings.Repeat("═", boxWidth) + "╝"
}

func (r *Renderer) coloredCenterLine(text, color string) string {
	visualLen := visualLength(text)
	padding := maxInt(0, boxWidth-visualLen)
	left := padding / 2
	right := padding - left
	return r.c(Cyan) + "║" + r.c(Reset) + strings.Repeat(" ", left) +
		r.c(color) + text + r.c(Reset) +
		strings.Repeat(" ", right) + r.c(Cyan) + "║" + r.c(Reset)
}

func (r *Renderer) paddedLine(text string) string {
	visualLen := visualLength(text)
	padding := maxInt(0, boxWidth-visualLen-1)
	return r.c(Cyan) + "║" + r.c(Reset) + " " + text + strings.Repeat(" ", padding) + r.c(Cyan) + "║" + r.c(Reset)
}

func (r *Renderer) paddedLineRaw(text string) string {
	// For lines with ANSI codes, we need to calculate visible length differently
	visualLen := visualLengthWithANSI(text)
	padding := maxInt(0, boxWidth-visualLen-1)
	return r.c(Cyan) + "║" + r.c(Reset) + " " + text + strings.Repeat(" ", padding) + r.c(Cyan) + "║" + r.c(Reset)
}

func (r *Renderer) coloredPaddedLine(text, color string) string {
	visualLen := visualLength(text)
	padding := maxInt(0, boxWidth-visualLen-1)
	return r.c(Cyan) + "║" + r.c(Reset) + " " + r.c(color) + text + r.c(Reset) + strings.Repeat(" ", padding) + r.c(Cyan) + "║" + r.c(Reset)
}

func visualLength(s string) int {
	length := 0
	for _, r := range s {
		if r >= 0x1F300 && r <= 0x1FAFF {
			length += 2
		} else if r >= 0x2600 && r <= 0x27BF {
			length += 2
		} else {
			length++
		}
	}
	return length
}

// visualLengthWithANSI calculates visible length ignoring ANSI escape codes.
func visualLengthWithANSI(s string) int {
	// Strip ANSI codes for length calculation
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		if r >= 0x1F300 && r <= 0x1FAFF {
			length += 2
		} else if r >= 0x2600 && r <= 0x27BF {
			length += 2
		} else {
			length++
		}
	}
	return length
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

//nolint:unparam // a is 0 for clamping, standard utility function pattern
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
