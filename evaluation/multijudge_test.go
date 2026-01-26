package evaluation

import "testing"

func TestAggregateEvaluations_Mean(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.WeightedScore = 8.0
	eval1.Decision = Decision{Status: DecisionPass}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.WeightedScore = 6.0
	eval2.Decision = Decision{Status: DecisionPass}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationMean)

	expectedMean := 7.0
	if result.AggregatedScore != expectedMean {
		t.Errorf("Expected mean score %f, got %f", expectedMean, result.AggregatedScore)
	}

	if result.AggregationMethod != AggregationMean {
		t.Errorf("Expected method 'mean', got %s", result.AggregationMethod)
	}
}

func TestAggregateEvaluations_Median(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.WeightedScore = 8.0

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.WeightedScore = 6.0

	eval3 := NewEvaluationReport("test", "doc.md")
	eval3.WeightedScore = 2.0

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2, eval3}, AggregationMedian)

	expectedMedian := 6.0
	if result.AggregatedScore != expectedMedian {
		t.Errorf("Expected median score %f, got %f", expectedMedian, result.AggregatedScore)
	}
}

func TestAggregateEvaluations_Conservative(t *testing.T) {
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.WeightedScore = 8.0
	eval1.Decision = Decision{Status: DecisionPass}

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.WeightedScore = 6.0
	eval2.Decision = Decision{Status: DecisionFail}

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationConservative)

	expectedMin := 6.0
	if result.AggregatedScore != expectedMin {
		t.Errorf("Expected min score %f, got %f", expectedMin, result.AggregatedScore)
	}

	// Conservative should fail if any judge fails
	if result.ConsolidatedDecision.Status != DecisionFail {
		t.Errorf("Expected fail decision, got %s", result.ConsolidatedDecision.Status)
	}
}

func TestAggregateEvaluations_Agreement(t *testing.T) {
	// High agreement (same scores)
	eval1 := NewEvaluationReport("test", "doc.md")
	eval1.WeightedScore = 7.0

	eval2 := NewEvaluationReport("test", "doc.md")
	eval2.WeightedScore = 7.0

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationMean)

	if result.Agreement != 1.0 {
		t.Errorf("Expected perfect agreement (1.0), got %f", result.Agreement)
	}

	// Lower agreement (different scores)
	eval3 := NewEvaluationReport("test", "doc.md")
	eval3.WeightedScore = 2.0

	result2 := AggregateEvaluations([]*EvaluationReport{eval1, eval3}, AggregationMean)

	if result2.Agreement >= 1.0 {
		t.Errorf("Expected lower agreement, got %f", result2.Agreement)
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

	result := AggregateEvaluations([]*EvaluationReport{eval1, eval2}, AggregationMean)

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
	result := AggregateEvaluations([]*EvaluationReport{}, AggregationMean)

	if result.AggregatedScore != 0 {
		t.Errorf("Expected 0 score for empty input, got %f", result.AggregatedScore)
	}
}
