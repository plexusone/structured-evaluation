package rubric

import (
	"encoding/json"
	"testing"
)

func TestCompositeCategoryRoundTrips(t *testing.T) {
	rs := &RubricSet{
		Name:    "Test Rubric",
		Version: "1.0",
		PassCriteria: RubricPassCriteria{
			ScoreThresholds: &ScoreThresholds{Pass: 80, Partial: 60},
		},
		Categories: []Category{
			{
				Name:        "Discovery Clarity",
				Weight:      20,
				Description: "Row 1 quality",
				Criteria: []Criterion{
					{
						ID:     "DC1",
						Name:   "Problem Definition",
						Weight: 8,
						Pass: CriterionLevel{
							Description: "Specific, evidence-based problem statement",
							Indicators:  []string{"Quantified impact", "Evidence from research"},
						},
						Fail: CriterionLevel{Description: "Unclear or missing"},
					},
				},
			},
		},
	}

	if !rs.Categories[0].IsComposite() {
		t.Fatal("category with criteria should be composite")
	}

	data, err := json.Marshal(rs)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back RubricSet
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !back.Categories[0].IsComposite() {
		t.Error("composite state lost in round-trip")
	}
	if back.PassCriteria.ScoreThresholds == nil || back.PassCriteria.ScoreThresholds.Pass != 80 {
		t.Error("score thresholds lost in round-trip")
	}
	got := back.Categories[0].Criteria[0].Pass.Indicators
	if len(got) != 2 || got[0] != "Quantified impact" {
		t.Errorf("indicators lost or wrong: %v", got)
	}
}

func TestSimpleCategoryIsNotComposite(t *testing.T) {
	// A category scored directly via Scale (the existing form) is not composite.
	c := Category{Name: "Clarity", Scale: Scale{Type: ScaleTypeCategorical}}
	if c.IsComposite() {
		t.Error("category without criteria should not be composite")
	}
}
