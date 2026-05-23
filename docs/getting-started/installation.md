# Installation

## Go Module

Add structured-evaluation to your Go project:

```bash
go get github.com/plexusone/structured-evaluation
```

## CLI Tool

Install the `sevaluation` command-line tool:

```bash
go install github.com/plexusone/structured-evaluation/cmd/sevaluation@latest
```

Verify installation:

```bash
sevaluation version
```

## Requirements

- Go 1.21 or later
- No external dependencies required

## Package Structure

| Package | Import Path | Description |
|---------|-------------|-------------|
| `evaluation` | `github.com/plexusone/structured-evaluation/evaluation` | Core evaluation types |
| `summary` | `github.com/plexusone/structured-evaluation/summary` | GO/NO-GO summary reports |
| `combine` | `github.com/plexusone/structured-evaluation/combine` | DAG-based aggregation |
| `render/terminal` | `github.com/plexusone/structured-evaluation/render/terminal` | ANSI terminal renderer |
| `render/markdown` | `github.com/plexusone/structured-evaluation/render/markdown` | Markdown renderer |
| `render/detailed` | `github.com/plexusone/structured-evaluation/render/detailed` | Detailed terminal renderer |
| `render/box` | `github.com/plexusone/structured-evaluation/render/box` | Box-format renderer |
| `schema` | `github.com/plexusone/structured-evaluation/schema` | JSON Schema definitions |

## Next Steps

- [Quick Start](quickstart.md) - Create your first evaluation report
- [Report Types](../concepts/report-types.md) - Understand the evaluation model
