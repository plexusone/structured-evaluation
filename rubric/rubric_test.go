package rubric

import "testing"

func TestNewRubricSet(t *testing.T) {
	rs := NewRubricSet("test-v1", "Test Rubric", "1.0.0")

	if rs.ID != "test-v1" {
		t.Errorf("Expected ID 'test-v1', got %s", rs.ID)
	}
	if rs.Name != "Test Rubric" {
		t.Errorf("Expected name 'Test Rubric', got %s", rs.Name)
	}
	if rs.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", rs.Version)
	}
	if rs.EvaluationType != EvaluationTypeAnalytic {
		t.Errorf("Expected evaluation type 'analytic', got %s", rs.EvaluationType)
	}
}

func TestNewCategory(t *testing.T) {
	cat := NewCategory("quality", "Quality", "Measures output quality")

	if cat.ID != "quality" {
		t.Errorf("Expected ID 'quality', got %s", cat.ID)
	}
	if cat.Name != "Quality" {
		t.Errorf("Expected name 'Quality', got %s", cat.Name)
	}
	if cat.Scale.Type != ScaleTypeCategorical {
		t.Errorf("Expected scale type 'categorical', got %s", cat.Scale.Type)
	}
}

func TestCategory_WithPassPartialFail(t *testing.T) {
	cat := NewCategory("quality", "Quality", "Test").
		WithPassPartialFail(
			[]string{"Excellent output", "No errors"},
			[]string{"Good output", "Minor issues"},
			[]string{"Poor output", "Major errors"},
		)

	if len(cat.Scale.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(cat.Scale.Options))
	}

	if cat.Scale.Options[0].Value != "pass" {
		t.Errorf("Expected first option value 'pass', got %s", cat.Scale.Options[0].Value)
	}
	if cat.Scale.Options[1].Value != "partial" {
		t.Errorf("Expected second option value 'partial', got %s", cat.Scale.Options[1].Value)
	}
	if cat.Scale.Options[2].Value != "fail" {
		t.Errorf("Expected third option value 'fail', got %s", cat.Scale.Options[2].Value)
	}
}

func TestCategory_AddOption(t *testing.T) {
	cat := NewCategory("quality", "Quality", "Test").
		AddOption("pass", "Pass", "Criterion 1", "Criterion 2").
		AddOption("fail", "Fail", "Issue 1")

	if len(cat.Scale.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(cat.Scale.Options))
	}

	if len(cat.Scale.Options[0].Criteria) != 2 {
		t.Errorf("Expected 2 criteria for pass, got %d", len(cat.Scale.Options[0].Criteria))
	}
}

func TestCategory_GetOptionForValue(t *testing.T) {
	cat := NewCategory("quality", "Quality", "Test").
		WithPassPartialFail(
			[]string{"Excellent"},
			[]string{"Good"},
			[]string{"Poor"},
		)

	passOpt := cat.GetOptionForValue("pass")
	if passOpt == nil {
		t.Error("Expected to find pass option")
	} else if passOpt.Label != "Pass" {
		t.Errorf("Expected label 'Pass', got %s", passOpt.Label)
	}

	nonExistent := cat.GetOptionForValue("nonexistent")
	if nonExistent != nil {
		t.Error("Expected nil for nonexistent option")
	}
}

func TestRubricSet_AddCategory(t *testing.T) {
	rs := NewRubricSet("test-v1", "Test", "1.0.0")

	cat1 := NewCategory("cat1", "Category 1", "First category")
	cat2 := NewCategory("cat2", "Category 2", "Second category")

	rs.AddCategory(*cat1).AddCategory(*cat2)

	if len(rs.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(rs.Categories))
	}
}

func TestRubricSet_GetCategory(t *testing.T) {
	rs := NewRubricSet("test-v1", "Test", "1.0.0")
	cat := NewCategory("quality", "Quality", "Test")
	rs.AddCategory(*cat)

	found := rs.GetCategory("quality")
	if found == nil {
		t.Error("Expected to find category 'quality'")
	}

	notFound := rs.GetCategory("nonexistent")
	if notFound != nil {
		t.Error("Expected nil for nonexistent category")
	}
}

func TestRubricSet_GetRequiredCategories(t *testing.T) {
	rs := NewRubricSet("test-v1", "Test", "1.0.0")

	cat1 := NewCategory("cat1", "Category 1", "Required").SetRequired(true)
	cat2 := NewCategory("cat2", "Category 2", "Optional")
	cat3 := NewCategory("cat3", "Category 3", "Required").SetRequired(true)

	rs.AddCategory(*cat1).AddCategory(*cat2).AddCategory(*cat3)

	required := rs.GetRequiredCategories()
	if len(required) != 2 {
		t.Errorf("Expected 2 required categories, got %d", len(required))
	}
}

func TestRubricSet_Validate(t *testing.T) {
	// Valid rubric
	rs := NewRubricSet("test-v1", "Test", "1.0.0")
	cat := NewCategory("quality", "Quality", "Test").
		WithPassPartialFail([]string{"Good"}, []string{"OK"}, []string{"Bad"})
	rs.AddCategory(*cat)

	issues := rs.Validate()
	if len(issues) > 0 {
		t.Errorf("Expected no issues, got: %v", issues)
	}

	// Invalid rubric - missing ID
	invalid := &RubricSet{Name: "Test", Version: "1.0"}
	issues = invalid.Validate()
	if len(issues) == 0 {
		t.Error("Expected validation issues for missing ID")
	}

	// Invalid rubric - categorical without options
	rs2 := NewRubricSet("test-v1", "Test", "1.0.0")
	cat2 := NewCategory("empty", "Empty", "No options")
	rs2.AddCategory(*cat2)

	issues = rs2.Validate()
	if len(issues) == 0 {
		t.Error("Expected validation issues for categorical without options")
	}
}

func TestRubricSet_ToJSON(t *testing.T) {
	rs := NewRubricSet("test-v1", "Test Rubric", "1.0.0")
	cat := NewCategory("quality", "Quality", "Test").
		SetRequired(true).
		SetWeight(1.5).
		WithPassPartialFail(
			[]string{"Excellent output"},
			[]string{"Acceptable output"},
			[]string{"Poor output"},
		)
	rs.AddCategory(*cat)

	jsonData, err := rs.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize rubric: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Expected non-empty JSON")
	}

	// Verify it contains expected fields
	jsonStr := string(jsonData)
	if !contains(jsonStr, "test-v1") {
		t.Error("JSON should contain rubric ID")
	}
	if !contains(jsonStr, "categorical") {
		t.Error("JSON should contain scale type")
	}
	if !contains(jsonStr, "pass") {
		t.Error("JSON should contain pass option")
	}
}

func TestCategory_WithChecklist(t *testing.T) {
	cat := NewCategory("completeness", "Completeness", "Test").
		WithChecklist(
			[]string{"overview", "mitigations"},
			[]string{"tldr", "detection"},
			&ChecklistThreshold{Required: "all", Optional: 1},
		)

	if cat.Scale.Type != ScaleTypeChecklist {
		t.Errorf("Expected scale type 'checklist', got %s", cat.Scale.Type)
	}
	if len(cat.Scale.RequiredItems) != 2 {
		t.Errorf("Expected 2 required items, got %d", len(cat.Scale.RequiredItems))
	}
	if len(cat.Scale.OptionalItems) != 2 {
		t.Errorf("Expected 2 optional items, got %d", len(cat.Scale.OptionalItems))
	}
	if cat.Scale.PassingThreshold.Optional != 1 {
		t.Errorf("Expected optional threshold 1, got %d", cat.Scale.PassingThreshold.Optional)
	}
}

func TestCategory_SetExamples(t *testing.T) {
	cat := NewCategory("quality", "Quality", "Test").
		SetExamples(&CategoryExamples{
			Pass: &Example{
				Excerpt:   "This is excellent work.",
				Reasoning: "Clear, concise, accurate.",
			},
			Fail: &Example{
				Excerpt:   "This is poor.",
				Reasoning: "Vague, incomplete.",
			},
		})

	if cat.Examples == nil {
		t.Error("Expected examples to be set")
	}
	if cat.Examples.Pass == nil {
		t.Error("Expected pass example to be set")
	}
	if cat.Examples.Pass.Reasoning != "Clear, concise, accurate." {
		t.Errorf("Unexpected reasoning: %s", cat.Examples.Pass.Reasoning)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
