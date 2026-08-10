package claims

import (
	"testing"
	"time"
)

func extClaim(id string, verdict Verdict, url, quote string, stat *StatisticalDetail) Claim {
	c := Claim{
		ID:       id,
		Text:     id,
		Category: ClaimStatistical,
		Verdict:  verdict,
		Validation: &Validation{
			Type: SourceExternal,
			External: &ExternalValidation{
				URL:        url,
				SourceType: ExternalReputableVendor,
				QuotedText: quote,
			},
		},
		Statistical: stat,
	}
	return c
}

func TestLint_CriteriaCorroborationRequirement(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	// Primary role — the SourceRole-driven rule (RMI-002) would not flag
	// this; only the criteria-driven threshold (RMI-003) should.
	c.Validation.External.Role = SourceRolePrimary

	r := &ClaimsReport{
		Claims:   []Claim{c},
		Criteria: ClaimsCriteria{MinCorroboratingSources: 2},
	}
	f := Lint(r)
	if !hasRule(f, "verified-insufficient-corroboration") {
		t.Errorf("expected verified-insufficient-corroboration, got %+v", f)
	}
	if hasRule(f, "verified-role-needs-corroboration") {
		t.Error("primary role should not trigger the role-based rule")
	}
}

func TestLint_CriteriaCorroborationSatisfied(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	c.RelatedClaimIDs = []string{"x-corroborating-1"}

	r := &ClaimsReport{
		Claims:   []Claim{c},
		Criteria: ClaimsCriteria{MinCorroboratingSources: 2},
	}
	if f := Lint(r); hasRule(f, "verified-insufficient-corroboration") {
		t.Errorf("corroborated claim should not be flagged, got %+v", f)
	}
}

func TestLint_CriteriaStaleAsOfDateFlagged(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	old := time.Now().AddDate(-3, 0, 0)
	c.Statistical.AsOfDate = &old

	r := &ClaimsReport{
		Claims:   []Claim{c},
		Criteria: ClaimsCriteria{MaxClaimAge: 365 * 24 * time.Hour},
	}
	f := Lint(r)
	if !hasRule(f, "verified-stale-as-of-date") {
		t.Errorf("expected verified-stale-as-of-date, got %+v", f)
	}
	if !HasErrors(f) {
		t.Error("stale claim should be an error")
	}
}

func TestLint_CriteriaFreshAsOfDatePasses(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	recent := time.Now().AddDate(0, -1, 0)
	c.Statistical.AsOfDate = &recent

	r := &ClaimsReport{
		Claims:   []Claim{c},
		Criteria: ClaimsCriteria{MaxClaimAge: 365 * 24 * time.Hour},
	}
	if f := Lint(r); hasRule(f, "verified-stale-as-of-date") {
		t.Errorf("recent claim should not be flagged, got %+v", f)
	}
}

func TestLint_CriteriaStaleDisabledByDefault(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	old := time.Now().AddDate(-10, 0, 0)
	c.Statistical.AsOfDate = &old

	// No Criteria set on the report (zero value) — requirement is off.
	if f := Lint(&ClaimsReport{Claims: []Claim{c}}); hasRule(f, "verified-stale-as-of-date") {
		t.Errorf("staleness requirement should be disabled by default, got %+v", f)
	}
}

func TestLint_CriteriaStaleUnsetAsOfDateNeverFlagged(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	// AsOfDate left unset — age is unknown, not stale.

	r := &ClaimsReport{
		Claims:   []Claim{c},
		Criteria: ClaimsCriteria{MaxClaimAge: 24 * time.Hour},
	}
	if f := Lint(r); hasRule(f, "verified-stale-as-of-date") {
		t.Errorf("claim with unset AsOfDate should never be flagged, got %+v", f)
	}
}

func TestLint_CriteriaCorroborationDisabledByDefault(t *testing.T) {
	c := extClaim("x", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	// No Criteria set on the report (zero value) — requirement is off.
	if f := Lint(&ClaimsReport{Claims: []Claim{c}}); hasRule(f, "verified-insufficient-corroboration") {
		t.Errorf("corroboration requirement should be disabled by default, got %+v", f)
	}
}

func TestLint_SecondaryAnalysisRequiresCorroboration(t *testing.T) {
	c := extClaim("arr", VerdictVerified, "https://example.com", "grew to $500 million",
		NewStatisticalDetail(500, "million USD", PrecisionApproximate))
	c.Validation.External.Role = SourceRoleSecondaryAnalysis

	f := Lint(&ClaimsReport{Claims: []Claim{c}})
	if !hasRule(f, "verified-role-needs-corroboration") {
		t.Errorf("expected verified-role-needs-corroboration, got %+v", f)
	}
	if !HasErrors(f) {
		t.Error("missing corroboration for secondary-analysis should be an error")
	}
}

func TestLint_SelfReportedRequiresCorroboration(t *testing.T) {
	c := extClaim("f500", VerdictVerified, "https://cursor.com/enterprise", "64% of Fortune 500",
		NewStatisticalDetail(64, "% of Fortune 500", PrecisionExact))
	c.Validation.External.Role = SourceRoleSelfReported

	f := Lint(&ClaimsReport{Claims: []Claim{c}})
	if !hasRule(f, "verified-role-needs-corroboration") {
		t.Errorf("expected verified-role-needs-corroboration, got %+v", f)
	}
}

func TestLint_CorroboratedSecondaryAnalysisPasses(t *testing.T) {
	c := extClaim("arr", VerdictVerified, "https://example.com", "grew to $500 million",
		NewStatisticalDetail(500, "million USD", PrecisionApproximate))
	c.Validation.External.Role = SourceRoleSecondaryAnalysis
	c.RelatedClaimIDs = []string{"arr-corroborating-1"}

	f := Lint(&ClaimsReport{Claims: []Claim{c}})
	if hasRule(f, "verified-role-needs-corroboration") {
		t.Errorf("corroborated secondary-analysis claim should not be flagged, got %+v", f)
	}
}

func TestLint_PrimaryAndUnsetRoleNeedNoCorroboration(t *testing.T) {
	primary := extClaim("p", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	primary.Validation.External.Role = SourceRolePrimary

	unset := extClaim("u", VerdictVerified, "https://example.com", "1,500+ PRs merged",
		NewStatisticalDetail(1500, "PRs", PrecisionApproximate))
	// Role left unset — reports predating this field must not be penalized.

	f := Lint(&ClaimsReport{Claims: []Claim{primary, unset}})
	if hasRule(f, "verified-role-needs-corroboration") {
		t.Errorf("primary and unset-role claims should not require corroboration, got %+v", f)
	}
}

func hasRule(findings []LintFinding, rule string) bool {
	for _, f := range findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

func TestLint_VerifiedRequiresQuoteAndURL(t *testing.T) {
	r := &ClaimsReport{Claims: []Claim{
		extClaim("no-quote", VerdictVerified, "https://example.com", "", nil),
		extClaim("no-url", VerdictVerified, "", "some quote", nil),
	}}
	f := Lint(r)
	if !hasRule(f, "verified-requires-quote") {
		t.Error("expected verified-requires-quote for the quoteless claim")
	}
	if !hasRule(f, "verified-requires-url") {
		t.Error("expected verified-requires-url for the urlless claim")
	}
	if !HasErrors(f) {
		t.Error("expected errors")
	}
}

func TestLint_VerifiedWithQuoteAndValuePasses(t *testing.T) {
	stat := NewStatisticalDetail(1500, "PRs", PrecisionApproximate)
	r := &ClaimsReport{Claims: []Claim{
		extClaim("ok", VerdictVerified, "https://example.com", "1,500+ PRs merged", stat),
	}}
	f := Lint(r)
	if len(f) != 0 {
		t.Errorf("expected no findings, got %+v", f)
	}
}

func TestLint_ValueNotInQuoteWarns(t *testing.T) {
	stat := NewStatisticalDetail(3, "billion USD", PrecisionApproximate)
	r := &ClaimsReport{Claims: []Claim{
		extClaim("arr", VerdictVerified, "https://example.com", "grew to $500 million", stat),
	}}
	f := Lint(r)
	if !hasRule(f, "verified-value-in-quote") {
		t.Error("expected verified-value-in-quote warning")
	}
	if HasErrors(f) {
		t.Error("value-in-quote should be a warning, not an error")
	}
	if !HasWarnings(f) {
		t.Error("expected a warning")
	}
}

func TestLint_NonVerifiedClaimsSkipped(t *testing.T) {
	r := &ClaimsReport{Claims: []Claim{
		extClaim("nr", VerdictNeedsReview, "", "", nil),
		extClaim("rej", VerdictRejected, "", "", nil),
		extClaim("unv", VerdictUnverified, "", "", nil),
	}}
	if f := Lint(r); len(f) != 0 {
		t.Errorf("non-verified claims should not be gated, got %+v", f)
	}
}

func TestLint_DerivedAndInternal(t *testing.T) {
	derivedNoSources := Claim{ID: "d", Verdict: VerdictVerified, Category: ClaimStatistical,
		Validation: &Validation{Type: SourceDerived, Derived: &DerivedValidation{}}}
	internalNoEvidence := Claim{ID: "i", Verdict: VerdictVerified, Category: ClaimStatistical,
		Validation: &Validation{Type: SourceInternal, Internal: &InternalValidation{}}}
	f := Lint(&ClaimsReport{Claims: []Claim{derivedNoSources, internalNoEvidence}})
	if !hasRule(f, "verified-derived-needs-sources") {
		t.Error("expected verified-derived-needs-sources")
	}
	if !hasRule(f, "verified-internal-needs-evidence") {
		t.Error("expected verified-internal-needs-evidence")
	}
}

func TestLint_MissingAndDuplicateIDs(t *testing.T) {
	r := &ClaimsReport{Claims: []Claim{
		{ID: "", Verdict: VerdictRejected},
		{ID: "dup", Verdict: VerdictRejected},
		{ID: "dup", Verdict: VerdictRejected},
	}}
	f := Lint(r)
	if !hasRule(f, "claim-missing-id") {
		t.Error("expected claim-missing-id")
	}
	if !hasRule(f, "claim-duplicate-id") {
		t.Error("expected claim-duplicate-id")
	}
}

func TestAddThousands(t *testing.T) {
	cases := map[string]string{
		"1500": "1,500", "20000": "20,000", "77000": "77,000",
		"1000000": "1,000,000", "999": "999", "-1500": "-1,500", "1234.5": "1,234.5",
	}
	for in, want := range cases {
		if got := addThousands(in); got != want {
			t.Errorf("addThousands(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValueInText(t *testing.T) {
	if !valueInText(1500, "1,500+ PRs merged") {
		t.Error("expected 1500 to match '1,500'")
	}
	if !valueInText(20, "crossed 20 million users") {
		t.Error("expected 20 to match")
	}
	if !valueInText(4.7, "4.7 million subscribers") {
		t.Error("expected 4.7 to match")
	}
	if valueInText(1000000, "over 1 million daily active users") {
		t.Error("1000000 should not match '1 million' (unit-scaled) — warning is expected")
	}
}
