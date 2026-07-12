# Structured Evaluation

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/structured-evaluation/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/structured-evaluation
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/structured-evaluation
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/structured-evaluation
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fstructured-evaluation
 [loc-svg]: https://tokei.rs/b1/github/plexusone/structured-evaluation
 [repo-url]: https://github.com/plexusone/structured-evaluation
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/structured-evaluation/blob/main/LICENSE

A reusable evaluation framework for LLM-as-Judge and multi-agent workflows.

## Overview

`structured-evaluation` provides standardized types for evaluation reports, enabling:

- ⚖️ **LLM-as-Judge assessments** with categorical and 1-5 integer scoring
- 🔧 **Automated repair** via reason codes with repair prompts
- 📊 **Coverage tracking** for spec completeness metrics
- 📈 **Confidence & routing** for human review of low-confidence evaluations
- ✅ **GO/NO-GO summary reports** for deterministic checks (CI, tests, validation)
- 🔗 **Multi-agent coordination** with DAG-based report aggregation
- 📋 **Claims validation** for factual claim extraction and source verification

## Architecture

```
┌───────────────────────────────────────────────────────────┐
│                    SummaryReport (GO/NO-GO)               │
│  ┌──────────────────────┐  ┌──────────────────────┐       │
│  │  Embedded Reports    │  │   Team Sections      │       │
│  │  (Full-Fidelity)     │  │   (Task Results)     │       │
│  └──────────────────────┘  └──────────────────────┘       │
└───────────────────────────────────────────────────────────┘
                              ▲
              ┌───────────────┴───────────────┐
              │                               │
┌─────────────┴─────────────┐   ┌─────────────┴─────────────┐
│     Rubric (rubric/)      │   │   ClaimsReport (claims/)  │
│  ┌─────────────────────┐  │   │  ┌─────────────────────┐  │
│  │ Category Results    │  │   │  │ Claims + Validation │  │
│  │ (pass/partial/fail) │  │   │  │ (verified/rejected) │  │
│  ├─────────────────────┤  │   │  ├─────────────────────┤  │
│  │ Findings            │  │   │  │ Sources             │  │
│  │ (severity-based)    │  │   │  │ (external/internal) │  │
│  └─────────────────────┘  │   │  └─────────────────────┘  │
│  LLM-as-Judge scoring     │   │  Fact verification        │
└───────────────────────────┘   └───────────────────────────┘
```

**Three complementary report types:**

| Package | Purpose | Evaluation Type |
|---------|---------|-----------------|
| `rubric/` | Categorical scoring with findings | Subjective (LLM-as-Judge) |
| `claims/` | Fact verification with sources | Objective (source-backed) |
| `summary/` | GO/NO-GO aggregation | Deterministic |

## Installation

```bash
go get github.com/plexusone/structured-evaluation
```

## Packages

| Package | Description |
|---------|-------------|
| `rubric` | Rubric, CategoryResult, Finding, Severity types for LLM-as-Judge |
| `claims` | ClaimsReport, Claim, Validation, Verdict for source verification |
| `summary` | SummaryReport, TeamSection, TaskResult for GO/NO-GO checks |
| `combine` | DAG-based report aggregation using Kahn's algorithm |
| `render/box` | ASCII box renderer for deterministic TUI output |
| `render/detailed` | Detailed terminal renderer for rubric reports |
| `render/terminal` | ANSI-colored terminal renderer with UTF8 icons |
| `render/markdown` | Markdown report renderer |
| `schema` | JSON Schema generation and embedding |

## Report Types

### Rubric (LLM-as-Judge)

For subjective quality assessments with detailed findings:

```go
import "github.com/plexusone/structured-evaluation/rubric"

report := rubric.NewRubric("prd", "document.md")

// Add category with 1-5 integer score and confidence
result := rubric.NewCategoryResultWithIntScore(
    "problem_definition",
    rubric.ScoreGood,  // 4/5
    0.9,               // High confidence
    "Clear problem statement with measurable goals",
)
report.AddCategoryResult(*result)

// Add finding with reason code for automated repair
finding := rubric.NewFindingWithCode(
    "f1", "metrics",
    rubric.CodeMETRICNoBaseline,
    "Missing baseline metrics",
    "No baseline measurements defined for success metrics",
)
finding.SetRecommendation("Add current baseline measurements")
report.AddFinding(*finding)

// Set minimum score threshold
report.PassCriteria.MinIntScore = rubric.ScoreGood  // Require 4+

report.Finalize(nil, "sevaluation check document.md")
// report.IntScore, report.Confidence, report.Blocking are computed
```

### Summary Report (GO/NO-GO)

For deterministic checks with pass/fail status:

```go
import "github.com/plexusone/structured-evaluation/summary"

report := summary.NewSummaryReport("my-service", "v1.0.0", "Release Validation")
report.AddTeam(summary.TeamSection{
    ID:   "qa",
    Name: "Quality Assurance",
    Tasks: []summary.TaskResult{
        {ID: "unit-tests", Status: summary.StatusGo, Detail: "Coverage: 92%"},
        {ID: "e2e-tests", Status: summary.StatusWarn, Detail: "2 flaky tests"},
    },
})
```

### Claims Report (v0.6.0)

For factual claim extraction and source validation:

```go
import "github.com/plexusone/structured-evaluation/claims"

report := claims.NewClaimsReport("security-advisory.md")

// External source: CVE from NVD
claim := claims.NewClaim("cvss", "CVSS 8.8 High", claims.ClaimRiskAssessment,
    claims.Location{Section: "severity"})
claim.SetValidation(claims.NewExternalValidation(
    "https://nvd.nist.gov/vuln/detail/CVE-2026-25253",
    claims.ExternalNVD,
))
report.AddClaim(*claim)

// Internal validation: exploit confirmed via code
exploit := claims.NewClaim("exploit", "RCE confirmed", claims.ClaimTechnicalFinding,
    claims.Location{Section: "impact"})
exploit.SetValidation(claims.NewInternalValidation(
    claims.MethodCodeExecution, "poc.py", true,
))
report.AddClaim(*exploit)

report.Finalize()
// report.Decision.Passed, report.Summary.Counts
```

## Severity Levels

Following InfoSec conventions:

| Severity | Icon | Blocking | Description |
|----------|------|----------|-------------|
| Critical | 🔴 | Yes | Must fix before approval |
| High | 🔴 | Yes | Must fix before approval |
| Medium | 🟡 | No | Should fix, tracked |
| Low | 🟢 | No | Nice to fix |
| Info | ⚪ | No | Informational only |

## Pass Criteria

Default criteria (zero blocking findings, all categories passing):

```go
criteria := rubric.DefaultPassCriteria()
// MaxCritical: 0, MaxHigh: 0, MaxMedium: -1 (unlimited), RequireAllPass: false

criteria := rubric.StrictPassCriteria()
// MaxCritical: 0, MaxHigh: 0, MaxMedium: 3, RequireAllPass: true
```

## Report Validation (v0.7.0)

Validate evaluation reports for correctness:

```go
result := rubric.ValidateReport(&report)

if !result.Valid {
    fmt.Printf("Invalid: %d errors, %d warnings\n", result.ErrorCount, result.WarningCount)
    for _, issue := range result.Issues {
        fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Path, issue.Message)
    }
}

// Get valid enum values for tooling
scores := rubric.ValidScoreValues()         // ["pass", "partial", "fail"]
severities := rubric.ValidSeverityValues()  // ["critical", "high", "medium", "low", "info"]
```

## CLI Tool

```bash
# Install
go install github.com/plexusone/structured-evaluation/cmd/sevaluation@latest

# Render reports
sevaluation render report.json --format=detailed
sevaluation render report.json --format=terminal   # ANSI colors + UTF8 icons
sevaluation render report.json --format=markdown   # Markdown output
sevaluation render report.json --format=box
sevaluation render report.json --format=json

# Lint reports for correctness (v0.7.0)
sevaluation lint report.json              # Basic validation
sevaluation lint report.json --strict     # Warnings are errors
sevaluation lint report.json --format=json

# Check pass/fail (exit code 0/1)
sevaluation check report.json

# Validate structure
sevaluation validate report.json

# Generate JSON Schema
sevaluation schema generate -o ./schema/
```

## DAG-Based Aggregation

For multi-agent workflows with dependencies:

```go
import "github.com/plexusone/structured-evaluation/combine"

results := []combine.AgentResult{
    {TeamID: "qa", Tasks: qaTasks},
    {TeamID: "security", Tasks: secTasks, DependsOn: []string{"qa"}},
    {TeamID: "release", Tasks: relTasks, DependsOn: []string{"qa", "security"}},
}

report := combine.AggregateResults(results, "my-project", "v1.0.0", "Release")
// Teams are topologically sorted: qa → security → release
```

## JSON Schema

Schemas are embedded for runtime validation:

```go
import "github.com/plexusone/structured-evaluation/schema"

rubricSchema := schema.RubricSchemaJSON
claimsSchema := schema.ClaimsSchemaJSON
summarySchema := schema.SummarySchemaJSON
```

## RubricSet (v0.4.0)

Define explicit criteria for consistent categorical evaluations:

```go
cat := rubric.NewCategory("quality", "Output Quality", "Overall quality assessment").
    WithPassPartialFail(
        []string{"Meets all requirements, no significant issues"},
        []string{"Meets most requirements, minor issues"},
        []string{"Missing key requirements or major issues"},
    )

// Use default PRD rubric
rubricSet := rubric.DefaultPRDRubricSet()
```

## Judge Metadata (v0.2.0)

Track LLM judge configuration for reproducibility:

```go
judge := rubric.NewJudgeMetadata("claude-3-opus").
    WithProvider("anthropic").
    WithPrompt("prd-eval-v1", "1.0").
    WithTemperature(0.0).
    WithTokenUsage(1500, 800)

report.SetJudge(judge)
```

## Pairwise Comparison (v0.2.0)

Compare two outputs instead of absolute scoring:

```go
comparison := rubric.NewPairwiseComparison(input, outputA, outputB)
comparison.SetWinner(rubric.WinnerA, "A is more accurate", 0.9)

// Aggregate multiple comparisons
result := rubric.ComputePairwiseResult(comparisons)
// result.WinRateA, result.OverallWinner
```

## Multi-Judge Aggregation (v0.4.0)

Combine evaluations from multiple judges:

```go
result := rubric.AggregateEvaluations(evaluations, rubric.AggregationMajority)

// Methods: AggregationMajority, AggregationConservative, AggregationOptimistic
// result.Agreement - inter-judge agreement (0-1)
// result.Disagreements - categories with significant disagreement
// result.ConsolidatedDecision - final aggregated decision
```

## Likert Scales (v0.5.0)

Use 1-5 numeric scales for human comparison studies:

```go
// Create a Likert-scale category
cat := rubric.NewCategory("quality", "Content Quality", "Overall quality").
    WithLikert5(rubric.StandardLikert5Anchors())

// Record a Likert score (automatically maps to categorical)
result := rubric.NewCategoryResultFromLikert("quality", 4, config, "Good quality")
// result.Score = ScorePass, result.NumericScore = 4.0

// Or record both categorical and numeric
result := rubric.NewCategoryResultWithNumeric("quality", rubric.ScorePass, 4.5, "reasoning")
```

## Inter-Rater Reliability (v0.5.0)

Compare LLM evaluations with human ground truth:

```go
// Compute IRR metrics
metrics := rubric.ComputeIRRFromResults(humanResults, llmResults)

fmt.Printf("Exact Agreement: %.1f%%\n", metrics.ExactAgreement*100)
fmt.Printf("Adjacent Agreement: %.1f%%\n", metrics.AdjacentAgreement*100)
fmt.Printf("Pearson r: %.3f\n", metrics.PearsonCorrelation)

// Categorical agreement with confusion matrix
agreement := rubric.ComputeCategoricalAgreement(humanResults, llmResults)
```

## Claims Validation (v0.6.0)

Validate factual claims have proper source backing:

```go
import "github.com/plexusone/structured-evaluation/claims"

report := claims.NewClaimsReport("article.md")

// Source types: external (URL), internal (code/lab), derived, subjective
// Reliability tiers: authoritative, high, medium, low
// Verdicts: verified, unverified, needs-review, rejected

// Configure pass criteria
report.SetCriteria(claims.ClaimsCriteria{
    RequireAllVerified:           true,
    AllowSubjectiveWithDisclaimer: false,
    MinReliabilityTier:           claims.ReliabilityHigh,
})

report.Finalize()
if report.IsPassing() {
    fmt.Println("Ready for publication")
}
```

## Embedded Reports (v0.6.0)

Archive full-fidelity reports within SummaryReport:

```go
report := summary.NewSummaryReport("project", "v1.0.0", "RELEASE")

// Embed detailed reports
report.EmbedRubricReport("quality-review", rubricReport)
report.EmbedClaimsReport("source-validation", claimsReport)

// Retrieve later
var r rubric.Rubric
report.GetEmbeddedRubricReport("quality-review", &r)
```

## OmniObserve Integration

Export evaluations to Opik, Phoenix, or Langfuse:

```go
import "github.com/plexusone/omniobserve/integrations/sevaluation"

// Export to observability platform
err := sevaluation.Export(ctx, provider, traceID, report)
```

## Integration

Designed to work with:

- `github.com/plexusone/omniobserve` - LLM observability (Opik, Phoenix, Langfuse)
- `github.com/grokify/structured-requirements` - PRD evaluation templates
- `github.com/plexusone/multi-agent-spec` - Agent coordination
- `github.com/grokify/structured-changelog` - Release validation

## License

MIT License - see [LICENSE](LICENSE) for details.
