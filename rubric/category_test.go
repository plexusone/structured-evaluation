package rubric

import "testing"

func TestCategoryResult_AddFinding_ComputesSeverity(t *testing.T) {
	cr := NewCategoryResult("clarity", ScorePartial, "some ambiguity")
	if cr.Severity != "" {
		t.Fatalf("new category should have no severity, got %q", cr.Severity)
	}

	cr.AddFinding(Finding{Severity: SeverityLow})
	if cr.Severity != SeverityLow {
		t.Errorf("after one low finding: got %q, want %q", cr.Severity, SeverityLow)
	}

	cr.AddFinding(Finding{Severity: SeverityCritical})
	if cr.Severity != SeverityCritical {
		t.Errorf("after adding a critical finding: got %q, want %q", cr.Severity, SeverityCritical)
	}
}

func TestRubric_AddCategoryResult_ComputesSeverityIfUnset(t *testing.T) {
	r := NewRubric("prd", "prd.md")

	// Category built via struct literal (bypassing AddFinding) still gets
	// its severity computed on add.
	r.AddCategoryResult(CategoryResult{
		Category: "clarity",
		Score:    ScoreFail,
		Findings: []Finding{{Severity: SeverityHigh}, {Severity: SeverityLow}},
	})

	if got := r.Categories[0].Severity; got != SeverityHigh {
		t.Errorf("AddCategoryResult: severity = %q, want %q", got, SeverityHigh)
	}
}

func TestRubric_AddCategoryResult_RespectsExplicitSeverity(t *testing.T) {
	r := NewRubric("prd", "prd.md")

	// An explicitly-set severity is not overwritten, even if it disagrees
	// with what the findings alone would compute.
	r.AddCategoryResult(CategoryResult{
		Category: "clarity",
		Score:    ScorePass,
		Severity: SeverityInfo,
		Findings: []Finding{{Severity: SeverityCritical}},
	})

	if got := r.Categories[0].Severity; got != SeverityInfo {
		t.Errorf("AddCategoryResult: severity = %q, want %q (explicit value should survive)", got, SeverityInfo)
	}
}

func TestRubric_Evaluate_ComputesSeverityAsSafetyNet(t *testing.T) {
	r := NewRubric("prd", "prd.md")

	// Simulate a category appended directly to the slice, bypassing both
	// AddFinding and AddCategoryResult.
	r.Categories = append(r.Categories, CategoryResult{
		Category: "completeness",
		Score:    ScoreFail,
		Findings: []Finding{{Severity: SeverityMedium}},
	})

	r.Evaluate(nil)

	if got := r.Categories[0].Severity; got != SeverityMedium {
		t.Errorf("Evaluate safety net: severity = %q, want %q", got, SeverityMedium)
	}
}

func TestRubric_Evaluate_NoFindingsMeansNoSeverity(t *testing.T) {
	r := NewRubric("prd", "prd.md")
	r.AddCategoryResult(CategoryResult{Category: "clarity", Score: ScorePass})
	r.Evaluate(nil)

	if got := r.Categories[0].Severity; got != "" {
		t.Errorf("category with no findings: severity = %q, want empty", got)
	}
}
