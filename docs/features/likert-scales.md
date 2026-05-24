# Likert Scales

Structured-evaluation supports Likert scales (1-5 numeric ratings) alongside categorical scoring. This enables human comparison studies and inter-rater reliability analysis.

## When to Use Likert Scales

| Use Case | Recommended Scale |
|----------|-------------------|
| LLM-as-Judge decisions | Categorical (pass/partial/fail) |
| Human comparison studies | Likert (1-5) |
| Inter-rater reliability | Likert (1-5) |
| Calibration analysis | Likert (1-5) |
| Simple automation | Categorical |

## Creating Likert Categories

```go
import "github.com/plexusone/structured-evaluation/evaluation"

// Using standard 1-5 anchors
cat := evaluation.NewCategory("quality", "Content Quality", "Overall quality assessment").
    WithLikert5(evaluation.StandardLikert5Anchors())

// Custom anchors
cat := evaluation.NewCategory("clarity", "Clarity", "How clear is the writing").
    WithLikert5([]evaluation.LikertAnchor{
        {Value: 5, Label: "Crystal Clear", Description: "No ambiguity, easy to understand"},
        {Value: 4, Label: "Clear", Description: "Minor clarifications needed"},
        {Value: 3, Label: "Adequate", Description: "Understandable with effort"},
        {Value: 2, Label: "Unclear", Description: "Significant confusion"},
        {Value: 1, Label: "Incomprehensible", Description: "Cannot understand"},
    })
```

## Standard Anchors

The `StandardLikert5Anchors()` helper provides:

| Score | Label | Description |
|-------|-------|-------------|
| 5 | Excellent | Exceeds all expectations |
| 4 | Good | Meets expectations with minor improvements possible |
| 3 | Adequate | Meets minimum requirements |
| 2 | Needs Improvement | Below expectations |
| 1 | Poor | Does not meet requirements |

## Automatic Categorical Mapping

Likert scores are automatically mapped to categorical for decisions:

| Likert Score | Categorical |
|--------------|-------------|
| 4-5 | Pass |
| 3 | Partial |
| 1-2 | Fail |

Thresholds are configurable:

```go
passThreshold := 4
partialThreshold := 3
config := &evaluation.LikertConfig{
    Min:              1,
    Max:              5,
    PassThreshold:    &passThreshold,
    PartialThreshold: &partialThreshold,
}
cat.WithLikert(config)
```

## Recording Results

### From Likert Score

```go
// Categorical score is derived automatically
result := evaluation.NewCategoryResultFromLikert(
    "quality",   // category ID
    4,           // Likert score
    config,      // LikertConfig
    "Good overall quality with minor issues",
)
// result.Score = ScorePass
// result.NumericScore = 4.0
```

### Dual Scores

```go
// Record both categorical and numeric
result := evaluation.NewCategoryResultWithNumeric(
    "quality",
    evaluation.ScorePass,
    4.5,  // numeric for human comparison
    "Reasoning here",
)
```

### Adding Numeric to Existing

```go
result := evaluation.NewCategoryResult("quality", evaluation.ScorePass, "Good").
    SetNumericScore(4.5)
```

## Accessing Numeric Scores

```go
if result.HasNumericScore() {
    score := result.GetNumericScore()
    fmt.Printf("Numeric: %.1f\n", score)
}
```

## Validation

Rubric validation checks Likert configurations:

```go
rs := evaluation.NewRubricSet("test", "Test", "1.0")
rs.AddCategory(*cat)

issues := rs.Validate()
// Checks:
// - LikertConfig is present for likert scale type
// - Min < Max
```

## Best Practices

1. **Use categorical for automation** - Pass/partial/fail is cleaner for decision-making
2. **Use Likert for calibration** - Compare LLM ratings with human ground truth
3. **Include anchors** - Detailed anchor descriptions improve rater consistency
4. **Store both when needed** - Use `NumericScore` field for analysis while keeping categorical for decisions
