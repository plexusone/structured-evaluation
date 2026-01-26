package evaluation

import (
	"math"
	"sort"
)

// MultiJudgeResult aggregates evaluations from multiple judges.
// This improves reliability by combining perspectives and detecting disagreement.
type MultiJudgeResult struct {
	// Evaluations are the individual judge evaluations.
	Evaluations []*EvaluationReport `json:"evaluations"`

	// Judges contains metadata for each judge.
	Judges []*JudgeMetadata `json:"judges"`

	// AggregatedScore is the combined score (mean, median, or weighted).
	AggregatedScore float64 `json:"aggregated_score"`

	// AggregationMethod describes how scores were combined.
	AggregationMethod AggregationMethod `json:"aggregation_method"`

	// Agreement measures inter-judge agreement (0-1, higher = more agreement).
	Agreement float64 `json:"agreement"`

	// Disagreements lists categories where judges significantly disagreed.
	Disagreements []JudgeDisagreement `json:"disagreements,omitempty"`

	// ConsolidatedDecision is the final decision after aggregation.
	ConsolidatedDecision Decision `json:"consolidated_decision"`

	// ConsolidatedFindings merges findings from all judges.
	ConsolidatedFindings []Finding `json:"consolidated_findings"`
}

// AggregationMethod specifies how to combine multiple judge scores.
type AggregationMethod string

const (
	// AggregationMean uses the arithmetic mean of scores.
	AggregationMean AggregationMethod = "mean"

	// AggregationMedian uses the median score.
	AggregationMedian AggregationMethod = "median"

	// AggregationWeighted uses weighted average based on judge confidence.
	AggregationWeighted AggregationMethod = "weighted"

	// AggregationMajority uses majority vote for pass/fail.
	AggregationMajority AggregationMethod = "majority"

	// AggregationConservative uses the lowest/most critical score.
	AggregationConservative AggregationMethod = "conservative"
)

// JudgeDisagreement captures where judges had significantly different scores.
type JudgeDisagreement struct {
	// Category is the evaluation dimension.
	Category string `json:"category"`

	// Scores are the individual judge scores.
	Scores []JudgeScore `json:"scores"`

	// Range is the difference between max and min scores.
	Range float64 `json:"range"`

	// StandardDeviation measures score spread.
	StandardDeviation float64 `json:"standard_deviation"`
}

// JudgeScore is a score from a specific judge.
type JudgeScore struct {
	// JudgeID identifies the judge.
	JudgeID string `json:"judge_id"`

	// Score is the judge's score.
	Score float64 `json:"score"`
}

// AggregateEvaluations combines multiple evaluation reports.
func AggregateEvaluations(evaluations []*EvaluationReport, method AggregationMethod) *MultiJudgeResult {
	if len(evaluations) == 0 {
		return &MultiJudgeResult{}
	}

	result := &MultiJudgeResult{
		Evaluations:       evaluations,
		AggregationMethod: method,
		Judges:            make([]*JudgeMetadata, 0),
	}

	// Collect all weighted scores
	scores := make([]float64, len(evaluations))
	for i, eval := range evaluations {
		scores[i] = eval.WeightedScore
	}

	// Compute aggregated score
	switch method {
	case AggregationMean:
		result.AggregatedScore = mean(scores)
	case AggregationMedian:
		result.AggregatedScore = median(scores)
	case AggregationConservative:
		result.AggregatedScore = min(scores)
	default:
		result.AggregatedScore = mean(scores)
	}

	// Compute agreement
	result.Agreement = computeAgreement(scores)

	// Find disagreements per category
	result.Disagreements = findDisagreements(evaluations)

	// Consolidate findings (deduplicate similar ones)
	result.ConsolidatedFindings = consolidateFindings(evaluations)

	// Compute consolidated decision
	result.ConsolidatedDecision = consolidateDecision(evaluations, method)

	return result
}

// computeAgreement calculates inter-judge agreement using normalized standard deviation.
func computeAgreement(scores []float64) float64 {
	if len(scores) <= 1 {
		return 1.0
	}

	stdDev := standardDeviation(scores)
	// Normalize: max possible std dev for 0-10 scale is 5
	// Agreement = 1 - (stdDev / 5)
	agreement := 1 - (stdDev / 5)
	if agreement < 0 {
		agreement = 0
	}
	return agreement
}

// findDisagreements identifies categories with significant disagreement.
func findDisagreements(evaluations []*EvaluationReport) []JudgeDisagreement {
	if len(evaluations) <= 1 {
		return nil
	}

	// Map category -> scores
	categoryScores := make(map[string][]float64)
	for _, eval := range evaluations {
		for _, cat := range eval.Categories {
			categoryScores[cat.Category] = append(categoryScores[cat.Category], cat.Score)
		}
	}

	var disagreements []JudgeDisagreement
	for category, scores := range categoryScores {
		stdDev := standardDeviation(scores)
		scoreRange := maxVal(scores) - minVal(scores)

		// Threshold: std dev > 1.5 or range > 3 indicates disagreement
		if stdDev > 1.5 || scoreRange > 3 {
			judgeScores := make([]JudgeScore, len(scores))
			for i, s := range scores {
				judgeID := ""
				if i < len(evaluations) && evaluations[i].Metadata.ReviewerID != "" {
					judgeID = evaluations[i].Metadata.ReviewerID
				}
				judgeScores[i] = JudgeScore{JudgeID: judgeID, Score: s}
			}

			disagreements = append(disagreements, JudgeDisagreement{
				Category:          category,
				Scores:            judgeScores,
				Range:             scoreRange,
				StandardDeviation: stdDev,
			})
		}
	}

	return disagreements
}

// consolidateFindings merges findings from all evaluations.
func consolidateFindings(evaluations []*EvaluationReport) []Finding {
	// Use map to deduplicate by title+category
	seen := make(map[string]Finding)
	for _, eval := range evaluations {
		for _, f := range eval.Findings {
			key := f.Category + ":" + f.Title
			if existing, ok := seen[key]; ok {
				// Keep higher severity
				if f.Severity.Weight() > existing.Severity.Weight() {
					seen[key] = f
				}
			} else {
				seen[key] = f
			}
		}
	}

	findings := make([]Finding, 0, len(seen))
	for _, f := range seen {
		findings = append(findings, f)
	}

	// Sort by severity (highest first)
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Severity.Weight() > findings[j].Severity.Weight()
	})

	return findings
}

// consolidateDecision determines the final decision.
func consolidateDecision(evaluations []*EvaluationReport, method AggregationMethod) Decision {
	if len(evaluations) == 0 {
		return Decision{Status: DecisionHumanReview}
	}

	// Count decision types
	counts := make(map[DecisionStatus]int)
	for _, eval := range evaluations {
		counts[eval.Decision.Status]++
	}

	// For conservative method, any fail means fail
	if method == AggregationConservative {
		if counts[DecisionFail] > 0 {
			return Decision{Status: DecisionFail}
		}
		if counts[DecisionConditional] > 0 {
			return Decision{Status: DecisionConditional}
		}
		if counts[DecisionHumanReview] > 0 {
			return Decision{Status: DecisionHumanReview}
		}
		return Decision{Status: DecisionPass}
	}

	// For majority, use most common decision
	var maxCount int
	var majorityDecision DecisionStatus
	for status, count := range counts {
		if count > maxCount {
			maxCount = count
			majorityDecision = status
		}
	}

	// Require true majority (>50%)
	if float64(maxCount) > float64(len(evaluations))/2 {
		return Decision{Status: majorityDecision}
	}

	// No clear majority, recommend human review
	return Decision{Status: DecisionHumanReview}
}

// Helper functions

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func standardDeviation(values []float64) float64 {
	if len(values) <= 1 {
		return 0
	}
	m := mean(values)
	var sumSquares float64
	for _, v := range values {
		sumSquares += (v - m) * (v - m)
	}
	return math.Sqrt(sumSquares / float64(len(values)))
}

func minVal(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func maxVal(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func min(values []float64) float64 {
	return minVal(values)
}
