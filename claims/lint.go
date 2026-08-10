package claims

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// LintSeverity classifies a lint finding.
type LintSeverity string

const (
	// LintError is a finding that should block publication (non-zero exit).
	LintError LintSeverity = "error"

	// LintWarning is an advisory finding for human review.
	LintWarning LintSeverity = "warning"
)

// LintFinding is a single problem found by Lint.
type LintFinding struct {
	ClaimID  string       `json:"claimId,omitempty"`
	Rule     string       `json:"rule"`
	Severity LintSeverity `json:"severity"`
	Message  string       `json:"message"`
}

// Lint checks a claims report for evidence-integrity problems, centered on the
// invariant that a "verified" verdict must actually be backed by evidence:
// an external claim needs a resolving URL and a verbatim quote; a derived
// claim needs source claims; an internal claim needs evidence. The
// value-appears-in-quote check is advisory (a warning), because legitimate
// verified claims quote a rule rather than the number (e.g. a "0%" target) or
// a range, or scale the value into a unit ("20 million", "$60 billion").
//
// A verified external claim whose SourceRole is secondary-analysis or
// self-reported (SourceRole.RequiresCorroboration) must carry at least one
// entry in RelatedClaimIDs — an independent corroborating claim — or Lint
// errors. An unset Role is not flagged, so existing reports that predate the
// Role field are unaffected.
//
// Separately, if the report's own Criteria.MinCorroboratingSources is set
// (> 1), every verified claim in scope (all categories, or only
// Criteria.CorroborationCategories when set) must have at least that many
// independent sources (itself plus RelatedClaimIDs), regardless of
// validation type or SourceRole. Disabled by default (0), so existing
// reports are unaffected unless they opt in via Criteria.
//
// Lint does not mutate the report. It complements DetermineVerdict: because a
// verdict can be hand-authored (bypassing DetermineVerdict), Lint re-checks
// that a stated "verified" is earned.
func Lint(r *ClaimsReport) []LintFinding {
	var out []LintFinding
	seen := map[string]bool{}

	for i := range r.Claims {
		c := &r.Claims[i]

		if c.ID == "" {
			out = append(out, LintFinding{Rule: "claim-missing-id", Severity: LintError,
				Message: "claim has no id"})
		} else if seen[c.ID] {
			out = append(out, LintFinding{ClaimID: c.ID, Rule: "claim-duplicate-id", Severity: LintError,
				Message: "duplicate claim id"})
		}
		seen[c.ID] = true

		// Only verified claims must clear the evidence bar. needs-review,
		// rejected, and unverified carry no such obligation.
		if c.Verdict != VerdictVerified {
			continue
		}

		// General, configurable corroboration requirement (r.Criteria).
		// Independent of validation type: corroboration is about having
		// multiple independent claims behind a fact, not about how any one
		// of them was validated. Distinct from the SourceRole-driven check
		// below, which is a fixed policy (secondary-analysis/self-reported
		// always need corroboration) rather than a criteria-configured one.
		if !IsSufficientlyCorroborated(*c, r.Criteria) {
			out = append(out, finding(c, "verified-insufficient-corroboration", LintError,
				fmt.Sprintf("verified claim has %d source(s), fewer than the required %d",
					1+len(c.RelatedClaimIDs), r.Criteria.MinCorroboratingSources)))
		}

		v := c.Validation
		if v == nil {
			out = append(out, finding(c, "verified-requires-validation", LintError,
				"verified claim has no validation"))
			continue
		}

		switch v.Type {
		case SourceExternal:
			out = append(out, lintVerifiedExternal(c, v.External)...)
		case SourceDerived:
			if v.Derived == nil || len(v.Derived.SourceClaimIDs) == 0 {
				out = append(out, finding(c, "verified-derived-needs-sources", LintError,
					"verified derived claim lists no source claim ids"))
			}
		case SourceInternal:
			if v.Internal == nil || (v.Internal.EvidencePath == "" && v.Internal.Output == "") {
				out = append(out, finding(c, "verified-internal-needs-evidence", LintError,
					"verified internal claim has no evidence path or output"))
			}
		case SourceSubjective:
			out = append(out, finding(c, "verified-subjective", LintWarning,
				"verified claim is subjective — confirm it should be published as verified"))
		}
	}
	return out
}

func lintVerifiedExternal(c *Claim, e *ExternalValidation) []LintFinding {
	var out []LintFinding
	if e == nil {
		return []LintFinding{finding(c, "verified-requires-external", LintError,
			"verified external claim has no external validation")}
	}
	if strings.TrimSpace(e.URL) == "" {
		out = append(out, finding(c, "verified-requires-url", LintError,
			"verified claim has no source URL"))
	}
	if strings.TrimSpace(e.QuotedText) == "" {
		out = append(out, finding(c, "verified-requires-quote", LintError,
			"verified claim has no quoted text from the source"))
		return out // can't check value-in-quote without a quote
	}
	if c.Statistical != nil && !valueInText(c.Statistical.Value, e.QuotedText) {
		out = append(out, finding(c, "verified-value-in-quote", LintWarning,
			fmt.Sprintf("value %s not found in the quoted text — confirm the quote supports it",
				formatLintValue(c.Statistical.Value))))
	}
	if e.Role.RequiresCorroboration() && len(c.RelatedClaimIDs) == 0 {
		out = append(out, finding(c, "verified-role-needs-corroboration", LintError,
			fmt.Sprintf("verified claim sourced as %q has no corroborating relatedClaimIds", e.Role)))
	}
	return out
}

func finding(c *Claim, rule string, sev LintSeverity, msg string) LintFinding {
	return LintFinding{ClaimID: c.ID, Rule: rule, Severity: sev, Message: msg}
}

// HasErrors reports whether any finding is error severity.
func HasErrors(findings []LintFinding) bool {
	for _, f := range findings {
		if f.Severity == LintError {
			return true
		}
	}
	return false
}

// HasWarnings reports whether any finding is warning severity.
func HasWarnings(findings []LintFinding) bool {
	for _, f := range findings {
		if f.Severity == LintWarning {
			return true
		}
	}
	return false
}

// valueInText reports whether the numeric value appears in the text, trying the
// canonical rendering and a thousands-separated form. Intentionally
// conservative: a false result is a warning, not proof of a bad claim.
func valueInText(v float64, text string) bool {
	canon := formatLintValue(v)
	t := strings.ToLower(text)
	if strings.Contains(t, canon) {
		return true
	}
	if v == math.Trunc(v) && math.Abs(v) >= 1000 {
		if strings.Contains(t, addThousands(canon)) {
			return true
		}
	}
	return false
}

// formatLintValue renders a float without trailing zeros (1500, 4.7, 0, 77000).
func formatLintValue(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// addThousands inserts comma grouping into an integer string ("77000" ->
// "77,000"). Handles an optional leading minus; leaves any fractional part.
func addThousands(s string) string {
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	intPart, frac := s, ""
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		intPart, frac = s[:dot], s[dot:]
	}
	var b strings.Builder
	n := len(intPart)
	for i, ch := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	res := b.String() + frac
	if neg {
		res = "-" + res
	}
	return res
}
