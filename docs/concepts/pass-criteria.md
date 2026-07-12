# Pass Criteria

Pass criteria define the thresholds for evaluation decisions. They control how findings and category results translate to pass/fail outcomes.

## PassCriteria Structure

```go
type PassCriteria struct {
    MinCategoriesPassing string        `json:"minCategoriesPassing"` // "all", "all_required", or number
    MaxFindings          *FindingLimits `json:"maxFindingsSeverity"`  // Max findings by severity
    MinIntScore          IntegerScore  `json:"minIntScore"`          // Minimum overall score (1-5)
}

type FindingLimits struct {
    Critical int `json:"critical"` // Max critical findings (-1 = unlimited)
    High     int `json:"high"`     // Max high findings (-1 = unlimited)
    Medium   int `json:"medium"`   // Max medium findings (-1 = unlimited)
    Low      int `json:"low"`      // Max low findings (-1 = unlimited)
}
```

## Default Criteria

```go
func DefaultPassCriteria() PassCriteria {
    return PassCriteria{
        MinCategoriesPassing: "all_required",
        MaxFindings: &FindingLimits{
            Critical: 0,  // No critical findings allowed
            High:     0,  // No high findings allowed
            Medium:   -1, // Unlimited medium findings
            Low:      -1, // Unlimited low findings
        },
    }
}
```

With default criteria:

- ❌ Any critical finding → Fail
- ❌ Any high finding → Fail
- ✅ Medium/low/info findings → Allowed
- ✅ Partial categories → Allowed (if not required)

## Strict Criteria

```go
func StrictPassCriteria() PassCriteria {
    return PassCriteria{
        MaxCritical:    0,
        MaxHigh:        0,
        MaxMedium:      3,    // Max 3 medium findings
        RequireAllPass: true, // All categories must pass
    }
}
```

With strict criteria:

- ❌ Any critical finding → Fail
- ❌ Any high finding → Fail
- ❌ More than 3 medium findings → Fail
- ❌ Any partial or fail category → Fail

## Custom Criteria

```go
criteria := rubric.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        1,  // Allow 1 high finding
    MaxMedium:      5,  // Allow up to 5 medium findings
    RequireAllPass: false,
}

report.SetPassCriteria(criteria)
report.Finalize(nil, "reviewer")
```

## Decision Status

Based on criteria evaluation:

```go
const (
    DecisionPass        DecisionStatus = "pass"         // All criteria met
    DecisionFail        DecisionStatus = "fail"         // Blocking issues found
    DecisionConditional DecisionStatus = "conditional"  // Partial pass, needs attention
    DecisionHumanReview DecisionStatus = "human_review" // Uncertain, needs human review
)
```

## Decision Logic

```go
// Pseudocode for decision computation
func computeDecision(report *Rubric, criteria PassCriteria) Decision {
    counts := report.Decision.FindingCounts
    catCounts := report.Decision.CategoryCounts

    // Check blocking findings
    if counts.Critical > criteria.MaxCritical {
        return Decision{Status: DecisionFail, Rationale: "Too many critical findings"}
    }
    if counts.High > criteria.MaxHigh {
        return Decision{Status: DecisionFail, Rationale: "Too many high findings"}
    }
    if criteria.MaxMedium >= 0 && counts.Medium > criteria.MaxMedium {
        return Decision{Status: DecisionFail, Rationale: "Too many medium findings"}
    }

    // Check category requirements
    if criteria.RequireAllPass && catCounts.Fail > 0 {
        return Decision{Status: DecisionFail, Rationale: "Some categories failed"}
    }
    if criteria.RequireAllPass && catCounts.Partial > 0 {
        return Decision{Status: DecisionConditional, Rationale: "Some categories partial"}
    }

    // All checks passed
    if catCounts.Fail == 0 && catCounts.Partial == 0 {
        return Decision{Status: DecisionPass}
    }

    return Decision{Status: DecisionConditional}
}
```

## Example Configurations

### Security Review

```go
// Zero tolerance for security issues
criteria := rubric.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        0,
    MaxMedium:      0,  // Even medium security findings block
    RequireAllPass: true,
}
```

### Documentation Review

```go
// More lenient for docs
criteria := rubric.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        2,   // Allow some high-priority gaps
    MaxMedium:      -1,  // Unlimited medium
    RequireAllPass: false,
}
```

### Release Gate

```go
// Strict for releases
criteria := rubric.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        0,
    MaxMedium:      5,
    RequireAllPass: true,
}
```

## CLI Integration

```bash
# Check with default criteria
sevaluation check report.json

# Exit codes:
# 0 = pass
# 1 = fail or conditional
```

## MinIntScore (v0.9.0)

Require a minimum overall integer score (1-5) for approval:

```go
criteria := rubric.PassCriteria{
    MinCategoriesPassing: "all_required",
    MaxFindings: &rubric.FindingLimits{
        Critical: 0,
        High:     0,
        Medium:   -1,
        Low:      -1,
    },
    MinIntScore: rubric.ScoreGood,  // Require 4+ overall
}

report.SetPassCriteria(criteria)
report.Evaluate(rubricSet)

// If computed IntScore < 4, report.Pass = false
```

### Use Cases

```go
// Quality gate - require "Good" (4) or better
criteria.MinIntScore = rubric.ScoreGood

// Strict gate - require "Excellent" (5)
criteria.MinIntScore = rubric.ScoreExcellent

// Minimum viable - require "Acceptable" (3) or better
criteria.MinIntScore = rubric.ScoreAcceptable
```

### Combined with Finding Limits

MinIntScore is checked after IntScore is computed:

```go
// Both must pass:
// 1. No critical/high findings
// 2. Overall score >= 4
criteria := rubric.PassCriteria{
    MaxFindings: &rubric.FindingLimits{Critical: 0, High: 0, Medium: -1, Low: -1},
    MinIntScore: rubric.ScoreGood,
}
```

## Next Steps

- [Findings & Severity](findings.md) - Understanding severity levels
- [Rubrics](../features/rubrics.md) - Define evaluation criteria
