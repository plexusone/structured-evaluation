package claims

import (
	"reflect"
	"testing"
)

func quotedClaim(id string, conf Confidence, refs ...string) Claim {
	c := Claim{ID: id, Text: "A claim.", Verdict: VerdictVerified, Confidence: conf}
	for _, r := range refs {
		c.Evidence = append(c.Evidence, Evidence{SourceRef: r, Stance: StanceSupports, Quote: "A passage."})
	}
	return c
}

func TestSupportShortfalls(t *testing.T) {
	w := WikipediaClaimsCriteria()
	unquoted := Claim{ID: "c", Verdict: VerdictVerified, Confidence: ConfidenceHigh,
		Evidence: []Evidence{{SourceRef: "src-a", Stance: StanceSupports}}}
	exceptional := quotedClaim("c", ConfidenceHigh, "src-a", "src-a")
	exceptional.Exceptional = true
	tests := []struct {
		name string
		c    Claim
		want []string
	}{
		{"meets all", quotedClaim("c", ConfidenceHigh, "src-a"), nil},
		{"no quote", unquoted, []string{"requireQuotedSupport"}},
		{"low confidence", quotedClaim("c", ConfidenceLow, "src-a"), []string{"minSupportConfidence"}},
		{"no confidence", quotedClaim("c", "", "src-a"), []string{"minSupportConfidence"}},
		{"exceptional with one source", exceptional, []string{"exceptionalMinSources"}},
		{"not verified", Claim{ID: "c", Verdict: VerdictRejected}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SupportShortfalls(tt.c, w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shortfalls = %v, want %v", got, tt.want)
			}
		})
	}
	exceptional.Evidence[1].SourceRef = "src-b"
	if got := SupportShortfalls(exceptional, w); got != nil {
		t.Errorf("exceptional claim with two sources: shortfalls = %v", got)
	}
	// Off by default: the default criteria ask for none of it.
	if got := SupportShortfalls(unquoted, DefaultClaimsCriteria()); got != nil {
		t.Errorf("default criteria shortfalls = %v, want none", got)
	}
}

func TestEvaluateClaimsSupportCriteria(t *testing.T) {
	cs := []Claim{quotedClaim("c1", ConfidenceHigh, "src-a"), quotedClaim("c2", ConfidenceLow, "src-a")}
	if d := EvaluateClaims(cs, WikipediaClaimsCriteria()); d.Passed || d.Status != ClaimsDecisionFail {
		t.Errorf("a low-confidence claim should fail the Wikipedia criteria, got %+v", d)
	}
	if d := EvaluateClaims(cs, DefaultClaimsCriteria()); !d.Passed || d.Status != ClaimsDecisionPass {
		t.Errorf("default criteria should pass, got %+v", d)
	}
}

func TestRefuted(t *testing.T) {
	refuted := Claim{ID: "myth", Text: "It began in a garage.", Verdict: VerdictRefuted,
		Evidence: []Evidence{{SourceRef: "src-history", Stance: StanceDisputes, Quote: "The first premises were a hall."}}}
	if !VerdictRefuted.IsBlocking() || VerdictRefuted.IsPassing() {
		t.Error("refuted should block and not pass")
	}
	counts := CountClaims([]Claim{refuted})
	if counts.Refuted != 1 {
		t.Errorf("counts = %+v", counts)
	}
	if d := EvaluateClaims([]Claim{refuted}, DefaultClaimsCriteria()); d.Passed {
		t.Errorf("a document stating a refuted claim should fail, got %+v", d)
	}
	r := &ClaimsReport{Claims: []Claim{refuted}}
	if got := Lint(r); len(got) != 0 {
		t.Errorf("lint = %+v, want none", got)
	}
	refuted.Evidence[0].Quote = ""
	r = &ClaimsReport{Claims: []Claim{refuted}}
	if got := Lint(r); len(got) != 1 || got[0].Rule != "refuted-requires-quoted-dispute" {
		t.Errorf("lint = %+v, want refuted-requires-quoted-dispute", got)
	}
}

func TestLintSupportCriteria(t *testing.T) {
	c := quotedClaim("c1", ConfidenceLow, "src-a")
	c.Exceptional = true
	r := &ClaimsReport{Claims: []Claim{c}, Criteria: WikipediaClaimsCriteria()}
	var rules []string
	for _, f := range Lint(r) {
		rules = append(rules, f.Rule)
	}
	want := []string{"verified-low-confidence", "verified-exceptional-needs-sources"}
	if !reflect.DeepEqual(rules, want) {
		t.Errorf("rules = %v, want %v", rules, want)
	}
}
