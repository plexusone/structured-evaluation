package evaluation

import "testing"

func TestNewRubric(t *testing.T) {
	r := NewRubric("quality", "Measures output quality")

	if r.Category != "quality" {
		t.Errorf("Expected category 'quality', got %s", r.Category)
	}
	if r.Description != "Measures output quality" {
		t.Errorf("Expected description 'Measures output quality', got %s", r.Description)
	}
}

func TestRubric_AddAnchor(t *testing.T) {
	r := NewRubric("quality", "Test").
		AddAnchor(10, "Excellent", "Perfect output", "Criterion 1", "Criterion 2").
		AddAnchor(5, "Average", "Acceptable output")

	if len(r.Anchors) != 2 {
		t.Errorf("Expected 2 anchors, got %d", len(r.Anchors))
	}

	if r.Anchors[0].Score != 10 {
		t.Errorf("Expected score 10, got %f", r.Anchors[0].Score)
	}
	if r.Anchors[0].Label != "Excellent" {
		t.Errorf("Expected label 'Excellent', got %s", r.Anchors[0].Label)
	}
	if len(r.Anchors[0].Criteria) != 2 {
		t.Errorf("Expected 2 criteria, got %d", len(r.Anchors[0].Criteria))
	}
}

func TestRubric_AddRangeAnchor(t *testing.T) {
	r := NewRubric("quality", "Test").
		AddRangeAnchor(8, 10, "Excellent", "Top tier")

	if r.Anchors[0].MinScore != 8 {
		t.Errorf("Expected min score 8, got %f", r.Anchors[0].MinScore)
	}
	if r.Anchors[0].MaxScore != 10 {
		t.Errorf("Expected max score 10, got %f", r.Anchors[0].MaxScore)
	}
}

func TestRubric_GetAnchorForScore(t *testing.T) {
	r := NewRubric("quality", "Test").
		AddRangeAnchor(8, 10, "Excellent", "Top tier").
		AddRangeAnchor(5, 7.9, "Good", "Acceptable").
		AddRangeAnchor(0, 4.9, "Poor", "Needs work")

	tests := []struct {
		score    float64
		expected string
	}{
		{10, "Excellent"},
		{9, "Excellent"},
		{8, "Excellent"},
		{7.5, "Good"},
		{5, "Good"},
		{4, "Poor"},
		{0, "Poor"},
	}

	for _, tt := range tests {
		anchor := r.GetAnchorForScore(tt.score)
		if anchor == nil {
			t.Errorf("No anchor found for score %f", tt.score)
			continue
		}
		if anchor.Label != tt.expected {
			t.Errorf("Score %f: expected label %s, got %s", tt.score, tt.expected, anchor.Label)
		}
	}
}

func TestDefaultPRDRubricSet(t *testing.T) {
	rubricSet := DefaultPRDRubricSet()

	if rubricSet.ID != "prd-evaluation-v1" {
		t.Errorf("Expected ID 'prd-evaluation-v1', got %s", rubricSet.ID)
	}

	if len(rubricSet.Rubrics) < 3 {
		t.Errorf("Expected at least 3 rubrics, got %d", len(rubricSet.Rubrics))
	}

	// Check that each rubric has anchors
	for _, r := range rubricSet.Rubrics {
		if len(r.Anchors) == 0 {
			t.Errorf("Rubric %s has no anchors", r.Category)
		}
	}
}
