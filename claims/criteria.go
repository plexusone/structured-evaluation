package claims

import (
	"strings"
	"time"
)

// ClaimsCriteria defines requirements for claim validation approval.
type ClaimsCriteria struct {
	// RequireAllVerified requires all claims to be verified.
	RequireAllVerified bool `json:"requireAllVerified"`

	// AllowSubjectiveWithDisclaimer permits subjective claims if acknowledged.
	AllowSubjectiveWithDisclaimer bool `json:"allowSubjectiveWithDisclaimer"`

	// AllowNeedsReview permits claims that need review (conditional pass).
	AllowNeedsReview bool `json:"allowNeedsReview"`

	// MinReliabilityTier is the minimum acceptable reliability for external sources.
	MinReliabilityTier ReliabilityTier `json:"minReliabilityTier"`

	// RequireReproducible requires internal validations to be reproducible.
	RequireReproducible bool `json:"requireReproducible"`

	// RequiredCategories are claim categories that must have at least one verified claim.
	RequiredCategories []ClaimCategory `json:"requiredCategories,omitempty"`

	// MinCorroboratingSources is the minimum number of independent sources
	// required for a verified claim to be treated as sufficiently
	// corroborated, counting the claim's own source plus each entry in its
	// RelatedClaimIDs. 0 or 1 disables the requirement — a single source is
	// always sufficient. See IsSufficientlyCorroborated.
	MinCorroboratingSources int `json:"minCorroboratingSources,omitempty"`

	// CorroborationCategories restricts MinCorroboratingSources to specific
	// claim categories (e.g. only ClaimStatistical). An empty slice applies
	// the requirement to every category.
	CorroborationCategories []ClaimCategory `json:"corroborationCategories,omitempty"`

	// MaxClaimAge is how old a verified statistical claim's
	// StatisticalDetail.AsOfDate may be, relative to now, before it is
	// treated as stale (presented as current when it no longer is). Zero or
	// negative disables the requirement — no claim is ever too old. A claim
	// with no Statistical detail, or one whose AsOfDate is unset, is never
	// flagged: age is unknown, not stale. See IsFresh.
	MaxClaimAge time.Duration `json:"maxClaimAge,omitempty"`
}

// DefaultClaimsCriteria returns standard criteria.
func DefaultClaimsCriteria() ClaimsCriteria {
	return ClaimsCriteria{
		RequireAllVerified:            false,
		AllowSubjectiveWithDisclaimer: true,
		AllowNeedsReview:              true,
		MinReliabilityTier:            ReliabilityMedium,
		RequireReproducible:           false,
	}
}

// StrictClaimsCriteria returns strict criteria for high-stakes publications.
func StrictClaimsCriteria() ClaimsCriteria {
	return ClaimsCriteria{
		RequireAllVerified:            true,
		AllowSubjectiveWithDisclaimer: false,
		AllowNeedsReview:              false,
		MinReliabilityTier:            ReliabilityHigh,
		RequireReproducible:           true,
	}
}

// ClaimsDecision represents the validation decision for claims.
type ClaimsDecision struct {
	// Status is the decision outcome.
	Status ClaimsDecisionStatus `json:"status"`

	// Passed indicates if the claims validation passed.
	Passed bool `json:"passed"`

	// Rationale explains the decision.
	Rationale string `json:"rationale"`

	// Counts summarizes claims by verdict.
	Counts ClaimsCounts `json:"counts"`
}

// ClaimsDecisionStatus represents the decision outcome.
type ClaimsDecisionStatus string

const (
	// ClaimsDecisionPass indicates all claims are verified.
	ClaimsDecisionPass ClaimsDecisionStatus = "pass"

	// ClaimsDecisionConditional indicates some claims need review.
	ClaimsDecisionConditional ClaimsDecisionStatus = "conditional"

	// ClaimsDecisionFail indicates unverified or rejected claims exist.
	ClaimsDecisionFail ClaimsDecisionStatus = "fail"
)

// ClaimsCounts tracks claims by verdict.
type ClaimsCounts struct {
	Total       int `json:"total"`
	Verified    int `json:"verified"`
	Unverified  int `json:"unverified"`
	NeedsReview int `json:"needsReview"`
	Rejected    int `json:"rejected"`
}

// CountClaims counts claims by verdict.
func CountClaims(claims []Claim) ClaimsCounts {
	counts := ClaimsCounts{}
	for _, c := range claims {
		counts.Total++
		switch c.Verdict {
		case VerdictVerified:
			counts.Verified++
		case VerdictUnverified:
			counts.Unverified++
		case VerdictNeedsReview:
			counts.NeedsReview++
		case VerdictRejected:
			counts.Rejected++
		}
	}
	return counts
}

// EvaluateClaims checks claims against criteria.
func EvaluateClaims(claims []Claim, criteria ClaimsCriteria) ClaimsDecision {
	now := time.Now().UTC()
	counts := CountClaims(claims)
	insufficientCorroboration := countInsufficientlyCorroborated(claims, criteria)
	stale := countStaleClaims(claims, criteria, now)

	decision := ClaimsDecision{
		Counts: counts,
	}

	// Check for rejected claims
	if counts.Rejected > 0 {
		decision.Status = ClaimsDecisionFail
		decision.Passed = false
		decision.Rationale = formatRejectedRationale(counts.Rejected)
		return decision
	}

	// Check for unverified claims
	if counts.Unverified > 0 {
		decision.Status = ClaimsDecisionFail
		decision.Passed = false
		decision.Rationale = formatUnverifiedRationale(counts.Unverified)
		return decision
	}

	// Check needs-review claims, insufficiently corroborated claims, and
	// stale claims. A verified claim that doesn't clear
	// criteria.MinCorroboratingSources, or whose AsOfDate is older than
	// criteria.MaxClaimAge, is treated the same way as needs-review for
	// decision purposes: a conditional pass if criteria allows it,
	// otherwise a fail. This does not change the claim's own stored Verdict
	// or the Counts breakdown — only the report-level decision.
	if counts.NeedsReview > 0 || insufficientCorroboration > 0 || stale > 0 {
		if !criteria.AllowNeedsReview {
			decision.Status = ClaimsDecisionFail
			decision.Passed = false
			decision.Rationale = formatNeedsReviewRationale(counts.NeedsReview, insufficientCorroboration, stale)
			return decision
		}
		decision.Status = ClaimsDecisionConditional
		decision.Passed = true
		decision.Rationale = formatConditionalRationale(counts.NeedsReview, insufficientCorroboration, stale)
		return decision
	}

	// Check required categories
	if len(criteria.RequiredCategories) > 0 {
		missing := findMissingCategories(claims, criteria.RequiredCategories)
		if len(missing) > 0 {
			decision.Status = ClaimsDecisionFail
			decision.Passed = false
			decision.Rationale = formatMissingCategoriesRationale(missing)
			return decision
		}
	}

	decision.Status = ClaimsDecisionPass
	decision.Passed = true
	decision.Rationale = "All claims verified"
	return decision
}

func findMissingCategories(claims []Claim, required []ClaimCategory) []ClaimCategory {
	found := make(map[ClaimCategory]bool)
	for _, c := range claims {
		if c.Verdict == VerdictVerified {
			found[c.Category] = true
		}
	}

	var missing []ClaimCategory
	for _, cat := range required {
		if !found[cat] {
			missing = append(missing, cat)
		}
	}
	return missing
}

func formatRejectedRationale(count int) string {
	if count == 1 {
		return "1 claim rejected"
	}
	return itoa(count) + " claims rejected"
}

func formatUnverifiedRationale(count int) string {
	if count == 1 {
		return "1 claim unverified"
	}
	return itoa(count) + " claims unverified"
}

func formatNeedsReviewRationale(needsReview, insufficientCorroboration, stale int) string {
	return joinFailReasons(needsReview, insufficientCorroboration, stale) + " (not allowed by criteria)"
}

func joinFailReasons(needsReview, insufficientCorroboration, stale int) string {
	var parts []string
	if p := countedPhrase(needsReview, "1 claim needs review", "claims need review"); p != "" {
		parts = append(parts, p)
	}
	if p := countedPhrase(insufficientCorroboration, "1 claim lacks sufficient corroboration", "claims lack sufficient corroboration"); p != "" {
		parts = append(parts, p)
	}
	if p := countedPhrase(stale, "1 claim has a stale statistic", "claims have stale statistics"); p != "" {
		parts = append(parts, p)
	}
	return joinReasons(parts)
}

func formatConditionalRationale(needsReview, insufficientCorroboration, stale int) string {
	return "Passed with " + joinConditionalReasons(needsReview, insufficientCorroboration, stale)
}

func joinConditionalReasons(needsReview, insufficientCorroboration, stale int) string {
	var parts []string
	if p := countedPhrase(needsReview, "1 claim needing review", "claims needing review"); p != "" {
		parts = append(parts, p)
	}
	if p := countedPhrase(insufficientCorroboration, "1 claim lacking sufficient corroboration", "claims lacking sufficient corroboration"); p != "" {
		parts = append(parts, p)
	}
	if p := countedPhrase(stale, "1 claim with a stale statistic", "claims with stale statistics"); p != "" {
		parts = append(parts, p)
	}
	return joinReasons(parts)
}

// countedPhrase renders a count as a phrase, or "" if the count is zero.
// one is the exact 1-count phrase; many is appended after the count for
// 2+ (e.g. countedPhrase(2, "1 claim needs review", "claims need review")
// -> "2 claims need review").
func countedPhrase(n int, one, many string) string {
	switch {
	case n == 0:
		return ""
	case n == 1:
		return one
	default:
		return itoa(n) + " " + many
	}
}

// joinReasons joins 1-3 non-empty reason phrases into a single clause,
// preserving the exact single-reason and "X and Y" two-reason wording used
// before multi-reason rationales existed.
func joinReasons(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
	}
}

// IsSufficientlyCorroborated reports whether c has enough independent
// sources to satisfy criteria.MinCorroboratingSources, counting c's own
// source plus each entry in c.RelatedClaimIDs. Always true when the
// requirement is disabled (MinCorroboratingSources <= 1), or when
// criteria.CorroborationCategories is set and does not include c.Category.
func IsSufficientlyCorroborated(c Claim, criteria ClaimsCriteria) bool {
	if criteria.MinCorroboratingSources <= 1 {
		return true
	}
	if len(criteria.CorroborationCategories) > 0 && !categoryInScope(c.Category, criteria.CorroborationCategories) {
		return true
	}
	return 1+len(c.RelatedClaimIDs) >= criteria.MinCorroboratingSources
}

// IsFresh reports whether c's statistical value is recent enough to satisfy
// criteria.MaxClaimAge, measured from now. Always true when: the
// requirement is disabled (MaxClaimAge <= 0), c has no Statistical detail,
// or its AsOfDate is unset — age is unknown, not stale, so an unset date is
// never penalized.
func IsFresh(c Claim, criteria ClaimsCriteria, now time.Time) bool {
	if criteria.MaxClaimAge <= 0 {
		return true
	}
	if c.Statistical == nil || c.Statistical.AsOfDate == nil {
		return true
	}
	return now.Sub(*c.Statistical.AsOfDate) <= criteria.MaxClaimAge
}

// countStaleClaims counts verified claims whose statistic is not IsFresh.
// Only verified claims are checked, matching countInsufficientlyCorroborated:
// staleness governs whether a claim earns the verified label for decision
// purposes, so a claim that already isn't verified has nothing to add here.
func countStaleClaims(claims []Claim, criteria ClaimsCriteria, now time.Time) int {
	n := 0
	for _, c := range claims {
		if c.Verdict == VerdictVerified && !IsFresh(c, criteria, now) {
			n++
		}
	}
	return n
}

func categoryInScope(cat ClaimCategory, scope []ClaimCategory) bool {
	for _, s := range scope {
		if s == cat {
			return true
		}
	}
	return false
}

// countInsufficientlyCorroborated counts verified claims that do not meet
// criteria's corroboration requirement. Only verified claims are checked:
// corroboration governs whether a claim earns the verified label for
// decision purposes, so a claim that is already needs-review, rejected, or
// unverified has nothing to add here.
func countInsufficientlyCorroborated(claims []Claim, criteria ClaimsCriteria) int {
	n := 0
	for _, c := range claims {
		if c.Verdict == VerdictVerified && !IsSufficientlyCorroborated(c, criteria) {
			n++
		}
	}
	return n
}

func formatMissingCategoriesRationale(missing []ClaimCategory) string {
	if len(missing) == 1 {
		return "Required category missing verified claim: " + string(missing[0])
	}
	result := "Required categories missing verified claims:"
	for _, cat := range missing {
		result += " " + string(cat)
	}
	return result
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}
