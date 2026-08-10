package claims

import "testing"

func TestDefaultReliabilityForSourceType(t *testing.T) {
	tests := []struct {
		name     string
		input    ExternalSourceType
		expected ReliabilityTier
	}{
		{"nvd is authoritative", ExternalNVD, ReliabilityAuthoritative},
		{"vendor advisory is authoritative", ExternalVendorAdvisory, ReliabilityAuthoritative},
		{"framework official is authoritative", ExternalFramework, ReliabilityAuthoritative},
		{"peer reviewed is high", ExternalPeerReviewed, ReliabilityHigh},
		{"reputable vendor is high", ExternalReputableVendor, ReliabilityHigh},
		{"api is high", ExternalAPI, ReliabilityHigh},
		{"community is medium", ExternalCommunity, ReliabilityMedium},
		{"aggregator is low", ExternalAggregator, ReliabilityLow},
		{"unknown defaults to low", ExternalSourceType("unknown"), ReliabilityLow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultReliabilityForSourceType(tt.input)
			if got != tt.expected {
				t.Errorf("DefaultReliabilityForSourceType(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestReliabilityTier_AggregatorIsRejectedByDefault(t *testing.T) {
	tier := DefaultReliabilityForSourceType(ExternalAggregator)
	if tier.IsAcceptable() {
		t.Errorf("aggregator sources should not be auto-acceptable, got tier %q", tier)
	}
	if tier.RequiresReview() {
		t.Errorf("aggregator sources should reject outright, not just require review, got tier %q", tier)
	}
}
