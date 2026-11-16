# Structured Evaluation — Roadmap

**Repository:** `github.com/plexusone/structured-evaluation`
**Status:** In progress — 0 of 1 phases completed

> RMI IDs are stable and permanent. Commits implementing an item carry the trailer
> `Refs: RMI-STRUCTUREDEVALUATION-NNN`. A phase is complete only when all of its RMIs are complete.

## Phase 1 — Encyclopedic Claims Standard

**Theme:** Claims reports name the standard they were judged against, record claims shown false, and can be checked
against criteria at Wikipedia's standard.
**Status:** Planned — 0 of 1 items completed

- [ ] `RMI-STRUCTUREDEVALUATION-001` Refuted verdict, report standard and encyclopedic criteria
  - Acceptance: verdict `refuted` (a reliable source shows the claim false), distinct from `contradicted` (a cited source disagrees); `ClaimsReport.standard` and `ClaimsMetadata.standard` name the standard and version applied; `ClaimsCriteria` gains `requireQuotedSupport`, `minSupportConfidence` and `exceptionalMinSources`, and `Claim.exceptional` marks a claim needing more sources; `WikipediaClaimsCriteria()` sets high reliability, quoted support, at least medium confidence and two sources for exceptional claims; `EvaluateClaims` applies the new fields; schema and TypeScript types regenerated
