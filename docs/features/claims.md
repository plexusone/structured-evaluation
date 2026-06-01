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

- [Report Types](../concepts/report-types.md) - Compare EvaluationReport, SummaryReport, ClaimsReport
- [v0.6.0 Release Notes](../releases/v0.6.0.md) - Full changelog
