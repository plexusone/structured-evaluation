# Structured Evaluation

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/structured-evaluation
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/structured-evaluation
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/structured-evaluation
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/structured-evaluation
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fstructured-evaluation
 [loc-svg]: https://tokei.rs/b1/github/plexusone/structured-evaluation
 [repo-url]: https://github.com/plexusone/structured-evaluation
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/structured-evaluation/blob/master/LICENSE

A reusable evaluation framework for LLM-as-Judge and multi-agent workflows.

## Overview

`structured-evaluation` provides standardized types for evaluation reports, enabling:

- **LLM-as-Judge assessments** with weighted category scores and severity-based findings
- **GO/NO-GO summary reports** for deterministic checks (CI, tests, validation)
- **Multi-agent coordination** with DAG-based report aggregation

## Installation

```bash
go get github.com/plexusone/structured-evaluation
```

## Packages

| Package | Description |
|---------|-------------|
| `evaluation` | EvaluationReport, CategoryScore, Finding, Severity types |
| `summary` | SummaryReport, TeamSection, TaskResult for GO/NO-GO checks |
| `combine` | DAG-based report aggregation using Kahn's algorithm |
| `render/box` | Box-format terminal renderer for summary reports |
| `render/detailed` | Detailed terminal renderer for evaluation reports |
| `schema` | JSON Schema generation and embedding |

## Report Types

### Evaluation Report (LLM-as-Judge)

For subjective quality assessments with detailed findings:

```go
import "github.com/plexusone/structured-evaluation/evaluation"

report := evaluation.NewEvaluationReport("prd", "document.md")
report.AddCategory(evaluation.NewCategoryScore("problem_definition", 0.20, 8.5, "Clear problem statement"))
report.AddFinding(evaluation.Finding{
    Severity:       evaluation.SeverityMedium,
    Category:       "metrics",
    Title:          "Missing baseline metrics",
    Recommendation: "Add current baseline measurements",
})
report.Finalize("sevaluation check document.md")
```

### Summary Report (GO/NO-GO)

For deterministic checks with pass/fail status:

```go
import "github.com/plexusone/structured-evaluation/summary"

report := summary.NewSummaryReport("my-service", "v1.0.0", "Release Validation")
report.AddTeam(summary.TeamSection{
    ID:   "qa",
    Name: "Quality Assurance",
    Tasks: []summary.TaskResult{
        {ID: "unit-tests", Status: summary.StatusGo, Detail: "Coverage: 92%"},
        {ID: "e2e-tests", Status: summary.StatusWarn, Detail: "2 flaky tests"},
    },
})
```

## Severity Levels

Following InfoSec conventions:

| Severity | Icon | Blocking | Description |
|----------|------|----------|-------------|
| Critical | 🔴 | Yes | Must fix before approval |
| High | 🔴 | Yes | Must fix before approval |
| Medium | 🟡 | No | Should fix, tracked |
| Low | 🟢 | No | Nice to fix |
| Info | ⚪ | No | Informational only |

## Pass Criteria

Default criteria (zero blocking findings, minimum score):

```go
criteria := evaluation.DefaultPassCriteria()
// MaxCritical: 0, MaxHigh: 0, MaxMedium: -1 (unlimited), MinScore: 7.0

criteria := evaluation.StrictPassCriteria()
// MaxCritical: 0, MaxHigh: 0, MaxMedium: 3, MinScore: 8.0
```

## CLI Tool

```bash
# Install
go install github.com/plexusone/structured-evaluation/cmd/sevaluation@latest

# Render reports
sevaluation render report.json --format=detailed
sevaluation render report.json --format=box
sevaluation render report.json --format=json

# Check pass/fail (exit code 0/1)
sevaluation check report.json

# Validate structure
sevaluation validate report.json

# Generate JSON Schema
sevaluation schema generate -o ./schema/
```

## DAG-Based Aggregation

For multi-agent workflows with dependencies:

```go
import "github.com/plexusone/structured-evaluation/combine"

results := []combine.AgentResult{
    {TeamID: "qa", Tasks: qaTasks},
    {TeamID: "security", Tasks: secTasks, DependsOn: []string{"qa"}},
    {TeamID: "release", Tasks: relTasks, DependsOn: []string{"qa", "security"}},
}

report := combine.AggregateResults(results, "my-project", "v1.0.0", "Release")
// Teams are topologically sorted: qa → security → release
```

## JSON Schema

Schemas are embedded for runtime validation:

```go
import "github.com/plexusone/structured-evaluation/schema"

evalSchema := schema.EvaluationSchemaJSON
summarySchema := schema.SummarySchemaJSON
```

## Rubrics (v0.2.0)

Define explicit scoring criteria for consistent evaluations:

```go
rubric := evaluation.NewRubric("quality", "Output quality").
    AddRangeAnchor(8, 10, "Excellent", "Near perfect").
    AddRangeAnchor(5, 7.9, "Good", "Acceptable").
    AddRangeAnchor(0, 4.9, "Poor", "Needs work")

// Use default PRD rubric
rubricSet := evaluation.DefaultPRDRubricSet()
```

## Judge Metadata (v0.2.0)

Track LLM judge configuration for reproducibility:

```go
judge := evaluation.NewJudgeMetadata("claude-3-opus").
    WithProvider("anthropic").
    WithPrompt("prd-eval-v1", "1.0").
    WithTemperature(0.0).
    WithTokenUsage(1500, 800)

report.SetJudge(judge)
```

## Pairwise Comparison (v0.2.0)

Compare two outputs instead of absolute scoring:

```go
comparison := evaluation.NewPairwiseComparison(input, outputA, outputB)
comparison.SetWinner(evaluation.WinnerA, "A is more accurate", 0.9)

// Aggregate multiple comparisons
result := evaluation.ComputePairwiseResult(comparisons)
// result.WinRateA, result.OverallWinner
```

## Multi-Judge Aggregation (v0.2.0)

Combine evaluations from multiple judges:

```go
result := evaluation.AggregateEvaluations(evaluations, evaluation.AggregationMean)

// Methods: AggregationMean, AggregationMedian, AggregationConservative, AggregationMajority
// result.Agreement - inter-judge agreement (0-1)
// result.Disagreements - categories with significant disagreement
// result.ConsolidatedDecision - final aggregated decision
```

## OmniObserve Integration

Export evaluations to Opik, Phoenix, or Langfuse:

```go
import "github.com/plexusone/omniobserve/integrations/sevaluation"

// Export to observability platform
err := sevaluation.Export(ctx, provider, traceID, report)
```

## Integration

Designed to work with:

- `github.com/plexusone/omniobserve` - LLM observability (Opik, Phoenix, Langfuse)
- `github.com/grokify/structured-requirements` - PRD evaluation templates
- `github.com/plexusone/multi-agent-spec` - Agent coordination
- `github.com/grokify/structured-changelog` - Release validation

## License

MIT License - see [LICENSE](LICENSE) for details.
