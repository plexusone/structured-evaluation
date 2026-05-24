package evaluation

import "math"

// IRRMetrics contains inter-rater reliability metrics.
// These metrics are useful when comparing LLM and human ratings.
type IRRMetrics struct {
	// ExactAgreement is the percentage of exact score matches.
	ExactAgreement float64 `json:"exactAgreement"`

	// AdjacentAgreement is the percentage within ±1 of each other.
	AdjacentAgreement float64 `json:"adjacentAgreement"`

	// MeanAbsoluteDifference is the average absolute difference.
	MeanAbsoluteDifference float64 `json:"meanAbsoluteDifference"`

	// PearsonCorrelation measures linear correlation (-1 to 1).
	PearsonCorrelation float64 `json:"pearsonCorrelation"`

	// SampleSize is the number of paired ratings.
	SampleSize int `json:"sampleSize"`
}

// RatingPair represents a pair of ratings for the same item.
type RatingPair struct {
	// Rater1 is the first rater's score (e.g., human).
	Rater1 float64

	// Rater2 is the second rater's score (e.g., LLM).
	Rater2 float64

	// Category is the category being rated.
	Category string

	// ItemID identifies the item being rated.
	ItemID string
}

// ComputeIRR calculates inter-rater reliability metrics from paired ratings.
func ComputeIRR(pairs []RatingPair) *IRRMetrics {
	if len(pairs) == 0 {
		return &IRRMetrics{}
	}

	n := float64(len(pairs))
	var exactMatches, adjacentMatches int
	var sumDiff, sumR1, sumR2, sumR1Sq, sumR2Sq, sumR1R2 float64

	for _, p := range pairs {
		diff := math.Abs(p.Rater1 - p.Rater2)

		if diff == 0 {
			exactMatches++
		}
		if diff <= 1 {
			adjacentMatches++
		}

		sumDiff += diff
		sumR1 += p.Rater1
		sumR2 += p.Rater2
		sumR1Sq += p.Rater1 * p.Rater1
		sumR2Sq += p.Rater2 * p.Rater2
		sumR1R2 += p.Rater1 * p.Rater2
	}

	// Pearson correlation
	numerator := n*sumR1R2 - sumR1*sumR2
	denominator := math.Sqrt((n*sumR1Sq - sumR1*sumR1) * (n*sumR2Sq - sumR2*sumR2))

	var pearson float64
	if denominator > 0 {
		pearson = numerator / denominator
	}

	return &IRRMetrics{
		ExactAgreement:         float64(exactMatches) / n,
		AdjacentAgreement:      float64(adjacentMatches) / n,
		MeanAbsoluteDifference: sumDiff / n,
		PearsonCorrelation:     pearson,
		SampleSize:             len(pairs),
	}
}

// ComputeIRRFromResults computes IRR metrics from two sets of category results.
// Useful for comparing LLM evaluation with human ground truth.
func ComputeIRRFromResults(results1, results2 []CategoryResult) *IRRMetrics {
	// Build map of results2 by category
	r2Map := make(map[string]*CategoryResult)
	for i := range results2 {
		r2Map[results2[i].Category] = &results2[i]
	}

	var pairs []RatingPair
	for _, r1 := range results1 {
		r2, ok := r2Map[r1.Category]
		if !ok {
			continue
		}

		// Use numeric scores if available, otherwise convert categorical
		score1 := getNumericOrCategorical(r1)
		score2 := getNumericOrCategorical(*r2)

		pairs = append(pairs, RatingPair{
			Rater1:   score1,
			Rater2:   score2,
			Category: r1.Category,
		})
	}

	return ComputeIRR(pairs)
}

// getNumericOrCategorical returns numeric score if available,
// otherwise converts categorical to numeric (pass=5, partial=3, fail=1).
func getNumericOrCategorical(cr CategoryResult) float64 {
	if cr.NumericScore != nil {
		return *cr.NumericScore
	}

	// Default mapping for categorical to 1-5 scale
	switch cr.Score {
	case ScorePass:
		return 5.0
	case ScorePartial:
		return 3.0
	case ScoreFail:
		return 1.0
	default:
		return 3.0
	}
}

// CategoricalAgreement computes agreement between categorical scores.
type CategoricalAgreement struct {
	// ExactAgreement is percentage of exact categorical matches.
	ExactAgreement float64 `json:"exactAgreement"`

	// ConfusionMatrix shows disagreement patterns.
	// Keys are "rater1_score:rater2_score" (e.g., "pass:partial").
	ConfusionMatrix map[string]int `json:"confusionMatrix"`

	// SampleSize is the number of paired ratings.
	SampleSize int `json:"sampleSize"`
}

// ComputeCategoricalAgreement computes agreement between categorical scores.
func ComputeCategoricalAgreement(results1, results2 []CategoryResult) *CategoricalAgreement {
	r2Map := make(map[string]*CategoryResult)
	for i := range results2 {
		r2Map[results2[i].Category] = &results2[i]
	}

	confusion := make(map[string]int)
	var exactMatches, total int

	for _, r1 := range results1 {
		r2, ok := r2Map[r1.Category]
		if !ok {
			continue
		}

		total++
		key := string(r1.Score) + ":" + string(r2.Score)
		confusion[key]++

		if r1.Score == r2.Score {
			exactMatches++
		}
	}

	var exactAgreement float64
	if total > 0 {
		exactAgreement = float64(exactMatches) / float64(total)
	}

	return &CategoricalAgreement{
		ExactAgreement:  exactAgreement,
		ConfusionMatrix: confusion,
		SampleSize:      total,
	}
}
