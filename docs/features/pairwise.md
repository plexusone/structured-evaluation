# Pairwise Comparison

Pairwise comparison evaluates two outputs against each other rather than absolute scoring. This is useful when relative quality is more important than absolute scores.

## Overview

```go
type PairwiseComparison struct {
    ID           string        `json:"id"`
    Input        string        `json:"input"`
    OutputA      string        `json:"output_a"`
    OutputB      string        `json:"output_b"`
    Winner       Winner        `json:"winner"`
    Reasoning    string        `json:"reasoning"`
    Confidence   float64       `json:"confidence"`
    SwapPosition bool          `json:"swap_position"` // For bias detection
    Judge        *JudgeMetadata `json:"judge,omitempty"`
}

const (
    WinnerA   Winner = "a"
    WinnerB   Winner = "b"
    WinnerTie Winner = "tie"
)
```

## Creating a Comparison

```go
comparison := rubric.NewPairwiseComparison(
    "Summarize this article about climate change",
    "Output A: Climate change is causing global temperatures...",
    "Output B: The article discusses how climate change...",
)

// Set the winner after evaluation
comparison.SetWinner(
    rubric.WinnerA,
    "Output A provides more specific details and better structure",
    0.85, // Confidence 0-1
)
```

## Position Bias Detection

LLMs can exhibit position bias (preferring the first or second option). Detect this by running comparisons with swapped positions:

```go
// Original comparison
comp1 := rubric.NewPairwiseComparison(input, outputA, outputB)
comp1.SetWinner(rubric.WinnerA, "...", 0.9)

// Swapped comparison
comp2 := rubric.NewPairwiseComparison(input, outputB, outputA)
comp2.SwapPosition = true
comp2.SetWinner(rubric.WinnerB, "...", 0.85) // Should also prefer original A

// Check for consistency
if comp1.Winner == "a" && comp2.Winner == "b" {
    // Consistent: both prefer the same output
} else {
    // Position bias detected
}
```

## Aggregating Comparisons

Compute overall results from multiple comparisons:

```go
comparisons := []rubric.PairwiseComparison{
    comp1, comp2, comp3, // Multiple comparisons
}

result := rubric.ComputePairwiseResult(comparisons)

fmt.Printf("Win rate A: %.1f%%\n", result.WinRateA*100)
fmt.Printf("Win rate B: %.1f%%\n", result.WinRateB*100)
fmt.Printf("Tie rate: %.1f%%\n", result.TieRate*100)
fmt.Printf("Overall winner: %s\n", result.OverallWinner)
```

### PairwiseResult Structure

```go
type PairwiseResult struct {
    TotalComparisons int     `json:"total_comparisons"`
    WinsA            int     `json:"wins_a"`
    WinsB            int     `json:"wins_b"`
    Ties             int     `json:"ties"`
    WinRateA         float64 `json:"win_rate_a"`
    WinRateB         float64 `json:"win_rate_b"`
    TieRate          float64 `json:"tie_rate"`
    OverallWinner    Winner  `json:"overall_winner"`
    PositionBias     float64 `json:"position_bias"` // 0 = no bias
}
```

## Use Cases

### Model Comparison

Compare outputs from different models:

```go
// Compare GPT-4 vs Claude responses
comp := rubric.NewPairwiseComparison(
    "Explain quantum computing",
    gpt4Response,
    claudeResponse,
)
```

### A/B Testing Prompts

Compare different prompt strategies:

```go
// Compare concise vs detailed prompts
comp := rubric.NewPairwiseComparison(
    userQuery,
    responseFromConcisePrompt,
    responseFromDetailedPrompt,
)
```

### Human Preference Alignment

Collect human preferences for RLHF:

```go
// Store human preference
comp := rubric.NewPairwiseComparison(instruction, outputA, outputB)
comp.SetWinner(rubric.WinnerB, "Preferred by human annotator", 1.0)
comp.Judge = &rubric.JudgeMetadata{
    Model:    "human",
    Provider: "internal-annotation",
}
```

## Best Practices

### Randomize Position

Always randomize which output appears first to detect bias:

```go
if rand.Float64() < 0.5 {
    comp = NewPairwiseComparison(input, outputA, outputB)
} else {
    comp = NewPairwiseComparison(input, outputB, outputA)
    comp.SwapPosition = true
}
```

### Use Multiple Judges

Reduce individual judge bias:

```go
results := []PairwiseComparison{}

for _, judge := range judges {
    comp := runComparison(input, outputA, outputB, judge)
    results = append(results, comp)
}

aggregate := ComputePairwiseResult(results)
```

### Track Confidence

Use confidence scores to weight comparisons:

```go
comp.SetWinner(WinnerA, reasoning, 0.95) // High confidence
comp.SetWinner(WinnerTie, reasoning, 0.55) // Low confidence, essentially a tie
```

## Next Steps

- [Multi-Judge Aggregation](multi-judge.md) - Combine multiple judges
- [Rubrics](rubrics.md) - Define evaluation criteria
