package evaluation

// PassCriteria defines the requirements for approval.
type PassCriteria struct {
	// MaxCritical is the maximum allowed critical findings (default 0).
	MaxCritical int `json:"max_critical"`

	// MaxHigh is the maximum allowed high severity findings (default 0).
	MaxHigh int `json:"max_high"`

	// MaxMedium is the maximum allowed medium findings (-1 = unlimited).
	MaxMedium int `json:"max_medium,omitempty"`

	// MinScore is the minimum weighted score required.
	MinScore float64 `json:"min_score"`
}

// DefaultPassCriteria returns standard pass criteria.
// Zero Critical/High, minimum score 7.0.
func DefaultPassCriteria() PassCriteria {
	return PassCriteria{
		MaxCritical: 0,
		MaxHigh:     0,
		MaxMedium:   -1, // Unlimited
		MinScore:    7.0,
	}
}

// StrictPassCriteria returns strict pass criteria.
// Zero Critical/High, max 3 Medium, minimum score 8.0.
func StrictPassCriteria() PassCriteria {
	return PassCriteria{
		MaxCritical: 0,
		MaxHigh:     0,
		MaxMedium:   3,
		MinScore:    8.0,
	}
}

// Decision represents the evaluation decision.
type Decision struct {
	// Status is the decision outcome.
	Status DecisionStatus `json:"status"`

	// Passed indicates if the evaluation passed.
	Passed bool `json:"passed"`

	// Rationale explains the decision.
	Rationale string `json:"rationale"`

	// FindingCounts summarizes findings by severity.
	FindingCounts FindingCounts `json:"finding_counts"`

	// WeightedScore is the final weighted score.
	WeightedScore float64 `json:"weighted_score"`
}

// DecisionStatus represents the decision outcome.
type DecisionStatus string

const (
	DecisionPass        DecisionStatus = "pass"         // Meets all criteria
	DecisionConditional DecisionStatus = "conditional"  // Meets score but has findings
	DecisionFail        DecisionStatus = "fail"         // Has blocking findings
	DecisionHumanReview DecisionStatus = "human_review" // Requires human judgment
)

// Evaluate checks findings and score against criteria.
func Evaluate(findings []Finding, weightedScore float64, criteria PassCriteria) Decision {
	counts := CountFindings(findings)

	decision := Decision{
		FindingCounts: counts,
		WeightedScore: weightedScore,
	}

	// Check for blocking findings
	criticalExceeded := counts.Critical > criteria.MaxCritical
	highExceeded := counts.High > criteria.MaxHigh
	mediumExceeded := criteria.MaxMedium >= 0 && counts.Medium > criteria.MaxMedium
	scoreBelowMin := weightedScore < criteria.MinScore

	if criticalExceeded || highExceeded {
		decision.Status = DecisionFail
		decision.Passed = false
		decision.Rationale = formatFailRationale(counts, criteria)
		return decision
	}

	if scoreBelowMin {
		decision.Status = DecisionHumanReview
		decision.Passed = false
		decision.Rationale = formatScoreRationale(weightedScore, criteria.MinScore)
		return decision
	}

	if mediumExceeded {
		decision.Status = DecisionConditional
		decision.Passed = false
		decision.Rationale = formatMediumRationale(counts.Medium, criteria.MaxMedium)
		return decision
	}

	if counts.Medium > 0 || counts.Low > 0 {
		decision.Status = DecisionConditional
		decision.Passed = true
		decision.Rationale = "Passed with non-blocking findings"
		return decision
	}

	decision.Status = DecisionPass
	decision.Passed = true
	decision.Rationale = "Meets all criteria"
	return decision
}

func formatFailRationale(counts FindingCounts, criteria PassCriteria) string {
	if counts.Critical > criteria.MaxCritical {
		return "Blocked: " + itoa(counts.Critical) + " critical findings (max " + itoa(criteria.MaxCritical) + ")"
	}
	return "Blocked: " + itoa(counts.High) + " high severity findings (max " + itoa(criteria.MaxHigh) + ")"
}

func formatScoreRationale(score, min float64) string {
	return "Score " + ftoa(score) + " below minimum " + ftoa(min)
}

func formatMediumRationale(count, max int) string {
	return itoa(count) + " medium findings exceeds limit of " + itoa(max)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}

func ftoa(f float64) string {
	// Simple float to string for scores
	whole := int(f)
	frac := int((f - float64(whole)) * 10)
	return itoa(whole) + "." + itoa(frac)
}
