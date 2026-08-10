package claims

import "testing"

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
