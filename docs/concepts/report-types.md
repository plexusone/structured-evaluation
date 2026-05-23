# Report Types

Structured Evaluation provides two report types for different use cases.

## EvaluationReport

For **subjective quality assessments** using LLM-as-Judge or human reviewers.

### Structure

```go
type EvaluationReport struct {
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
report := evaluation.NewEvaluationReport("prd-review", "requirements.md")
report.AddCategory(evaluation.CategoryResult{
    Category:  "clarity",
    Score:     evaluation.ScorePass,
    Reasoning: "Requirements are clearly written",
})
report.AddFinding(evaluation.Finding{
    Severity: evaluation.SeverityMedium,
    Title:    "Missing edge case",
})
report.Finalize("reviewer")
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

## Comparison

| Aspect | EvaluationReport | SummaryReport |
|--------|------------------|---------------|
| **Purpose** | Subjective assessment | Deterministic checks |
| **Scoring** | Categorical (pass/partial/fail) | Binary (go/warn/nogo) |
| **Structure** | Categories + Findings | Teams + Tasks |
| **Source** | LLM or human reviewer | Automated systems |
| **Findings** | Detailed with severity | Simple status + detail |

## When to Use Which

### Use EvaluationReport when:

- Assessment requires judgment
- Multiple criteria need scoring
- Detailed findings with recommendations needed
- Reproducibility via rubrics matters

### Use SummaryReport when:

- Checks are pass/fail
- Results come from automated systems
- Aggregating across teams/domains
- CI/CD pipeline integration

## Next Steps

- [Categorical Scoring](scoring.md) - Understand pass/partial/fail
- [Findings & Severity](findings.md) - Issue tracking
- [DAG Aggregation](../features/dag-aggregation.md) - Multi-agent coordination
