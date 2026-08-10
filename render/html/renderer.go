// Package html renders claims reports as self-contained HTML pages.
//
// The output is a single standalone document with inline CSS and no external
// dependencies, so it can be written to a file and opened directly or emailed
// as an artifact. Claims are grouped by verdict (verified, needs-review,
// rejected, unverified) with each claim's value, source, reliability tier,
// and rationale (commentary) shown.
package html

import (
	"io"
	"strconv"

	"github.com/plexusone/structured-evaluation/claims"
)

// Renderer renders a claims report as a standalone HTML page.
type Renderer struct {
	w io.Writer
}

// New creates a new HTML renderer writing to w.
func New(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

// RenderClaims writes the claims report as a self-contained HTML document.
func (r *Renderer) RenderClaims(report *claims.ClaimsReport) error {
	return claimsTemplate.Execute(r.w, buildView(report))
}

// --- view model -------------------------------------------------------------

type reportView struct {
	Title         string
	Document      string
	GeneratedAt   string
	ValidatedBy   string
	DecisionLabel string
	DecisionClass string
	DecisionText  string
	Counts        claims.ClaimsCounts
	Groups        []groupView
}

type groupView struct {
	Verdict string
	Label   string
	Icon    string
	Class   string
	Claims  []claimView
}

type claimView struct {
	ID          string
	Text        string
	Category    string
	Rationale   string
	Icon        string
	Class       string
	HasStat     bool
	Value       string
	Unit        string
	Precision   string
	AsOfDate    string
	HasSource   bool
	URL         string
	SourceType  string
	Reliability string
	RelClass    string
	QuotedText  string
	Verified    bool
	AccessedAt  string
	IsDerived   bool
	DerivedFrom string
}

// verdictOrder controls the display order of verdict groups.
var verdictOrder = []struct {
	verdict claims.Verdict
	label   string
	icon    string
	class   string
}{
	{claims.VerdictVerified, "Verified", "✓", "verified"},
	{claims.VerdictNeedsReview, "Needs review", "?", "needs-review"},
	{claims.VerdictRejected, "Rejected", "✗", "rejected"},
	{claims.VerdictUnverified, "Unverified", "✗", "unverified"},
}

func buildView(report *claims.ClaimsReport) reportView {
	title := report.Metadata.DocumentTitle
	if title == "" {
		title = report.Metadata.Document
	}
	if title == "" {
		title = "Claims Report"
	}

	v := reportView{
		Title:         title,
		Document:      report.Metadata.Document,
		ValidatedBy:   report.Metadata.ValidatedBy,
		DecisionLabel: decisionLabel(report.Decision.Status),
		DecisionClass: string(report.Decision.Status),
		DecisionText:  report.Decision.Rationale,
		Counts:        report.Summary.Counts,
	}
	if !report.Metadata.GeneratedAt.IsZero() {
		v.GeneratedAt = report.Metadata.GeneratedAt.Format("2006-01-02")
	}

	for _, g := range verdictOrder {
		group := groupView{Verdict: string(g.verdict), Label: g.label, Icon: g.icon, Class: g.class}
		for i := range report.Claims {
			c := &report.Claims[i]
			if c.Verdict != g.verdict {
				continue
			}
			group.Claims = append(group.Claims, buildClaim(c, g.icon, g.class))
		}
		if len(group.Claims) > 0 {
			v.Groups = append(v.Groups, group)
		}
	}
	return v
}

func buildClaim(c *claims.Claim, icon, class string) claimView {
	cv := claimView{
		ID:        c.ID,
		Text:      c.Text,
		Category:  string(c.Category),
		Rationale: c.Rationale,
		Icon:      icon,
		Class:     class,
	}

	if s := c.Statistical; s != nil {
		cv.HasStat = true
		cv.Value = formatValue(s.Value)
		cv.Unit = s.Unit
		cv.Precision = string(s.Precision)
		if s.AsOfDate != nil && !s.AsOfDate.IsZero() {
			cv.AsOfDate = s.AsOfDate.Format("2006-01-02")
		}
	}

	if val := c.Validation; val != nil {
		if e := val.External; e != nil {
			cv.HasSource = true
			cv.URL = e.URL
			cv.SourceType = string(e.SourceType)
			cv.Reliability = string(e.Reliability)
			cv.RelClass = reliabilityClass(e.Reliability)
			cv.QuotedText = e.QuotedText
			cv.Verified = e.VerifiedMatch
			if e.AccessedAt != nil {
				cv.AccessedAt = e.AccessedAt.Format("2006-01-02")
			}
		}
		if d := val.Derived; d != nil {
			cv.IsDerived = true
			for i, id := range d.SourceClaimIDs {
				if i > 0 {
					cv.DerivedFrom += ", "
				}
				cv.DerivedFrom += id
			}
		}
	}
	return cv
}

func decisionLabel(s claims.ClaimsDecisionStatus) string {
	switch s {
	case claims.ClaimsDecisionPass:
		return "PASS"
	case claims.ClaimsDecisionConditional:
		return "CONDITIONAL"
	case claims.ClaimsDecisionFail:
		return "FAIL"
	default:
		return string(s)
	}
}

func reliabilityClass(t claims.ReliabilityTier) string {
	switch t {
	case claims.ReliabilityAuthoritative, claims.ReliabilityHigh:
		return "rel-high"
	case claims.ReliabilityMedium:
		return "rel-medium"
	default:
		return "rel-low"
	}
}

// formatValue renders a float without trailing zeros (1500, 4.7, 0, 77000).
func formatValue(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}
