# Report Rendering

Structured Evaluation provides multiple renderers for displaying reports in different formats.

## Available Renderers

| Package | Format | Use Case |
|---------|--------|----------|
| `render/terminal` | ANSI terminal | CLI output with colors |
| `render/markdown` | Markdown | Documentation, GitHub |
| `render/detailed` | Detailed terminal | Verbose CLI output |
| `render/box` | Box format | Summary reports |

## Terminal Renderer

ANSI-colored output with UTF8 icons:

```go
import "github.com/plexusone/structured-evaluation/render/terminal"

renderer := terminal.New(os.Stdout)
err := renderer.Render(report)
```

### Output Example

```
┌──────────────────────────────────────────────────────────────────────────┐
│ EVALUATION REPORT: requirements.md                                        │
├──────────────────────────────────────────────────────────────────────────┤
│ Results:  2 pass, 1 partial, 1 fail                                       │
│ Decision: CONDITIONAL (1 Critical, 0 High, 2 Medium)                      │
├──────────────────────────────────────────────────────────────────────────┤
│ RESULTS BY CATEGORY                                                       │
├──────────────────────────────────────────────────────────────────────────┤
│   problem_definition     🟢 PASS     Clear problem with measurable impact │
│   user_stories           🟡 PARTIAL  Some lack acceptance criteria        │
│   success_metrics        🔴 FAIL     No quantitative metrics defined      │
│   scope_definition       🟢 PASS     Clear boundaries established         │
├──────────────────────────────────────────────────────────────────────────┤
│ ⚠️ PRD-REVIEW CONDITIONAL (2/4 categories passed)                         │
└──────────────────────────────────────────────────────────────────────────┘
```

### Color Options

```go
// Disable colors (e.g., for piping)
renderer := terminal.New(os.Stdout, terminal.WithNoColor())

// Custom width
renderer := terminal.New(os.Stdout, terminal.WithWidth(100))
```

## Markdown Renderer

Generates Markdown suitable for documentation or GitHub:

```go
import "github.com/plexusone/structured-evaluation/render/markdown"

renderer := markdown.New(os.Stdout)
err := renderer.Render(report)
```

### Output Example

```markdown
## Evaluation Report: requirements.md

### Summary

**Overall Decision: CONDITIONAL** ⚠️

| Metric | Value |
|--------|-------|
| Categories | 2 pass, 1 partial, 1 fail |
| Findings | 1 critical, 0 high, 2 medium |
| Decision | Conditional |

---

### Category Results

| Category | Score | Reasoning |
|----------|-------|-----------|
| problem_definition | 🟢 PASS | Clear problem with measurable impact |
| user_stories | 🟡 PARTIAL | Some lack acceptance criteria |
| success_metrics | 🔴 FAIL | No quantitative metrics defined |

---

### Findings

#### 🔴 Critical

**Missing success metrics**
- Category: success_metrics
- The PRD does not define how success will be measured
- *Recommendation:* Add 2-3 quantitative KPIs with target values
```

## Detailed Renderer

Verbose terminal output with full finding details:

```go
import "github.com/plexusone/structured-evaluation/render/detailed"

renderer := detailed.NewTerminal(os.Stdout)
err := renderer.Render(report)
```

## Box Renderer

For summary reports (GO/NO-GO):

```go
import "github.com/plexusone/structured-evaluation/render/box"

renderer := box.NewRenderer(os.Stdout)
err := renderer.Render(&summaryReport)
```

### Output Example

```
╔══════════════════════════════════════════════════════════════════════════╗
║                    RELEASE VALIDATION: my-service v2.0.0                  ║
╠══════════════════════════════════════════════════════════════════════════╣
║ TESTING                                                                   ║
║   ✅ unit-tests         Coverage: 92%                                     ║
║   ✅ integration-tests  All 47 tests pass                                 ║
║   ⚠️ e2e-tests          2 flaky tests skipped                             ║
╠══════════════════════════════════════════════════════════════════════════╣
║ SECURITY                                                                  ║
║   ✅ sast               No critical findings                              ║
║   ✅ dependency-scan    All deps up to date                               ║
╠══════════════════════════════════════════════════════════════════════════╣
║                              ✅ GO FOR RELEASE                            ║
╚══════════════════════════════════════════════════════════════════════════╝
```

## JSON Output

For programmatic use:

```go
import "encoding/json"

output, err := json.MarshalIndent(&report, "", "  ")
if err != nil {
    return err
}
fmt.Println(string(output))
```

## CLI Usage

```bash
# Terminal format (default)
sevaluation render report.json --format=terminal

# Markdown
sevaluation render report.json --format=markdown > report.md

# Detailed
sevaluation render report.json --format=detailed

# Box (for summary reports)
sevaluation render summary.json --format=box

# JSON
sevaluation render report.json --format=json
```

## Writing to Files

```go
file, err := os.Create("report.md")
if err != nil {
    return err
}
defer file.Close()

renderer := markdown.New(file)
return renderer.Render(report)
```

## Custom Renderers

Implement the Renderer interface:

```go
type Renderer interface {
    Render(report *rubric.Rubric) error
}

// Example: HTML renderer
type HTMLRenderer struct {
    w io.Writer
}

func (r *HTMLRenderer) Render(report *rubric.Rubric) error {
    // Generate HTML output
    tmpl := template.Must(template.ParseFiles("report.html"))
    return tmpl.Execute(r.w, report)
}
```

## Next Steps

- [CLI Commands](../cli/commands.md) - Full CLI reference
- [Report Types](../concepts/report-types.md) - Understanding reports
