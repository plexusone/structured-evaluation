package evaluation

import "testing"

func TestEvaluate_Pass(t *testing.T) {
	findings := []Finding{}
	score := 8.0
	criteria := DefaultPassCriteria()

	decision := Evaluate(findings, score, criteria)

	if !decision.Passed {
		t.Errorf("expected pass, got fail: %s", decision.Rationale)
	}
	if decision.Status != DecisionPass {
		t.Errorf("expected status pass, got %s", decision.Status)
	}
}

func TestEvaluate_FailCritical(t *testing.T) {
	findings := []Finding{
		{Severity: SeverityCritical, Title: "Critical issue"},
	}
	score := 9.0
	criteria := DefaultPassCriteria()

	decision := Evaluate(findings, score, criteria)

	if decision.Passed {
		t.Error("expected fail due to critical finding")
	}
	if decision.Status != DecisionFail {
		t.Errorf("expected status fail, got %s", decision.Status)
	}
}

func TestEvaluate_FailHigh(t *testing.T) {
	findings := []Finding{
		{Severity: SeverityHigh, Title: "High issue"},
	}
	score := 9.0
	criteria := DefaultPassCriteria()

	decision := Evaluate(findings, score, criteria)

	if decision.Passed {
		t.Error("expected fail due to high finding")
	}
	if decision.Status != DecisionFail {
		t.Errorf("expected status fail, got %s", decision.Status)
	}
}

func TestEvaluate_Conditional(t *testing.T) {
	findings := []Finding{
		{Severity: SeverityMedium, Title: "Medium issue"},
	}
	score := 8.0
	criteria := DefaultPassCriteria()

	decision := Evaluate(findings, score, criteria)

	if !decision.Passed {
		t.Error("expected conditional pass")
	}
	if decision.Status != DecisionConditional {
		t.Errorf("expected status conditional, got %s", decision.Status)
	}
}

func TestEvaluate_BelowMinScore(t *testing.T) {
	findings := []Finding{}
	score := 5.0
	criteria := DefaultPassCriteria()

	decision := Evaluate(findings, score, criteria)

	if decision.Passed {
		t.Error("expected fail due to low score")
	}
	if decision.Status != DecisionHumanReview {
		t.Errorf("expected status human_review, got %s", decision.Status)
	}
}

func TestFindingCounts_BlockingCount(t *testing.T) {
	counts := FindingCounts{
		Critical: 1,
		High:     2,
		Medium:   3,
		Low:      4,
	}

	if counts.BlockingCount() != 3 {
		t.Errorf("expected blocking count 3, got %d", counts.BlockingCount())
	}
}
