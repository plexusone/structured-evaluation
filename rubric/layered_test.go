package rubric

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// TestLayeredFieldsRoundTrip verifies Class, Blocking, Evaluation on Category
// and Criterion, and RubricSet.JudgeInstructions, survive a YAML round-trip.
func TestLayeredFieldsRoundTrip(t *testing.T) {
	rs := RubricSet{
		ID:      "test-v1",
		Name:    "Test",
		Version: "1.0.0",
		JudgeInstructions: []string{
			"Cite the relevant section and requirement IDs",
			"Do not reward length",
		},
		Categories: []Category{
			{
				ID:         "leadership",
				Name:       "Leadership",
				Class:      ClassLeadershipPrinciple,
				Evaluation: EvalMethodSemantic,
				Scale:      Scale{Type: ScaleTypeCategorical, Options: []ScaleOption{{Value: "pass", Label: "Pass", Criteria: []string{"x"}}}},
			},
			{
				ID:         "readiness",
				Name:       "Readiness",
				Class:      ClassImplementationReadiness,
				Blocking:   true,
				Evaluation: EvalMethodDeterministic,
				Scale:      Scale{Type: ScaleTypeCategorical, Options: []ScaleOption{{Value: "pass", Label: "Pass", Criteria: []string{"x"}}}},
				Criteria: []Criterion{
					{
						ID:         "trace",
						Name:       "Traceability",
						Class:      ClassImplementationReadiness,
						Blocking:   true,
						Evaluation: EvalMethodDeterministic,
						Pass:       CriterionLevel{Description: "traced"},
						Fail:       CriterionLevel{Description: "untraced"},
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(&rs)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got RubricSet
	if err := yaml.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(got.JudgeInstructions) != 2 {
		t.Fatalf("JudgeInstructions = %v, want 2 entries", got.JudgeInstructions)
	}

	lp := got.GetCategory("leadership")
	if lp == nil {
		t.Fatal("leadership category missing after round-trip")
	}
	if lp.Class != ClassLeadershipPrinciple || lp.Evaluation != EvalMethodSemantic {
		t.Errorf("leadership category class/evaluation = %q/%q, want %q/%q", lp.Class, lp.Evaluation, ClassLeadershipPrinciple, EvalMethodSemantic)
	}
	if lp.Blocking {
		t.Error("leadership category should not be blocking")
	}

	ir := got.GetCategory("readiness")
	if ir == nil {
		t.Fatal("readiness category missing after round-trip")
	}
	if ir.Class != ClassImplementationReadiness || !ir.Blocking || ir.Evaluation != EvalMethodDeterministic {
		t.Errorf("readiness category = %+v, want implementation_readiness/blocking/deterministic", ir)
	}
	if len(ir.Criteria) != 1 || ir.Criteria[0].Class != ClassImplementationReadiness || !ir.Criteria[0].Blocking {
		t.Errorf("readiness criterion = %+v, want implementation_readiness/blocking", ir.Criteria)
	}
}

// TestLegacyRubricParsesUnchanged verifies a v0.13.0-shaped rubric (no layer
// fields at all) parses identically: the new fields default to their zero
// values and every pre-existing field is unaffected (TRD-003).
func TestLegacyRubricParsesUnchanged(t *testing.T) {
	const legacyYAML = `
id: prd-rubric
name: Enterprise PRD Rubric
version: "1.0"
description: Comprehensive PRD evaluation
evaluationType: analytic
passCriteria:
  minCategoriesPassing: all_required
  maxFindingsSeverity:
    critical: 0
    high: 0
    medium: 2
    low: -1
categories:
  - id: problem-definition
    name: Problem Definition
    description: Is the problem clearly defined?
    weight: 2
    required: true
    scale:
      type: categorical
      options:
        - value: pass
          label: Pass
          criteria:
            - Problem clearly stated
`
	var rs RubricSet
	if err := yaml.Unmarshal([]byte(legacyYAML), &rs); err != nil {
		t.Fatalf("Unmarshal legacy rubric: %v", err)
	}

	if rs.ID != "prd-rubric" || rs.Name != "Enterprise PRD Rubric" {
		t.Errorf("legacy fields not preserved: id=%q name=%q", rs.ID, rs.Name)
	}
	if len(rs.Categories) != 1 || rs.Categories[0].ID != "problem-definition" {
		t.Fatalf("legacy categories not preserved: %+v", rs.Categories)
	}

	cat := rs.Categories[0]
	if cat.Class != "" || cat.Evaluation != "" || cat.Blocking {
		t.Errorf("legacy category should have zero-value layer fields, got class=%q evaluation=%q blocking=%v", cat.Class, cat.Evaluation, cat.Blocking)
	}
	if len(rs.JudgeInstructions) != 0 {
		t.Errorf("legacy rubric should have no JudgeInstructions, got %v", rs.JudgeInstructions)
	}

	if issues := rs.Validate(); len(issues) != 0 {
		t.Errorf("legacy rubric should validate cleanly, got issues: %v", issues)
	}
}

// TestValidate_LeadershipPrincipleMustNotBlock guards INV-3: a
// leadership_principle category or criterion must never be blocking, because
// advisory judgment cannot gate implementation.
func TestValidate_LeadershipPrincipleMustNotBlock(t *testing.T) {
	t.Run("category", func(t *testing.T) {
		rs := NewRubricSet("test-v1", "Test", "1.0.0")
		cat := NewCategory("think_big", "Think Big", "LP judgment").
			WithPassPartialFail([]string{"Good"}, []string{"OK"}, []string{"Bad"})
		cat.Class = ClassLeadershipPrinciple
		cat.Blocking = true
		rs.AddCategory(*cat)

		issues := rs.Validate()
		if !containsSubstring(issues, "must not be blocking") {
			t.Errorf("expected a 'must not be blocking' issue, got: %v", issues)
		}
	})

	t.Run("criterion", func(t *testing.T) {
		rs := NewRubricSet("test-v1", "Test", "1.0.0")
		cat := NewCategory("mixed", "Mixed", "Composite category")
		cat.Criteria = []Criterion{
			{
				ID:       "earn_trust",
				Name:     "Earn Trust",
				Class:    ClassLeadershipPrinciple,
				Blocking: true,
				Pass:     CriterionLevel{Description: "trusted"},
				Fail:     CriterionLevel{Description: "untrusted"},
			},
		}
		rs.AddCategory(*cat)

		issues := rs.Validate()
		if !containsSubstring(issues, "must not be blocking") {
			t.Errorf("expected a 'must not be blocking' issue, got: %v", issues)
		}
	})

	t.Run("non-blocking leadership principle is valid", func(t *testing.T) {
		rs := NewRubricSet("test-v1", "Test", "1.0.0")
		cat := NewCategory("think_big", "Think Big", "LP judgment").
			WithPassPartialFail([]string{"Good"}, []string{"OK"}, []string{"Bad"})
		cat.Class = ClassLeadershipPrinciple
		cat.Blocking = false
		rs.AddCategory(*cat)

		issues := rs.Validate()
		if containsSubstring(issues, "must not be blocking") {
			t.Errorf("non-blocking leadership_principle should not raise the invariant issue, got: %v", issues)
		}
	})
}

func containsSubstring(issues []string, substr string) bool {
	for _, issue := range issues {
		if len(issue) >= len(substr) {
			for i := 0; i+len(substr) <= len(issue); i++ {
				if issue[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}
