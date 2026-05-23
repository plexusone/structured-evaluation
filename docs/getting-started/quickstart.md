# Quick Start

This guide walks you through creating and rendering evaluation reports.

## Creating an Evaluation Report

```go
package main

import (
    "os"

    "github.com/plexusone/structured-evaluation/evaluation"
    "github.com/plexusone/structured-evaluation/render/terminal"
)

func main() {
    // Create a new evaluation report
    report := evaluation.NewEvaluationReport("prd-review", "requirements.md")

    // Add category results with categorical scores
    report.AddCategory(evaluation.CategoryResult{
        Category:  "problem_definition",
        Score:     evaluation.ScorePass,
        Reasoning: "Problem is clearly defined with measurable impact",
    })

    report.AddCategory(evaluation.CategoryResult{
        Category:  "user_stories",
        Score:     evaluation.ScorePartial,
        Reasoning: "Stories present but some lack acceptance criteria",
    })

    report.AddCategory(evaluation.CategoryResult{
        Category:  "success_metrics",
        Score:     evaluation.ScoreFail,
        Reasoning: "No quantitative success metrics defined",
    })

    // Add findings for issues discovered
    report.AddFinding(evaluation.Finding{
        Severity:       evaluation.SeverityHigh,
        Category:       "success_metrics",
        Title:          "Missing success metrics",
        Description:    "The PRD does not define how success will be measured",
        Recommendation: "Add 2-3 quantitative KPIs with target values",
    })

    report.AddFinding(evaluation.Finding{
        Severity:       evaluation.SeverityMedium,
        Category:       "user_stories",
        Title:          "Incomplete acceptance criteria",
        Description:    "3 of 8 user stories lack testable acceptance criteria",
        Recommendation: "Add Given/When/Then criteria for each story",
    })

    // Finalize computes the decision
    report.Finalize("sevaluation check requirements.md")

    // Render to terminal with colors
    renderer := terminal.New(os.Stdout)
    renderer.Render(&report)
}
```

## Understanding the Decision

After calling `Finalize()`, the report computes a decision based on:

1. **Category Results** - Count of pass/partial/fail categories
2. **Finding Severity** - Count of critical/high/medium/low findings
3. **Pass Criteria** - Configurable thresholds

```go
// Check the decision
if report.Decision.Passed {
    fmt.Println("✅ Evaluation passed")
} else {
    fmt.Printf("❌ Evaluation failed: %s\n", report.Decision.Rationale)
}

// Access category counts
counts := report.Decision.CategoryCounts
fmt.Printf("Results: %d pass, %d partial, %d fail\n",
    counts.Pass, counts.Partial, counts.Fail)
```

## Creating a Summary Report

For deterministic GO/NO-GO checks:

```go
import "github.com/plexusone/structured-evaluation/summary"

report := summary.NewSummaryReport("my-service", "v1.0.0", "Release Validation")

report.AddTeam(summary.TeamSection{
    ID:   "qa",
    Name: "Quality Assurance",
    Tasks: []summary.TaskResult{
        {ID: "unit-tests", Status: summary.StatusGo, Detail: "Coverage: 92%"},
        {ID: "integration-tests", Status: summary.StatusGo, Detail: "All 47 tests pass"},
        {ID: "e2e-tests", Status: summary.StatusWarn, Detail: "2 flaky tests skipped"},
    },
})

report.AddTeam(summary.TeamSection{
    ID:   "security",
    Name: "Security",
    Tasks: []summary.TaskResult{
        {ID: "sast", Status: summary.StatusGo, Detail: "No critical findings"},
        {ID: "dependency-scan", Status: summary.StatusGo, Detail: "All deps up to date"},
    },
})
```

## Rendering Reports

### Terminal Output (ANSI Colors)

```go
import "github.com/plexusone/structured-evaluation/render/terminal"

renderer := terminal.New(os.Stdout)
renderer.Render(&report)
```

### Markdown Output

```go
import "github.com/plexusone/structured-evaluation/render/markdown"

renderer := markdown.New(os.Stdout)
renderer.Render(&report)
```

### JSON Output

```go
import "encoding/json"

output, _ := json.MarshalIndent(&report, "", "  ")
fmt.Println(string(output))
```

## Using the CLI

```bash
# Render a report
sevaluation render report.json --format=terminal

# Check pass/fail (exit code 0 or 1)
sevaluation check report.json

# Validate JSON structure
sevaluation validate report.json
```

## Next Steps

- [Categorical Scoring](../concepts/scoring.md) - Understand pass/partial/fail
- [Rubrics](../features/rubrics.md) - Define evaluation criteria
- [Multi-Judge](../features/multi-judge.md) - Aggregate multiple evaluations
