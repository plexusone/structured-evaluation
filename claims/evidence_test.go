package claims

import (
	"encoding/json"
	"strings"
	"testing"
)

func multiSourceClaim(id string, v Verdict) Claim {
	return Claim{
		ID: id, Text: "The company was founded in 1929.", Category: ClaimTimeline, Verdict: v,
		CheckedBy: "checker",
		Evidence: []Evidence{
			{SourceRef: "src-history", Stance: StanceSupports, Quote: "fondata nel 1929", Language: "it",
				Translation: &Translation{Text: "founded in 1929", Language: "en"}, Role: SourceRolePrimary},
			{URL: "https://example.org/other", Stance: StanceDisputes, Quote: "founded in 1930"},
		},
	}
}

func TestNewVerdictsCountAndDecide(t *testing.T) {
	tests := []struct {
		name   string
		v      Verdict
		status ClaimsDecisionStatus
	}{
		{"contested asks a person", VerdictContested, ClaimsDecisionConditional},
		{"citation gap asks a person", VerdictCitationGap, ClaimsDecisionConditional},
		{"contradicted fails", VerdictContradicted, ClaimsDecisionFail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewClaimsReport("doc")
			r.AddClaim(multiSourceClaim("c1", VerdictVerified))
			r.AddClaim(multiSourceClaim("c2", tt.v))
			d := r.Evaluate()
			if d.Status != tt.status {
				t.Errorf("status = %s, want %s", d.Status, tt.status)
			}
			c := d.Counts
			if c.Total != 2 || c.Verified != 1 || c.Contested+c.Contradicted+c.CitationGap != 1 {
				t.Errorf("counts = %+v", c)
			}
		})
	}
	if !VerdictContradicted.IsBlocking() || VerdictContested.IsBlocking() || !VerdictContested.NeedsPerson() {
		t.Error("verdict helpers disagree with the decision")
	}
}

func TestContestedFailsWhenReviewNotAllowed(t *testing.T) {
	r := NewClaimsReport("doc")
	r.SetCriteria(StrictClaimsCriteria())
	r.AddClaim(multiSourceClaim("c1", VerdictContested))
	if d := r.Evaluate(); d.Status != ClaimsDecisionFail {
		t.Errorf("status = %s, want fail", d.Status)
	}
}

func TestAgreement(t *testing.T) {
	agreed := multiSourceClaim("c1", VerdictVerified)
	agreed.Review = &Review{Judge: "skeptic", Agrees: true}
	overruled := multiSourceClaim("c2", VerdictContested)
	overruled.Review = &Review{Judge: "skeptic", Agrees: false, FirstVerdict: VerdictVerified, Rationale: "the 1930 date"}
	unreviewed := multiSourceClaim("c3", VerdictVerified)

	if a := AgreementOf([]Claim{unreviewed}); a != nil {
		t.Errorf("agreement without reviews = %+v, want nil", a)
	}
	r := NewClaimsReport("doc")
	for _, c := range []Claim{agreed, overruled, unreviewed} {
		r.AddClaim(c)
	}
	r.Evaluate()
	a := r.Summary.SecondReview
	if a == nil || a.Reviewed != 2 || a.Agreed != 1 || a.Rate != 0.5 {
		t.Errorf("agreement = %+v", a)
	}
	if len(r.Summary.NeedsReviewClaims) != 1 || r.Summary.NeedsReviewClaims[0] != "c2" {
		t.Errorf("needs review = %v", r.Summary.NeedsReviewClaims)
	}
}

func TestLintAcceptsQuotedEvidence(t *testing.T) {
	quoted := multiSourceClaim("c1", VerdictVerified)
	unquoted := multiSourceClaim("c2", VerdictVerified)
	unquoted.Evidence = []Evidence{{SourceRef: "src-history", Stance: StanceSupports}}
	disputedOnly := multiSourceClaim("c3", VerdictVerified)
	disputedOnly.Evidence = disputedOnly.Evidence[1:]

	r := NewClaimsReport("doc")
	for _, c := range []Claim{quoted, unquoted, disputedOnly} {
		r.AddClaim(c)
	}
	got := map[string]string{}
	for _, f := range Lint(r) {
		got[f.ClaimID] = f.Rule
	}
	if _, ok := got["c1"]; ok {
		t.Errorf("quoted supporting evidence was flagged: %v", got)
	}
	for _, id := range []string{"c2", "c3"} {
		if got[id] != "verified-requires-validation" {
			t.Errorf("%s: rule = %q, want verified-requires-validation", id, got[id])
		}
	}
}

func TestEvidenceFieldsOmittedWhenUnused(t *testing.T) {
	r := NewClaimsReport("doc")
	c := NewClaim("c1", "A fact.", ClaimTimeline, Location{})
	r.AddClaim(*c)
	r.Evaluate()
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"evidence"`, `"checkedBy"`, `"review"`, `"inaccessible"`, `"secondReview"`, `"contested"`, `"citationGap"`} {
		if strings.Contains(string(data), field) {
			t.Errorf("report without the new fields contains %s", field)
		}
	}
}
