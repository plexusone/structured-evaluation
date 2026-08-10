# Claims Validation

The `claims` package provides types for extracting and validating factual claims in documents. This enables verification that claims are properly sourced (external references) or objectively validated (internal evidence).

## Overview

When publishing content—especially technical documentation, security advisories, or research—it's important to verify that factual claims have proper backing. The claims package provides:

- **Claim extraction**: Identify and categorize claims in documents
- **Source validation**: Track how each claim is validated
- **Reliability tiers**: Categorize source trustworthiness
- **Verdict system**: Determine if claims are verified, unverified, or rejected
- **Pass criteria**: Configure requirements for approval

## Source Types

Claims can be validated through four approaches:

### External Sources

Claims backed by URL references:

```go
validation := claims.NewExternalValidation(
    "https://nvd.nist.gov/vuln/detail/CVE-2026-25253",
    claims.ExternalNVD,
)
validation.External.QuotedText = "CVSS 3.1 Base Score: 8.8"
validation.External.VerifiedMatch = true
```

External source types:

| Type | Description | Reliability |
|------|-------------|-------------|
| `ExternalNVD` | NIST National Vulnerability Database | Authoritative |
| `ExternalVendorAdvisory` | Official vendor advisory | Authoritative |
| `ExternalFramework` | MITRE, OWASP, CWE official docs | Authoritative |
| `ExternalPeerReviewed` | Peer-reviewed publications | High |
| `ExternalReputableVendor` | Research firms, established vendors | High |
| `ExternalAPI` | Public APIs (FIRST.org EPSS) | High |
| `ExternalCommunity` | Blogs, forums, community sources | Medium |
| `ExternalAggregator` | Third-party "stats roundup" sites with no original reporting (e.g. AI-generated SEO stats pages) | Low |

`ExternalAggregator` is distinct from `ExternalCommunity`: a community source
(a named blog post, a forum thread) is still someone's own account of
something. An aggregator has no original reporting to fall back on — it
just reposts numbers, sometimes distorted, from elsewhere — so it defaults
to auto-reject rather than requires-review. Classify a source this way based
on the *domain*, not the label it happens to use: a known aggregator citing
itself as "WHO" or similar should still be classified `ExternalAggregator`.

### Internal Evidence

Claims validated through code, testing, or observation:

```go
validation := claims.NewInternalValidation(
    claims.MethodCodeExecution,
    "exploits/poc.py",
    true, // reproducible
)
validation.Internal.ValidatedBy = "Security Team"
validation.Internal.Output = "Shell access obtained"
validation.Internal.Environment = &claims.ValidationEnvironment{
    Product: "OpenClaw",
    Version: "2.4.1",
    Platform: "Ubuntu 22.04",
}
```

Internal validation methods:

| Method | Description |
|--------|-------------|
| `MethodCodeExecution` | Validated by running code |
| `MethodLabTesting` | Controlled lab testing |
| `MethodCodeReview` | Code inspection |
| `MethodLogAnalysis` | Log file analysis |
| `MethodCalculation` | Mathematical calculation |
| `MethodObservation` | Direct observation |

### Derived Claims

Claims calculated from other validated claims:

```go
validation := claims.NewDerivedValidation(
    []string{"claim-1", "claim-2"}, // source claim IDs
    "aggregation",                   // derivation method
    "risk = cvss * exploitability",  // formula (optional)
)
validation.Derived.Reasoning = "Combined risk based on CVSS and exploit availability"
```

### Subjective Estimates

Estimates that lack objective backing:

```go
validation := claims.NewSubjectiveValidation(
    true,  // acknowledged as estimate in document
    claims.RecommendKeepWithDisclaimer,
)
validation.Subjective.Methodology = "Expert judgment"
validation.Subjective.Rationale = "Based on historical incident data"
```

Recommendations for subjective claims:

- `RecommendKeepWithDisclaimer`: Accept with explicit disclaimer
- `RecommendRemove`: Remove from document
- `RecommendFindSource`: Find external source
- `RecommendConvertToInternal`: Validate internally

## Reliability Tiers

External sources are categorized by trustworthiness:

| Tier | Auto-Accept | Description |
|------|-------------|-------------|
| **Authoritative** | Yes | Official, authoritative source |
| **High** | Yes | Highly reputable source |
| **Medium** | Review | Requires human review |
| **Low** | Reject | Unverified, reject by default |

Default reliability by source type:

```go
claims.DefaultReliabilityForSourceType(claims.ExternalNVD)
// → ReliabilityAuthoritative

claims.DefaultReliabilityForSourceType(claims.ExternalCommunity)
// → ReliabilityMedium
```

## Verdicts

Each claim receives a verdict based on its validation:

| Verdict | Blocking | Description |
|---------|----------|-------------|
| **Verified** | No | Properly validated |
| **Unverified** | Yes | No validation provided |
| **NeedsReview** | Configurable | Requires human review |
| **Rejected** | Yes | Failed validation |

Verdict determination:

```go
verdict := claims.DetermineVerdict(validation)

// External: verified if reliability is acceptable
// Internal: verified if reproducible
// Derived: verified if all source claims are verified
// Subjective: needs-review (never auto-verified)
```

## Claims Report

Create a report with multiple claims:

```go
report := claims.NewClaimsReport("security-advisory.md")
report.Metadata.DocumentTitle = "CVE-2026-25253 Advisory"
report.Metadata.DocumentVersion = "1.0.0"

// Add claims
report.AddClaim(*claim1)
report.AddClaim(*claim2)
report.AddClaim(*claim3)

// Configure criteria
report.SetCriteria(claims.ClaimsCriteria{
    RequireAllVerified:           false,
    AllowSubjectiveWithDisclaimer: true,
    AllowNeedsReview:             true,
    MinReliabilityTier:           claims.ReliabilityMedium,
})

// Evaluate
report.Finalize()

// Check result
if report.IsPassing() {
    fmt.Println("Ready for publication")
} else {
    fmt.Println("Issues:", report.Decision.Rationale)
}
```

## Pass Criteria

Configure what's required for approval:

```go
// Strict: all claims must be verified
strict := claims.ClaimsCriteria{
    RequireAllVerified:           true,
    AllowSubjectiveWithDisclaimer: false,
    AllowNeedsReview:             false,
    MinReliabilityTier:           claims.ReliabilityHigh,
}

// Permissive: allow some flexibility
permissive := claims.ClaimsCriteria{
    RequireAllVerified:           false,
    AllowSubjectiveWithDisclaimer: true,
    AllowNeedsReview:             true,
    MinReliabilityTier:           claims.ReliabilityMedium,
}
```

`MinCorroboratingSources`/`CorroborationCategories` and `MaxClaimAge` (v0.13.0)
are additional, opt-in criteria — see
[Evidence-Integrity Linting](#evidence-integrity-linting-v0130). When either
is set, a verified claim that fails it is treated the same way as a
needs-review claim by `EvaluateClaims`: a conditional pass if
`AllowNeedsReview`, otherwise a fail. This only affects the report-level
`Decision`, never the claim's own stored `Verdict`.

## Summary Statistics

After finalization, review statistics:

```go
report.Finalize()

counts := report.Summary.Counts
fmt.Printf("Total: %d, Verified: %d, Unverified: %d\n",
    counts.Total, counts.Verified, counts.Unverified)

// Claims by category
for cat, count := range report.Summary.ByCategory {
    fmt.Printf("%s: %d\n", cat, count)
}

// Claims needing attention
for _, id := range report.Summary.UnverifiedClaims {
    fmt.Printf("Unverified: %s\n", id)
}
```

## Claim Categories

Categorize claims for better organization:

| Category | Description |
|----------|-------------|
| `ClaimMetadata` | Identifiers, versions, dates |
| `ClaimTechnicalFinding` | Technical observations |
| `ClaimFrameworkMapping` | MITRE, OWASP mappings |
| `ClaimRiskAssessment` | Risk and impact assessments |
| `ClaimTimeline` | Temporal claims |
| `ClaimStatistical` | Numeric/statistical claims |
| `ClaimGuidance` | Recommendations |
| `ClaimAttribution` | Source credits |

## Statistical Claims

`ClaimStatistical` claims can carry structured numeric detail in addition to
the free-text `Text` field, via `Claim.Statistical`. This matters because
`Text` is a rendered string (e.g. `"4.7M paid subscribers"`) — without a
separate structured value, there's no way to query, compare, or
programmatically re-check a claim's number against its
`Validation.External.QuotedText` excerpt. `StatisticalDetail` keeps the two
independent:

```go
claim := claims.NewClaim("stat-1", "4.7M paid subscribers", claims.ClaimStatistical,
    claims.Location{Section: "metrics"})

claim.SetValidation(claims.NewExternalValidation(
    "https://example.com/earnings",
    claims.ExternalReputableVendor,
))
claim.Validation.External.QuotedText = "4.7 million paid subscribers"

asOf := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
claim.SetStatistical(
    claims.NewStatisticalDetail(4_700_000, "subscribers", claims.PrecisionExact).
        WithAsOfDate(asOf),
)
```

`AsOfDate` is deliberately separate from `Validation.External.AccessedAt`:
`AccessedAt` is when the source URL was crawled; `AsOfDate` is when the
underlying fact was true (e.g. the earnings period a figure describes, or
the date a public statement was made). A claim crawled today can describe a
fact from months earlier — conflating the two dates misdates the claim.

`AccessedAt` (and `InternalValidation.ValidatedAt`) are `*time.Time`, nil
when unset — set them via the chainable `WithAccessedAt`/`WithValidatedAt`,
or let `NewExternalValidation`/`NewInternalValidation` default them to
`time.Now().UTC()`:

```go
validation.External.WithAccessedAt(time.Now().UTC())
```

### Precision

| Value | Meaning |
|-------|---------|
| `PrecisionExact` | The source states this value directly |
| `PrecisionApproximate` | The source itself qualifies it (e.g. "1M+", "~50%", "over 20,000") |
| `PrecisionEstimated` | A third-party or analyst estimate, not source-stated |
| `PrecisionRange` | A bounded range rather than a single point value — put the more notable end in `Value` and describe the range in `Text` |

`Statistical` is nil for non-statistical claims, and setting it on a claim
whose `Category` isn't `ClaimStatistical` is harmless but not rendered by
convention — it's meant for numeric claims specifically.

## Evidence-Integrity Linting (v0.13.0)

`DetermineVerdict` computes a verdict from a `Validation`, but a verdict can
also be hand-authored directly on a `Claim` — bypassing that computation
entirely. `claims.Lint` re-checks the result: it re-verifies that every claim
stated as `verified` actually earns the label, regardless of how the verdict
was set. This is the exact gap that let a claim marked `verified` cite a
figure synthesized from secondary reporting, when the primary source said
something substantially different.

```go
findings := claims.Lint(report)
if claims.HasErrors(findings) {
    for _, f := range findings {
        fmt.Printf("[%s] %s: %s\n", f.Severity, f.ClaimID, f.Message)
    }
}
```

Or from the CLI:

```bash
sevaluation lint report.json           # errors fail (exit 1), warnings are advisory
sevaluation lint report.json --strict  # warnings fail too
```

`Lint` only gates `Verdict == VerdictVerified` claims — needs-review,
rejected, and unverified claims carry no evidence obligation. It never
mutates the report.

### Baseline checks (always on)

| Rule | Severity | Applies to | Checks |
|------|----------|------------|--------|
| `claim-missing-id` / `claim-duplicate-id` | Error | Any claim | Every claim has a unique `ID` |
| `verified-requires-validation` | Error | Verified | `Validation` is not nil |
| `verified-requires-url` / `verified-requires-quote` | Error | External | `URL` and `QuotedText` are both set |
| `verified-value-in-quote` | Warning | External, with `Statistical` | `Statistical.Value` appears in `QuotedText` — advisory because a legitimate quote can state a rule ("0% target") or a range, or scale the value into a unit ("20 million") rather than repeat the literal number |
| `verified-derived-needs-sources` | Error | Derived | `SourceClaimIDs` is non-empty |
| `verified-internal-needs-evidence` | Error | Internal | `EvidencePath` or `Output` is set |
| `verified-subjective` | Warning | Subjective | Flags that a subjective estimate is published as verified — confirm intentionally |

### `SourceRole`: how directly a source speaks for the claim

`ExternalSourceType` (NVD, reputable vendor, community, aggregator, ...)
categorizes a source's general authority. `SourceRole` is a separate,
orthogonal axis: how directly *this particular citation* speaks for the
claim, independent of how authoritative the outlet is overall.

```go
type SourceRole string

const (
    SourceRolePrimary          SourceRole = "primary"           // the claim's own subject states it
    SourceRoleSecondaryRelay   SourceRole = "secondary-relay"    // a wire report quoting the primary source directly
    SourceRoleSecondaryAnalysis SourceRole = "secondary-analysis" // an outlet's own synthesis/estimate across other reporting
    SourceRoleSelfReported     SourceRole = "self-reported"      // the claim's subject's own marketing/PR material
)
```

Two "reputable vendor" citations can carry very different trust: a wire
report relaying an earnings call is near-primary, but a research outlet's own
synthesized estimate across several other reports is `secondary-analysis` —
exactly the kind of citation that can drift from the underlying fact. Set it
on `ExternalValidation.Role` (`sourceRole` in JSON), optional and empty by
default so existing reports are unaffected:

```go
validation.External.Role = claims.SourceRoleSecondaryAnalysis
```

`SourceRole.RequiresCorroboration()` is `true` for `secondary-analysis` and
`self-reported` (`false` for `primary`, `secondary-relay`, and an unset
role). When it's true, `Lint` requires at least one corroborating claim in
`RelatedClaimIDs`, or errors with `verified-role-needs-corroboration`:

```go
claim.RelatedClaimIDs = []string{"arr-independent-estimate"}
```

### `MinCorroboratingSources`: a general corroboration threshold

`SourceRole` encodes a fixed policy (secondary-analysis and self-reported
always need corroboration). `ClaimsCriteria.MinCorroboratingSources` is a
separate, configurable threshold that applies regardless of role or
validation type — useful for high-stakes categories where even a `primary`
source shouldn't stand alone:

```go
report.Criteria = claims.ClaimsCriteria{
    MinCorroboratingSources: 2,
    // Optional: restrict the requirement to specific categories.
    // An empty slice applies it to every category.
    CorroborationCategories: []claims.ClaimCategory{claims.ClaimStatistical},
}
```

A claim counts its own source plus each entry in `RelatedClaimIDs` toward
the threshold (`IsSufficientlyCorroborated`). Disabled by default
(`MinCorroboratingSources <= 1`), so opting in is required. `Lint` flags a
shortfall as `verified-insufficient-corroboration`; `EvaluateClaims` folds
it into the report-level decision the same way as needs-review claims (see
[Pass Criteria](#pass-criteria) below).

### `MaxClaimAge`: staleness

A verified statistic can quietly go stale — presented as current years after
the underlying fact was true (e.g. a 2022–23 study cited as if still
representative). `ClaimsCriteria.MaxClaimAge` bounds how old a verified
claim's `Statistical.AsOfDate` may be, relative to now:

```go
report.Criteria = claims.ClaimsCriteria{
    MaxClaimAge: 365 * 24 * time.Hour,
}
```

A claim with no `Statistical` detail, or an unset `AsOfDate`, is never
flagged — age unknown is not the same as stale. Disabled by default
(`MaxClaimAge <= 0`). `Lint` flags a violation as `verified-stale-as-of-date`;
`EvaluateClaims` folds it into the report-level decision like needs-review.

## Example: Security Advisory

Complete example validating a security advisory:

```go
package main

import (
    "fmt"
    "github.com/plexusone/structured-evaluation/claims"
)

func main() {
    report := claims.NewClaimsReport("CVE-2026-25253-advisory.md")

    // CVE ID from NVD
    cve := claims.NewClaim("cve-id", "CVE-2026-25253", claims.ClaimMetadata,
        claims.Location{Section: "header"})
    cve.SetValidation(claims.NewExternalValidation(
        "https://nvd.nist.gov/vuln/detail/CVE-2026-25253",
        claims.ExternalNVD,
    ))
    report.AddClaim(*cve)

    // CVSS score from NVD
    cvss := claims.NewClaim("cvss", "CVSS 8.8 High", claims.ClaimRiskAssessment,
        claims.Location{Section: "severity"})
    cvss.SetValidation(claims.NewExternalValidation(
        "https://nvd.nist.gov/vuln/detail/CVE-2026-25253",
        claims.ExternalNVD,
    ))
    report.AddClaim(*cvss)

    // Exploit confirmed via testing
    exploit := claims.NewClaim("exploit", "Remote code execution confirmed",
        claims.ClaimTechnicalFinding, claims.Location{Section: "impact"})
    exploit.SetValidation(claims.NewInternalValidation(
        claims.MethodCodeExecution, "exploits/poc.py", true,
    ))
    report.AddClaim(*exploit)

    // Evaluate
    report.Finalize()

    fmt.Printf("Decision: %s\n", report.Decision.Status)
    fmt.Printf("Claims: %d verified, %d unverified\n",
        report.Summary.Counts.Verified,
        report.Summary.Counts.Unverified,
    )
}
```

## Next Steps

- [Report Types](../concepts/report-types.md) - Compare Rubric, SummaryReport, ClaimsReport
- [CLI Commands](../cli/commands.md#lint) - `sevaluation lint` reference for claims reports
- [v0.13.0 Release Notes](../releases/v0.13.0.md) - Evidence-integrity linting, SourceRole, corroboration, staleness
- [v0.6.0 Release Notes](../releases/v0.6.0.md) - Original claims validation release
