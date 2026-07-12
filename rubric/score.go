package rubric

// IntegerScore represents a 1-5 integer evaluation score.
// This scale is preferred for LLM judges as research shows they are
// unreliable at finer granularity than 5 levels.
type IntegerScore int

const (
	// ScoreUnacceptable indicates the spec does not meet requirements.
	ScoreUnacceptable IntegerScore = 1

	// ScoreMajorRevisions indicates significant work is needed.
	ScoreMajorRevisions IntegerScore = 2

	// ScoreAcceptable indicates minimum requirements are met.
	ScoreAcceptable IntegerScore = 3

	// ScoreGood indicates the spec meets expectations well.
	ScoreGood IntegerScore = 4

	// ScoreExcellent indicates the spec exceeds expectations.
	ScoreExcellent IntegerScore = 5
)

// String returns the human-readable label for the score.
func (s IntegerScore) String() string {
	labels := []string{"", "Unacceptable", "Major Revisions", "Acceptable", "Good", "Excellent"}
	if s >= 1 && s <= 5 {
		return labels[s]
	}
	return "Unknown"
}

// ToCategorical converts the integer score to a categorical ScoreValue.
// 1-2 = fail, 3 = partial, 4-5 = pass
func (s IntegerScore) ToCategorical() ScoreValue {
	switch {
	case s <= 2:
		return ScoreFail
	case s == 3:
		return ScorePartial
	default:
		return ScorePass
	}
}

// IsValid returns true if the score is in the valid 1-5 range.
func (s IntegerScore) IsValid() bool {
	return s >= 1 && s <= 5
}

// IsPassing returns true if the score is considered passing (4 or higher).
func (s IntegerScore) IsPassing() bool {
	return s >= ScoreGood
}

// Icon returns the emoji icon for the score.
func (s IntegerScore) Icon() string {
	switch s {
	case ScoreExcellent:
		return "🌟"
	case ScoreGood:
		return "🟢"
	case ScoreAcceptable:
		return "🟡"
	case ScoreMajorRevisions:
		return "🟠"
	case ScoreUnacceptable:
		return "🔴"
	default:
		return "⚪"
	}
}

// ParseIntegerScore converts an integer to IntegerScore, clamping to valid range.
func ParseIntegerScore(score int) IntegerScore {
	if score < 1 {
		return ScoreUnacceptable
	}
	if score > 5 {
		return ScoreExcellent
	}
	return IntegerScore(score)
}

// AllIntegerScores returns all valid integer scores in descending order.
func AllIntegerScores() []IntegerScore {
	return []IntegerScore{
		ScoreExcellent,
		ScoreGood,
		ScoreAcceptable,
		ScoreMajorRevisions,
		ScoreUnacceptable,
	}
}
