package evaluation

import "testing"

func TestDefaultPassCriteria(t *testing.T) {
	criteria := DefaultPassCriteria()

	if criteria.MinCategoriesPassing != "all_required" {
		t.Errorf("Expected 'all_required', got %s", criteria.MinCategoriesPassing)
	}
	if criteria.MaxFindings == nil {
		t.Fatal("Expected MaxFindings to be set")
	}
	if criteria.MaxFindings.Critical != 0 {
		t.Errorf("Expected 0 critical allowed, got %d", criteria.MaxFindings.Critical)
	}
	if criteria.MaxFindings.High != 0 {
		t.Errorf("Expected 0 high allowed, got %d", criteria.MaxFindings.High)
	}
}

func TestStrictPassCriteria(t *testing.T) {
	criteria := StrictPassCriteria()

	if criteria.MinCategoriesPassing != "all" {
		t.Errorf("Expected 'all', got %s", criteria.MinCategoriesPassing)
	}
	if criteria.MaxFindings.Medium != 3 {
		t.Errorf("Expected max 3 medium, got %d", criteria.MaxFindings.Medium)
	}
}

func TestEvaluate_Pass(t *testing.T) {
	rubric := NewRubricSet("test", "Test", "1.0")
	cat := NewCategory("quality", "Quality", "Test").SetRequired(true)
	rubric.AddCategory(*cat)

	results := []CategoryResult{
		{Category: "quality", Score: ScorePass, Reasoning: "Good"},
	}
	findings := []Finding{}
	criteria := DefaultPassCriteria()

	decision := Evaluate(results, findings, criteria, rubric)

	if decision.Status != DecisionPass {
		t.Errorf("Expected pass, got %s", decision.Status)
	}
	if !decision.Passed {
		t.Error("Expected Passed to be true")
	}
}

func TestEvaluate_FailOnCriticalFindings(t *testing.T) {
	results := []CategoryResult{
		{Category: "quality", Score: ScorePass},
	}
	findings := []Finding{
		{ID: "F1", Severity: SeverityCritical, Title: "Critical issue"},
	}
	criteria := DefaultPassCriteria()

	decision := Evaluate(results, findings, criteria, nil)

	if decision.Status != DecisionFail {
		t.Errorf("Expected fail, got %s", decision.Status)
	}
	if decision.Passed {
		t.Error("Expected Passed to be false")
	}
}

func TestEvaluate_FailOnHighFindings(t *testing.T) {
	results := []CategoryResult{
		{Category: "quality", Score: ScorePass},
	}
	findings := []Finding{
		{ID: "F1", Severity: SeverityHigh, Title: "High severity issue"},
	}
	criteria := DefaultPassCriteria()

	decision := Evaluate(results, findings, criteria, nil)

	if decision.Status != DecisionFail {
		t.Errorf("Expected fail, got %s", decision.Status)
	}
}

func TestEvaluate_FailOnRequiredCategoryFail(t *testing.T) {
	rubric := NewRubricSet("test", "Test", "1.0")
	cat := NewCategory("quality", "Quality", "Test").SetRequired(true)
	rubric.AddCategory(*cat)

	results := []CategoryResult{
		{Category: "quality", Score: ScoreFail, Reasoning: "Poor"},
	}
	findings := []Finding{}
	criteria := DefaultPassCriteria()

	decision := Evaluate(results, findings, criteria, rubric)

	if decision.Status != DecisionFail {
		t.Errorf("Expected fail, got %s", decision.Status)
	}
}

func TestEvaluate_ConditionalOnPartial(t *testing.T) {
	rubric := NewRubricSet("test", "Test", "1.0")
	cat := NewCategory("quality", "Quality", "Test").SetRequired(true)
	rubric.AddCategory(*cat)

	results := []CategoryResult{
		{Category: "quality", Score: ScorePartial, Reasoning: "OK"},
	}
	findings := []Finding{}
	criteria := PassCriteria{
		MinCategoriesPassing: "all_required",
		MaxFindings:          &FindingLimits{Critical: 0, High: 0, Medium: -1},
	}

	// With "all_required", partial on required should fail
	decision := Evaluate(results, findings, criteria, rubric)
	if decision.Status != DecisionFail {
		t.Errorf("Expected fail for partial on required category, got %s", decision.Status)
	}
}

func TestEvaluate_ConditionalOnMediumFindings(t *testing.T) {
	results := []CategoryResult{
		{Category: "quality", Score: ScorePass},
	}
	findings := []Finding{
		{ID: "F1", Severity: SeverityMedium, Title: "Medium issue"},
	}
	criteria := DefaultPassCriteria()

	decision := Evaluate(results, findings, criteria, nil)

	if decision.Status != DecisionConditional {
		t.Errorf("Expected conditional, got %s", decision.Status)
	}
	if !decision.Passed {
		t.Error("Expected Passed to be true for conditional with medium findings")
	}
}

func TestEvaluate_AllCategoriesMustPass(t *testing.T) {
	results := []CategoryResult{
		{Category: "cat1", Score: ScorePass},
		{Category: "cat2", Score: ScorePartial},
	}
	findings := []Finding{}
	criteria := PassCriteria{
		MinCategoriesPassing: "all",
		MaxFindings:          &FindingLimits{Critical: 0, High: 0, Medium: -1},
	}

	decision := Evaluate(results, findings, criteria, nil)

	if decision.Status != DecisionFail {
		t.Errorf("Expected fail when not all categories pass, got %s", decision.Status)
	}
}

func TestEvaluate_NumericThreshold(t *testing.T) {
	results := []CategoryResult{
		{Category: "cat1", Score: ScorePass},
		{Category: "cat2", Score: ScorePass},
		{Category: "cat3", Score: ScoreFail},
	}
	findings := []Finding{}
	criteria := PassCriteria{
		MinCategoriesPassing: "2",
		MaxFindings:          &FindingLimits{Critical: 0, High: 0, Medium: -1},
	}

	decision := Evaluate(results, findings, criteria, nil)

	if decision.Status != DecisionPass {
		t.Errorf("Expected pass with 2/3 categories passing, got %s", decision.Status)
	}
}

func TestScoreValue_Methods(t *testing.T) {
	if !ScorePass.IsPassing() {
		t.Error("ScorePass should be passing")
	}
	if ScoreFail.IsPassing() {
		t.Error("ScoreFail should not be passing")
	}
	if !ScorePartial.IsPartial() {
		t.Error("ScorePartial should be partial")
	}
	if !ScoreFail.IsFailing() {
		t.Error("ScoreFail should be failing")
	}
}

func TestScoreValue_Icon(t *testing.T) {
	tests := []struct {
		score    ScoreValue
		expected string
	}{
		{ScorePass, "🟢"},
		{ScorePartial, "🟡"},
		{ScoreFail, "🔴"},
		{ScoreValue("unknown"), "⚪"},
	}

	for _, tt := range tests {
		if got := tt.score.Icon(); got != tt.expected {
			t.Errorf("Score %s: expected icon %s, got %s", tt.score, tt.expected, got)
		}
	}
}

func TestCountResults(t *testing.T) {
	results := []CategoryResult{
		{Category: "cat1", Score: ScorePass},
		{Category: "cat2", Score: ScorePass},
		{Category: "cat3", Score: ScorePartial},
		{Category: "cat4", Score: ScoreFail},
	}

	counts := CountResults(results)

	if counts.Pass != 2 {
		t.Errorf("Expected 2 pass, got %d", counts.Pass)
	}
	if counts.Partial != 1 {
		t.Errorf("Expected 1 partial, got %d", counts.Partial)
	}
	if counts.Fail != 1 {
		t.Errorf("Expected 1 fail, got %d", counts.Fail)
	}
	if counts.Total != 4 {
		t.Errorf("Expected 4 total, got %d", counts.Total)
	}
}

func TestAllRequiredPassing(t *testing.T) {
	rubric := NewRubricSet("test", "Test", "1.0")
	cat1 := NewCategory("required1", "Required 1", "Test").SetRequired(true)
	cat2 := NewCategory("optional1", "Optional 1", "Test")
	rubric.AddCategory(*cat1).AddCategory(*cat2)

	// All required passing
	results1 := []CategoryResult{
		{Category: "required1", Score: ScorePass},
		{Category: "optional1", Score: ScoreFail}, // Optional can fail
	}
	if !AllRequiredPassing(results1, rubric) {
		t.Error("Expected all required to be passing")
	}

	// Required failing
	results2 := []CategoryResult{
		{Category: "required1", Score: ScoreFail},
		{Category: "optional1", Score: ScorePass},
	}
	if AllRequiredPassing(results2, rubric) {
		t.Error("Expected all required to not be passing")
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
