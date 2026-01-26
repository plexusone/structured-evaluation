# Structured Evaluation

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

A reusable evaluation framework for LLM-as-Judge and multi-agent workflows.

## Overview

`structured-evaluation` provides standardized types for evaluation reports, enabling:

- **LLM-as-Judge assessments** with weighted category scores and severity-based findings
- **GO/NO-GO summary reports** for deterministic checks (CI, tests, validation)
- **Multi-agent coordination** with DAG-based report aggregation

## Installation

```bash
go get github.com/agentplexus/structured-evaluation
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
import "github.com/agentplexus/structured-evaluation/evaluation"

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
import "github.com/agentplexus/structured-evaluation/summary"

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
go install github.com/agentplexus/structured-evaluation/cmd/sevaluation@latest

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
import "github.com/agentplexus/structured-evaluation/combine"

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
import "github.com/agentplexus/structured-evaluation/schema"

evalSchema := schema.EvaluationSchemaJSON
summarySchema := schema.SummarySchemaJSON
```

## Integration

Designed to work with:

- `github.com/grokify/structured-requirements` - PRD evaluation templates
- `github.com/agentplexus/multi-agent-spec` - Agent coordination
- `github.com/grokify/structured-changelog` - Release validation

## License

MIT License - see [LICENSE](LICENSE) for details.

 [build-status-svg]: https://github.com/agentplexus/structured-evaluation/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/agentplexus/structured-evaluation/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/agentplexus/structured-evaluation/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/agentplexus/structured-evaluation/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/agentplexus/structured-evaluation
 [goreport-url]: https://goreportcard.com/report/github.com/agentplexus/structured-evaluation
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/agentplexus/structured-evaluation
 [docs-godoc-url]: https://pkg.go.dev/github.com/agentplexus/structured-evaluation
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/agentplexus/structured-evaluation/blob/master/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/agentplexus/structured-evaluation/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/agentplexus/structured-evaluation?badge
