package claims

import "testing"

func verifiedClaim(id string, category ClaimCategory, related ...string) Claim {
	return Claim{
		ID:              id,
		Category:        category,
		Verdict:         VerdictVerified,
		RelatedClaimIDs: related,
	}
}

func TestIsSufficientlyCorroborated(t *testing.T) {
	single := verifiedClaim("c1", ClaimStatistical)
	double := verifiedClaim("c2", ClaimStatistical, "c2-corroborating-1")

	tests := []struct {
		name     string
		claim    Claim
		criteria ClaimsCriteria
		want     bool
	}{
		{"disabled (0) always sufficient", single, ClaimsCriteria{MinCorroboratingSources: 0}, true},
		{"disabled (1) always sufficient", single, ClaimsCriteria{MinCorroboratingSources: 1}, true},
		{"threshold 2, single source fails", single, ClaimsCriteria{MinCorroboratingSources: 2}, false},
		{"threshold 2, one related source passes", double, ClaimsCriteria{MinCorroboratingSources: 2}, true},
		{"threshold 3, one related source still fails", double, ClaimsCriteria{MinCorroboratingSources: 3}, false},
		{
			"category out of scope is always sufficient", single,
			ClaimsCriteria{MinCorroboratingSources: 2, CorroborationCategories: []ClaimCategory{ClaimTimeline}},
			true,
		},
		{
			"category in scope is checked", single,
			ClaimsCriteria{MinCorroboratingSources: 2, CorroborationCategories: []ClaimCategory{ClaimStatistical}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSufficientlyCorroborated(tt.claim, tt.criteria); got != tt.want {
				t.Errorf("IsSufficientlyCorroborated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateClaims_CorroborationDisabledByDefault(t *testing.T) {
	claims := []Claim{verifiedClaim("c1", ClaimStatistical)}
	decision := EvaluateClaims(claims, DefaultClaimsCriteria())
	if decision.Status != ClaimsDecisionPass {
		t.Errorf("expected pass with corroboration disabled, got %s: %s", decision.Status, decision.Rationale)
	}
}

func TestEvaluateClaims_InsufficientCorroborationIsConditionalWhenAllowed(t *testing.T) {
	claims := []Claim{verifiedClaim("c1", ClaimStatistical)}
	criteria := DefaultClaimsCriteria()
	criteria.MinCorroboratingSources = 2
	criteria.AllowNeedsReview = true

	decision := EvaluateClaims(claims, criteria)
	if decision.Status != ClaimsDecisionConditional || !decision.Passed {
		t.Errorf("expected conditional pass, got %s (passed=%v): %s", decision.Status, decision.Passed, decision.Rationale)
	}
	if decision.Rationale != "Passed with 1 claim lacking sufficient corroboration" {
		t.Errorf("unexpected rationale: %q", decision.Rationale)
	}
}

func TestEvaluateClaims_InsufficientCorroborationFailsWhenNeedsReviewDisallowed(t *testing.T) {
	claims := []Claim{verifiedClaim("c1", ClaimStatistical)}
	criteria := DefaultClaimsCriteria()
	criteria.MinCorroboratingSources = 2
	criteria.AllowNeedsReview = false

	decision := EvaluateClaims(claims, criteria)
	if decision.Status != ClaimsDecisionFail || decision.Passed {
		t.Errorf("expected fail, got %s (passed=%v): %s", decision.Status, decision.Passed, decision.Rationale)
	}
	if decision.Rationale != "1 claim lacks sufficient corroboration (not allowed by criteria)" {
		t.Errorf("unexpected rationale: %q", decision.Rationale)
	}
}

func TestEvaluateClaims_CorroboratedClaimPasses(t *testing.T) {
	claims := []Claim{verifiedClaim("c1", ClaimStatistical, "c1-corroborating-1")}
	criteria := DefaultClaimsCriteria()
	criteria.MinCorroboratingSources = 2

	decision := EvaluateClaims(claims, criteria)
	if decision.Status != ClaimsDecisionPass {
		t.Errorf("expected pass, got %s: %s", decision.Status, decision.Rationale)
	}
}

func TestEvaluateClaims_NeedsReviewAndCorroborationCombine(t *testing.T) {
	claims := []Claim{
		{ID: "nr", Category: ClaimStatistical, Verdict: VerdictNeedsReview},
		verifiedClaim("uc", ClaimStatistical),
	}
	criteria := DefaultClaimsCriteria()
	criteria.MinCorroboratingSources = 2

	decision := EvaluateClaims(claims, criteria)
	want := "Passed with 1 claim needing review and 1 claim lacking sufficient corroboration"
	if decision.Rationale != want {
		t.Errorf("rationale = %q, want %q", decision.Rationale, want)
	}
}

// Regression: the single-reason rationale text predates this RMI and must
// not change now that a second reason (corroboration) exists.
func TestEvaluateClaims_SingleReasonRationaleUnchanged(t *testing.T) {
	needsReviewOnly := []Claim{{ID: "nr", Category: ClaimStatistical, Verdict: VerdictNeedsReview}}
	d := EvaluateClaims(needsReviewOnly, DefaultClaimsCriteria())
	if d.Rationale != "Passed with 1 claim needing review" {
		t.Errorf("rationale = %q", d.Rationale)
	}

	criteria := DefaultClaimsCriteria()
	criteria.AllowNeedsReview = false
	d = EvaluateClaims(needsReviewOnly, criteria)
	if d.Rationale != "1 claim needs review (not allowed by criteria)" {
		t.Errorf("rationale = %q", d.Rationale)
	}
}
