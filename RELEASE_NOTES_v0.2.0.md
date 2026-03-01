# Release Notes - v0.2.0

**Release Date:** 2026-01-26

## Overview

v0.2.0 adds LLM-as-Judge best practices including rubric definitions, judge metadata tracking, pairwise comparison mode, and multi-judge aggregation.

## Highlights

- **Rubric Definitions** with score anchors for consistent evaluations
- **Judge Metadata** tracking (model, prompt, temperature, tokens, observability links)
- **Pairwise Comparison** mode as alternative to absolute scoring
- **Multi-Judge Aggregation** with agreement metrics and disagreement detection

## New Types

### Rubric & ScoreAnchor

Define explicit scoring criteria for each evaluation category:

```go
rubric := evaluation.NewRubric("quality", "Output quality").
    AddRangeAnchor(8, 10, "Excellent", "Near perfect output").
    AddRangeAnchor(5, 7.9, "Good", "Acceptable with minor issues").
    AddRangeAnchor(0, 4.9, "Poor", "Significant problems")

// Get anchor for a score
anchor := rubric.GetAnchorForScore(8.5) // Returns "Excellent" anchor

// Use default PRD rubric set
rubricSet := evaluation.DefaultPRDRubricSet()
```

### JudgeMetadata

Track LLM judge configuration for reproducibility:

```go
judge := evaluation.NewJudgeMetadata("claude-3-opus-20240229").
    WithProvider("anthropic").
    WithPrompt("prd-eval-v1", "1.0").
    WithTemperature(0.0).
    WithRubric("prd-evaluation-v1", "1.0").
    WithTokenUsage(1500, 800)

report.SetJudge(judge)
```

### PairwiseComparison

Compare two outputs instead of absolute scoring:

```go
comparison := evaluation.NewPairwiseComparison(input, outputA, outputB)
comparison.SetWinner(evaluation.WinnerA, "A is more accurate", 0.9)
comparison.AddCategoryScore("clarity", evaluation.WinnerB, "B is clearer", 0.3)

// Detect position bias with swapped comparison
swapped := comparison.SwappedComparison()

// Aggregate multiple comparisons
result := evaluation.ComputePairwiseResult(comparisons)
// result.WinRateA, result.WinRateB, result.OverallWinner
```

### ReferenceData

Store gold/expected data for reference-based evaluation:

```go
ref := evaluation.NewReferenceData(input, expectedOutput).
    WithContext("document1.txt", "document2.txt").
    WithAnnotation("quality", 9.0, "human-reviewer-1")

report.SetReference(ref)
```

### MultiJudgeResult

Aggregate evaluations from multiple judges:

```go
result := evaluation.AggregateEvaluations(evaluations, evaluation.AggregationMean)

// Available methods:
// - AggregationMean: arithmetic mean
// - AggregationMedian: median score
// - AggregationConservative: minimum score (most critical)
// - AggregationMajority: majority vote on pass/fail

// Check agreement
if result.Agreement < 0.7 {
    log.Println("Low inter-judge agreement:", result.Disagreements)
}

// Use consolidated results
finalDecision := result.ConsolidatedDecision
findings := result.ConsolidatedFindings
```

## Updated EvaluationReport

New fields added to `EvaluationReport`:

```go
type EvaluationReport struct {
    // ... existing fields ...

    // v0.2.0 additions
    Judge     *JudgeMetadata `json:"judge,omitempty"`
    RubricID  string         `json:"rubric_id,omitempty"`
    Reference *ReferenceData `json:"reference,omitempty"`
}
```

## OmniObserve Integration

The `omniobserve` package now includes `integrations/sevaluation` for exporting evaluation reports to Opik, Phoenix, and Langfuse:

```go
import "github.com/plexusone/omniobserve/integrations/sevaluation"

// Export to observability platform
err := sevaluation.Export(ctx, provider, traceID, report)

// Import platform results
report := sevaluation.ImportEvalResult(evalResult)
```

## Aggregation Methods

| Method | Description | Use Case |
|--------|-------------|----------|
| `mean` | Arithmetic mean of scores | General aggregation |
| `median` | Median score | Robust to outliers |
| `conservative` | Minimum score | Safety-critical evaluations |
| `majority` | Majority vote on pass/fail | Binary decisions |
| `weighted` | Weighted by judge confidence | Variable judge quality |

## Migration from v0.1.0

v0.2.0 is fully backward compatible. Existing code works without changes. New fields are optional:

```go
// v0.1.0 code still works
report := evaluation.NewEvaluationReport("prd", "doc.md")
report.AddCategory(...)
report.Finalize("sevaluation check doc.md")

// Optionally add v0.2.0 features
report.SetJudge(judge)
report.SetRubric("prd-evaluation-v1")
```

## Installation

```bash
go get github.com/plexusone/structured-evaluation@v0.2.0
```

## Contributors

- PlexusOne Team
- Claude Opus 4.5 (Co-Author)
