package rubric

import (
	"encoding/json"
	"testing"
)

func TestCoverageSection(t *testing.T) {
	section := CoverageSection{
		Total:      10,
		Complete:   8,
		Percentage: 80,
		Missing:    []string{"item1", "item2"},
	}

	if section.Total != 10 {
		t.Errorf("expected Total 10, got %d", section.Total)
	}
	if section.Complete != 8 {
		t.Errorf("expected Complete 8, got %d", section.Complete)
	}
	if section.Percentage != 80 {
		t.Errorf("expected Percentage 80, got %d", section.Percentage)
	}
	if len(section.Missing) != 2 {
		t.Errorf("expected 2 missing items, got %d", len(section.Missing))
	}
}

func TestCoverageReport(t *testing.T) {
	t.Run("NewCoverageReport", func(t *testing.T) {
		cr := NewCoverageReport()
		if cr.Sections == nil {
			t.Error("expected Sections to be initialized")
		}
		if len(cr.Sections) != 0 {
			t.Error("expected empty Sections")
		}
	})

	t.Run("SetSection", func(t *testing.T) {
		cr := NewCoverageReport()
		cr.SetSection("components", 10, 8, []string{"card", "dialog"})

		section := cr.GetSection("components")
		if section.Total != 10 {
			t.Errorf("expected Total 10, got %d", section.Total)
		}
		if section.Complete != 8 {
			t.Errorf("expected Complete 8, got %d", section.Complete)
		}
		if section.Percentage != 80 {
			t.Errorf("expected Percentage 80, got %d", section.Percentage)
		}
	})

	t.Run("ComputeOverall", func(t *testing.T) {
		cr := NewCoverageReport()
		cr.SetSection("components", 10, 10, nil) // 100%
		cr.SetSection("foundations", 4, 2, nil)  // 50%
		cr.SetSection("patterns", 5, 5, nil)     // 100%

		overall := cr.ComputeOverall()
		// (100 + 50 + 100) / 3 = 83
		if overall != 83 {
			t.Errorf("expected overall 83, got %d", overall)
		}
	})

	t.Run("ComputeOverallWeighted", func(t *testing.T) {
		cr := NewCoverageReport()
		cr.SetSection("components", 10, 10, nil) // 100%
		cr.SetSection("foundations", 4, 2, nil)  // 50%

		weights := map[string]float64{
			"components":  2.0,
			"foundations": 1.0,
		}
		overall := cr.ComputeOverallWeighted(weights)
		// (100 * 2 + 50 * 1) / 3 = 83
		if overall != 83 {
			t.Errorf("expected overall 83, got %d", overall)
		}
	})

	t.Run("AllComplete", func(t *testing.T) {
		cr := NewCoverageReport()
		cr.SetSection("a", 10, 10, nil)
		cr.SetSection("b", 5, 5, nil)

		if !cr.AllComplete() {
			t.Error("expected AllComplete to be true")
		}

		cr.SetSection("c", 10, 5, nil)
		if cr.AllComplete() {
			t.Error("expected AllComplete to be false")
		}
	})

	t.Run("MeetsThreshold", func(t *testing.T) {
		cr := NewCoverageReport()
		cr.Overall = 80

		if !cr.MeetsThreshold(80) {
			t.Error("expected MeetsThreshold(80) to be true")
		}
		if !cr.MeetsThreshold(75) {
			t.Error("expected MeetsThreshold(75) to be true")
		}
		if cr.MeetsThreshold(85) {
			t.Error("expected MeetsThreshold(85) to be false")
		}
	})

	t.Run("SectionsAboveThreshold", func(t *testing.T) {
		cr := NewCoverageReport()
		cr.SetSection("high", 10, 9, nil) // 90%
		cr.SetSection("low", 10, 5, nil)  // 50%

		above := cr.SectionsAboveThreshold(80)
		if len(above) != 1 || above[0] != "high" {
			t.Errorf("expected [high], got %v", above)
		}
	})
}

func TestRubricCoverageExtension(t *testing.T) {
	t.Run("SetAndGetCoverage", func(t *testing.T) {
		r := NewRubric("test", "test.md")

		cr := NewCoverageReport()
		cr.SetSection("components", 10, 8, []string{"card"})
		cr.ComputeOverall()

		r.SetCoverage(cr)

		retrieved := r.GetCoverage()
		if retrieved == nil {
			t.Fatal("expected coverage to be retrieved")
		}
		if retrieved.Overall != 80 {
			t.Errorf("expected Overall 80, got %d", retrieved.Overall)
		}
	})

	t.Run("JSONRoundTrip", func(t *testing.T) {
		r := NewRubric("test", "test.md")

		cr := NewCoverageReport()
		cr.SetSection("components", 10, 8, []string{"card", "dialog"})
		cr.ComputeOverall()

		r.SetCoverage(cr)

		// Marshal to JSON
		data, err := json.Marshal(r)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		// Unmarshal back
		var r2 Rubric
		if err := json.Unmarshal(data, &r2); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		// Check coverage survived
		retrieved := r2.GetCoverage()
		if retrieved == nil {
			t.Fatal("expected coverage after JSON round-trip")
		}
		if retrieved.Overall != 80 {
			t.Errorf("expected Overall 80, got %d", retrieved.Overall)
		}

		section := retrieved.GetSection("components")
		if section.Total != 10 {
			t.Errorf("expected Total 10, got %d", section.Total)
		}
		if len(section.Missing) != 2 {
			t.Errorf("expected 2 missing items, got %d", len(section.Missing))
		}
	})
}

func TestRubricExtensions(t *testing.T) {
	t.Run("SetAndGetExtension", func(t *testing.T) {
		r := NewRubric("test", "test.md")

		r.SetExtension("custom", "value")
		if r.GetExtension("custom") != "value" {
			t.Error("expected extension value")
		}
	})

	t.Run("HasExtension", func(t *testing.T) {
		r := NewRubric("test", "test.md")

		if r.HasExtension("missing") {
			t.Error("expected HasExtension to return false")
		}

		r.SetExtension("present", true)
		if !r.HasExtension("present") {
			t.Error("expected HasExtension to return true")
		}
	})

	t.Run("NilExtensions", func(t *testing.T) {
		r := &Rubric{}

		if r.GetExtension("anything") != nil {
			t.Error("expected nil from empty rubric")
		}
		if r.HasExtension("anything") {
			t.Error("expected false from empty rubric")
		}
	})
}

func TestMinIntScoreCriteria(t *testing.T) {
	t.Run("PassesWhenScoreMeetsMinimum", func(t *testing.T) {
		r := NewRubric("test", "test.md")
		r.PassCriteria.MinIntScore = ScoreGood // 4

		// Add a passing category
		r.AddCategoryResult(CategoryResult{
			Category: "quality",
			Score:    ScorePass,
			IntScore: ScoreGood,
		})

		r.IntScore = ScoreGood
		r.Evaluate(nil)

		if !r.Pass {
			t.Error("expected pass when score meets minimum")
		}
	})

	t.Run("FailsWhenScoreBelowMinimum", func(t *testing.T) {
		r := NewRubric("test", "test.md")
		r.PassCriteria.MinIntScore = ScoreGood // 4

		// Add a category with lower score
		r.AddCategoryResult(CategoryResult{
			Category: "quality",
			Score:    ScorePartial,
			IntScore: ScoreAcceptable, // 3
		})

		r.Evaluate(nil)

		if r.Pass {
			t.Error("expected fail when score below minimum")
		}
		if r.Decision.Status != DecisionFail {
			t.Errorf("expected DecisionFail, got %s", r.Decision.Status)
		}
	})

	t.Run("NoMinIntScoreCheck", func(t *testing.T) {
		r := NewRubric("test", "test.md")
		// MinIntScore is 0 (disabled)

		r.AddCategoryResult(CategoryResult{
			Category: "quality",
			Score:    ScorePartial,
			IntScore: ScoreAcceptable,
		})

		r.Evaluate(nil)

		// Should pass because MinIntScore is not set
		if !r.Pass {
			t.Error("expected pass when MinIntScore is not set")
		}
	})
}
