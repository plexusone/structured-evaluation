# Structured Evaluation

**Reusable evaluation framework for LLM-as-Judge and multi-agent workflows**

Structured Evaluation provides standardized Go types for evaluation reports, enabling consistent quality assessment across LLM-based and deterministic workflows.

## Features

- ⚖️ **LLM-as-Judge Assessments** - Categorical scoring (pass/partial/fail) with severity-based findings
- ✅ **GO/NO-GO Summary Reports** - Deterministic checks for CI, tests, and validation
- 🔗 **Multi-Agent Coordination** - DAG-based report aggregation using Kahn's algorithm
- 📊 **Rubric Definitions** - Explicit criteria for consistent evaluations
- 🔄 **Pairwise Comparison** - Compare outputs instead of absolute scoring
- 👥 **Multi-Judge Aggregation** - Combine evaluations from multiple judges with agreement metrics

## Quick Example

```go
package main

import (
    "fmt"
    "os"

    "github.com/plexusone/structured-evaluation/evaluation"
    "github.com/plexusone/structured-evaluation/render/terminal"
)

func main() {
    report := evaluation.NewEvaluationReport("prd", "document.md")

    // Add category results (pass/partial/fail)
    report.AddCategory(evaluation.CategoryResult{
        Category:  "problem_definition",
        Score:     evaluation.ScorePass,
        Reasoning: "Clear problem statement with measurable goals",
    })
    report.AddCategory(evaluation.CategoryResult{
        Category:  "user_stories",
        Score:     evaluation.ScorePartial,
        Reasoning: "Stories present but missing acceptance criteria",
    })

    // Add findings
    report.AddFinding(evaluation.Finding{
        Severity:       evaluation.SeverityMedium,
        Category:       "metrics",
        Title:          "Missing baseline metrics",
        Recommendation: "Add current baseline measurements",
    })

    report.Finalize("sevaluation check document.md")

    // Render to terminal
    renderer := terminal.New(os.Stdout)
    renderer.Render(&report)
}
```

## Report Types

| Type | Purpose | Use Case |
|------|---------|----------|
| **EvaluationReport** | LLM-as-Judge assessments | PRD reviews, code quality, content evaluation |
| **SummaryReport** | GO/NO-GO deterministic checks | CI pipelines, release validation, test results |

## Severity Levels

Following InfoSec conventions:

| Severity | Icon | Blocking | Description |
|----------|------|----------|-------------|
| Critical | 🔴 | Yes | Must fix before approval |
| High | 🔴 | Yes | Must fix before approval |
| Medium | 🟡 | No | Should fix, tracked |
| Low | 🟢 | No | Nice to fix |
| Info | ⚪ | No | Informational only |

## Next Steps

- [Installation](getting-started/installation.md) - Get started with structured-evaluation
- [Quick Start](getting-started/quickstart.md) - Create your first evaluation report
- [Concepts](concepts/report-types.md) - Understand the evaluation model
- [CLI](cli/index.md) - Use the command-line tool
