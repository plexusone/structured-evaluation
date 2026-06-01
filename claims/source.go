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
	default:
		return ReliabilityLow
	}
}
