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

## Creating Findings

```go
// Critical finding - blocks approval
report.AddFinding(evaluation.Finding{
    Severity:       evaluation.SeverityCritical,
    Category:       "security",
    Title:          "SQL injection vulnerability",
    Description:    "User input is concatenated directly into SQL query",
    Location:       "src/db/queries.go:142",
    Recommendation: "Use parameterized queries",
    References:     []string{"OWASP-A03:2021"},
})

// Medium finding - tracked but doesn't block
report.AddFinding(evaluation.Finding{
    Severity:       evaluation.SeverityMedium,
    Category:       "documentation",
    Title:          "Missing API documentation",
    Description:    "Public endpoints lack OpenAPI annotations",
    Recommendation: "Add swagger comments to all public handlers",
})

// Info finding - observation
report.AddFinding(evaluation.Finding{
    Severity:    evaluation.SeverityInfo,
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
severity := evaluation.SeverityHigh

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

## Next Steps

- [Pass Criteria](pass-criteria.md) - Configure blocking thresholds
- [Report Rendering](../features/rendering.md) - Display findings
