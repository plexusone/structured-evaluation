# CLI Commands

Detailed reference for all `sevaluation` commands.

## render

Render evaluation or summary reports in various formats.

### Usage

```bash
sevaluation render <file> [--format=<format>]
```

### Formats

| Format | Description | Report Type |
|--------|-------------|-------------|
| `terminal` | ANSI colors + UTF8 icons | Rubric |
| `markdown` | Markdown output | Rubric |
| `detailed` | Verbose terminal output | Rubric |
| `box` | ASCII box format (deterministic) | Rubric, Summary |
| `json` | Pretty-printed JSON | Rubric, Summary, Claims |
| `html` | Self-contained HTML audit page, grouped by verdict | Claims |

### Examples

```bash
# Terminal output (default)
sevaluation render eval.json

# Explicit format
sevaluation render eval.json --format=terminal

# Markdown for documentation
sevaluation render eval.json --format=markdown > report.md

# JSON for programmatic use
sevaluation render eval.json --format=json | jq '.decision'

# Box format for summary reports
sevaluation render summary.json --format=box

# Claims report -> self-contained HTML audit page, grouped by verdict
sevaluation render claims.json --format=html > report.html
```

### Auto-Detection

The command auto-detects report type from top-level JSON keys — `categories`
(rubric/evaluation), `teams` (summary), or `claims` (claims report) — and
uses the appropriate renderer.

---

## lint

Validate report correctness. Auto-detects evaluation reports (`categories`)
vs. claims reports (`claims`) from the file's top-level keys and applies the
matching checks. Added in v0.7.0; claims-report support added in v0.13.0.

### Usage

```bash
sevaluation lint <file> [--strict] [--format=<format>]
```

### Flags

| Flag | Description |
|------|-------------|
| `--strict` | Treat warnings as errors |
| `--format` | Output format: `text` (default), `json` |

### Validation Checks — Evaluation Reports

| Check | Level | Description |
|-------|-------|-------------|
| Enum Values | Error | Score, severity, decision status must be valid |
| Required Fields | Error | `metadata.document` and `reviewType` are required |
| Finding Title | Error | Each finding must have a title |
| Count Accuracy | Warning | Reported counts must match actual data |
| Decision Consistency | Warning | Decision should align with findings |
| OverallDecision Match | Warning | `overallDecision` should match `decision.status` |

### Validation Checks — Claims Reports

Evidence-integrity checks for `verified` claims — see
[Evidence-Integrity Linting](../features/claims.md#evidence-integrity-linting-v0130)
for the full explanation of each rule.

| Rule | Level | Description |
|------|-------|-------------|
| `claim-missing-id` / `claim-duplicate-id` | Error | Every claim has a unique `id` |
| `verified-requires-validation` | Error | A verified claim has a `validation` object |
| `verified-requires-url` / `verified-requires-quote` | Error | A verified external claim has a source URL and a verbatim quote |
| `verified-value-in-quote` | Warning | The claim's statistical value appears in the quoted text |
| `verified-derived-needs-sources` | Error | A verified derived claim lists source claim ids |
| `verified-internal-needs-evidence` | Error | A verified internal claim has an evidence path or output |
| `verified-subjective` | Warning | A subjective estimate is published as verified — confirm intentionally |
| `verified-role-needs-corroboration` | Error | A `secondary-analysis`/`self-reported` source has a corroborating `relatedClaimId` |
| `verified-insufficient-corroboration` | Error | Opt-in via `criteria.minCorroboratingSources` — a verified claim has fewer independent sources than required |
| `verified-stale-as-of-date` | Error | Opt-in via `criteria.maxClaimAge` — a verified statistic's `asOfDate` is older than the freshness threshold |

### Exit Codes

| Code | Status | Meaning |
|------|--------|---------|
| 0 | Valid | No errors (warnings allowed unless `--strict`) |
| 1 | Invalid | Has errors, or has warnings with `--strict` |

### Examples

```bash
# Basic validation
sevaluation lint report.json

# Strict mode (warnings are errors)
sevaluation lint report.json --strict

# JSON output for programmatic use
sevaluation lint report.json --format=json

# Claims report — fails if any verified claim lacks a quote/URL/corroboration
sevaluation lint claims.json --strict

# CI pipeline validation
for report in reports/*.json; do
    if ! sevaluation lint "$report" --strict; then
        echo "Invalid report: $report"
        exit 1
    fi
done
```

### Output

```
✅ Valid: 0 errors, 0 warnings
```

or

```
❌ Invalid: 2 errors, 1 warning
  [error] categories[0].score: invalid score value "passed", must be one of: pass, partial, fail
  [error] findings[1].title: finding must have a title
  [warning] summary.categoryCount: reported 5, actual 4
```

---

## check

Check if a report passes evaluation criteria. Useful for CI/CD gates.
Auto-detects evaluation, summary, and claims reports; for a claims report,
this checks `report.Decision.Passed` (the criteria decision, e.g. from
`MinCorroboratingSources`/`MaxClaimAge`) — it does not run `claims.Lint`, so
use `lint` as a separate, stricter gate.

### Usage

```bash
sevaluation check <file>
```

### Exit Codes

| Code | Status | Meaning |
|------|--------|---------|
| 0 | Pass | All criteria met |
| 1 | Fail/Conditional | Blocking issues or conditional pass |

### Examples

```bash
# Basic check
sevaluation check report.json

# Use in CI
if sevaluation check report.json; then
    echo "✅ Evaluation passed"
    deploy_to_production
else
    echo "❌ Evaluation failed"
    exit 1
fi

# Capture output
result=$(sevaluation check report.json 2>&1)
```

### Output

```
✅ PASSED: prd-review (4/4 categories)
```

or

```
❌ FAILED: prd-review - 2 blocking issues
```

---

## validate

Validate a JSON file against the appropriate schema.

### Usage

```bash
sevaluation validate <file>
```

### Examples

```bash
# Validate rubric report
sevaluation validate eval.json

# Validate summary report
sevaluation validate summary.json
```

### Output

```
✅ Valid rubric report
```

or

```
❌ Invalid: missing required field "review_type"
```

---

## schema

Work with JSON schemas.

### Subcommands

| Subcommand | Description |
|------------|-------------|
| `generate` | Generate schema files |
| `show` | Display embedded schema |

### generate

```bash
sevaluation schema generate -o ./schema/
```

Generates:

- `rubric.schema.json`
- `summary.schema.json`

### show

```bash
# Show rubric schema
sevaluation schema show rubric

# Show summary schema
sevaluation schema show summary
```

---

## version

Print version information.

### Usage

```bash
sevaluation version
```

### Output

```
sevaluation v0.7.0
```

---

## Common Workflows

### CI Pipeline Gate

```yaml
# .github/workflows/pr-check.yaml
- name: Run evaluation
  run: |
    sevaluation check eval-report.json
```

### Generate Documentation

```bash
# Generate markdown reports
for f in reports/*.json; do
    name=$(basename "$f" .json)
    sevaluation render "$f" --format=markdown > "docs/reports/${name}.md"
done
```

### Validate Before Commit

```bash
# Pre-commit hook
#!/bin/bash
for report in reports/*.json; do
    if ! sevaluation validate "$report"; then
        echo "Invalid report: $report"
        exit 1
    fi
done
```

### Lint Reports in CI

```bash
# Validate all reports with strict mode
for report in reports/*.json; do
    if ! sevaluation lint "$report" --strict; then
        echo "Invalid report: $report"
        exit 1
    fi
done
```

### Compare Reports

```bash
# Extract decisions from multiple reports
for f in reports/*.json; do
    echo -n "$f: "
    sevaluation render "$f" --format=json | jq -r '.decision.status'
done
```

## Next Steps

- [Report Rendering](../features/rendering.md) - Format details
- [Pass Criteria](../concepts/pass-criteria.md) - Check thresholds
