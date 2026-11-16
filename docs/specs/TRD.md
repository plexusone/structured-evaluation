# Structured Evaluation — Technical Requirements

This document records design decisions for features planned in the [roadmap](ROADMAP.md). The package documentation
and the [site](../index.md) describe what exists.

## 1. Encyclopedic Claims Standard (Phase 1)

### Refuted

`refuted` records that a reliable source shows a claim false. It differs from `contradicted`, which records that a
source the claim cites disagrees with it, a fault the author fixes. A refuted claim is kept as a record that it is
false, for example a misconception that circulates widely. Like the other verdicts, `refuted` needs quoted evidence:
at least one `disputes` passage.

### Standard

A report names the standard its judges applied, as `<name>@<version>`, in `ClaimsMetadata.standard`. Tools combining
reports can then refuse to mix standards, and a reader knows which rules a verdict answered to.

### Criteria

New `ClaimsCriteria` fields, each off at its zero value:

| Field | Requires |
|---|---|
| `requireQuotedSupport` | Every verified claim has a `supports` passage with a quote |
| `minSupportConfidence` | A verified claim's confidence is at least this level |
| `exceptionalMinSources` | A claim marked `exceptional` has supporting passages from at least this many distinct sources |

`Claim.exceptional` marks a superlative, a first, a surprising claim or an allegation about a person. Judges set it.

`WikipediaClaimsCriteria()` combines high source reliability, quoted support, at least medium confidence and two
sources for exceptional claims, matching an encyclopedic standard of verifiability. `EvaluateClaims` reports each
failing claim with the criterion it fails.
