# Release Notes - v0.3.0

**Release Date:** 2026-03-01

## Overview

This release updates the module path following the GitHub organization rename from `agentplexus` to `plexusone`.

## Breaking Changes

### Module Path Changed

The Go module path has changed from `github.com/agentplexus/structured-evaluation` to `github.com/plexusone/structured-evaluation`.

**Before:**

```go
import "github.com/agentplexus/structured-evaluation/evaluation"
import "github.com/agentplexus/structured-evaluation/summary"
import "github.com/agentplexus/structured-evaluation/combine"
```

**After:**

```go
import "github.com/plexusone/structured-evaluation/evaluation"
import "github.com/plexusone/structured-evaluation/summary"
import "github.com/plexusone/structured-evaluation/combine"
```

### JSON Schema URLs Updated

The `$id` fields in JSON Schema files have been updated:

- `evaluation.schema.json`: `https://github.com/plexusone/structured-evaluation/schema/evaluation.schema.json`
- `summary.schema.json`: `https://github.com/plexusone/structured-evaluation/schema/summary.schema.json`

## Upgrade Guide

Update all import statements in your code:

```bash
# Using sed (macOS)
find . -name "*.go" -exec sed -i '' 's|github.com/agentplexus/structured-evaluation|github.com/plexusone/structured-evaluation|g' {} +

# Using sed (Linux)
find . -name "*.go" -exec sed -i 's|github.com/agentplexus/structured-evaluation|github.com/plexusone/structured-evaluation|g' {} +
```

Then update your dependencies:

```bash
go get github.com/plexusone/structured-evaluation@v0.3.0
go mod tidy
```

## Other Changes

### Documentation

- Updated README with new module path and badge URLs
- Updated all example imports in documentation
- Updated related package references (`omniobserve`, `multi-agent-spec`)

## Contributors

- PlexusOne Team
- Claude Opus 4.5 (Co-Author)
