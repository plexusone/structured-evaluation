# CLAUDE.md - structured-evaluation

Project-specific instructions for Claude Code.

## Project Overview

`structured-evaluation` is a reusable Go library defining standardized types
for evaluation reports: `rubric` (LLM-as-Judge categorical scoring),
`claims` (factual claim extraction and source verification), and `summary`
(deterministic GO/NO-GO aggregation). Consumed as an importable Go module
and, via `cmd/sevaluation`, as a CLI. Downstream consumers include
`github.com/plexusone/agent-team-stats` (produces `ClaimsReport`s as
evidence) and other org repos wanting evaluation/decision types without
redefining them.

## Architecture

- **Library-first.** Report types and logic live in `rubric/`, `claims/`,
  `summary/`, `combine/` (DAG aggregation); `cmd/sevaluation` is a thin
  Cobra adapter over one shared layer, `cmd/genschema` regenerates schemas.
- **`render/*`** are independent renderer packages (`box`, `detailed`,
  `terminal`, `markdown`, `html`) — one report type in, one output format
  out, no shared renderer state.
- **`schema/`** holds JSON Schema files, generated (never hand-edited) from
  the Go structs via `invopop/jsonschema` reflection — Go-first, per org
  convention. Regenerate with `go run ./cmd/genschema` after any exported
  type change in `claims`, `rubric`, or `summary`, and check the diff is
  additive/expected before committing (unrelated schema drift signals an
  accidental type change elsewhere).
- **`ts/`** is a generated TypeScript/Zod package downstream of the same
  JSON Schema — regenerate via `npm run generate` in `ts/`, never hand-edit.
- **`docs/`** is an MkDocs Material site (`mkdocs.yml`); `docs/releases/`
  holds one file per version, indexed in `mkdocs.yml`'s nav and in
  `docs/releases/changelog.md`.

## Claims Package: Design Intent

The `claims` package exists to make "verified" mean something specific and
checkable, not just an asserted label. Two complementary layers:

- **`DetermineVerdict`** computes a verdict from a `Validation` — but a
  verdict can also be hand-authored directly on a `Claim`, bypassing that
  computation.
- **`claims.Lint`** re-checks the result regardless of how the verdict was
  set: a verified claim must carry real evidence (URL + verbatim quote,
  source claim ids, evidence path/output), sufficient independent
  corroboration when required, and a statistic that isn't stale. `Lint`
  never mutates the report and only gates `VerdictVerified` claims.

This distinction (compute vs. re-check) exists because a real false
positive motivated it: a claim marked `verified` cited a figure synthesized
from secondary reporting, and a later primary-source check found a
substantially different number. `SourceRole`, `MinCorroboratingSources`, and
`MaxClaimAge` (all v0.13.0) are the direct fixes for that failure mode — see
[docs/features/claims.md](docs/features/claims.md#evidence-integrity-linting-v0130).

New criteria/lint additions to this package should stay **additive and
opt-in by default**: verify any change against real ledgers (e.g. in a
consuming repo's case-study data) to confirm identical lint output unless
the new field is explicitly set.

## Conventions

- Follow the org Go standards: Cobra for CLI, conventional commits, gofmt +
  `golangci-lint`, error handling per the priority order (return, never
  discard).
- Verify dependency versions before bumping (`go list -m -versions ...`).
- Every exported type change: regenerate `schema/*.schema.json` via
  `cmd/genschema` and consider whether `ts/` needs regenerating too.
- Releases: update `CHANGELOG.json` (compact single-line leaf objects — the
  existing style, not `json.dump(indent=2)`), regenerate `CHANGELOG.md` via
  `schangelog generate`, add `docs/releases/vX.Y.Z.md`, and add the release
  to `mkdocs.yml`'s nav. Tag only after commits are pushed and CI is green.
