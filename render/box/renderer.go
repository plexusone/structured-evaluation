// Package box provides box-format terminal rendering for summary reports.
package box

import (
	"fmt"
	"io"
	"strings"

	"github.com/agentplexus/structured-evaluation/combine"
	"github.com/agentplexus/structured-evaluation/summary"
)

const boxWidth = 78

// Renderer renders summary reports in box format.
type Renderer struct {
	w io.Writer
}

// New creates a new box renderer.
func New(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

// Render outputs the summary report in box format.
// Automatically sorts teams by DAG order.
func (r *Renderer) Render(report *summary.SummaryReport) error {
	combine.SortReportByDAG(report)

	var b strings.Builder

	// Header
	b.WriteString(header())
	b.WriteString("\n")
	b.WriteString(centerLine("TEAM STATUS REPORT"))
	b.WriteString("\n")
	b.WriteString(separator())
	b.WriteString("\n")

	// Project info
	b.WriteString(paddedLine(fmt.Sprintf("Project: %s", report.Project)))
	b.WriteString("\n")
	if report.Version != "" {
		b.WriteString(paddedLine(fmt.Sprintf("Version: %s", report.Version)))
		b.WriteString("\n")
	}
	if report.Target != "" {
		b.WriteString(paddedLine(fmt.Sprintf("Target:  %s", report.Target)))
		b.WriteString("\n")
	}

	// Phase
	if report.Phase != "" {
		b.WriteString(separator())
		b.WriteString("\n")
		b.WriteString(paddedLine(report.Phase))
		b.WriteString("\n")
	}

	// Teams
	for _, team := range report.Teams {
		b.WriteString(separator())
		b.WriteString("\n")
		b.WriteString(paddedLine(fmt.Sprintf("%s (%s)", team.ID, team.Name)))
		b.WriteString("\n")

		for _, task := range team.Tasks {
			b.WriteString(paddedLine(formatTaskLine(task)))
			b.WriteString("\n")
		}
	}

	// Final message
	b.WriteString(separator())
	b.WriteString("\n")
	b.WriteString(centerLine(report.FinalMessage()))
	b.WriteString("\n")
	b.WriteString(footer())
	b.WriteString("\n")

	_, err := fmt.Fprint(r.w, b.String())
	return err
}

func formatTaskLine(task summary.TaskResult) string {
	id := task.ID
	if len(id) > 24 {
		id = id[:21] + "..."
	}

	icon := task.Status.Icon()
	statusText := string(task.Status)

	detail := task.Detail
	if len(detail) > 38 {
		detail = detail[:35] + "..."
	}

	return fmt.Sprintf("  %-24s %s %-5s %s", id, icon, statusText, detail)
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

func centerLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen)
	left := padding / 2
	right := padding - left
	return "║" + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "║"
}

func paddedLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen-1)
	return "║ " + text + strings.Repeat(" ", padding) + "║"
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
