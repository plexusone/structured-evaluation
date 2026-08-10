package html

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/plexusone/structured-evaluation/claims"
)

func sampleReport() *claims.ClaimsReport {
	r := claims.NewClaimsReport("case-study.md")
	r.Metadata.DocumentTitle = "Sample Case Study"
	r.Metadata.GeneratedAt = time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)

	// Verified statistical claim with a source and a quote.
	verified := claims.NewClaim("stat-1", "1,500+ PRs merged by agents",
		claims.ClaimStatistical, claims.Location{Section: "Metrics"})
	verified.SetValidation((&claims.Validation{
		Type: claims.SourceExternal,
		External: &claims.ExternalValidation{
			URL:           "https://example.com/blog",
			SourceType:    claims.ExternalReputableVendor,
			Reliability:   claims.ReliabilityHigh,
			QuotedText:    "1,500+ PRs Later",
			VerifiedMatch: true,
		},
	}))
	verified.SetStatistical(claims.NewStatisticalDetail(1500, "PRs", claims.PrecisionApproximate))
	verified.Rationale = "Primary source."
	r.AddClaim(*verified)

	// Rejected claim from an aggregator, with characters that must be escaped.
	rejected := claims.NewClaim("stat-2", "9,900% growth <script>",
		claims.ClaimStatistical, claims.Location{Section: "Metrics"})
	rejected.SetValidation((&claims.Validation{
		Type: claims.SourceExternal,
		External: &claims.ExternalValidation{
			URL:         "https://aggregator.example/stats",
			SourceType:  claims.ExternalAggregator,
			Reliability: claims.ReliabilityLow,
		},
	}))
	rejected.SetStatistical(claims.NewStatisticalDetail(9900, "%", claims.PrecisionEstimated))
	rejected.Rationale = "Aggregator-only; do not use."
	r.AddClaim(*rejected)

	r.Finalize()
	return r
}

func TestRenderClaims(t *testing.T) {
	var buf bytes.Buffer
	if err := New(&buf).RenderClaims(sampleReport()); err != nil {
		t.Fatalf("RenderClaims failed: %v", err)
	}
	out := buf.String()

	checks := []string{
		"<!DOCTYPE html>",
		"Sample Case Study",
		`class="decision fail"`, // one rejected claim -> fail
		"1 claim rejected",
		`<section class="group verified"`,
		`<section class="group rejected"`,
		"PRs merged by agents",
		"Primary source.",
		"reputable-vendor",
		"aggregator",
		"rel-low",
		"1500 PRs",
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("output missing %q", c)
		}
	}

	// The needs-review group has no claims here and must be omitted.
	if strings.Contains(out, `class="group needs-review"`) {
		t.Error("empty needs-review group should be omitted")
	}

	// Script tag in claim text must be escaped, not emitted raw.
	if strings.Contains(out, "<script>") {
		t.Error("claim text was not HTML-escaped")
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("expected escaped script tag in output")
	}
}

func TestRenderClaims_Deterministic(t *testing.T) {
	var a, b bytes.Buffer
	if err := New(&a).RenderClaims(sampleReport()); err != nil {
		t.Fatal(err)
	}
	if err := New(&b).RenderClaims(sampleReport()); err != nil {
		t.Fatal(err)
	}
	if a.String() != b.String() {
		t.Error("render output is not deterministic")
	}
}
