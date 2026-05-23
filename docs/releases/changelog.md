# Changelog

All notable changes to structured-evaluation are documented here.

For the canonical changelog, see [CHANGELOG.md](https://github.com/plexusone/structured-evaluation/blob/main/CHANGELOG.md) in the repository.

## [v0.4.0](v0.4.0.md) - 2026-05-23

**Breaking Changes**: Switch from numeric scores to categorical pass/partial/fail values.

### Highlights

- Switch from numeric scores to categorical pass/partial/fail values
- New terminal and markdown renderers for evaluation reports

### Changed

- **Breaking**: `CategoryScore` renamed to `CategoryResult` with `Score` field (pass/partial/fail)
- **Breaking**: `ScoreStatus` renamed to `ScoreValue` with values pass/partial/fail
- **Breaking**: Removed `WeightedScore` from `EvaluationReport`
- Added `CategoryCounts` to `Decision` (pass/partial/fail counts)
- Updated detailed renderer to display category counts

### Added

- `render/terminal` package with ANSI-colored output and UTF8 icons
- `render/markdown` package for Markdown report generation
- CLI `terminal` and `markdown` render format options

---

## [v0.3.1](https://github.com/plexusone/structured-evaluation/releases/tag/v0.3.1) - 2026-05-19

### Build

- Updated CI workflows to use shared workflows
- Renamed workflow files to standard filenames

### Dependencies

- Bump github.com/invopop/jsonschema from 0.13.0 to 0.14.0

---

## [v0.3.0](v0.3.0.md) - 2026-03-01

### Changed

- **Breaking**: Module path changed from `github.com/agentplexus/structured-evaluation` to `github.com/plexusone/structured-evaluation`

---

## [v0.2.0](v0.2.0.md) - 2026-01-26

### Highlights

- Rubric definitions with score anchors
- Judge metadata tracking
- Pairwise comparison mode
- Multi-judge aggregation

### Added

- Rubric and ScoreAnchor types
- JudgeMetadata for tracking LLM configuration
- PairwiseComparison for relative evaluation
- MultiJudgeResult for aggregating evaluations

---

## [v0.1.0](v0.1.0.md) - 2026-01-26

### Highlights

- Initial release
- LLM-as-Judge evaluation reports
- GO/NO-GO summary reports
- DAG-based aggregation

### Added

- `evaluation` package with core types
- `summary` package for deterministic checks
- `combine` package with DAG aggregation
- `render` packages for terminal and box output
- `sevaluation` CLI tool
