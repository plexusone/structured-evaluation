# Inter-Rater Reliability (IRR)

Structured-evaluation provides metrics for comparing ratings between different evaluators (e.g., LLM vs human, or multiple humans).

## Use Cases

- **LLM Calibration**: Compare LLM ratings with human ground truth
- **Quality Assurance**: Verify LLM evaluations align with expert judgment
- **Research**: Measure agreement in evaluation studies
- **Continuous Improvement**: Track LLM alignment over time

## Available Metrics

### IRRMetrics

```go
type IRRMetrics struct {
    ExactAgreement         float64  // % exact matches
    AdjacentAgreement      float64  // % within ±1
    MeanAbsoluteDifference float64  // average |diff|
    PearsonCorrelation     float64  // -1 to 1
    SampleSize             int      // number of pairs
}
```

### CategoricalAgreement

```go
type CategoricalAgreement struct {
    ExactAgreement  float64         // % exact categorical matches
    ConfusionMatrix map[string]int  // disagreement patterns
    SampleSize      int
}
```

## Computing IRR

### From Rating Pairs

```go
pairs := []evaluation.RatingPair{
    {Rater1: 5, Rater2: 4, Category: "quality", ItemID: "doc1"},
    {Rater1: 3, Rater2: 3, Category: "quality", ItemID: "doc2"},
    {Rater1: 4, Rater2: 5, Category: "quality", ItemID: "doc3"},
}

metrics := evaluation.ComputeIRR(pairs)

fmt.Printf("Exact Agreement: %.1f%%\n", metrics.ExactAgreement*100)
fmt.Printf("Adjacent Agreement: %.1f%%\n", metrics.AdjacentAgreement*100)
fmt.Printf("Mean Absolute Diff: %.2f\n", metrics.MeanAbsoluteDifference)
fmt.Printf("Pearson r: %.3f\n", metrics.PearsonCorrelation)
```

### From Category Results

Compare two sets of evaluation results directly:

```go
// Human ratings
humanResults := []evaluation.CategoryResult{
    *evaluation.NewCategoryResultWithNumeric("quality", evaluation.ScorePass, 5.0, ""),
    *evaluation.NewCategoryResultWithNumeric("clarity", evaluation.ScorePartial, 3.0, ""),
}

// LLM ratings
llmResults := []evaluation.CategoryResult{
    *evaluation.NewCategoryResultWithNumeric("quality", evaluation.ScorePass, 4.0, ""),
    *evaluation.NewCategoryResultWithNumeric("clarity", evaluation.ScorePartial, 3.0, ""),
}

metrics := evaluation.ComputeIRRFromResults(humanResults, llmResults)
```

### Categorical Agreement

For pass/partial/fail comparisons:

```go
agreement := evaluation.ComputeCategoricalAgreement(humanResults, llmResults)

fmt.Printf("Exact Agreement: %.1f%%\n", agreement.ExactAgreement*100)

// Analyze disagreement patterns
for pattern, count := range agreement.ConfusionMatrix {
    fmt.Printf("  %s: %d\n", pattern, count)
}
// Output: pass:pass: 1, partial:partial: 1
```

## Interpreting Metrics

### Exact Agreement

| Value | Interpretation |
|-------|----------------|
| > 80% | Excellent agreement |
| 60-80% | Good agreement |
| 40-60% | Moderate agreement |
| < 40% | Poor agreement |

### Adjacent Agreement

| Value | Interpretation |
|-------|----------------|
| > 90% | Excellent (ratings within 1 point) |
| 70-90% | Acceptable |
| < 70% | Needs calibration |

### Pearson Correlation

| Value | Interpretation |
|-------|----------------|
| > 0.8 | Strong positive correlation |
| 0.5-0.8 | Moderate correlation |
| 0.3-0.5 | Weak correlation |
| < 0.3 | Little to no correlation |

## Automatic Score Conversion

When comparing categorical-only results, scores are converted:

| Categorical | Numeric |
|-------------|---------|
| Pass | 5.0 |
| Partial | 3.0 |
| Fail | 1.0 |

This allows IRR computation even without explicit numeric scores.

## Example: LLM Calibration Workflow

```go
// 1. Collect human ground truth
humanResults := runHumanEvaluation(documents)

// 2. Run LLM evaluation
llmResults := runLLMEvaluation(documents)

// 3. Compute IRR
metrics := evaluation.ComputeIRRFromResults(humanResults, llmResults)
catAgreement := evaluation.ComputeCategoricalAgreement(humanResults, llmResults)

// 4. Report
fmt.Println("=== LLM Calibration Report ===")
fmt.Printf("Sample Size: %d evaluations\n", metrics.SampleSize)
fmt.Printf("Exact Agreement: %.1f%%\n", metrics.ExactAgreement*100)
fmt.Printf("Adjacent Agreement: %.1f%%\n", metrics.AdjacentAgreement*100)
fmt.Printf("Correlation: %.3f\n", metrics.PearsonCorrelation)
fmt.Printf("Categorical Agreement: %.1f%%\n", catAgreement.ExactAgreement*100)

// 5. Identify systematic biases
for pattern, count := range catAgreement.ConfusionMatrix {
    if strings.Contains(pattern, ":") && pattern[:strings.Index(pattern, ":")] != pattern[strings.Index(pattern, ":")+1:] {
        fmt.Printf("Disagreement %s: %d cases\n", pattern, count)
    }
}
```

## Best Practices

1. **Collect sufficient samples** - Aim for 30+ pairs for reliable metrics
2. **Use numeric scores when possible** - More granular than categorical
3. **Track over time** - Monitor calibration drift
4. **Analyze confusion patterns** - Identify systematic biases (e.g., LLM always harsher)
5. **Iterate on prompts** - Use IRR feedback to improve LLM evaluation prompts
