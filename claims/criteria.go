package claims

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
	counts := CountClaims(claims)

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

	// Check needs-review claims
	if counts.NeedsReview > 0 {
		if !criteria.AllowNeedsReview {
			decision.Status = ClaimsDecisionFail
			decision.Passed = false
			decision.Rationale = formatNeedsReviewRationale(counts.NeedsReview)
			return decision
		}
		decision.Status = ClaimsDecisionConditional
		decision.Passed = true
		decision.Rationale = formatConditionalRationale(counts.NeedsReview)
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

func formatNeedsReviewRationale(count int) string {
	if count == 1 {
		return "1 claim needs review (not allowed by criteria)"
	}
	return itoa(count) + " claims need review (not allowed by criteria)"
}

func formatConditionalRationale(count int) string {
	if count == 1 {
		return "Passed with 1 claim needing review"
	}
	return "Passed with " + itoa(count) + " claims needing review"
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
