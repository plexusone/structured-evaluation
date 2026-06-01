package claims

// Verdict represents the validation result for a claim.
type Verdict string

const (
	// VerdictVerified indicates the claim is verified.
	VerdictVerified Verdict = "verified"

	// VerdictUnverified indicates the claim could not be verified.
	VerdictUnverified Verdict = "unverified"

	// VerdictNeedsReview indicates the claim requires human review.
	VerdictNeedsReview Verdict = "needs-review"

	// VerdictRejected indicates the claim should be removed.
	VerdictRejected Verdict = "rejected"
)

// IsPassing returns true if the verdict is acceptable for publication.
func (v Verdict) IsPassing() bool {
	return v == VerdictVerified
}

// IsBlocking returns true if the verdict blocks publication.
func (v Verdict) IsBlocking() bool {
	return v == VerdictUnverified || v == VerdictRejected
}

// Icon returns the emoji icon for the verdict.
func (v Verdict) Icon() string {
	switch v {
	case VerdictVerified:
		return "✓"
	case VerdictUnverified:
		return "✗"
	case VerdictNeedsReview:
		return "?"
	case VerdictRejected:
		return "✗"
	default:
		return "○"
	}
}

// DetermineVerdict computes the verdict based on validation.
func DetermineVerdict(v *Validation) Verdict {
	if v == nil {
		return VerdictUnverified
	}

	switch v.Type {
	case SourceExternal:
		if v.External == nil {
			return VerdictUnverified
		}
		if v.External.URL == "" {
			return VerdictUnverified
		}
		if v.External.Reliability.RequiresReview() {
			return VerdictNeedsReview
		}
		if v.External.Reliability == ReliabilityLow {
			return VerdictRejected
		}
		return VerdictVerified

	case SourceInternal:
		if v.Internal == nil {
			return VerdictUnverified
		}
		if v.Internal.EvidencePath == "" && v.Internal.Output == "" {
			return VerdictUnverified
		}
		if !v.Internal.Reproducible {
			return VerdictNeedsReview
		}
		return VerdictVerified

	case SourceDerived:
		if v.Derived == nil || len(v.Derived.SourceClaimIDs) == 0 {
			return VerdictUnverified
		}
		// Derived claims need source claims to be verified (checked at report level)
		return VerdictVerified

	case SourceSubjective:
		if v.Subjective == nil {
			return VerdictRejected
		}
		switch v.Subjective.Recommendation {
		case RecommendKeepWithDisclaimer:
			if v.Subjective.Acknowledged {
				return VerdictVerified
			}
			return VerdictNeedsReview
		case RecommendRemove:
			return VerdictRejected
		default:
			return VerdictNeedsReview
		}

	default:
		return VerdictUnverified
	}
}
