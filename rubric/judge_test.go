package rubric

import (
	"encoding/json"
	"testing"
)

func TestNewJudgeMetadata_SetsSemanticMethod(t *testing.T) {
	j := NewJudgeMetadata("claude-3-opus-20240229")
	if j.Method != EvalMethodSemantic {
		t.Fatalf("expected Method=semantic, got %q", j.Method)
	}
	if j.Model != "claude-3-opus-20240229" {
		t.Fatalf("expected Model to be set, got %q", j.Model)
	}
	if issues := j.Validate(); len(issues) != 0 {
		t.Fatalf("expected no validation issues, got %v", issues)
	}
}

func TestNewDeterministicJudgeMetadata(t *testing.T) {
	j := NewDeterministicJudgeMetadata("golangci-lint")
	if j.Method != EvalMethodDeterministic {
		t.Fatalf("expected Method=deterministic, got %q", j.Method)
	}
	if j.CheckID != "golangci-lint" {
		t.Fatalf("expected CheckID to be set, got %q", j.CheckID)
	}
	j.WithCheckVersion("v2.7.0")
	if j.CheckVersion != "v2.7.0" {
		t.Fatalf("expected CheckVersion to be set, got %q", j.CheckVersion)
	}
	if issues := j.Validate(); len(issues) != 0 {
		t.Fatalf("expected no validation issues, got %v", issues)
	}
}

func TestNewHumanJudgeMetadata(t *testing.T) {
	j := NewHumanJudgeMetadata("jsmith")
	if j.Method != EvalMethodHuman {
		t.Fatalf("expected Method=human, got %q", j.Method)
	}
	if j.ReviewerID != "jsmith" {
		t.Fatalf("expected ReviewerID to be set, got %q", j.ReviewerID)
	}
	j.WithReviewerName("Jane Smith")
	if j.ReviewerName != "Jane Smith" {
		t.Fatalf("expected ReviewerName to be set, got %q", j.ReviewerName)
	}
	if issues := j.Validate(); len(issues) != 0 {
		t.Fatalf("expected no validation issues, got %v", issues)
	}
}

func TestJudgeMetadata_Validate_MissingRequiredField(t *testing.T) {
	tests := []struct {
		name string
		j    *JudgeMetadata
		want string
	}{
		{"semantic missing model", &JudgeMetadata{Method: EvalMethodSemantic}, "semantic judge requires Model"},
		{"empty method missing model", &JudgeMetadata{}, "semantic judge requires Model"},
		{"deterministic missing check id", &JudgeMetadata{Method: EvalMethodDeterministic}, "deterministic judge requires CheckID"},
		{"human missing reviewer id", &JudgeMetadata{Method: EvalMethodHuman}, "human judge requires ReviewerID"},
		{"unknown method", &JudgeMetadata{Method: EvaluationMethod("bogus")}, "unknown Method: bogus"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := tt.j.Validate()
			if len(issues) != 1 || issues[0] != tt.want {
				t.Fatalf("expected issues [%q], got %v", tt.want, issues)
			}
		})
	}
}

// A JudgeMetadata JSON payload shaped exactly like what v0.14.0 (before
// Method/CheckID/ReviewerID existed) would have produced must still
// unmarshal cleanly and validate as a semantic judge — the whole point of
// treating an empty Method as EvalMethodSemantic.
func TestJudgeMetadata_BackwardCompatible_PreMethodJSON(t *testing.T) {
	legacyJSON := `{
		"judge_id": "j-1",
		"model": "gpt-4-turbo",
		"model_provider": "openai",
		"evaluated_at": "2026-08-15T00:00:00Z"
	}`
	var j JudgeMetadata
	if err := json.Unmarshal([]byte(legacyJSON), &j); err != nil {
		t.Fatalf("unmarshal legacy JudgeMetadata: %v", err)
	}
	if j.Method != "" {
		t.Fatalf("expected Method to be empty for legacy JSON, got %q", j.Method)
	}
	if issues := j.Validate(); len(issues) != 0 {
		t.Fatalf("expected legacy JudgeMetadata to validate as semantic, got issues %v", issues)
	}
}

func TestJudgeMetadata_JSONRoundTrip_DeterministicAndHuman(t *testing.T) {
	original := NewDeterministicJudgeMetadata("schema-validate").WithCheckVersion("v1.2.3")
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var round JudgeMetadata
	if err := json.Unmarshal(data, &round); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if round.Method != EvalMethodDeterministic || round.CheckID != "schema-validate" || round.CheckVersion != "v1.2.3" {
		t.Fatalf("round-trip mismatch: got %+v", round)
	}
	// Model must not be present for a deterministic judge — it was made
	// omitempty precisely so non-LLM judges don't carry an empty model field.
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}
	if _, present := raw["model"]; present {
		t.Fatalf("expected no \"model\" key for a deterministic judge, got %v", raw)
	}
}

func TestValidEvaluationMethodValues(t *testing.T) {
	got := ValidEvaluationMethodValues()
	want := []string{"deterministic", "semantic", "human"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
