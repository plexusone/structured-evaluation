package claims

import "strings"

// Stance says how a quoted passage bears on a claim.
type Stance string

const (
	// StanceSupports means the passage says what the claim says.
	StanceSupports Stance = "supports"

	// StanceDisputes means the passage says otherwise.
	StanceDisputes Stance = "disputes"

	// StanceContext means the passage is relevant but neither supports nor
	// disputes the claim.
	StanceContext Stance = "context"
)

// Translation renders a quote in another language, usually English, so a
// reviewer can read evidence in a language they do not know.
type Translation struct {
	// Text is the translated passage.
	Text string `json:"text"`

	// Language is the BCP 47 tag of Text, such as "en".
	Language string `json:"language,omitempty"`
}

// Evidence is one passage from one source that bears on a claim. A claim
// can carry several, from several sources, each with its own stance, so a
// report can record that sources disagree rather than only a single
// validating source.
type Evidence struct {
	// SourceRef is the ID of the source in the caller's source registry,
	// such as the source list a claim set was exported with. Set it, URL,
	// or both.
	SourceRef string `json:"sourceRef,omitempty"`

	// URL locates the source when it is not in a registry, or to make the
	// report readable on its own.
	URL string `json:"url,omitempty"`

	// Title names the source.
	Title string `json:"title,omitempty"`

	// Stance is how the passage bears on the claim.
	Stance Stance `json:"stance"`

	// Quote is the passage, verbatim, in its original language.
	Quote string `json:"quote,omitempty"`

	// Locator points within the source: a page, section or timestamp.
	Locator string `json:"locator,omitempty"`

	// Language is the BCP 47 tag of Quote when it is not English.
	Language string `json:"language,omitempty"`

	// Translation renders a quote that is not in English.
	Translation *Translation `json:"translation,omitempty"`

	// Role is how directly the source speaks for the claim.
	Role SourceRole `json:"sourceRole,omitempty"`

	// HowChecked records how the passage was obtained, such as "opened the
	// page" or "read the scanned book".
	HowChecked string `json:"howChecked,omitempty"`
}

// HasQuote reports whether the evidence quotes its source.
func (e Evidence) HasQuote() bool {
	return strings.TrimSpace(e.Quote) != ""
}

// HasSource reports whether the evidence names its source.
func (e Evidence) HasSource() bool {
	return strings.TrimSpace(e.SourceRef) != "" || strings.TrimSpace(e.URL) != ""
}

// SupportingQuotes returns the claim's evidence that supports it with a
// quote from a named source.
func (c *Claim) SupportingQuotes() []Evidence {
	var out []Evidence
	for _, e := range c.Evidence {
		if e.Stance == StanceSupports && e.HasQuote() && e.HasSource() {
			out = append(out, e)
		}
	}
	return out
}

// Confidence is a judge's confidence in a verdict.
type Confidence string

const (
	// ConfidenceHigh means the sources settle the claim.
	ConfidenceHigh Confidence = "high"

	// ConfidenceMedium means the sources mostly settle it, with a gap such
	// as a detail no source states.
	ConfidenceMedium Confidence = "medium"

	// ConfidenceLow means the verdict rests on weak or partial evidence.
	ConfidenceLow Confidence = "low"
)

// Review is a second judge's check of a claim's verdict. Producers set the
// claim's Verdict to the final verdict: when the second judge disagrees and
// gives a verdict of its own, that verdict is the claim's Verdict and the
// first judge's is recorded in FirstVerdict.
type Review struct {
	// Judge identifies the second judge.
	Judge string `json:"judge"`

	// Agrees reports whether the second judge agreed with the first.
	Agrees bool `json:"agrees"`

	// FirstVerdict is the first judge's verdict when the second judge
	// disagreed and its verdict replaced it.
	FirstVerdict Verdict `json:"firstVerdict,omitempty"`

	// Rationale explains the second judge's view.
	Rationale string `json:"rationale,omitempty"`
}

// InaccessibleSource is a source a judge could not open, so a reader knows
// which evidence was never checked.
type InaccessibleSource struct {
	// SourceRef is the ID of the source in the caller's registry.
	SourceRef string `json:"sourceRef,omitempty"`

	// URL locates the source.
	URL string `json:"url,omitempty"`

	// Reason says why it could not be opened.
	Reason string `json:"reason"`
}

// ReviewAgreement summarizes how often second judges agreed with first
// judges, a measure of how consistent the judging is.
type ReviewAgreement struct {
	// Reviewed is the number of claims a second judge checked.
	Reviewed int `json:"reviewed"`

	// Agreed is the number of those where the second judge agreed.
	Agreed int `json:"agreed"`

	// Rate is Agreed divided by Reviewed.
	Rate float64 `json:"rate"`
}

// AgreementOf returns the second-judge agreement of claims, or nil when no
// claim has a second check.
func AgreementOf(claims []Claim) *ReviewAgreement {
	var a ReviewAgreement
	for _, c := range claims {
		if c.Review == nil {
			continue
		}
		a.Reviewed++
		if c.Review.Agrees {
			a.Agreed++
		}
	}
	if a.Reviewed == 0 {
		return nil
	}
	a.Rate = float64(a.Agreed) / float64(a.Reviewed)
	return &a
}
