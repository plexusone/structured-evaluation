package claims

import (
	"testing"
	"time"
)

func verifiedClaim(id string, related ...string) Claim {
	return Claim{
		ID:              id,
		Category:        ClaimStatistical,
		Verdict:         VerdictVerified,
		RelatedClaimIDs: related,
	}
}

func verifiedStatClaim(id string, asOfDate *time.Time) Claim {
	return Claim{
		ID:          id,
		Category:    ClaimStatistical,
		Verdict:     VerdictVerified,
		Statistical: &StatisticalDetail{Value: 1, AsOfDate: asOfDate},
	}
}

func TestIsSufficientlyCorroborated(t *testing.T) {
	single := verifiedClaim("c1")
	double := verifiedClaim("c2", "c2-corroborating-1")

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
	claims := []Claim{verifiedClaim("c1")}
	decision := EvaluateClaims(claims, DefaultClaimsCriteria())
	if decision.Status != ClaimsDecisionPass {
		t.Errorf("expected pass with corroboration disabled, got %s: %s", decision.Status, decision.Rationale)
	}
}

func TestEvaluateClaims_InsufficientCorroborationIsConditionalWhenAllowed(t *testing.T) {
	claims := []Claim{verifiedClaim("c1")}
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
	claims := []Claim{verifiedClaim("c1")}
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
	claims := []Claim{verifiedClaim("c1", "c1-corroborating-1")}
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
		verifiedClaim("uc"),
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

func TestIsFresh(t *testing.T) {
	now := time.Now()
	old := now.AddDate(-3, 0, 0)
	recent := now.AddDate(0, -1, 0)

	tests := []struct {
		name     string
		claim    Claim
		criteria ClaimsCriteria
		want     bool
	}{
		{"disabled (0) always fresh", verifiedStatClaim("c1", &old), ClaimsCriteria{}, true},
		{"no statistical detail always fresh", Claim{ID: "c2", Verdict: VerdictVerified},
			ClaimsCriteria{MaxClaimAge: 365 * 24 * time.Hour}, true},
		{"unset AsOfDate always fresh", verifiedStatClaim("c3", nil),
			ClaimsCriteria{MaxClaimAge: 365 * 24 * time.Hour}, true},
		{"old AsOfDate exceeds threshold", verifiedStatClaim("c4", &old),
			ClaimsCriteria{MaxClaimAge: 365 * 24 * time.Hour}, false},
		{"recent AsOfDate within threshold", verifiedStatClaim("c5", &recent),
			ClaimsCriteria{MaxClaimAge: 365 * 24 * time.Hour}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFresh(tt.claim, tt.criteria, now); got != tt.want {
				t.Errorf("IsFresh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateClaims_StaleDisabledByDefault(t *testing.T) {
	old := time.Now().AddDate(-5, 0, 0)
	claims := []Claim{verifiedStatClaim("c1", &old)}
	decision := EvaluateClaims(claims, DefaultClaimsCriteria())
	if decision.Status != ClaimsDecisionPass {
		t.Errorf("expected pass with staleness disabled, got %s: %s", decision.Status, decision.Rationale)
	}
}

func TestEvaluateClaims_StaleClaimIsConditionalWhenAllowed(t *testing.T) {
	old := time.Now().AddDate(-3, 0, 0)
	claims := []Claim{verifiedStatClaim("c1", &old)}
	criteria := DefaultClaimsCriteria()
	criteria.MaxClaimAge = 365 * 24 * time.Hour
	criteria.AllowNeedsReview = true

	decision := EvaluateClaims(claims, criteria)
	if decision.Status != ClaimsDecisionConditional || !decision.Passed {
		t.Errorf("expected conditional pass, got %s (passed=%v): %s", decision.Status, decision.Passed, decision.Rationale)
	}
	if decision.Rationale != "Passed with 1 claim with a stale statistic" {
		t.Errorf("unexpected rationale: %q", decision.Rationale)
	}
}

func TestEvaluateClaims_StaleClaimFailsWhenNeedsReviewDisallowed(t *testing.T) {
	old := time.Now().AddDate(-3, 0, 0)
	claims := []Claim{verifiedStatClaim("c1", &old)}
	criteria := DefaultClaimsCriteria()
	criteria.MaxClaimAge = 365 * 24 * time.Hour
	criteria.AllowNeedsReview = false

	decision := EvaluateClaims(claims, criteria)
	if decision.Status != ClaimsDecisionFail || decision.Passed {
		t.Errorf("expected fail, got %s (passed=%v): %s", decision.Status, decision.Passed, decision.Rationale)
	}
	if decision.Rationale != "1 claim has a stale statistic (not allowed by criteria)" {
		t.Errorf("unexpected rationale: %q", decision.Rationale)
	}
}

func TestEvaluateClaims_AllThreeReasonsCombine(t *testing.T) {
	old := time.Now().AddDate(-3, 0, 0)
	staleButCorroborated := verifiedStatClaim("stale", &old)
	staleButCorroborated.RelatedClaimIDs = []string{"stale-corroborating-1"}
	claims := []Claim{
		{ID: "nr", Category: ClaimStatistical, Verdict: VerdictNeedsReview},
		verifiedClaim("uc"),
		staleButCorroborated,
	}
	criteria := DefaultClaimsCriteria()
	criteria.MinCorroboratingSources = 2
	criteria.MaxClaimAge = 365 * 24 * time.Hour

	decision := EvaluateClaims(claims, criteria)
	want := "Passed with 1 claim needing review, 1 claim lacking sufficient corroboration, and 1 claim with a stale statistic"
	if decision.Rationale != want {
		t.Errorf("rationale = %q, want %q", decision.Rationale, want)
	}
}
