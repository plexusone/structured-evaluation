package claims

import (
	"fmt"
	"time"
)

// ClaimsReport is the report for claim validation.
type ClaimsReport struct {
	// Schema is the JSON Schema URL.
	Schema string `json:"$schema,omitempty"`

	// Metadata contains report identification.
	Metadata ClaimsMetadata `json:"metadata"`

	// Claims are the extracted and validated claims.
	Claims []Claim `json:"claims"`

	// Summary provides aggregated statistics.
	Summary ClaimsSummary `json:"summary"`

	// Criteria defines the pass requirements.
	Criteria ClaimsCriteria `json:"criteria"`

	// Decision is the validation outcome.
	Decision ClaimsDecision `json:"decision"`
}

// ClaimsMetadata contains report identification.
type ClaimsMetadata struct {
	// Document is the filename or path being validated.
	Document string `json:"document"`

	// DocumentID is the document identifier.
	DocumentID string `json:"documentId,omitempty"`

	// DocumentTitle is the document title.
	DocumentTitle string `json:"documentTitle,omitempty"`

	// DocumentVersion is the document version.
	DocumentVersion string `json:"documentVersion,omitempty"`

	// GeneratedAt is when the report was created.
	GeneratedAt time.Time `json:"generatedAt"`

	// GeneratedBy identifies what created this report.
	GeneratedBy string `json:"generatedBy,omitempty"`

	// ValidatedBy identifies who validated the claims.
	ValidatedBy string `json:"validatedBy,omitempty"`
}

// ClaimsSummary provides aggregated statistics.
type ClaimsSummary struct {
	// Counts by verdict.
	Counts ClaimsCounts `json:"counts"`

	// ByCategory counts claims by category.
	ByCategory map[ClaimCategory]int `json:"byCategory,omitempty"`

	// BySourceType counts claims by validation source type.
	BySourceType map[SourceType]int `json:"bySourceType,omitempty"`

	// ByReliability counts external claims by reliability tier.
	ByReliability map[ReliabilityTier]int `json:"byReliability,omitempty"`

	// UnverifiedClaims lists IDs of unverified claims.
	UnverifiedClaims []string `json:"unverifiedClaims,omitempty"`

	// NeedsReviewClaims lists IDs of claims needing review.
	NeedsReviewClaims []string `json:"needsReviewClaims,omitempty"`

	// RejectedClaims lists IDs of rejected claims.
	RejectedClaims []string `json:"rejectedClaims,omitempty"`
}

// NewClaimsReport creates a new claims report.
func NewClaimsReport(document string) *ClaimsReport {
	return &ClaimsReport{
		Metadata: ClaimsMetadata{
			Document:    document,
			GeneratedAt: time.Now().UTC(),
			GeneratedBy: "structured-evaluation",
		},
		Claims:   []Claim{},
		Criteria: DefaultClaimsCriteria(),
	}
}

// AddClaim adds a claim to the report.
func (r *ClaimsReport) AddClaim(c Claim) {
	r.Claims = append(r.Claims, c)
}

// SetCriteria sets the validation criteria.
func (r *ClaimsReport) SetCriteria(criteria ClaimsCriteria) {
	r.Criteria = criteria
}

// Evaluate computes the summary and decision.
func (r *ClaimsReport) Evaluate() ClaimsDecision {
	r.computeSummary()
	r.Decision = EvaluateClaims(r.Claims, r.Criteria)
	return r.Decision
}

func (r *ClaimsReport) computeSummary() {
	r.Summary = ClaimsSummary{
		Counts:            CountClaims(r.Claims),
		ByCategory:        make(map[ClaimCategory]int),
		BySourceType:      make(map[SourceType]int),
		ByReliability:     make(map[ReliabilityTier]int),
		UnverifiedClaims:  []string{},
		NeedsReviewClaims: []string{},
		RejectedClaims:    []string{},
	}

	for _, c := range r.Claims {
		r.Summary.ByCategory[c.Category]++

		if c.Validation != nil {
			r.Summary.BySourceType[c.Validation.Type]++

			if c.Validation.Type == SourceExternal && c.Validation.External != nil {
				r.Summary.ByReliability[c.Validation.External.Reliability]++
			}
		}

		switch c.Verdict {
		case VerdictUnverified:
			r.Summary.UnverifiedClaims = append(r.Summary.UnverifiedClaims, c.ID)
		case VerdictNeedsReview:
			r.Summary.NeedsReviewClaims = append(r.Summary.NeedsReviewClaims, c.ID)
		case VerdictRejected:
			r.Summary.RejectedClaims = append(r.Summary.RejectedClaims, c.ID)
		}
	}
}

// Finalize computes all derived fields.
func (r *ClaimsReport) Finalize() {
	r.Evaluate()
}

// GetClaim returns a claim by ID, or nil if not found.
func (r *ClaimsReport) GetClaim(claimID string) *Claim {
	for i := range r.Claims {
		if r.Claims[i].ID == claimID {
			return &r.Claims[i]
		}
	}
	return nil
}

// GenerateSummaryText creates a human-readable summary.
func (r *ClaimsReport) GenerateSummaryText() string {
	counts := r.Summary.Counts
	summary := fmt.Sprintf("Claims: %d total, %d verified, %d unverified, %d needs review, %d rejected. ",
		counts.Total, counts.Verified, counts.Unverified, counts.NeedsReview, counts.Rejected)

	summary += "Decision: " + string(r.Decision.Status) + "."
	return summary
}

// IsPassing returns true if the report passes validation.
func (r *ClaimsReport) IsPassing() bool {
	return r.Decision.Passed
}

// ValidateDerivedClaims checks that all derived claims reference verified source claims.
func (r *ClaimsReport) ValidateDerivedClaims() {
	claimMap := make(map[string]*Claim)
	for i := range r.Claims {
		claimMap[r.Claims[i].ID] = &r.Claims[i]
	}

	for i := range r.Claims {
		c := &r.Claims[i]
		if c.Validation != nil && c.Validation.Type == SourceDerived && c.Validation.Derived != nil {
			allSourcesVerified := true
			for _, sourceID := range c.Validation.Derived.SourceClaimIDs {
				source, ok := claimMap[sourceID]
				if !ok || !source.IsVerified() {
					allSourcesVerified = false
					break
				}
			}
			if !allSourcesVerified {
				c.Verdict = VerdictUnverified
				c.Rationale = "One or more source claims are not verified"
			}
		}
	}
}
