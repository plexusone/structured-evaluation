package rubric

import "testing"

func TestIntegerScore_String(t *testing.T) {
	tests := []struct {
		score IntegerScore
		want  string
	}{
		{ScoreUnacceptable, "Unacceptable"},
		{ScoreMajorRevisions, "Major Revisions"},
		{ScoreAcceptable, "Acceptable"},
		{ScoreGood, "Good"},
		{ScoreExcellent, "Excellent"},
		{IntegerScore(0), "Unknown"},
		{IntegerScore(6), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.score.String(); got != tt.want {
				t.Errorf("IntegerScore(%d).String() = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

func TestIntegerScore_ToCategorical(t *testing.T) {
	tests := []struct {
		score IntegerScore
		want  ScoreValue
	}{
		{ScoreUnacceptable, ScoreFail},
		{ScoreMajorRevisions, ScoreFail},
		{ScoreAcceptable, ScorePartial},
		{ScoreGood, ScorePass},
		{ScoreExcellent, ScorePass},
	}

	for _, tt := range tests {
		t.Run(tt.score.String(), func(t *testing.T) {
			if got := tt.score.ToCategorical(); got != tt.want {
				t.Errorf("IntegerScore(%d).ToCategorical() = %s, want %s", tt.score, got, tt.want)
			}
		})
	}
}

func TestIntegerScore_IsValid(t *testing.T) {
	tests := []struct {
		score IntegerScore
		want  bool
	}{
		{IntegerScore(0), false},
		{ScoreUnacceptable, true},
		{ScoreMajorRevisions, true},
		{ScoreAcceptable, true},
		{ScoreGood, true},
		{ScoreExcellent, true},
		{IntegerScore(6), false},
	}

	for _, tt := range tests {
		if got := tt.score.IsValid(); got != tt.want {
			t.Errorf("IntegerScore(%d).IsValid() = %v, want %v", tt.score, got, tt.want)
		}
	}
}

func TestIntegerScore_IsPassing(t *testing.T) {
	tests := []struct {
		score IntegerScore
		want  bool
	}{
		{ScoreUnacceptable, false},
		{ScoreMajorRevisions, false},
		{ScoreAcceptable, false},
		{ScoreGood, true},
		{ScoreExcellent, true},
	}

	for _, tt := range tests {
		t.Run(tt.score.String(), func(t *testing.T) {
			if got := tt.score.IsPassing(); got != tt.want {
				t.Errorf("IntegerScore(%d).IsPassing() = %v, want %v", tt.score, got, tt.want)
			}
		})
	}
}

func TestParseIntegerScore(t *testing.T) {
	tests := []struct {
		input int
		want  IntegerScore
	}{
		{-1, ScoreUnacceptable},
		{0, ScoreUnacceptable},
		{1, ScoreUnacceptable},
		{2, ScoreMajorRevisions},
		{3, ScoreAcceptable},
		{4, ScoreGood},
		{5, ScoreExcellent},
		{6, ScoreExcellent},
		{100, ScoreExcellent},
	}

	for _, tt := range tests {
		if got := ParseIntegerScore(tt.input); got != tt.want {
			t.Errorf("ParseIntegerScore(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestIntegerScore_Icon(t *testing.T) {
	// Just verify icons are non-empty
	for _, score := range AllIntegerScores() {
		if icon := score.Icon(); icon == "" {
			t.Errorf("IntegerScore(%d).Icon() should not be empty", score)
		}
	}
}
