package claims

import (
	"testing"
	"time"
)

func TestNewStatisticalDetail(t *testing.T) {
	d := NewStatisticalDetail(4700000, "subscribers", PrecisionExact)
	if d.Value != 4700000 {
		t.Errorf("Value = %v, want 4700000", d.Value)
	}
	if d.Unit != "subscribers" {
		t.Errorf("Unit = %q, want %q", d.Unit, "subscribers")
	}
	if d.Precision != PrecisionExact {
		t.Errorf("Precision = %q, want %q", d.Precision, PrecisionExact)
	}
	if d.AsOfDate != nil {
		t.Errorf("AsOfDate = %v, want nil", d.AsOfDate)
	}
}

func TestStatisticalDetail_WithAsOfDate(t *testing.T) {
	asOf := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
	d := NewStatisticalDetail(10, "x", PrecisionApproximate).WithAsOfDate(asOf)

	if d.AsOfDate == nil || !d.AsOfDate.Equal(asOf) {
		t.Errorf("AsOfDate = %v, want %v", d.AsOfDate, asOf)
	}
}

func TestClaim_SetStatistical(t *testing.T) {
	claim := NewClaim("stat-1", "4.7M paid subscribers", ClaimStatistical, Location{Section: "metrics"})
	detail := NewStatisticalDetail(4700000, "subscribers", PrecisionExact)

	claim.SetStatistical(detail)

	if claim.Statistical == nil {
		t.Fatal("Statistical is nil after SetStatistical")
	}
	if claim.Statistical.Value != 4700000 {
		t.Errorf("Statistical.Value = %v, want 4700000", claim.Statistical.Value)
	}
}
