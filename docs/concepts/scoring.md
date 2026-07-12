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
report.AddCategoryResult(rubric.CategoryResult{
    Category:  "problem_definition",
    Score:     rubric.ScorePass,
    Reasoning: "Problem is clearly stated with measurable business impact",
})

report.AddCategoryResult(rubric.CategoryResult{
    Category:  "user_stories",
    Score:     rubric.ScorePartial,
    Reasoning: "Stories present but 2 of 5 lack acceptance criteria",
})

report.AddCategoryResult(rubric.CategoryResult{
    Category:  "success_metrics",
    Score:     rubric.ScoreFail,
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
score := rubric.ScorePass

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

## Integer Scores (v0.9.0)

For LLM-as-Judge evaluations, use the 1-5 integer scale. Research shows LLMs are most reliable at this granularity.

### IntegerScore Type

```go
const (
    ScoreUnacceptable  IntegerScore = 1  // Does not meet requirements
    ScoreMajorRevisions IntegerScore = 2  // Significant work needed
    ScoreAcceptable    IntegerScore = 3  // Minimum requirements met
    ScoreGood          IntegerScore = 4  // Meets expectations well
    ScoreExcellent     IntegerScore = 5  // Exceeds expectations
)
```

### Using IntegerScore

```go
// Create category result with integer score
result := rubric.NewCategoryResultWithIntScore(
    "quality",
    rubric.ScoreGood,  // 4
    0.85,              // Confidence
    "Meets quality standards with minor improvements possible",
)

// Or set on existing result
result.SetIntScore(rubric.ScoreGood)
result.SetConfidence(0.85)
```

### Conversion Methods

```go
score := rubric.ScoreGood

score.String()        // "Good"
score.ToCategorical() // ScorePass (4-5 = pass, 3 = partial, 1-2 = fail)
score.IsValid()       // true

// Parse from int
score = rubric.ParseIntegerScore(4) // ScoreGood
```

### IntegerScore on Rubric

The overall rubric also has an IntegerScore:

```go
report := rubric.NewRubric("eval", "doc.md")
report.SetIntScore(rubric.ScoreGood)
report.SetConfidence(0.9)

// Or computed automatically from category scores
report.Evaluate(rubricSet)
// report.IntScore is now the weighted average of category IntScores
```

## Confidence (v0.9.0)

Track evaluator confidence for human review routing:

```go
result.SetConfidence(0.85)  // 0.0-1.0

// Check if human review needed
if result.HasLowConfidence() {  // Default threshold: 0.7
    // Route to human reviewer
}

// Custom threshold
if result.HasLowConfidence(0.8) {
    // Higher bar for confidence
}
```

### Overall Confidence

```go
report.SetConfidence(0.9)

// Or computed as minimum across categories
report.Evaluate(rubricSet)
// report.Confidence is the weakest category confidence

// Check if any category needs human review
if report.NeedsHumanReview() {
    // Route entire report to human
}
```

## Next Steps

- [Pass Criteria](pass-criteria.md) - Configure decision thresholds
- [Rubrics](../features/rubrics.md) - Define scoring criteria
- [Multi-Judge](../features/multi-judge.md) - Aggregate multiple evaluations
