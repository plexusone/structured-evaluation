# Report Types

Structured Evaluation provides three report types for different use cases.

## Rubric

For **subjective quality assessments** using LLM-as-Judge or human reviewers.

### Structure

```go
type Rubric struct {
    ReviewType  string           // e.g., "prd-review", "code-quality"
    Metadata    ReportMetadata   // Document info, timestamps
    Categories  []CategoryResult // Pass/partial/fail per category
    Findings    []Finding        // Issues discovered
    Decision    Decision         // Overall pass/fail decision
    NextSteps   []string         // Recommended actions
    Judge       *JudgeMetadata   // LLM judge info (optional)
    RubricID    string           // Rubric used (optional)
}
```

### Use Cases

- PRD/MRD quality review
- Code review assessments
- Content quality evaluation
- Design document review
- Security assessment reports

### Example

```go
report := rubric.NewRubric("prd-review", "requirements.md")
report.AddCategoryResult(rubric.CategoryResult{
    Category:  "clarity",
    Score:     rubric.ScorePass,
    Reasoning: "Requirements are clearly written",
})
report.AddFinding(rubric.Finding{
    Severity: rubric.SeverityMedium,
    Title:    "Missing edge case",
})
report.Finalize(nil, "reviewer")
```

## SummaryReport

For **deterministic GO/NO-GO checks** in CI/CD pipelines and release validation.

### Structure

```go
type SummaryReport struct {
    Project     string        // Project name
    Version     string        // Version being validated
    Title       string        // Report title
    Teams       []TeamSection // Organized by team/domain
    OverallStatus Status      // Computed overall status
    GeneratedAt time.Time     // Timestamp
}

type TeamSection struct {
    ID    string       // Team identifier
    Name  string       // Display name
    Tasks []TaskResult // Individual check results
}

type TaskResult struct {
    ID     string // Task identifier
    Status Status // go/warn/nogo
    Detail string // Additional context
}
```

### Use Cases

- Release readiness validation
- CI/CD pipeline gates
- Deployment checklists
- Compliance verification
- Test result aggregation

### Example

```go
report := summary.NewSummaryReport("my-service", "v2.0.0", "Release Validation")
report.AddTeam(summary.TeamSection{
    ID:   "testing",
    Name: "Testing",
    Tasks: []summary.TaskResult{
        {ID: "unit-tests", Status: summary.StatusGo, Detail: "100% pass"},
        {ID: "coverage", Status: summary.StatusGo, Detail: "92% coverage"},
    },
})
```

## ClaimsReport

For **factual claim validation** with source tracking and verification.

### Structure

```go
type ClaimsReport struct {
    Metadata   ClaimsMetadata  // Document info, timestamps
    Claims     []Claim         // Extracted claims with validation
    Summary    ClaimsSummary   // Counts by verdict, category
    Criteria   ClaimsCriteria  // Pass requirements
    Decision   ClaimsDecision  // Overall pass/fail decision
}

type Claim struct {
    ID          string       // Unique identifier
    Text        string       // Claim text
    Location    Location     // Where in document
    Category    ClaimCategory // metadata, technical-finding, etc.
    Validation  *Validation  // How validated
    Verdict     Verdict      // verified, unverified, needs-review, rejected
    Rationale   string       // Explanation
}
```

### Use Cases

- Security advisory validation
- Blog post fact-checking
- Technical documentation review
- Research paper claims
- Compliance verification

### Example

```go
report := claims.NewClaimsReport("security-advisory.md")

// External source validation
claim := claims.NewClaim("cvss", "CVSS 8.8", claims.ClaimRiskAssessment,
    claims.Location{Section: "severity"})
claim.SetValidation(claims.NewExternalValidation(
    "https://nvd.nist.gov/vuln/detail/CVE-2026-25253",
    claims.ExternalNVD,
))
report.AddClaim(*claim)

// Internal validation via code
exploit := claims.NewClaim("exploit", "RCE confirmed", claims.ClaimTechnicalFinding,
    claims.Location{Section: "impact"})
exploit.SetValidation(claims.NewInternalValidation(
    claims.MethodCodeExecution, "poc.py", true,
))
report.AddClaim(*exploit)

report.Finalize()
```

## Comparison

| Aspect | Rubric | SummaryReport | ClaimsReport |
|--------|--------|---------------|--------------|
| **Purpose** | Subjective assessment | Deterministic checks | Source validation |
| **Scoring** | Categorical (pass/partial/fail) | Binary (go/warn/nogo) | Verdict (verified/unverified) |
| **Structure** | Categories + Findings | Teams + Tasks | Claims + Validation |
| **Source** | LLM or human reviewer | Automated systems | External URLs or internal evidence |
| **Output** | Quality assessment | GO/NO-GO decision | Publication readiness |

## When to Use Which

### Use Rubric when:

- Assessment requires judgment
- Multiple criteria need scoring
- Detailed findings with recommendations needed
- Reproducibility via rubrics matters

### Use SummaryReport when:

- Checks are pass/fail
- Results come from automated systems
- Aggregating across teams/domains
- CI/CD pipeline integration

### Use ClaimsReport when:

- Document contains factual claims
- Claims need source references (URLs)
- Internal evidence validates claims (code, testing)
- Publishing requires fact-checking

## Next Steps

- [Categorical Scoring](scoring.md) - Understand pass/partial/fail
- [Findings & Severity](findings.md) - Issue tracking
- [Claims Validation](../features/claims.md) - Source verification
- [DAG Aggregation](../features/dag-aggregation.md) - Multi-agent coordination
