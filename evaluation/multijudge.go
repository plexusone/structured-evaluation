package evaluation

import (
	"sort"
)

// MultiJudgeResult aggregates evaluations from multiple judges.
// This improves reliability by combining perspectives and detecting disagreement.
type MultiJudgeResult struct {
	// Evaluations are the individual judge evaluations.
	Evaluations []*EvaluationReport `json:"evaluations"`

	// Judges contains metadata for each judge.
	Judges []*JudgeMetadata `json:"judges"`

	// AggregatedCategories are the combined category results.
	AggregatedCategories []CategoryResult `json:"aggregatedCategories"`

	// AggregationMethod describes how scores were combined.
	AggregationMethod AggregationMethod `json:"aggregationMethod"`

	// Agreement measures inter-judge agreement (0-1, higher = more agreement).
	Agreement float64 `json:"agreement"`

	// Disagreements lists categories where judges significantly disagreed.
	Disagreements []JudgeDisagreement `json:"disagreements,omitempty"`

	// ConsolidatedDecision is the final decision after aggregation.
	ConsolidatedDecision Decision `json:"consolidatedDecision"`

	// ConsolidatedFindings merges findings from all judges.
	ConsolidatedFindings []Finding `json:"consolidatedFindings"`
}

// AggregationMethod specifies how to combine multiple judge scores.
type AggregationMethod string

const (
	// AggregationMajority uses majority vote for pass/partial/fail.
	AggregationMajority AggregationMethod = "majority"

	// AggregationConservative uses the lowest/most critical score.
	AggregationConservative AggregationMethod = "conservative"

	// AggregationOptimistic uses the highest/most lenient score.
	AggregationOptimistic AggregationMethod = "optimistic"

	// AggregationUnanimous requires all judges to agree.
	AggregationUnanimous AggregationMethod = "unanimous"
)

// JudgeDisagreement captures where judges had significantly different scores.
type JudgeDisagreement struct {
	// Category is the evaluation dimension.
	Category string `json:"category"`

	// Scores are the individual judge scores.
	Scores []JudgeCategoricalScore `json:"scores"`

	// UniqueScores is the number of distinct scores given.
	UniqueScores int `json:"uniqueScores"`
}

// JudgeCategoricalScore is a categorical score from a specific judge.
type JudgeCategoricalScore struct {
	// JudgeID identifies the judge.
	JudgeID string `json:"judgeId"`

	// Score is the judge's categorical score.
	Score ScoreValue `json:"score"`
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

	// Collect judges
	for _, eval := range evaluations {
		if eval.Judge != nil {
			result.Judges = append(result.Judges, eval.Judge)
		}
	}

	// Aggregate category results
	result.AggregatedCategories = aggregateCategoryResults(evaluations, method)

	// Compute agreement
	result.Agreement = computeCategoricalAgreement(evaluations)

	// Find disagreements
	result.Disagreements = findCategoricalDisagreements(evaluations)

	// Consolidate findings (deduplicate similar ones)
	result.ConsolidatedFindings = consolidateFindings(evaluations)

	// Compute consolidated decision
	result.ConsolidatedDecision = consolidateDecision(evaluations, method)

	return result
}

// aggregateCategoryResults combines category results from multiple judges.
func aggregateCategoryResults(evaluations []*EvaluationReport, method AggregationMethod) []CategoryResult {
	if len(evaluations) == 0 {
		return nil
	}

	// Map category -> list of scores
	categoryScores := make(map[string][]ScoreValue)
	categoryReasonings := make(map[string][]string)

	for _, eval := range evaluations {
		for _, cat := range eval.Categories {
			categoryScores[cat.Category] = append(categoryScores[cat.Category], cat.Score)
			if cat.Reasoning != "" {
				categoryReasonings[cat.Category] = append(categoryReasonings[cat.Category], cat.Reasoning)
			}
		}
	}

	var results []CategoryResult
	for category, scores := range categoryScores {
		aggregatedScore := aggregateScores(scores, method)
		reasoning := "Aggregated from " + itoa(len(scores)) + " judges using " + string(method)

		results = append(results, CategoryResult{
			Category:  category,
			Score:     aggregatedScore,
			Reasoning: reasoning,
		})
	}

	return results
}

// aggregateScores combines categorical scores using the specified method.
func aggregateScores(scores []ScoreValue, method AggregationMethod) ScoreValue {
	if len(scores) == 0 {
		return ScoreFail
	}

	counts := make(map[ScoreValue]int)
	for _, s := range scores {
		counts[s]++
	}

	switch method {
	case AggregationConservative:
		// Most critical: fail > partial > pass
		if counts[ScoreFail] > 0 {
			return ScoreFail
		}
		if counts[ScorePartial] > 0 {
			return ScorePartial
		}
		return ScorePass

	case AggregationOptimistic:
		// Most lenient: pass > partial > fail
		if counts[ScorePass] > 0 {
			return ScorePass
		}
		if counts[ScorePartial] > 0 {
			return ScorePartial
		}
		return ScoreFail

	case AggregationUnanimous:
		// All must agree
		if len(counts) == 1 {
			for score := range counts {
				return score
			}
		}
		// No unanimous agreement - return partial as middle ground
		return ScorePartial

	case AggregationMajority:
		fallthrough
	default:
		// Majority vote
		var maxCount int
		var majorityScore ScoreValue
		for score, count := range counts {
			if count > maxCount {
				maxCount = count
				majorityScore = score
			}
		}
		// Require true majority (>50%)
		if float64(maxCount) > float64(len(scores))/2 {
			return majorityScore
		}
		// No clear majority - return partial
		return ScorePartial
	}
}

// computeCategoricalAgreement calculates inter-judge agreement.
func computeCategoricalAgreement(evaluations []*EvaluationReport) float64 {
	if len(evaluations) <= 1 {
		return 1.0
	}

	// Map category -> list of scores
	categoryScores := make(map[string][]ScoreValue)
	for _, eval := range evaluations {
		for _, cat := range eval.Categories {
			categoryScores[cat.Category] = append(categoryScores[cat.Category], cat.Score)
		}
	}

	if len(categoryScores) == 0 {
		return 1.0
	}

	// Calculate agreement for each category
	var totalAgreement float64
	for _, scores := range categoryScores {
		counts := make(map[ScoreValue]int)
		for _, s := range scores {
			counts[s]++
		}

		// Find the most common score
		var maxCount int
		for _, count := range counts {
			if count > maxCount {
				maxCount = count
			}
		}

		// Agreement = proportion that agree with majority
		categoryAgreement := float64(maxCount) / float64(len(scores))
		totalAgreement += categoryAgreement
	}

	return totalAgreement / float64(len(categoryScores))
}

// findCategoricalDisagreements identifies categories with disagreement.
func findCategoricalDisagreements(evaluations []*EvaluationReport) []JudgeDisagreement {
	if len(evaluations) <= 1 {
		return nil
	}

	// Map category -> scores with judge IDs
	categoryScores := make(map[string][]JudgeCategoricalScore)
	for _, eval := range evaluations {
		judgeID := ""
		if eval.Metadata.ReviewerID != "" {
			judgeID = eval.Metadata.ReviewerID
		} else if eval.Judge != nil && eval.Judge.JudgeID != "" {
			judgeID = eval.Judge.JudgeID
		}

		for _, cat := range eval.Categories {
			categoryScores[cat.Category] = append(categoryScores[cat.Category], JudgeCategoricalScore{
				JudgeID: judgeID,
				Score:   cat.Score,
			})
		}
	}

	var disagreements []JudgeDisagreement
	for category, scores := range categoryScores {
		// Count unique scores
		uniqueScores := make(map[ScoreValue]bool)
		for _, s := range scores {
			uniqueScores[s.Score] = true
		}

		// Disagreement if more than 1 unique score
		if len(uniqueScores) > 1 {
			disagreements = append(disagreements, JudgeDisagreement{
				Category:     category,
				Scores:       scores,
				UniqueScores: len(uniqueScores),
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
			return Decision{Status: DecisionFail, Rationale: "Conservative: at least one judge failed"}
		}
		if counts[DecisionConditional] > 0 {
			return Decision{Status: DecisionConditional, Rationale: "Conservative: at least one judge conditional"}
		}
		if counts[DecisionHumanReview] > 0 {
			return Decision{Status: DecisionHumanReview, Rationale: "Conservative: at least one judge needs human review"}
		}
		return Decision{Status: DecisionPass, Passed: true, Rationale: "All judges passed"}
	}

	// For unanimous, all must agree
	if method == AggregationUnanimous {
		if len(counts) == 1 {
			for status := range counts {
				return Decision{Status: status, Passed: status == DecisionPass, Rationale: "Unanimous decision"}
			}
		}
		return Decision{Status: DecisionHumanReview, Rationale: "No unanimous agreement"}
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
		return Decision{
			Status:    majorityDecision,
			Passed:    majorityDecision == DecisionPass,
			Rationale: "Majority decision: " + itoa(maxCount) + "/" + itoa(len(evaluations)) + " judges",
		}
	}

	// No clear majority, recommend human review
	return Decision{Status: DecisionHumanReview, Rationale: "No clear majority among judges"}
}
