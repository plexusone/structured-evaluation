package claims

// ClaimCategory categorizes the type of claim.
type ClaimCategory string

const (
	// ClaimMetadata is metadata about the subject (e.g., identifiers, versions).
	ClaimMetadata ClaimCategory = "metadata"

	// ClaimTechnicalFinding is a technical observation or finding.
	ClaimTechnicalFinding ClaimCategory = "technical-finding"

	// ClaimFrameworkMapping maps to a standard framework (e.g., taxonomy IDs).
	ClaimFrameworkMapping ClaimCategory = "framework-mapping"

	// ClaimRiskAssessment is a risk or impact assessment.
	ClaimRiskAssessment ClaimCategory = "risk-assessment"

	// ClaimTimeline is a temporal claim (dates, events).
	ClaimTimeline ClaimCategory = "timeline"

	// ClaimStatistical is a statistical or numeric claim.
	ClaimStatistical ClaimCategory = "statistical"

	// ClaimGuidance is a recommendation or guidance claim.
	ClaimGuidance ClaimCategory = "guidance"

	// ClaimAttribution credits a source or author.
	ClaimAttribution ClaimCategory = "attribution"
)

// Location identifies where a claim appears in a document.
type Location struct {
	// Section is the section heading or identifier.
	Section string `json:"section,omitempty"`

	// Line is the line number (1-indexed).
	Line int `json:"line,omitempty"`

	// StartOffset is the character offset from start of document.
	StartOffset int `json:"startOffset,omitempty"`

	// EndOffset is the ending character offset.
	EndOffset int `json:"endOffset,omitempty"`
}

// Claim represents a factual claim extracted from a document.
type Claim struct {
	// ID is the unique identifier for this claim.
	ID string `json:"id"`

	// Text is the exact claim text from the document.
	Text string `json:"text"`

	// Location identifies where the claim appears.
	Location Location `json:"location"`

	// Category categorizes the type of claim.
	Category ClaimCategory `json:"category"`

	// Validation describes how the claim is validated.
	Validation *Validation `json:"validation,omitempty"`

	// Verdict is the validation result.
	Verdict Verdict `json:"verdict"`

	// Rationale explains the verdict.
	Rationale string `json:"rationale,omitempty"`

	// RelatedClaimIDs are IDs of related claims.
	RelatedClaimIDs []string `json:"relatedClaimIds,omitempty"`
}

// NewClaim creates a new claim with the given details.
func NewClaim(id, text string, category ClaimCategory, location Location) *Claim {
	return &Claim{
		ID:       id,
		Text:     text,
		Category: category,
		Location: location,
		Verdict:  VerdictUnverified,
	}
}

// SetValidation sets the validation and computes the verdict.
func (c *Claim) SetValidation(v *Validation) *Claim {
	c.Validation = v
	c.Verdict = DetermineVerdict(v)
	return c
}

// SetVerdict sets the verdict with rationale.
func (c *Claim) SetVerdict(verdict Verdict, rationale string) *Claim {
	c.Verdict = verdict
	c.Rationale = rationale
	return c
}

// AddRelatedClaim adds a related claim ID.
func (c *Claim) AddRelatedClaim(claimID string) *Claim {
	c.RelatedClaimIDs = append(c.RelatedClaimIDs, claimID)
	return c
}

// IsVerified returns true if the claim is verified.
func (c *Claim) IsVerified() bool {
	return c.Verdict == VerdictVerified
}

// IsBlocking returns true if the claim blocks publication.
func (c *Claim) IsBlocking() bool {
	return c.Verdict.IsBlocking()
}

// NeedsReview returns true if the claim requires human review.
func (c *Claim) NeedsReview() bool {
	return c.Verdict == VerdictNeedsReview
}
