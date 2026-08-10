package claims

import "time"

// Validation describes how a claim is validated.
// Exactly one of External, Internal, Derived, or Subjective should be set.
type Validation struct {
	// Type indicates the validation approach.
	Type SourceType `json:"type"`

	// External contains details for externally-sourced claims.
	External *ExternalValidation `json:"external,omitempty"`

	// Internal contains details for internally-validated claims.
	Internal *InternalValidation `json:"internal,omitempty"`

	// Derived contains details for claims derived from other claims.
	Derived *DerivedValidation `json:"derived,omitempty"`

	// Subjective contains details for subjective estimates.
	Subjective *SubjectiveValidation `json:"subjective,omitempty"`
}

// ExternalValidation describes validation via external source.
type ExternalValidation struct {
	// URL is the source URL.
	URL string `json:"url"`

	// SourceType categorizes the source.
	SourceType ExternalSourceType `json:"sourceType"`

	// Role distinguishes how directly this source speaks for the claim
	// (primary / secondary-relay / secondary-analysis / self-reported),
	// independent of SourceType's general authority category. Optional and
	// empty by default so existing claims are unaffected: an empty Role is
	// not flagged by RequiresCorroboration — only an explicitly-set
	// secondary-analysis or self-reported role is.
	Role SourceRole `json:"sourceRole,omitempty"`

	// Reliability indicates the trustworthiness of the source.
	Reliability ReliabilityTier `json:"reliability"`

	// AccessedAt is when the URL was accessed.
	AccessedAt time.Time `json:"accessedAt,omitempty"`

	// Archived indicates if the URL is archived (e.g., Wayback Machine).
	Archived bool `json:"archived,omitempty"`

	// ArchiveURL is the archive URL if different from primary URL.
	ArchiveURL string `json:"archiveUrl,omitempty"`

	// QuotedText is the exact text from the source supporting the claim.
	QuotedText string `json:"quotedText,omitempty"`

	// VerifiedMatch indicates the claim text matches the source.
	VerifiedMatch bool `json:"verifiedMatch,omitempty"`
}

// InternalValidation describes validation via internal evidence.
type InternalValidation struct {
	// Method describes how validation was performed.
	Method InternalValidationMethod `json:"method"`

	// EvidencePath is the path to code, logs, or artifacts.
	EvidencePath string `json:"evidencePath,omitempty"`

	// EvidenceHash is a hash of the evidence file for integrity.
	EvidenceHash string `json:"evidenceHash,omitempty"`

	// Reproducible indicates whether the validation can be reproduced.
	Reproducible bool `json:"reproducible"`

	// ReproductionSteps describes how to reproduce the validation.
	ReproductionSteps string `json:"reproductionSteps,omitempty"`

	// ValidatedBy identifies who performed the validation.
	ValidatedBy string `json:"validatedBy,omitempty"`

	// ValidatedAt is when validation was performed.
	ValidatedAt time.Time `json:"validatedAt,omitempty"`

	// Environment describes the validation environment.
	Environment *ValidationEnvironment `json:"environment,omitempty"`

	// Output is the observed output that validates the claim.
	Output string `json:"output,omitempty"`
}

// ValidationEnvironment describes the environment where validation occurred.
type ValidationEnvironment struct {
	// Product is the product being tested.
	Product string `json:"product,omitempty"`

	// Version is the version tested.
	Version string `json:"version,omitempty"`

	// Configuration describes the configuration used.
	Configuration string `json:"configuration,omitempty"`

	// Platform is the platform (os, arch).
	Platform string `json:"platform,omitempty"`
}

// DerivedValidation describes claims derived from other claims.
type DerivedValidation struct {
	// SourceClaimIDs are the IDs of claims this is derived from.
	SourceClaimIDs []string `json:"sourceClaimIds"`

	// DerivationMethod describes how the claim was derived.
	DerivationMethod string `json:"derivationMethod"`

	// Formula is the calculation formula if applicable.
	Formula string `json:"formula,omitempty"`

	// Reasoning explains the derivation logic.
	Reasoning string `json:"reasoning,omitempty"`
}

// SubjectiveValidation describes subjective estimates.
type SubjectiveValidation struct {
	// Acknowledged indicates if this is labeled as an estimate in the document.
	Acknowledged bool `json:"acknowledged"`

	// Methodology describes any methodology used for the estimate.
	Methodology string `json:"methodology,omitempty"`

	// Recommendation suggests what to do with this claim.
	Recommendation SubjectiveRecommendation `json:"recommendation"`

	// Rationale explains why this is acceptable or not.
	Rationale string `json:"rationale,omitempty"`
}

// SubjectiveRecommendation indicates what to do with a subjective claim.
type SubjectiveRecommendation string

const (
	// RecommendKeepWithDisclaimer keeps the claim with a disclaimer.
	RecommendKeepWithDisclaimer SubjectiveRecommendation = "keep-with-disclaimer"

	// RecommendRemove removes the claim.
	RecommendRemove SubjectiveRecommendation = "remove"

	// RecommendFindSource suggests finding an external source.
	RecommendFindSource SubjectiveRecommendation = "find-source"

	// RecommendConvertToInternal suggests validating internally.
	RecommendConvertToInternal SubjectiveRecommendation = "convert-to-internal"
)

// NewExternalValidation creates a validation for an external source.
func NewExternalValidation(url string, sourceType ExternalSourceType) *Validation {
	return &Validation{
		Type: SourceExternal,
		External: &ExternalValidation{
			URL:         url,
			SourceType:  sourceType,
			Reliability: DefaultReliabilityForSourceType(sourceType),
			AccessedAt:  time.Now().UTC(),
		},
	}
}

// NewInternalValidation creates a validation for internal evidence.
func NewInternalValidation(method InternalValidationMethod, evidencePath string, reproducible bool) *Validation {
	return &Validation{
		Type: SourceInternal,
		Internal: &InternalValidation{
			Method:       method,
			EvidencePath: evidencePath,
			Reproducible: reproducible,
			ValidatedAt:  time.Now().UTC(),
		},
	}
}

// NewDerivedValidation creates a validation for a derived claim.
func NewDerivedValidation(sourceClaimIDs []string, method, formula string) *Validation {
	return &Validation{
		Type: SourceDerived,
		Derived: &DerivedValidation{
			SourceClaimIDs:   sourceClaimIDs,
			DerivationMethod: method,
			Formula:          formula,
		},
	}
}

// NewSubjectiveValidation creates a validation for a subjective estimate.
func NewSubjectiveValidation(acknowledged bool, recommendation SubjectiveRecommendation) *Validation {
	return &Validation{
		Type: SourceSubjective,
		Subjective: &SubjectiveValidation{
			Acknowledged:   acknowledged,
			Recommendation: recommendation,
		},
	}
}
