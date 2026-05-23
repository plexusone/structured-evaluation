# Pass Criteria

Pass criteria define the thresholds for evaluation decisions. They control how findings and category results translate to pass/fail outcomes.

## PassCriteria Structure

```go
type PassCriteria struct {
    MaxCritical    int  `json:"max_critical"`     // Max critical findings allowed (-1 = unlimited)
    MaxHigh        int  `json:"max_high"`         // Max high findings allowed (-1 = unlimited)
    MaxMedium      int  `json:"max_medium"`       // Max medium findings allowed (-1 = unlimited)
    RequireAllPass bool `json:"require_all_pass"` // Require all categories to pass
}
```

## Default Criteria

```go
func DefaultPassCriteria() PassCriteria {
    return PassCriteria{
        MaxCritical:    0,  // No critical findings allowed
        MaxHigh:        0,  // No high findings allowed
        MaxMedium:      -1, // Unlimited medium findings
        RequireAllPass: false,
    }
}
```

With default criteria:

- ❌ Any critical finding → Fail
- ❌ Any high finding → Fail
- ✅ Medium/low/info findings → Allowed
- ✅ Partial categories → Allowed

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
criteria := evaluation.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        1,  // Allow 1 high finding
    MaxMedium:      5,  // Allow up to 5 medium findings
    RequireAllPass: false,
}

report.SetPassCriteria(criteria)
report.Finalize("reviewer")
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
func computeDecision(report *EvaluationReport, criteria PassCriteria) Decision {
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
criteria := evaluation.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        0,
    MaxMedium:      0,  // Even medium security findings block
    RequireAllPass: true,
}
```

### Documentation Review

```go
// More lenient for docs
criteria := evaluation.PassCriteria{
    MaxCritical:    0,
    MaxHigh:        2,   // Allow some high-priority gaps
    MaxMedium:      -1,  // Unlimited medium
    RequireAllPass: false,
}
```

### Release Gate

```go
// Strict for releases
criteria := evaluation.PassCriteria{
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

## Next Steps

- [Findings & Severity](findings.md) - Understanding severity levels
- [Rubrics](../features/rubrics.md) - Define evaluation criteria
