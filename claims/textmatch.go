package claims

import (
	"regexp"
	"strings"
	"unicode"
)

// NormalizeText prepares text for matching: lower case, with every run of
// characters other than letters and digits turned into one space, and a
// space between a digit and a letter. Quote styles, dashes, punctuation,
// thousands separators, line breaks and a page's missing space ("20febbr.")
// then cannot stop a match, so "1,142,462" and "1.142.462" both become
// "1 142 462". Letters keep their accents.
func NormalizeText(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	space := true
	var prev rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !space && unicode.IsDigit(r) != unicode.IsDigit(prev) {
				b.WriteByte(' ')
			}
			b.WriteRune(unicode.ToLower(r))
			space, prev = false, r
			continue
		}
		if !space {
			b.WriteByte(' ')
			space = true
		}
	}
	return strings.TrimSpace(b.String())
}

// ellipsis matches the marks that show words were left out of a quotation.
var ellipsis = regexp.MustCompile(`\[\s*(?:\.\.\.|…)\s*\]|\.\.\.|…`)

// QuoteFragments splits a quotation at its ellipses into the parts that
// must each appear in the source, in order, normalized with NormalizeText.
func QuoteFragments(quote string) []string {
	var out []string
	for _, part := range ellipsis.Split(quote, -1) {
		if n := NormalizeText(part); n != "" {
			out = append(out, n)
		}
	}
	return out
}

// QuoteMatch reports whether a quotation was found in a text.
type QuoteMatch struct {
	// Found reports whether every fragment appears in the text.
	Found bool `json:"found"`
	// InOrder reports whether the fragments appear in the order quoted.
	// A quotation that joins passages out of order is still found.
	InOrder bool `json:"inOrder"`
	// Fragments is the number of parts between ellipses.
	Fragments int `json:"fragments"`
	// Missing lists the fragments not found, normalized.
	Missing []string `json:"missing,omitempty"`
}

// MatchQuote reports whether quote appears in text: each part of the quote
// between ellipses must appear as whole words, after normalizing both with
// NormalizeText. InOrder says whether the parts appear in the order quoted.
// A quote with no letters or digits is not found.
func MatchQuote(quote, text string) QuoteMatch {
	frags := QuoteFragments(quote)
	m := QuoteMatch{Fragments: len(frags)}
	if len(frags) == 0 {
		return m
	}
	hay := " " + NormalizeText(text) + " "
	m.InOrder = true
	last := -1
	for _, f := range frags {
		i := strings.Index(hay, " "+f+" ")
		if i < 0 {
			m.Missing = append(m.Missing, f)
			continue
		}
		if next := strings.Index(hay[last+1:], " "+f+" "); next < 0 {
			m.InOrder = false
		} else {
			last += 1 + next + len(f)
		}
	}
	m.Found = len(m.Missing) == 0
	m.InOrder = m.InOrder && m.Found
	return m
}

// figure matches a number: thousands groups separated by a comma, point,
// apostrophe or narrow space ("1,142,462", "1.142.462"), or digits with an
// optional decimal part ("1947", "12.5"). A separator followed by a space
// ends the number, so "1947, 45" is two numbers.
var figure = regexp.MustCompile(`\d{1,3}(?:[,.'\x{00a0}\x{202f}]\d{3})+|\d+(?:[.,]\d+)?`)

// Figures returns the numbers in s, such as years, counts and amounts, with
// their separators removed, in order of appearance and without repeats.
func Figures(s string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range figure.FindAllString(s, -1) {
		d := digitsOnly(m)
		if d != "" && !seen[d] {
			seen[d] = true
			out = append(out, d)
		}
	}
	return out
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// MissingFigures returns the figures in statement that appear in none of
// texts. Separators are ignored on both sides, so "1,142,462" matches
// "1.142.462".
func MissingFigures(statement string, texts ...string) []string {
	have := map[string]bool{}
	for _, t := range texts {
		for _, f := range Figures(t) {
			have[f] = true
		}
	}
	var out []string
	for _, f := range Figures(statement) {
		if !have[f] {
			out = append(out, f)
		}
	}
	return out
}
