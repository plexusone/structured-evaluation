package claims

import (
	"encoding/json"
	"testing"
)

func TestClaimsReport_Basic(t *testing.T) {
	report := NewClaimsReport("document.md")
	report.Metadata.DocumentTitle = "Test Document"
	report.Metadata.DocumentVersion = "1.0.0"

	// Add a verified external claim
	claim1 := NewClaim("claim-1", "The score is 8.8", ClaimMetadata, Location{Section: "Header", Line: 1})
	claim1.SetValidation(NewExternalValidation("https://example.com/source", ExternalNVD))
	report.AddClaim(*claim1)

	// Add an internally validated claim
	claim2 := NewClaim("claim-2", "The code outputs X", ClaimTechnicalFinding, Location{Section: "Analysis", Line: 10})
	claim2.SetValidation(NewInternalValidation(MethodCodeExecution, "test.go", true))
	report.AddClaim(*claim2)

	// Add a derived claim
	claim3 := NewClaim("claim-3", "Combined result is Y", ClaimStatistical, Location{Section: "Summary", Line: 50})
	claim3.SetValidation(NewDerivedValidation([]string{"claim-1", "claim-2"}, "calculation", "claim-1 + claim-2"))
	report.AddClaim(*claim3)

	report.ValidateDerivedClaims()
	report.Finalize()

	if report.Summary.Counts.Total != 3 {
		t.Errorf("Expected 3 total claims, got %d", report.Summary.Counts.Total)
	}

	if report.Summary.Counts.Verified != 3 {
		t.Errorf("Expected 3 verified claims, got %d", report.Summary.Counts.Verified)
	}

	if !report.IsPassing() {
		t.Errorf("Expected report to pass, got %s: %s", report.Decision.Status, report.Decision.Rationale)
	}
}

func TestClaimsReport_Unverified(t *testing.T) {
	report := NewClaimsReport("document.md")

	// Add an unverified claim
	claim1 := NewClaim("claim-1", "Some unsourced claim", ClaimStatistical, Location{Section: "Stats", Line: 1})
	// No validation set
	report.AddClaim(*claim1)

	report.Finalize()

	if report.Summary.Counts.Unverified != 1 {
		t.Errorf("Expected 1 unverified claim, got %d", report.Summary.Counts.Unverified)
	}

	if report.IsPassing() {
		t.Errorf("Expected report to fail with unverified claims")
	}

	if report.Decision.Status != ClaimsDecisionFail {
		t.Errorf("Expected fail status, got %s", report.Decision.Status)
	}
}

func TestClaimsReport_SubjectiveWithDisclaimer(t *testing.T) {
	report := NewClaimsReport("document.md")

	// Add a subjective claim that is acknowledged
	claim1 := NewClaim("claim-1", "Approximately 50% of users", ClaimStatistical, Location{Section: "Stats", Line: 1})
	validation := NewSubjectiveValidation(true, RecommendKeepWithDisclaimer)
	validation.Subjective.Methodology = "Expert estimate"
	claim1.SetValidation(validation)
	report.AddClaim(*claim1)

	report.Finalize()

	if report.Summary.Counts.Verified != 1 {
		t.Errorf("Expected 1 verified claim (acknowledged subjective), got %d", report.Summary.Counts.Verified)
	}

	if !report.IsPassing() {
		t.Errorf("Expected report to pass with acknowledged subjective claim")
	}
}

func TestClaimsReport_DerivedFromUnverified(t *testing.T) {
	report := NewClaimsReport("document.md")

	// Add an unverified claim
	claim1 := NewClaim("claim-1", "Base value", ClaimMetadata, Location{Section: "Data", Line: 1})
	// No validation - unverified
	report.AddClaim(*claim1)

	// Add a derived claim that depends on the unverified claim
	claim2 := NewClaim("claim-2", "Derived value", ClaimStatistical, Location{Section: "Summary", Line: 10})
	claim2.SetValidation(NewDerivedValidation([]string{"claim-1"}, "calculation", "claim-1 * 2"))
	report.AddClaim(*claim2)

	report.ValidateDerivedClaims()
	report.Finalize()

	// The derived claim should now be unverified because its source is unverified
	derivedClaim := report.GetClaim("claim-2")
	if derivedClaim.Verdict != VerdictUnverified {
		t.Errorf("Expected derived claim to be unverified when source is unverified, got %s", derivedClaim.Verdict)
	}
}

func TestClaimsReport_JSON(t *testing.T) {
	report := NewClaimsReport("document.md")
	report.Metadata.DocumentTitle = "Test"

	claim := NewClaim("claim-1", "Test claim", ClaimMetadata, Location{Section: "Test", Line: 1})
	claim.SetValidation(NewExternalValidation("https://example.com", ExternalFramework))
	report.AddClaim(*claim)

	report.Finalize()

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal report: %v", err)
	}

	var parsed ClaimsReport
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal report: %v", err)
	}

	if len(parsed.Claims) != 1 {
		t.Errorf("Expected 1 claim after round-trip, got %d", len(parsed.Claims))
	}

	if parsed.Claims[0].Validation.Type != SourceExternal {
		t.Errorf("Expected external validation type, got %s", parsed.Claims[0].Validation.Type)
	}
}

func TestDetermineVerdict(t *testing.T) {
	tests := []struct {
		name       string
		validation *Validation
		expected   Verdict
	}{
		{
			name:       "nil validation",
			validation: nil,
			expected:   VerdictUnverified,
		},
		{
			name:       "external authoritative",
			validation: NewExternalValidation("https://nvd.nist.gov/...", ExternalNVD),
			expected:   VerdictVerified,
		},
		{
			name: "external low reliability",
			validation: &Validation{
				Type: SourceExternal,
				External: &ExternalValidation{
					URL:         "https://random-blog.com",
					SourceType:  ExternalCommunity,
					Reliability: ReliabilityLow,
				},
			},
			expected: VerdictRejected,
		},
		{
			name:       "internal reproducible",
			validation: NewInternalValidation(MethodCodeExecution, "test.go", true),
			expected:   VerdictVerified,
		},
		{
			name:       "internal not reproducible",
			validation: NewInternalValidation(MethodCodeExecution, "test.go", false),
			expected:   VerdictNeedsReview,
		},
		{
			name:       "subjective acknowledged",
			validation: NewSubjectiveValidation(true, RecommendKeepWithDisclaimer),
			expected:   VerdictVerified,
		},
		{
			name:       "subjective not acknowledged",
			validation: NewSubjectiveValidation(false, RecommendKeepWithDisclaimer),
			expected:   VerdictNeedsReview,
		},
		{
			name:       "subjective remove",
			validation: NewSubjectiveValidation(true, RecommendRemove),
			expected:   VerdictRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verdict := DetermineVerdict(tt.validation)
			if verdict != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, verdict)
			}
		})
	}
}
