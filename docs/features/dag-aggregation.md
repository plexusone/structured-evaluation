# DAG Aggregation

The `combine` package provides DAG-based (Directed Acyclic Graph) report aggregation for multi-agent workflows where teams have dependencies.

## Overview

In complex validation workflows, different teams may depend on each other:

```
QA → Security → Release
     ↓
     Compliance
```

DAG aggregation ensures reports are processed in dependency order using Kahn's algorithm.

## AgentResult

```go
type AgentResult struct {
    TeamID    string       `json:"team_id"`
    Tasks     []TaskResult `json:"tasks"`
    DependsOn []string     `json:"depends_on,omitempty"`
}

type TaskResult struct {
    ID     string `json:"id"`
    Status Status `json:"status"` // go/warn/nogo
    Detail string `json:"detail"`
}
```

## Basic Aggregation

Without dependencies:

```go
import "github.com/plexusone/structured-evaluation/combine"

results := []combine.AgentResult{
    {
        TeamID: "qa",
        Tasks: []combine.TaskResult{
            {ID: "unit-tests", Status: combine.StatusGo, Detail: "100% pass"},
            {ID: "coverage", Status: combine.StatusGo, Detail: "92% coverage"},
        },
    },
    {
        TeamID: "security",
        Tasks: []combine.TaskResult{
            {ID: "sast", Status: combine.StatusGo, Detail: "No critical findings"},
        },
    },
}

report := combine.AggregateResults(results, "my-service", "v2.0.0", "Release Validation")
```

## DAG Aggregation

With dependencies:

```go
results := []combine.AgentResult{
    {
        TeamID: "qa",
        Tasks:  qaTasks,
        // No dependencies - runs first
    },
    {
        TeamID:    "security",
        Tasks:     securityTasks,
        DependsOn: []string{"qa"}, // Waits for QA
    },
    {
        TeamID:    "compliance",
        Tasks:     complianceTasks,
        DependsOn: []string{"qa"}, // Also waits for QA
    },
    {
        TeamID:    "release",
        Tasks:     releaseTasks,
        DependsOn: []string{"security", "compliance"}, // Waits for both
    },
}

report, err := combine.AggregateWithDAG(results, "my-service", "v2.0.0", "Release")
if err != nil {
    // Handle cycle detection or missing dependencies
    log.Fatal(err)
}
```

## Topological Sorting

Teams are sorted using Kahn's algorithm:

```go
// Input order: qa, security, compliance, release
// Output order (topologically sorted): qa → security, compliance → release
```

### Algorithm

1. Build dependency graph
2. Find nodes with no incoming edges (no dependencies)
3. Process those nodes, remove their outgoing edges
4. Repeat until all nodes processed
5. If nodes remain, cycle detected

## Cycle Detection

Circular dependencies are detected and rejected:

```go
results := []combine.AgentResult{
    {TeamID: "a", DependsOn: []string{"b"}},
    {TeamID: "b", DependsOn: []string{"a"}}, // Cycle!
}

_, err := combine.AggregateWithDAG(results, ...)
// Error: cycle detected in dependency graph
```

## Missing Dependencies

References to non-existent teams are caught:

```go
results := []combine.AgentResult{
    {TeamID: "release", DependsOn: []string{"unknown"}}, // Error!
}

_, err := combine.AggregateWithDAG(results, ...)
// Error: unknown dependency "unknown" in team "release"
```

## Multi-Agent Workflow Example

```go
// Agent 1: QA runs tests
qaAgent := func() combine.AgentResult {
    return combine.AgentResult{
        TeamID: "qa",
        Tasks: []combine.TaskResult{
            runUnitTests(),
            runIntegrationTests(),
            checkCoverage(),
        },
    }
}

// Agent 2: Security scans (depends on QA passing)
securityAgent := func() combine.AgentResult {
    return combine.AgentResult{
        TeamID:    "security",
        DependsOn: []string{"qa"},
        Tasks: []combine.TaskResult{
            runSAST(),
            runDependencyScan(),
            runSecretsScan(),
        },
    }
}

// Agent 3: Release validation (depends on QA and Security)
releaseAgent := func() combine.AgentResult {
    return combine.AgentResult{
        TeamID:    "release",
        DependsOn: []string{"qa", "security"},
        Tasks: []combine.TaskResult{
            checkChangelog(),
            checkVersionBump(),
            checkMigrations(),
        },
    }
}

// Run agents and aggregate
results := []combine.AgentResult{
    qaAgent(),
    securityAgent(),
    releaseAgent(),
}

report, err := combine.AggregateWithDAG(results, "my-service", "v2.0.0", "Release")
```

## Parallel Execution

Agents without dependencies on each other can run in parallel:

```go
// These can run in parallel (both depend only on QA):
// - security
// - compliance
// - documentation

// This must wait for all of them:
// - release
```

## Status Propagation

If a dependency fails, dependent teams are blocked:

```go
// QA fails → Security, Compliance, Release all blocked
// Security fails → Release blocked (but Compliance can still run)
```

## Best Practices

### Keep DAGs Shallow

Deep dependency chains slow down workflows:

```go
// ❌ Deep chain
a → b → c → d → e → f

// ✅ Shallow with parallel branches
     ┌→ b ─┐
a →─┼→ c ─┼→ f
     └→ d ─┘
```

### Explicit Dependencies

Always declare dependencies explicitly:

```go
// ✅ Explicit
{TeamID: "release", DependsOn: []string{"qa", "security"}}

// ❌ Implicit ordering (fragile)
```

### Idempotent Tasks

Tasks should be safe to re-run:

```go
// ✅ Idempotent
{ID: "check-coverage", ...}  // Always produces same result

// ❌ Side effects
{ID: "deploy-to-prod", ...}  // Not idempotent!
```

## Next Steps

- [Report Types](../concepts/report-types.md) - Understanding SummaryReport
- [Report Rendering](rendering.md) - Display aggregated reports
