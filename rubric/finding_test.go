package rubric

import "testing"

func TestWorstSeverity(t *testing.T) {
	tests := []struct {
		name     string
		findings []Finding
		want     Severity
	}{
		{"empty", nil, ""},
		{"single", []Finding{{Severity: SeverityLow}}, SeverityLow},
		{
			"picks the highest weight, not declaration order",
			[]Finding{{Severity: SeverityLow}, {Severity: SeverityCritical}, {Severity: SeverityMedium}},
			SeverityCritical,
		},
		{
			"highest weight regardless of position",
			[]Finding{{Severity: SeverityHigh}, {Severity: SeverityInfo}},
			SeverityHigh,
		},
		{
			"tie keeps the first-seen severity at that weight",
			[]Finding{{Severity: SeverityMedium}, {Severity: SeverityMedium}},
			SeverityMedium,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WorstSeverity(tt.findings); got != tt.want {
				t.Errorf("WorstSeverity() = %q, want %q", got, tt.want)
			}
		})
	}
}
