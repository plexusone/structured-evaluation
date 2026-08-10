// Package claims provides types for claim extraction and source validation.
// This enables verification that factual claims in documents are properly
// sourced (external references) or objectively validated (internal evidence).
package claims

// SourceType identifies how a claim is validated.
type SourceType string

const (
	// SourceExternal indicates validation via external URL reference.
	SourceExternal SourceType = "external"

	// SourceInternal indicates validation via internal evidence (code, lab tests).
	SourceInternal SourceType = "internal"

	// SourceDerived indicates the claim is calculated from other validated claims.
	SourceDerived SourceType = "derived"

	// SourceSubjective indicates an estimate without objective backing.
	SourceSubjective SourceType = "subjective"
)

// ExternalSourceType categorizes the authority of external sources.
type ExternalSourceType string

const (
	// ExternalNVD is the NIST National Vulnerability Database.
	ExternalNVD ExternalSourceType = "nvd"

	// ExternalVendorAdvisory is an official vendor advisory.
	ExternalVendorAdvisory ExternalSourceType = "vendor-advisory"

	// ExternalFramework is official framework documentation (MITRE, OWASP, CWE).
	ExternalFramework ExternalSourceType = "framework-official"

	// ExternalPeerReviewed is a peer-reviewed publication.
	ExternalPeerReviewed ExternalSourceType = "peer-reviewed"

	// ExternalReputableVendor is from a reputable vendor (e.g., research firms).
	ExternalReputableVendor ExternalSourceType = "reputable-vendor"

	// ExternalCommunity is from community sources (blogs, forums).
	ExternalCommunity ExternalSourceType = "community"

	// ExternalAPI is from a public API (e.g., FIRST.org EPSS API).
	ExternalAPI ExternalSourceType = "api"

	// ExternalAggregator is a third-party content-roundup or "stats blog" site
	// that reposts figures without independent reporting or a traceable
	// primary source (e.g. AI-generated SEO stats pages). Distinct from
	// ExternalCommunity: a community source (a named blog post, a forum
	// thread) is still someone's own account: an aggregator has no original
	// reporting to fall back on, so it defaults to auto-reject rather than
	// requires-review.
	ExternalAggregator ExternalSourceType = "aggregator"
)

// ReliabilityTier indicates the trustworthiness of a source.
type ReliabilityTier string

const (
	// ReliabilityAuthoritative is an official, authoritative source (auto-accept).
	ReliabilityAuthoritative ReliabilityTier = "authoritative"

	// ReliabilityHigh is a highly reputable source (auto-accept).
	ReliabilityHigh ReliabilityTier = "high"

	// ReliabilityMedium is a moderately reputable source (requires review).
	ReliabilityMedium ReliabilityTier = "medium"

	// ReliabilityLow is an unverified or low-reputation source (reject).
	ReliabilityLow ReliabilityTier = "low"
)

// IsAcceptable returns true if the reliability tier is acceptable without review.
func (r ReliabilityTier) IsAcceptable() bool {
	return r == ReliabilityAuthoritative || r == ReliabilityHigh
}

// RequiresReview returns true if the reliability tier requires human review.
func (r ReliabilityTier) RequiresReview() bool {
	return r == ReliabilityMedium
}

// SourceRole distinguishes how directly an external source speaks for the
// claim, independent of ExternalSourceType (which is about the source's
// general authority/category). Two "reputable-vendor" sources can carry very
// different trust: a wire report quoting an earnings call verbatim is
// SourceRolePrimary-adjacent; an outlet's own synthesis across several other
// reports is SourceRoleSecondaryAnalysis and is exactly the shape that
// produced a false-positive "verified" figure ($3B ARR from a secondary
// synthesis, contradicted by the primary source's own $500M figure) — the
// incident that motivated this field.
type SourceRole string

const (
	// SourceRolePrimary is the entity the claim is about, speaking for
	// itself (a company blog post, an official filing, a direct quote from
	// an executive).
	SourceRolePrimary SourceRole = "primary"

	// SourceRoleSecondaryRelay is a reputable outlet directly reporting a
	// primary statement (e.g. a wire service relaying an earnings-call
	// quote) — one hop from primary, with low synthesis risk.
	SourceRoleSecondaryRelay SourceRole = "secondary-relay"

	// SourceRoleSecondaryAnalysis is an outlet's own synthesis, estimate, or
	// aggregation across multiple other reports — not a direct relay of a
	// single primary statement. Higher risk of drift from the underlying
	// fact; requires corroboration to be published as verified.
	SourceRoleSecondaryAnalysis SourceRole = "secondary-analysis"

	// SourceRoleSelfReported is the claim's subject describing itself in a
	// context with a promotional incentive (marketing page, press release
	// superlative) rather than an audited disclosure. Requires
	// corroboration to be published as verified.
	SourceRoleSelfReported SourceRole = "self-reported"
)

// RequiresCorroboration reports whether a claim sourced with this role needs
// at least one independent corroborating source before it can be published
// as verified. Primary and secondary-relay sources speak for themselves;
// secondary-analysis and self-reported sources carry synthesis or incentive
// risk that a second independent source mitigates.
func (r SourceRole) RequiresCorroboration() bool {
	return r == SourceRoleSecondaryAnalysis || r == SourceRoleSelfReported
}

// InternalValidationMethod describes how internal validation was performed.
type InternalValidationMethod string

const (
	// MethodCodeExecution indicates the claim was validated by running code.
	MethodCodeExecution InternalValidationMethod = "code-execution"

	// MethodLabTesting indicates validation via controlled lab testing.
	MethodLabTesting InternalValidationMethod = "lab-testing"

	// MethodCodeReview indicates validation via code review/inspection.
	MethodCodeReview InternalValidationMethod = "code-review"

	// MethodLogAnalysis indicates validation via log analysis.
	MethodLogAnalysis InternalValidationMethod = "log-analysis"

	// MethodCalculation indicates the claim is a calculation from other data.
	MethodCalculation InternalValidationMethod = "calculation"

	// MethodObservation indicates direct observation of behavior.
	MethodObservation InternalValidationMethod = "observation"
)

// DefaultReliabilityForSourceType returns the default reliability tier for a source type.
func DefaultReliabilityForSourceType(st ExternalSourceType) ReliabilityTier {
	switch st {
	case ExternalNVD, ExternalVendorAdvisory, ExternalFramework:
		return ReliabilityAuthoritative
	case ExternalPeerReviewed, ExternalReputableVendor, ExternalAPI:
		return ReliabilityHigh
	case ExternalCommunity:
		return ReliabilityMedium
	case ExternalAggregator:
		return ReliabilityLow
	default:
		return ReliabilityLow
	}
}
