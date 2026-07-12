package rubric

// PassCriteria defines the requirements for approval.
// Aligned with LLM-as-Judge best practices.
type PassCriteria struct {
	// MinCategoriesPassing specifies how many categories must pass.
	// Values: "all", "all_required", or a number like "3".
	MinCategoriesPassing string `json:"minCategoriesPassing"`

	// MaxFindings limits findings by severity.
	// Use -1 for unlimited.
	MaxFindings *FindingLimits `json:"maxFindingsSeverity,omitempty"`

	// MinIntScore is the minimum overall IntegerScore (1-5) required to pass.
	// If set, the overall score must be >= this value.
	// Use 0 to disable this check.
	MinIntScore IntegerScore `json:"minIntScore,omitempty"`
}

// DefaultPassCriteria returns standard pass criteria.
// All required categories must pass, 0 critical/high findings allowed.
func DefaultPassCriteria() PassCriteria {
	return PassCriteria{
		MinCategoriesPassing: "all_required",
		MaxFindings: &FindingLimits{
			Critical: 0,
			High:     0,
			Medium:   -1, // Unlimited
			Low:      -1, // Unlimited
		},
	}
}

// StrictPassCriteria returns strict pass criteria.
// All categories must pass, max 3 medium findings.
func StrictPassCriteria() PassCriteria {
	return PassCriteria{
		MinCategoriesPassing: "all",
		MaxFindings: &FindingLimits{
			Critical: 0,
			High:     0,
			Medium:   3,
			Low:      -1,
		},
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
	FindingCounts FindingCounts `json:"findingCounts"`

	// CategoryCounts summarizes category results.
	CategoryCounts CategoryResultCounts `json:"categoryCounts"`
}

// DecisionStatus represents the decision outcome.
type DecisionStatus string

const (
	DecisionPass        DecisionStatus = "pass"         // Meets all criteria
	DecisionConditional DecisionStatus = "conditional"  // Partial scores or non-blocking findings
	DecisionFail        DecisionStatus = "fail"         // Has blocking findings or required categories failed
	DecisionHumanReview DecisionStatus = "human_review" // Requires human judgment
)

// EvaluateResults checks category results and findings against criteria.
func EvaluateResults(results []CategoryResult, findings []Finding, criteria PassCriteria, rubricSet *RubricSet) Decision {
	findingCounts := CountFindings(findings)
	categoryCounts := CountResults(results)

	decision := Decision{
		FindingCounts:  findingCounts,
		CategoryCounts: categoryCounts,
	}

	// Check for blocking findings
	if criteria.MaxFindings != nil {
		if findingCounts.Critical > criteria.MaxFindings.Critical {
			decision.Status = DecisionFail
			decision.Passed = false
			decision.Rationale = formatCriticalRationale(findingCounts.Critical, criteria.MaxFindings.Critical)
			return decision
		}
		if findingCounts.High > criteria.MaxFindings.High {
			decision.Status = DecisionFail
			decision.Passed = false
			decision.Rationale = formatHighRationale(findingCounts.High, criteria.MaxFindings.High)
			return decision
		}
	}

	// Check category pass requirements
	switch criteria.MinCategoriesPassing {
	case "all":
		if categoryCounts.Fail > 0 || categoryCounts.Partial > 0 {
			decision.Status = DecisionFail
			decision.Passed = false
			decision.Rationale = "Not all categories passed: " + itoa(categoryCounts.Fail) + " failed, " + itoa(categoryCounts.Partial) + " partial"
			return decision
		}
	case "all_required":
		if rubricSet != nil && !AllRequiredPassing(results, rubricSet) {
			decision.Status = DecisionFail
			decision.Passed = false
			decision.Rationale = "One or more required categories did not pass"
			return decision
		}
	default:
		// Numeric threshold
		threshold := parseIntOrDefault(criteria.MinCategoriesPassing, 0)
		if categoryCounts.Pass < threshold {
			decision.Status = DecisionFail
			decision.Passed = false
			decision.Rationale = "Only " + itoa(categoryCounts.Pass) + " categories passed (minimum " + itoa(threshold) + " required)"
			return decision
		}
	}

	// Check medium finding limits
	if criteria.MaxFindings != nil && criteria.MaxFindings.Medium >= 0 {
		if findingCounts.Medium > criteria.MaxFindings.Medium {
			decision.Status = DecisionConditional
			decision.Passed = false
			decision.Rationale = formatMediumRationale(findingCounts.Medium, criteria.MaxFindings.Medium)
			return decision
		}
	}

	// Has partial scores but passes
	if categoryCounts.Partial > 0 {
		decision.Status = DecisionConditional
		decision.Passed = true
		decision.Rationale = "Passed with " + itoa(categoryCounts.Partial) + " partial score(s)"
		return decision
	}

	// Has non-blocking findings but passes
	if findingCounts.Medium > 0 || findingCounts.Low > 0 {
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

func formatCriticalRationale(count, max int) string {
	return "Blocked: " + itoa(count) + " critical findings (max " + itoa(max) + " allowed)"
}

func formatHighRationale(count, max int) string {
	return "Blocked: " + itoa(count) + " high severity findings (max " + itoa(max) + " allowed)"
}

func formatMediumRationale(count, max int) string {
	return itoa(count) + " medium findings exceeds limit of " + itoa(max)
}

func parseIntOrDefault(s string, def int) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return def
		}
	}
	if s == "" {
		return def
	}
	return result
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
