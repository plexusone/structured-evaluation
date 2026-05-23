# CLI Overview

The `sevaluation` CLI provides commands for working with evaluation reports.

## Installation

```bash
go install github.com/plexusone/structured-evaluation/cmd/sevaluation@latest
```

## Commands

| Command | Description |
|---------|-------------|
| `render` | Render a report in various formats |
| `check` | Check if a report passes criteria (exit code 0/1) |
| `validate` | Validate JSON structure against schema |
| `schema` | Generate or view JSON schemas |
| `version` | Print version information |

## Quick Examples

```bash
# Render report to terminal with colors
sevaluation render report.json --format=terminal

# Check pass/fail for CI
sevaluation check report.json
echo $?  # 0 = pass, 1 = fail

# Generate markdown report
sevaluation render report.json --format=markdown > report.md

# Validate report structure
sevaluation validate report.json
```

## Global Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help |
| `--version` | Show version |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success / Pass |
| 1 | Failure / Fail or Conditional |
| 2 | Invalid input or error |

## Next Steps

- [Commands Reference](commands.md) - Detailed command documentation
- [Report Rendering](../features/rendering.md) - Output formats
