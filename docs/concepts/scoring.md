# Categorical Scoring

As of v0.4.0, structured-evaluation uses **categorical scoring** instead of numeric scores. This aligns with how LLM judges naturally assess quality.

## Score Values

```go
const (
    ScorePass    ScoreValue = "pass"    // Meets requirements
    ScorePartial ScoreValue = "partial" // Partially meets requirements
    ScoreFail    ScoreValue = "fail"    // Does not meet requirements
)
```

## Why Categorical?

### Benefits over Numeric Scores

1. **Clearer semantics** - "pass" is unambiguous; "7.5" requires interpretation
2. **Better LLM alignment** - LLMs naturally reason in categories
3. **Simpler aggregation** - Majority voting vs. weighted averages
4. **Reduced bias** - No artificial precision (7.2 vs 7.3)

### Comparison

| Numeric | Categorical | Interpretation |
|---------|-------------|----------------|
| 8.0-10.0 | `pass` | Meets all requirements |
| 5.0-7.9 | `partial` | Meets most, minor issues |
| 0.0-4.9 | `fail` | Major gaps or issues |

## CategoryResult

Each evaluation category produces a result:

```go
type CategoryResult struct {
    Category  string     `json:"category"`  // Category ID
    Score     ScoreValue `json:"score"`     // pass/partial/fail
    Reasoning string     `json:"reasoning"` // Explanation
}
```

### Example

```go
report.AddCategory(evaluation.CategoryResult{
    Category:  "problem_definition",
    Score:     evaluation.ScorePass,
    Reasoning: "Problem is clearly stated with measurable business impact",
})

report.AddCategory(evaluation.CategoryResult{
    Category:  "user_stories",
    Score:     evaluation.ScorePartial,
    Reasoning: "Stories present but 2 of 5 lack acceptance criteria",
})

report.AddCategory(evaluation.CategoryResult{
    Category:  "success_metrics",
    Score:     evaluation.ScoreFail,
    Reasoning: "No quantitative success metrics defined",
})
```

## CategoryCounts

The decision includes category counts for quick assessment:

```go
type CategoryCounts struct {
    Pass    int `json:"pass"`
    Partial int `json:"partial"`
    Fail    int `json:"fail"`
    Total   int `json:"total"`
}
```

### Usage

```go
counts := report.Decision.CategoryCounts
fmt.Printf("Results: %d pass, %d partial, %d fail (of %d)\n",
    counts.Pass, counts.Partial, counts.Fail, counts.Total)
```

## Score Methods

```go
score := evaluation.ScorePass

score.IsPassing()  // true
score.IsPartial()  // false
score.IsFailing()  // false
score.Icon()       // "🟢"
```

| Score | Icon | IsPassing | IsPartial | IsFailing |
|-------|------|-----------|-----------|-----------|
| `pass` | 🟢 | true | false | false |
| `partial` | 🟡 | false | true | false |
| `fail` | 🔴 | false | false | true |

## Decision Logic

The overall decision is computed from category results:

```go
// All pass → DecisionPass
// Any fail with blocking findings → DecisionFail
// Mix of pass/partial → DecisionConditional
// Uncertain → DecisionHumanReview
```

## Migration from Numeric

If migrating from v0.3.x or earlier:

| Old API | New API |
|---------|---------|
| `CategoryScore` | `CategoryResult` |
| `Score float64` | `Score ScoreValue` |
| `MaxScore float64` | (removed) |
| `Status ScoreStatus` | (merged into Score) |
| `Justification string` | `Reasoning string` |
| `WeightedScore float64` | (removed from report) |

## Next Steps

- [Pass Criteria](pass-criteria.md) - Configure decision thresholds
- [Rubrics](../features/rubrics.md) - Define scoring criteria
- [Multi-Judge](../features/multi-judge.md) - Aggregate multiple evaluations
