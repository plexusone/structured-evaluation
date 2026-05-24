package evaluation

import "testing"

func TestLikertToCategorical(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		config   *LikertConfig
		expected ScoreValue
	}{
		// Default 1-5 scale (pass=4, partial=3)
		{"5 is pass", 5, nil, ScorePass},
		{"4 is pass", 4, nil, ScorePass},
		{"3 is partial", 3, nil, ScorePartial},
		{"2 is fail", 2, nil, ScoreFail},
		{"1 is fail", 1, nil, ScoreFail},

		// Custom thresholds
		{
			"custom pass threshold",
			3,
			&LikertConfig{Min: 1, Max: 5, PassThreshold: intPtr(3), PartialThreshold: intPtr(2)},
			ScorePass,
		},
		{
			"custom partial threshold",
			2,
			&LikertConfig{Min: 1, Max: 5, PassThreshold: intPtr(4), PartialThreshold: intPtr(2)},
			ScorePartial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LikertToCategorical(tt.score, tt.config)
			if result != tt.expected {
				t.Errorf("LikertToCategorical(%d) = %v, want %v", tt.score, result, tt.expected)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func TestWithLikert5(t *testing.T) {
	cat := NewCategory("test", "Test Category", "A test category").
		WithLikert5(StandardLikert5Anchors())

	if cat.Scale.Type != ScaleTypeLikert {
		t.Errorf("Scale.Type = %v, want %v", cat.Scale.Type, ScaleTypeLikert)
	}

	if cat.Scale.LikertConfig == nil {
		t.Fatal("LikertConfig is nil")
	}

	if cat.Scale.LikertConfig.Min != 1 {
		t.Errorf("Min = %d, want 1", cat.Scale.LikertConfig.Min)
	}

	if cat.Scale.LikertConfig.Max != 5 {
		t.Errorf("Max = %d, want 5", cat.Scale.LikertConfig.Max)
	}

	if len(cat.Scale.LikertConfig.Anchors) != 5 {
		t.Errorf("Anchors length = %d, want 5", len(cat.Scale.LikertConfig.Anchors))
	}
}

func TestStandardLikert5Anchors(t *testing.T) {
	anchors := StandardLikert5Anchors()

	if len(anchors) != 5 {
		t.Fatalf("Expected 5 anchors, got %d", len(anchors))
	}

	// Check values are 5, 4, 3, 2, 1
	expectedValues := []int{5, 4, 3, 2, 1}
	for i, anchor := range anchors {
		if anchor.Value != expectedValues[i] {
			t.Errorf("Anchor[%d].Value = %d, want %d", i, anchor.Value, expectedValues[i])
		}
		if anchor.Label == "" {
			t.Errorf("Anchor[%d].Label is empty", i)
		}
	}
}

func TestNewCategoryResultFromLikert(t *testing.T) {
	config := &LikertConfig{
		Min:              1,
		Max:              5,
		PassThreshold:    intPtr(4),
		PartialThreshold: intPtr(3),
	}

	tests := []struct {
		score    int
		expected ScoreValue
	}{
		{5, ScorePass},
		{4, ScorePass},
		{3, ScorePartial},
		{2, ScoreFail},
		{1, ScoreFail},
	}

	for _, tt := range tests {
		t.Run("score_"+itoa(tt.score), func(t *testing.T) {
			result := NewCategoryResultFromLikert("test", tt.score, config, "test reasoning")

			if result.Score != tt.expected {
				t.Errorf("Score = %v, want %v", result.Score, tt.expected)
			}

			if result.NumericScore == nil {
				t.Fatal("NumericScore is nil")
			}

			if *result.NumericScore != float64(tt.score) {
				t.Errorf("NumericScore = %v, want %v", *result.NumericScore, float64(tt.score))
			}
		})
	}
}

func TestNewCategoryResultWithNumeric(t *testing.T) {
	result := NewCategoryResultWithNumeric("test", ScorePass, 4.5, "good work")

	if result.Score != ScorePass {
		t.Errorf("Score = %v, want %v", result.Score, ScorePass)
	}

	if result.NumericScore == nil {
		t.Fatal("NumericScore is nil")
	}

	if *result.NumericScore != 4.5 {
		t.Errorf("NumericScore = %v, want 4.5", *result.NumericScore)
	}

	if result.Reasoning != "good work" {
		t.Errorf("Reasoning = %v, want 'good work'", result.Reasoning)
	}
}

func TestCategoryResultNumericHelpers(t *testing.T) {
	t.Run("without numeric score", func(t *testing.T) {
		result := NewCategoryResult("test", ScorePass, "reasoning")

		if result.HasNumericScore() {
			t.Error("HasNumericScore() = true, want false")
		}

		if result.GetNumericScore() != 0 {
			t.Errorf("GetNumericScore() = %v, want 0", result.GetNumericScore())
		}
	})

	t.Run("with numeric score", func(t *testing.T) {
		result := NewCategoryResult("test", ScorePass, "reasoning").
			SetNumericScore(4.5)

		if !result.HasNumericScore() {
			t.Error("HasNumericScore() = false, want true")
		}

		if result.GetNumericScore() != 4.5 {
			t.Errorf("GetNumericScore() = %v, want 4.5", result.GetNumericScore())
		}
	})
}

func TestRubricValidateLikert(t *testing.T) {
	t.Run("valid likert scale", func(t *testing.T) {
		rs := NewRubricSet("test", "Test", "1.0")
		rs.AddCategory(*NewCategory("cat1", "Category 1", "Desc").
			WithLikert5(StandardLikert5Anchors()))

		issues := rs.Validate()
		if len(issues) > 0 {
			t.Errorf("Validation failed: %v", issues)
		}
	})

	t.Run("likert scale without config", func(t *testing.T) {
		rs := NewRubricSet("test", "Test", "1.0")
		cat := NewCategory("cat1", "Category 1", "Desc")
		cat.Scale.Type = ScaleTypeLikert
		cat.Scale.LikertConfig = nil
		rs.AddCategory(*cat)

		issues := rs.Validate()
		hasLikertError := false
		for _, issue := range issues {
			if issue == "category cat1: likert scale requires LikertConfig" {
				hasLikertError = true
			}
		}
		if !hasLikertError {
			t.Error("Expected validation error for missing LikertConfig")
		}
	})

	t.Run("likert scale with invalid range", func(t *testing.T) {
		rs := NewRubricSet("test", "Test", "1.0")
		cat := NewCategory("cat1", "Category 1", "Desc")
		cat.Scale.Type = ScaleTypeLikert
		cat.Scale.LikertConfig = &LikertConfig{Min: 5, Max: 1}
		rs.AddCategory(*cat)

		issues := rs.Validate()
		hasRangeError := false
		for _, issue := range issues {
			if issue == "category cat1: likert scale min must be less than max" {
				hasRangeError = true
			}
		}
		if !hasRangeError {
			t.Error("Expected validation error for invalid range")
		}
	})
}
