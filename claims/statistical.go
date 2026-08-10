package claims

import "time"

// Precision indicates how exact a statistical claim's value is.
type Precision string

const (
	// PrecisionExact is a value stated exactly by the source (e.g. "4.7 million").
	PrecisionExact Precision = "exact"

	// PrecisionApproximate is a value the source itself qualifies as
	// approximate (e.g. "1M+", "~50%", "over 20,000").
	PrecisionApproximate Precision = "approximate"

	// PrecisionEstimated is a third-party or analyst estimate, not a figure
	// the primary source states directly.
	PrecisionEstimated Precision = "estimated"

	// PrecisionRange is a bounded range rather than a single point value
	// (e.g. "9.6 to 2.4 days"). Store the more notable end in Value and
	// describe the range in the claim Text.
	PrecisionRange Precision = "range"
)

// StatisticalDetail captures the structured numeric value behind a
// ClaimStatistical claim. It is kept separate from Claim.Text (which remains
// the free-text rendering, e.g. "4.7M paid subscribers") so the number
// itself is queryable and can be checked against
// ExternalValidation.QuotedText programmatically, instead of only existing
// inside a formatted string.
//
// AsOfDate is distinct from ExternalValidation.AccessedAt: AccessedAt is
// when the source URL was crawled; AsOfDate is when the underlying fact was
// true (e.g. the earnings period the figure describes, or the date a public
// statement was made). Case-study and market-signal claims routinely need
// both — a stat crawled today can describe a fact from months earlier.
type StatisticalDetail struct {
	// Value is the numeric value of the claim.
	Value float64 `json:"value"`

	// Unit is the unit of measurement (e.g. "%", "USD", "users"). Empty for
	// dimensionless counts.
	Unit string `json:"unit,omitempty"`

	// Precision indicates how exact the value is.
	Precision Precision `json:"precision,omitempty"`

	// AsOfDate is when the underlying fact was true, if known and different
	// from when the source was accessed.
	AsOfDate *time.Time `json:"asOfDate,omitempty"`
}

// NewStatisticalDetail creates a StatisticalDetail with the given value,
// unit, and precision. Set AsOfDate directly on the result if known.
func NewStatisticalDetail(value float64, unit string, precision Precision) *StatisticalDetail {
	return &StatisticalDetail{
		Value:     value,
		Unit:      unit,
		Precision: precision,
	}
}

// WithAsOfDate sets AsOfDate and returns the receiver for chaining.
func (d *StatisticalDetail) WithAsOfDate(t time.Time) *StatisticalDetail {
	d.AsOfDate = &t
	return d
}
