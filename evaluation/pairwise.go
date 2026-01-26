package evaluation

import "time"

// PairwiseComparison represents a comparison between two outputs.
// This is an alternative to absolute scoring that can reduce position bias
// and improve reliability of LLM-as-Judge evaluations.
type PairwiseComparison struct {
	// ID is the unique identifier for this comparison.
	ID string `json:"id,omitempty"`

	// Input is the shared input/prompt for both outputs.
	Input string `json:"input"`

	// OutputA is the first output being compared.
	OutputA string `json:"output_a"`

	// OutputB is the second output being compared.
	OutputB string `json:"output_b"`

	// Winner indicates which output won ("A", "B", or "tie").
	Winner PairwiseWinner `json:"winner"`

	// Confidence is the judge's confidence in the decision (0-1).
	Confidence float64 `json:"confidence,omitempty"`

	// Reasoning explains why this winner was chosen.
	Reasoning string `json:"reasoning"`

	// CategoryScores provides per-category comparisons if applicable.
	CategoryScores []PairwiseCategoryScore `json:"category_scores,omitempty"`

	// Judge contains metadata about the LLM judge.
	Judge *JudgeMetadata `json:"judge,omitempty"`

	// Metadata contains additional comparison context.
	Metadata map[string]any `json:"metadata,omitempty"`

	// CreatedAt is when this comparison was made.
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// PairwiseWinner indicates the winner of a pairwise comparison.
type PairwiseWinner string

const (
	// WinnerA indicates output A is better.
	WinnerA PairwiseWinner = "A"

	// WinnerB indicates output B is better.
	WinnerB PairwiseWinner = "B"

	// WinnerTie indicates both outputs are roughly equal.
	WinnerTie PairwiseWinner = "tie"

	// WinnerUncertain indicates the judge couldn't determine a winner.
	WinnerUncertain PairwiseWinner = "uncertain"
)

// PairwiseCategoryScore compares outputs on a specific dimension.
type PairwiseCategoryScore struct {
	// Category is the evaluation dimension.
	Category string `json:"category"`

	// Winner indicates which output won for this category.
	Winner PairwiseWinner `json:"winner"`

	// Margin indicates how much better the winner is (0-1, higher = larger gap).
	Margin float64 `json:"margin,omitempty"`

	// Reasoning explains the category-level comparison.
	Reasoning string `json:"reasoning,omitempty"`
}

// PairwiseResult aggregates multiple pairwise comparisons.
type PairwiseResult struct {
	// Comparisons are all the individual comparisons.
	Comparisons []PairwiseComparison `json:"comparisons"`

	// WinRateA is the percentage of comparisons won by A.
	WinRateA float64 `json:"win_rate_a"`

	// WinRateB is the percentage of comparisons won by B.
	WinRateB float64 `json:"win_rate_b"`

	// TieRate is the percentage of ties.
	TieRate float64 `json:"tie_rate"`

	// OverallWinner is the aggregated winner.
	OverallWinner PairwiseWinner `json:"overall_winner"`

	// Confidence is the overall confidence in the result.
	Confidence float64 `json:"confidence"`
}

// NewPairwiseComparison creates a new pairwise comparison.
func NewPairwiseComparison(input, outputA, outputB string) *PairwiseComparison {
	return &PairwiseComparison{
		Input:     input,
		OutputA:   outputA,
		OutputB:   outputB,
		CreatedAt: time.Now().UTC(),
	}
}

// SetWinner sets the comparison result.
func (p *PairwiseComparison) SetWinner(winner PairwiseWinner, reasoning string, confidence float64) {
	p.Winner = winner
	p.Reasoning = reasoning
	p.Confidence = confidence
}

// AddCategoryScore adds a per-category comparison.
func (p *PairwiseComparison) AddCategoryScore(category string, winner PairwiseWinner, reasoning string, margin float64) {
	p.CategoryScores = append(p.CategoryScores, PairwiseCategoryScore{
		Category:  category,
		Winner:    winner,
		Margin:    margin,
		Reasoning: reasoning,
	})
}

// ComputeResult aggregates multiple comparisons into a result.
func ComputePairwiseResult(comparisons []PairwiseComparison) *PairwiseResult {
	if len(comparisons) == 0 {
		return &PairwiseResult{}
	}

	var winsA, winsB, ties int
	for _, c := range comparisons {
		switch c.Winner {
		case WinnerA:
			winsA++
		case WinnerB:
			winsB++
		case WinnerTie:
			ties++
		}
	}

	total := float64(len(comparisons))
	result := &PairwiseResult{
		Comparisons: comparisons,
		WinRateA:    float64(winsA) / total,
		WinRateB:    float64(winsB) / total,
		TieRate:     float64(ties) / total,
	}

	// Determine overall winner
	switch {
	case winsA > winsB:
		result.OverallWinner = WinnerA
		result.Confidence = float64(winsA-winsB) / total
	case winsB > winsA:
		result.OverallWinner = WinnerB
		result.Confidence = float64(winsB-winsA) / total
	default:
		result.OverallWinner = WinnerTie
		result.Confidence = 0
	}

	return result
}

// SwappedComparison creates a comparison with A and B swapped.
// Running both orders helps detect position bias in the judge.
func (p *PairwiseComparison) SwappedComparison() *PairwiseComparison {
	return &PairwiseComparison{
		Input:     p.Input,
		OutputA:   p.OutputB,
		OutputB:   p.OutputA,
		CreatedAt: time.Now().UTC(),
		Metadata: map[string]any{
			"swapped_from": p.ID,
		},
	}
}
