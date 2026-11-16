package claims

import (
	"reflect"
	"testing"
)

func TestMatchQuote(t *testing.T) {
	text := `Nel 1947 le Ferrari vendute furono 7, nel 1948 45 e circa 80 l'anno.
	During a Board of Directors’ meeting of 30 December 1937, Ugo Gobbato announced the purchase of an 80% share.`
	tests := []struct {
		name  string
		quote string
		found bool
	}{
		{"exact", "Nel 1947 le Ferrari vendute furono 7", true},
		{"case, punctuation and line breaks", "nel 1948 45 e circa 80 l’anno", true},
		{"curly apostrophe in the source", "Board of Directors' meeting", true},
		{"ellipsis between parts in order", "During a Board ... Ugo Gobbato announced", true},
		{"bracketed ellipsis", "During a Board […] the purchase of an 80% share", true},
		{"parts out of order are found", "Ugo Gobbato … During a Board", true},
		{"missing space between digits and letters", "Nel 1947le Ferrari", true},
		{"changed figure", "Nel 1947 le Ferrari vendute furono 8", false},
		{"part of a word is not a match", "Ferrar", false},
		{"no words", "…", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if m := MatchQuote(tt.quote, text); m.Found != tt.found {
				t.Errorf("Found = %v, want %v (missing %v)", m.Found, tt.found, m.Missing)
			}
		})
	}
}

func TestMatchQuoteOrder(t *testing.T) {
	text := "first part here. second part there."
	for quote, inOrder := range map[string]bool{"first part … second part": true, "second part … first part": false} {
		t.Run(quote, func(t *testing.T) {
			if m := MatchQuote(quote, text); !m.Found || m.InOrder != inOrder {
				t.Errorf("Found %v InOrder %v, want found and InOrder %v", m.Found, m.InOrder, inOrder)
			}
		})
	}
}

func TestFigures(t *testing.T) {
	got := Figures("Ferrari sold 7 cars in 1947, 45 in 1948, and 1,142,462 lire (12.5%) in 1947.")
	want := []string{"7", "1947", "45", "1948", "1142462", "125"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Figures = %v, want %v", got, want)
	}
}

func TestMissingFigures(t *testing.T) {
	statement := "The liquidation paid 1,142,462 lire in 1937 for 8 cars."
	got := MissingFigures(statement, "liquidazione di 1.142.462 lire", "nel 1937")
	if !reflect.DeepEqual(got, []string{"8"}) {
		t.Errorf("MissingFigures = %v, want [8]", got)
	}
}

func TestNormalizeText(t *testing.T) {
	if got := NormalizeText("  L’Alfa—Romeo, “Portello”\n1.142.462 20febbr."); got != "l alfa romeo portello 1 142 462 20 febbr" {
		t.Errorf("NormalizeText = %q", got)
	}
}
