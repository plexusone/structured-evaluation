# Multi-Judge Aggregation

Combine evaluations from multiple judges (LLMs or humans) to improve reliability and reduce individual bias.

## Overview

```go
type MultiJudgeResult struct {
    Evaluations         []EvaluationReport `json:"evaluations"`
    AggregationMethod   AggregationMethod  `json:"aggregation_method"`
    Agreement           float64            `json:"agreement"` // 0-1
    Disagreements       []Disagreement     `json:"disagreements,omitempty"`
    ConsolidatedReport  EvaluationReport   `json:"consolidated_report"`
    ConsolidatedDecision Decision          `json:"consolidated_decision"`
}
```

## Aggregation Methods

```go
const (
    AggregationMajority     AggregationMethod = "majority"     // Most common score wins
    AggregationConservative AggregationMethod = "conservative" // Lowest score wins
    AggregationOptimistic   AggregationMethod = "optimistic"   // Highest score wins
)
```

### Majority Voting

Uses the most common score across judges:

```go
// 3 judges: pass, pass, partial → pass
// 3 judges: pass, partial, fail → no clear majority, use tie-breaker
```

### Conservative

Takes the lowest (most pessimistic) score:

```go
// 3 judges: pass, pass, partial → partial
// Useful for security-critical evaluations
```

### Optimistic

Takes the highest score:

```go
// 3 judges: pass, partial, partial → pass
// Useful when any judge passing is sufficient
```

## Aggregating Evaluations

```go
evaluations := []evaluation.EvaluationReport{
    judge1Report,
    judge2Report,
    judge3Report,
}

result := evaluation.AggregateEvaluations(
    evaluations,
    evaluation.AggregationMajority,
)

// Access results
fmt.Printf("Agreement: %.1f%%\n", result.Agreement*100)
fmt.Printf("Decision: %s\n", result.ConsolidatedDecision.Status)

// Check disagreements
for _, d := range result.Disagreements {
    fmt.Printf("Disagreement on %s: %v\n", d.Category, d.Scores)
}
```

## Agreement Metrics

Agreement measures how often judges agree:

```go
type Agreement struct {
    Overall    float64            `json:"overall"`    // 0-1, across all categories
    ByCategory map[string]float64 `json:"by_category"` // Per-category agreement
}
```

### Computing Agreement

```go
// Perfect agreement: all judges give same score → 1.0
// No agreement: all different scores → 0.0
// Partial: some agree → 0.33-0.66

// Fleiss' Kappa is used for >2 judges
```

## Disagreements

When judges disagree significantly:

```go
type Disagreement struct {
    Category string       `json:"category"`
    Scores   []ScoreValue `json:"scores"`     // What each judge scored
    Spread   float64      `json:"spread"`     // How much they disagree
    Resolved ScoreValue   `json:"resolved"`   // Final aggregated score
}
```

### Handling Disagreements

```go
for _, d := range result.Disagreements {
    if d.Spread > 0.5 {
        // Major disagreement - may need human review
        fmt.Printf("⚠️ Major disagreement on %s\n", d.Category)
    }
}
```

## Finding Consolidation

Findings from multiple judges are merged:

```go
// Same finding from multiple judges → deduplicated
// Different severity → highest severity used
// Different details → combined
```

## Example: Three-Judge Panel

```go
// Run evaluation with 3 different LLMs
judges := []JudgeConfig{
    {Model: "claude-3-opus", Provider: "anthropic"},
    {Model: "gpt-4", Provider: "openai"},
    {Model: "gemini-pro", Provider: "google"},
}

var evaluations []evaluation.EvaluationReport
for _, judge := range judges {
    report := runEvaluation(document, rubric, judge)
    evaluations = append(evaluations, report)
}

// Aggregate with majority voting
result := evaluation.AggregateEvaluations(evaluations, evaluation.AggregationMajority)

// Use consolidated report
if result.Agreement < 0.6 {
    // Low agreement - flag for human review
    result.ConsolidatedDecision.Status = evaluation.DecisionHumanReview
}
```

## Best Practices

### Judge Selection

- Use diverse judges (different models, providers)
- Consider cost/latency tradeoffs
- Include at least 3 judges for majority voting

### When to Use Multi-Judge

- High-stakes evaluations (security, compliance)
- Calibration of new rubrics
- Contentious or subjective assessments
- Building evaluation datasets

### Handling Low Agreement

```go
if result.Agreement < 0.5 {
    // Options:
    // 1. Flag for human review
    // 2. Use conservative aggregation
    // 3. Add more judges
    // 4. Refine rubric criteria
}
```

### Cost Optimization

```go
// Start with single judge
report := runEvaluation(doc, rubric, primaryJudge)

// Only use multi-judge for borderline cases
if report.Decision.Status == evaluation.DecisionConditional {
    reports := runMultiJudge(doc, rubric, allJudges)
    result := AggregateEvaluations(reports, AggregationMajority)
}
```

## Next Steps

- [Rubrics](rubrics.md) - Define evaluation criteria
- [Pairwise Comparison](pairwise.md) - Compare outputs
- [DAG Aggregation](dag-aggregation.md) - Multi-agent workflows
