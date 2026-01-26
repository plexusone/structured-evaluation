# Release Notes - v0.1.0

**Release Date:** 2026-01-26

## Overview

Initial release of `structured-evaluation`, a reusable evaluation framework for LLM-as-Judge and multi-agent workflows.

## Highlights

- **Evaluation Report type** for detailed LLM-as-Judge assessments with weighted category scores
- **Summary Report type** for GO/NO-GO deterministic checks
- **InfoSec severity levels** (Critical, High, Medium, Low, Info) with blocking thresholds
- **DAG-based aggregation** for multi-agent coordination using topological sort
- **Terminal renderers** with box-format and detailed output
- **CLI tool** (`sevaluation`) for rendering, validation, and pass/fail checks

## Packages

### evaluation/

Core types for LLM-as-Judge evaluations:

- `EvaluationReport` - Main report type with metadata, categories, findings, decision
- `CategoryScore` - Weighted scores (0-10) with pass/warn/fail status
- `Finding` - Issues with severity, recommendations, and ownership
- `Severity` - InfoSec levels with blocking rules
- `PassCriteria` - Configurable approval thresholds
- `Decision` - Pass/fail/conditional/human_review outcomes

### summary/

Types for deterministic GO/NO-GO checks:

- `SummaryReport` - Aggregated team results
- `TeamSection` - Agent/team outputs with dependencies
- `TaskResult` - Individual check outcomes
- `Status` - GO/WARN/NO-GO/SKIP with emoji icons

### combine/

Multi-agent coordination:

- `SortByDAG()` - Topological sort using Kahn's algorithm
- `AggregateResults()` - Combine agent outputs
- `AggregateWithDAG()` - Combine with explicit dependencies

### render/

Terminal output:

- `box.Renderer` - Box-format for summary reports
- `detailed.TerminalRenderer` - Detailed format for evaluation reports

### schema/

JSON Schema support:

- `GenerateEvaluationSchema()` - Generate from Go types
- `GenerateSummarySchema()` - Generate from Go types
- Embedded schemas via `//go:embed`

## CLI Commands

```bash
sevaluation render <file.json> --format=box|detailed|json
sevaluation check <file.json>      # Exit 0=pass, 1=fail
sevaluation validate <file.json>
sevaluation schema generate -o <dir>
```

## Pass Criteria

Default criteria for approval:

| Threshold | Value |
|-----------|-------|
| Max Critical | 0 |
| Max High | 0 |
| Max Medium | Unlimited |
| Min Score | 7.0 |

## Future Plans (v0.2.0)

Based on LLM-as-Judge best practices, planned additions:

- Rubric definitions with score anchors
- Judge model and prompt tracking
- Pairwise comparison mode
- Reference/gold data fields
- Multi-judge aggregation

## Installation

```bash
go get github.com/agentplexus/structured-evaluation@v0.1.0
```

## Contributors

- AgentPlexus Team
- Claude Opus 4.5 (Co-Author)
