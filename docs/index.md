# Structured Evaluation

**Reusable evaluation framework for LLM-as-Judge and multi-agent workflows**

Structured Evaluation provides standardized Go types for evaluation reports, enabling consistent quality assessment across LLM-based and deterministic workflows.

## Features

- ⚖️ **LLM-as-Judge Assessments** - Categorical scoring (pass/partial/fail) with severity-based findings
- ✅ **GO/NO-GO Summary Reports** - Deterministic checks for CI, tests, and validation
- 📋 **Claims Validation** - Factual claim extraction and source verification
- 🔗 **Multi-Agent Coordination** - DAG-based report aggregation using Kahn's algorithm
- 📊 **Rubric Definitions** - Explicit criteria for consistent evaluations
- 🔄 **Pairwise Comparison** - Compare outputs instead of absolute scoring
- 👥 **Multi-Judge Aggregation** - Combine evaluations from multiple judges with agreement metrics
- 🔍 **Report Validation** - Validate reports for enum correctness, counts, and consistency
- 🟦 **TypeScript / Zod Bindings** - Generated Zod schemas and TS types, downstream of the same Go structs

## Architecture

```
┌───────────────────────────────────────────────────────────┐
│                    SummaryReport (GO/NO-GO)               │
│  ┌──────────────────────┐  ┌──────────────────────┐       │
│  │  Embedded Reports    │  │   Team Sections      │       │
│  │  (Full-Fidelity)     │  │   (Task Results)     │       │
│  └──────────────────────┘  └──────────────────────┘       │
└───────────────────────────────────────────────────────────┘
                              ▲
              ┌───────────────┴───────────────┐
              │                               │
┌─────────────┴─────────────┐   ┌─────────────┴─────────────┐
│     Rubric (rubric/)      │   │   ClaimsReport (claims/)  │
│  ┌─────────────────────┐  │   │  ┌─────────────────────┐  │
│  │ Category Results    │  │   │  │ Claims + Validation │  │
│  │ (pass/partial/fail) │  │   │  │ (verified/rejected) │  │
│  ├─────────────────────┤  │   │  ├─────────────────────┤  │
│  │ Findings            │  │   │  │ Sources             │  │
│  │ (severity-based)    │  │   │  │ (external/internal) │  │
│  └─────────────────────┘  │   │  └─────────────────────┘  │
│  LLM-as-Judge scoring     │   │  Fact verification        │
└───────────────────────────┘   └───────────────────────────┘
```

**Three complementary report types:**

| Package | Purpose | Evaluation Type |
|---------|---------|-----------------|
| `rubric/` | Categorical scoring with findings | Subjective (LLM-as-Judge) |
| `claims/` | Fact verification with sources | Objective (source-backed) |
| `summary/` | GO/NO-GO aggregation | Deterministic |

## Quick Example

```go
package main

import (
    "os"

    "github.com/plexusone/structured-evaluation/rubric"
    "github.com/plexusone/structured-evaluation/render/terminal"
)

func main() {
    report := rubric.NewRubric("prd", "document.md")

    // Add category results (pass/partial/fail)
    report.AddCategoryResult(rubric.CategoryResult{
        Category:  "problem_definition",
        Score:     rubric.ScorePass,
        Reasoning: "Clear problem statement with measurable goals",
    })
    report.AddCategoryResult(rubric.CategoryResult{
        Category:  "user_stories",
        Score:     rubric.ScorePartial,
        Reasoning: "Stories present but missing acceptance criteria",
    })

    // Add findings
    report.AddFinding(rubric.Finding{
        Severity:       rubric.SeverityMedium,
        Category:       "metrics",
        Title:          "Missing baseline metrics",
        Recommendation: "Add current baseline measurements",
    })

    report.Finalize(nil, "sevaluation check document.md")

    // Render to terminal
    renderer := terminal.New(os.Stdout)
    renderer.Render(report)
}
```

## Report Types

| Type | Purpose | Use Case |
|------|---------|----------|
| **Rubric** | LLM-as-Judge assessments | PRD reviews, code quality, content evaluation |
| **SummaryReport** | GO/NO-GO deterministic checks | CI pipelines, release validation, test results |
| **ClaimsReport** | Factual claim validation | Security advisories, blog posts, documentation |

## Severity Levels

Following InfoSec conventions:

| Severity | Icon | Blocking | Description |
|----------|------|----------|-------------|
| Critical | 🔴 | Yes | Must fix before approval |
| High | 🔴 | Yes | Must fix before approval |
| Medium | 🟡 | No | Should fix, tracked |
| Low | 🟢 | No | Nice to fix |
| Info | ⚪ | No | Informational only |

## TypeScript / JavaScript

Consumers that read reports in TypeScript (a UI, a dashboard, a Node service)
don't need to hand-maintain a parallel type — `@plexusone/structured-evaluation`
generates Zod schemas and TS types from the same JSON Schema the Go library
embeds:

```ts
import { RubricSchema, type Rubric } from '@plexusone/structured-evaluation'

const report: Rubric = RubricSchema.parse(JSON.parse(rawJson))
```

See [Installation](getting-started/installation.md#typescript-javascript) for
setup and the [package README](https://github.com/plexusone/structured-evaluation/tree/main/ts)
for details.

## Next Steps

- [Installation](getting-started/installation.md) - Get started with structured-evaluation
- [Quick Start](getting-started/quickstart.md) - Create your first evaluation report
- [Concepts](concepts/report-types.md) - Understand the evaluation model
- [CLI](cli/index.md) - Use the command-line tool
