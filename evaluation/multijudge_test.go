package evaluation

import "testing"

func TestAggregateEvaluations_Majority(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})
	eval1.Decision = Decision{Status: DecisionPass, Passed: true}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})
	eval2.Decision = Decision{Status: DecisionPass, Passed: true}

	eval3 := NewEvaluationReport("test", "doc.md")
	eval3.AddCategoryResult(CategoryResult{Category: "quality", Score: ScoreFail})
	eval3.Decision = Decision{Status: DecisionFail}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2, eval3}, AggregationMajority)

	if result.AggregationMethod != AggregationMajority {
		t.Errorf("Expected method 'majority', got %s", result.AggregationMethod)
	}

	// 2/3 judges passed, so majority should pass
	if result.ConsolidatedDecision.Status != DecisionPass {
		t.Errorf("Expected pass decision, got %s", result.ConsolidatedDecision.Status)
	}

	// Quality category should aggregate to pass (2/3)
	if len(result.AggregatedCategories) != 1 {
		t.Fatalf("Expected 1 aggregated category, got %d", len(result.AggregatedCategories))
	}
	if result.AggregatedCategories[0].Score != ScorePass {
		t.Errorf("Expected aggregated score 'pass', got %s", result.AggregatedCategories[0].Score)
	}
}

func TestAggregateEvaluations_Conservative(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})
	eval1.Decision = Decision{Status: DecisionPass, Passed: true}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePartial})
	eval2.Decision = Decision{Status: DecisionConditional}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationConservative)

	// Conservative should take the worse score
	if len(result.AggregatedCategories) != 1 {
		t.Fatalf("Expected 1 aggregated category, got %d", len(result.AggregatedCategories))
	}
	if result.AggregatedCategories[0].Score != ScorePartial {
		t.Errorf("Expected aggregated score 'partial', got %s", result.AggregatedCategories[0].Score)
	}

	// Conservative should fail if any judge has conditional
	if result.ConsolidatedDecision.Status != DecisionConditional {
		t.Errorf("Expected conditional decision, got %s", result.ConsolidatedDecision.Status)
	}
}

func TestAggregateEvaluations_ConservativeFail(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.Decision = Decision{Status: DecisionPass, Passed: true}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.Decision = Decision{Status: DecisionFail}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationConservative)

	// Conservative should fail if any judge fails
	if result.ConsolidatedDecision.Status != DecisionFail {
		t.Errorf("Expected fail decision, got %s", result.ConsolidatedDecision.Status)
	}
}

func TestAggregateEvaluations_Unanimous(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})
	eval1.Decision = Decision{Status: DecisionPass, Passed: true}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})
	eval2.Decision = Decision{Status: DecisionPass, Passed: true}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationUnanimous)

	// Unanimous pass
	if result.ConsolidatedDecision.Status != DecisionPass {
		t.Errorf("Expected pass decision, got %s", result.ConsolidatedDecision.Status)
	}

	// Aggregated category should be pass (unanimous)
	if result.AggregatedCategories[0].Score != ScorePass {
		t.Errorf("Expected aggregated score 'pass', got %s", result.AggregatedCategories[0].Score)
	}
}

func TestAggregateEvaluations_UnanimousDisagreement(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})
	eval1.Decision = Decision{Status: DecisionPass, Passed: true}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddCategoryResult(CategoryResult{Category: "quality", Score: ScoreFail})
	eval2.Decision = Decision{Status: DecisionFail}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationUnanimous)

	// No unanimous agreement - should need human review
	if result.ConsolidatedDecision.Status != DecisionHumanReview {
		t.Errorf("Expected human_review decision, got %s", result.ConsolidatedDecision.Status)
	}

	// No unanimous score - should return partial
	if result.AggregatedCategories[0].Score != ScorePartial {
		t.Errorf("Expected aggregated score 'partial' for disagreement, got %s", result.AggregatedCategories[0].Score)
	}
}

func TestAggregateEvaluations_Agreement(t *testing.T) {
	// High agreement (same scores)
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationMajority)

	if result.Agreement != 1.0 {
		t.Errorf("Expected perfect agreement (1.0), got %f", result.Agreement)
	}

	// Lower agreement (different scores)
	eval3 := NewEvaluationReport("test", "doc.md")
	eval3.AddCategoryResult(CategoryResult{Category: "quality", Score: ScoreFail})

	result2 := AggregateEvaluations([]*EvaluationReport{eval1, eval3}, AggregationMajority)

	if result2.Agreement >= 1.0 {
		t.Errorf("Expected lower agreement, got %f", result2.Agreement)
	}
}

func TestAggregateEvaluations_Disagreements(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddCategoryResult(CategoryResult{Category: "quality", Score: ScorePass})

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddCategoryResult(CategoryResult{Category: "quality", Score: ScoreFail})

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationMajority)

	if len(result.Disagreements) != 1 {
		t.Errorf("Expected 1 disagreement, got %d", len(result.Disagreements))
	}
	if result.Disagreements[0].Category != "quality" {
		t.Errorf("Expected disagreement on 'quality', got %s", result.Disagreements[0].Category)
	}
	if result.Disagreements[0].UniqueScores != 2 {
		t.Errorf("Expected 2 unique scores, got %d", result.Disagreements[0].UniqueScores)
	}
}

func TestAggregateEvaluations_ConsolidateFindings(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.AddFinding(Finding{
		Category: "security",
		Title:    "SQL Injection",
		Severity: SeverityHigh,
	})

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.AddFinding(Finding{
		Category: "security",
		Title:    "SQL Injection",  // Same finding
		Severity: SeverityCritical, // But higher severity
	})
	eval2.AddFinding(Finding{
		Category: "performance",
		Title:    "Slow query",
		Severity: SeverityMedium,
	})

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationMajority)

	// Should have 2 unique findings
	if len(result.ConsolidatedFindings) != 2 {
		t.Errorf("Expected 2 consolidated findings, got %d", len(result.ConsolidatedFindings))
	}

	// SQL Injection should have the higher severity (critical)
	for _, f := range result.ConsolidatedFindings {
		if f.Title == "SQL Injection" && f.Severity != SeverityCritical {
			t.Errorf("Expected SQL Injection to have critical severity, got %s", f.Severity)
		}
	}
}

func TestAggregateEvaluations_Empty(t *testing.T) {
	result := AggregateEvaluations([]*EvaluationReport{}, AggregationMajority)

	if len(result.AggregatedCategories) != 0 {
		t.Errorf("Expected 0 categories for empty input, got %d", len(result.AggregatedCategories))
	}
}

func TestAggregateScores(t *testing.T) {
	tests := []struct {
		name     string
		scores   []ScoreValue
		method   AggregationMethod
		expected ScoreValue
	}{
		{
			name:     "majority pass",
			scores:   []ScoreValue{ScorePass, ScorePass, ScoreFail},
			method:   AggregationMajority,
			expected: ScorePass,
		},
		{
			name:     "majority fail",
			scores:   []ScoreValue{ScoreFail, ScoreFail, ScorePass},
			method:   AggregationMajority,
			expected: ScoreFail,
		},
		{
			name:     "majority no clear winner",
			scores:   []ScoreValue{ScorePass, ScoreFail},
			method:   AggregationMajority,
			expected: ScorePartial, // No majority
		},
		{
			name:     "conservative takes worst",
			scores:   []ScoreValue{ScorePass, ScorePartial, ScorePass},
			method:   AggregationConservative,
			expected: ScorePartial,
		},
		{
			name:     "conservative with fail",
			scores:   []ScoreValue{ScorePass, ScoreFail},
			method:   AggregationConservative,
			expected: ScoreFail,
		},
		{
			name:     "optimistic takes best",
			scores:   []ScoreValue{ScorePartial, ScoreFail},
			method:   AggregationOptimistic,
			expected: ScorePartial,
		},
		{
			name:     "optimistic with pass",
			scores:   []ScoreValue{ScorePass, ScoreFail},
			method:   AggregationOptimistic,
			expected: ScorePass,
		},
		{
			name:     "unanimous agrees",
			scores:   []ScoreValue{ScorePass, ScorePass},
			method:   AggregationUnanimous,
			expected: ScorePass,
		},
		{
			name:     "unanimous disagrees",
			scores:   []ScoreValue{ScorePass, ScoreFail},
			method:   AggregationUnanimous,
			expected: ScorePartial, // Middle ground
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aggregateScores(tt.scores, tt.method)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
