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
| `box` | Box-drawing format | Summary |
| `json` | Pretty-printed JSON | Both |

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
```

### Auto-Detection

The command auto-detects report type (rubric vs summary) and uses the appropriate renderer.

---

## check

Check if a report passes evaluation criteria. Useful for CI/CD gates.

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
sevaluation v0.6.0
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
