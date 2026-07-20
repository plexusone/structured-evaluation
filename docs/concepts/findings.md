# Findings & Severity

Findings capture specific issues discovered during evaluation, with severity levels following InfoSec conventions.

## Finding Structure

```go
type Finding struct {
    ID             string   `json:"id,omitempty"`
    Severity       Severity `json:"severity"`
    Category       string   `json:"category"`
    Title          string   `json:"title"`
    Description    string   `json:"description,omitempty"`
    Location       string   `json:"location,omitempty"`
    Recommendation string   `json:"recommendation,omitempty"`
    References     []string `json:"references,omitempty"`
}
```

## Severity Levels

| Severity | Icon | Blocking | Description |
|----------|------|----------|-------------|
| `Critical` | 🔴 | Yes | Must fix before approval; security vulnerabilities, data loss risks |
| `High` | 🔴 | Yes | Must fix before approval; major functionality gaps |
| `Medium` | 🟡 | No | Should fix; tracked issues that don't block |
| `Low` | 🟢 | No | Nice to fix; minor improvements |
| `Info` | ⚪ | No | Informational; observations, suggestions |

## Category Severity Rollup (v0.11.0)

`CategoryResult` carries a `Severity` field — the highest-severity finding in
that category, or empty if it has none. It's computed automatically from the
category's `Findings` via `WorstSeverity()`, never set independently by the
judge, so it can't drift out of sync with what was actually found:

```go
type CategoryResult struct {
    Category string     `json:"category"`
    Score    ScoreValue `json:"score"`
    Severity Severity   `json:"severity,omitempty"`
    Findings []Finding  `json:"findings,omitempty"`
    // ...
}
```

It's computed in two places, so every path is covered:

```go
// AddFinding recomputes Severity on every call.
cat := rubric.CategoryResult{Category: "security"}
cat.AddFinding(rubric.Finding{Severity: rubric.SeverityHigh, Title: "..."})
cat.Severity // SeverityHigh

// AddCategoryResult and Evaluate compute it as a safety net for categories
// built via struct literal or appended directly to Rubric.Categories.
report.AddCategoryResult(rubric.CategoryResult{
    Category: "security",
    Findings: []rubric.Finding{{Severity: rubric.SeverityCritical}},
})
// report.Categories[len(report.Categories)-1].Severity == SeverityCritical
```

An explicitly set `Severity` is never overwritten — the compute-if-unset
behavior only fills it in when it's still the zero value, matching the same
convention already used for `Rubric.IntScore` and `Rubric.Confidence`.

This is distinct from `Score`/`IntScore`, which measure quality. `Severity`
exists specifically to answer "which categories should I fix first?" —
useful for sorting or highlighting categories in a UI independent of their
pass/partial/fail status.

```go
// Sort categories by severity for a "fix these first" view
sort.Slice(report.Categories, func(i, j int) bool {
    return report.Categories[i].Severity.Weight() > report.Categories[j].Severity.Weight()
})
```

## Creating Findings

```go
// Critical finding - blocks approval
report.AddFinding(rubric.Finding{
    Severity:       rubric.SeverityCritical,
    Category:       "security",
    Title:          "SQL injection vulnerability",
    Description:    "User input is concatenated directly into SQL query",
    Location:       "src/db/queries.go:142",
    Recommendation: "Use parameterized queries",
    References:     []string{"OWASP-A03:2021"},
})

// Medium finding - tracked but doesn't block
report.AddFinding(rubric.Finding{
    Severity:       rubric.SeverityMedium,
    Category:       "documentation",
    Title:          "Missing API documentation",
    Description:    "Public endpoints lack OpenAPI annotations",
    Recommendation: "Add swagger comments to all public handlers",
})

// Info finding - observation
report.AddFinding(rubric.Finding{
    Severity:    rubric.SeverityInfo,
    Category:    "style",
    Title:       "Consider using table-driven tests",
    Description: "Test file has repetitive test cases",
})
```

## Finding Counts

The decision tracks finding counts by severity:

```go
type FindingCounts struct {
    Critical int `json:"critical"`
    High     int `json:"high"`
    Medium   int `json:"medium"`
    Low      int `json:"low"`
    Info     int `json:"info"`
    Total    int `json:"total"`
}

// BlockingCount returns critical + high
func (fc FindingCounts) BlockingCount() int {
    return fc.Critical + fc.High
}
```

### Usage

```go
counts := report.Decision.FindingCounts

if counts.BlockingCount() > 0 {
    fmt.Printf("❌ %d blocking issues found\n", counts.BlockingCount())
}

fmt.Printf("Findings: %d critical, %d high, %d medium, %d low\n",
    counts.Critical, counts.High, counts.Medium, counts.Low)
```

## Severity Methods

```go
severity := rubric.SeverityHigh

severity.Icon()        // "🔴"
severity.IsBlocking()  // true
severity.Priority()    // 2 (lower = more severe)
```

| Severity | Icon | IsBlocking | Priority |
|----------|------|------------|----------|
| Critical | 🔴 | true | 1 |
| High | 🔴 | true | 2 |
| Medium | 🟡 | false | 3 |
| Low | 🟢 | false | 4 |
| Info | ⚪ | false | 5 |

## Grouping Findings

Findings can be grouped by category or severity for display:

```go
// Group by category
byCategory := make(map[string][]Finding)
for _, f := range report.Findings {
    byCategory[f.Category] = append(byCategory[f.Category], f)
}

// Group by severity
bySeverity := make(map[Severity][]Finding)
for _, f := range report.Findings {
    bySeverity[f.Severity] = append(bySeverity[f.Severity], f)
}
```

## Best Practices

### Writing Good Finding Titles

- ✅ "Missing input validation on user ID parameter"
- ❌ "Bad code"
- ✅ "API rate limiting not implemented"
- ❌ "Needs work"

### Writing Recommendations

- Be specific and actionable
- Reference standards when applicable
- Provide code examples if helpful

```go
Finding{
    Title:          "Passwords stored in plaintext",
    Recommendation: "Use bcrypt with cost factor ≥12. See OWASP Password Storage Cheat Sheet.",
}
```

### Using Location

Include file paths and line numbers when available:

```go
Finding{
    Location: "src/auth/handler.go:89",  // File:line
    Location: "Section 3.2",              // Document section
    Location: "API: POST /users",         // Endpoint
}
```

## Reason Codes (v0.9.0)

Reason codes provide standardized finding identifiers for automated repair workflows.

### ReasonCode Format

Codes use `{CATEGORY}-{ISSUE}` format:

```go
// Category prefixes
CategoryREQ    = "REQ"    // Requirements
CategorySEC    = "SEC"    // Security
CategoryARCH   = "ARCH"   // Architecture
CategoryDOC    = "DOC"    // Documentation
CategoryMETRIC = "METRIC" // Metrics
CategoryUSER   = "USER"   // User personas
CategorySCALE  = "SCALE"  // Scalability
CategoryINFRA  = "INFRA"  // Infrastructure
CategorySCOPE  = "SCOPE"  // Scope
CategoryUX     = "UX"     // UX/Accessibility
```

### Pre-defined Codes

```go
// Requirements
CodeREQAmbiguous     // "REQ-AMBIGUOUS"
CodeREQIncomplete    // "REQ-INCOMPLETE"
CodeREQConflict      // "REQ-CONFLICT"

// Security
CodeSECMissingAuth   // "SEC-MISSING_AUTH"
CodeSECNoValidation  // "SEC-NO_VALIDATION"

// Metrics
CodeMETRICNoBaseline // "METRIC-NO_BASELINE"
CodeMETRICNoTarget   // "METRIC-NO_TARGET"
```

### Creating Findings with Codes

```go
// Using NewFindingWithCode (inherits severity from registry)
finding := rubric.NewFindingWithCode(
    "f1",
    "requirements",
    rubric.CodeREQAmbiguous,
    "Ambiguous requirement",
    "Requirement REQ-12 can be interpreted multiple ways",
)

// Or set code on existing finding
finding := rubric.NewFinding("f2", "security", rubric.SeverityHigh, "Missing auth", "...")
finding.SetCode(rubric.CodeSECMissingAuth)
finding.SetLocation("Section 3.2")
```

### Automated Repair

Each code has a repair prompt in the registry:

```go
// Get repair prompt for automated fixes
prompt := finding.GetRepairPrompt()
// Returns: "Rewrite this requirement to be specific and testable..."

// Or access code info directly
info := rubric.GetReasonCodeInfo(rubric.CodeREQAmbiguous)
info.Category        // "REQ"
info.DefaultSeverity // SeverityMedium
info.RepairPrompt    // "Rewrite this requirement..."
```

### Aggregating by Code

```go
// Count findings by reason code
counts := rubric.CountFindingsByCode(report.Findings)
// map[ReasonCode]int{"REQ-AMBIGUOUS": 3, "SEC-MISSING_AUTH": 1}

// Get blocking codes
blocking := rubric.GetBlockingCodes(report.Findings)
// []ReasonCode{"SEC-MISSING_AUTH"}
```

### Blocking Codes on Rubric

The rubric tracks which codes caused failure:

```go
report.Evaluate(rubricSet)

if !report.Pass {
    for _, code := range report.Blocking {
        fmt.Printf("Blocked by: %s\n", code)
    }
}
```

## Next Steps

- [Pass Criteria](pass-criteria.md) - Configure blocking thresholds
- [Report Rendering](../features/rendering.md) - Display findings
