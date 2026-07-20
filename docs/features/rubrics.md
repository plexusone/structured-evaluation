# Rubrics

Rubrics define explicit evaluation criteria for consistent assessments across evaluators (human or LLM).

## Overview

A category defines how one dimension of a document is scored. It can be scored two ways:

- **Flat** — a categorical scale (pass / partial / fail) with criteria per band.
- **Rich** — weighted sub-criteria, each with its own pass/partial/fail bands and concrete indicators (added in v0.10.0).

```go
type Category struct {
    ID          string            `json:"id" yaml:"id"`
    Name        string            `json:"name" yaml:"name"`
    Description string            `json:"description" yaml:"description"`
    Weight      float64           `json:"weight,omitempty" yaml:"weight,omitempty"`
    Required    bool              `json:"required,omitempty" yaml:"required,omitempty"`
    Scale       Scale             `json:"scale" yaml:"scale"`                             // flat scoring
    Criteria    []Criterion       `json:"criteria,omitempty" yaml:"criteria,omitempty"`   // rich scoring
    Examples    *CategoryExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
}
```

Use `cat.IsComposite()` to tell them apart — it returns true when the category
carries weighted `Criteria`.

## Flat Categories

A categorical scale with pass/partial/fail criteria (2-3 options is recommended
for LLM-as-Judge):

```go
cat := rubric.NewCategory("problem_definition", "Problem Definition",
    "Clarity and completeness of the problem statement").
    SetWeight(0.2).
    SetRequired(true).
    WithPassPartialFail(
        []string{"Problem is clearly stated with measurable business impact and affected users identified"},
        []string{"Problem is stated but lacks specificity or measurable impact"},
        []string{"Problem is vague, missing, or not actionable"},
    )
```

Other scale shapes: `WithBinary(pass, fail)`, `WithChecklist(required, optional, threshold)`, `WithLikert(config)`.

## Rich Weighted Criteria (v0.10.0)

When a dimension decomposes into independently scored parts, give the category
`Criteria` instead of a flat scale. Each criterion carries a weight and
pass/partial/fail bands, each with a description and concrete indicators an
evaluator can look for:

```go
cat := rubric.Category{
    ID:     "assumption_coverage",
    Name:   "Assumption Coverage",
    Weight: 25,
    Criteria: []rubric.Criterion{
        {
            ID:     "desirability",
            Name:   "Desirability",
            Weight: 25,
            Pass: rubric.CriterionLevel{
                Description: "Desirability assumptions are identified",
                Indicators:  []string{"customer demand cited", "willingness-to-pay evidence"},
            },
            Fail: rubric.CriterionLevel{
                Description: "No desirability assumptions surfaced",
            },
        },
    },
}
```

Rich categories pair with numeric `scoreThresholds` on the rubric's pass
criteria (weighted roll-up to a 0-100 score) — see
[Pass Criteria](../concepts/pass-criteria.md).

## Adding Examples

Few-shot examples calibrate an evaluator; including the reasoning improves LLM
alignment (chain-of-thought). Attach one per scoring band:

```go
cat.SetExamples(&rubric.CategoryExamples{
    Pass: &rubric.Example{
        Excerpt:   "Users spend 3+ hours/week manually reconciling invoices, costing $50k/year",
        Reasoning: "Quantifies impact, identifies users, and is actionable",
    },
    Fail: &rubric.Example{
        Excerpt:   "We need to improve the system",
        Reasoning: "Vague, no measurable impact, not actionable",
    },
})
```

## RubricSet

Group rubric categories for a specific review type:

```go
type RubricSet struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Description string     `json:"description"`
    Categories  []Category `json:"categories"`
}
```

### Creating a RubricSet

```go
rubricSet := rubric.NewRubricSet("prd-review", "PRD Review", "1.0.0").
    WithDescription("Evaluates Product Requirements Documents").
    AddCategory(problemDefinitionCategory).
    AddCategory(userStoriesCategory).
    AddCategory(successMetricsCategory).
    AddCategory(acceptanceCriteriaCategory)
```

## Authoring Rubrics in YAML (v0.10.0)

Rubric definitions carry `yaml` tags mirroring their `json` tags, so a rubric
can be authored as a YAML file and parsed directly into a `RubricSet`:

```yaml
id: prd-rubric
name: PRD Review
version: "1.0"
passCriteria:
  minCategoriesPassing: all_required
  maxFindingsSeverity: {critical: 0, high: 0, medium: -1, low: -1}
categories:
  - id: problem_definition
    name: Problem Definition
    weight: 0.2
    required: true
    scale:
      type: categorical
      options:
        - {value: pass, criteria: ["Problem is clear, measurable, and tied to users"]}
        - {value: partial, criteria: ["Problem is stated but lacks specificity"]}
        - {value: fail, criteria: ["Problem is vague or missing"]}
```

```go
var rs rubric.RubricSet
_ = yaml.Unmarshal(data, &rs)
```

### Definition Schema

The generated `RubricSet` definition schema is embedded for downstream tooling:

```go
import "github.com/plexusone/structured-evaluation/schema"

data := schema.RubricSetSchemaJSON // rubricset.schema.json
```

This is the **definition** schema (how a rubric is authored); `rubric.schema.json`
remains the **report** schema (an evaluation result).

## Using Rubrics with Reports

```go
// Create report with rubric reference
report := rubric.NewRubric("prd-review", "requirements.md")
report.RubricID = "prd-review-v1"

// Load a rubric definition for evaluation guidance
var rubricSet rubric.RubricSet
_ = yaml.Unmarshal(prdRubricYAML, &rubricSet)

// Evaluate each category using rubric criteria
for _, cat := range rubricSet.Categories {
    result := evaluateCategory(document, cat)
    report.AddCategoryResult(result)
}
```

## Rubric-Guided LLM Evaluation

When using LLM-as-Judge, include each scale option's criteria in the prompt:

```go
func buildPrompt(document string, cat rubric.Category) string {
    var b strings.Builder
    fmt.Fprintf(&b, "Evaluate the following document for %s.\n\nCriteria:\n", cat.Name)
    for _, opt := range cat.Scale.Options {
        fmt.Fprintf(&b, "- %s: %s\n", strings.ToUpper(opt.Value), strings.Join(opt.Criteria, "; "))
    }
    fmt.Fprintf(&b, "\nDocument:\n%s\n\nRespond with: score (pass/partial/fail) and reasoning.", document)
    return b.String()
}
```

For a rich category, iterate `cat.Criteria` and render each criterion's
`Pass.Description` and `Pass.Indicators` instead.

## Benefits

1. **Consistency** - Same criteria across evaluators
2. **Reproducibility** - Track which rubric version was used
3. **Transparency** - Clear expectations for authors
4. **Calibration** - Examples help align understanding

## Best Practices

### Writing Good Criteria

- Be specific and observable
- Use measurable language when possible
- Avoid subjective terms like "good" or "well-written"

```go
// ✅ Good criteria
WithPassPartialFail(
    []string{"All user stories follow Given/When/Then format with acceptance criteria"},
    []string{"Most stories follow the format; some lack acceptance criteria"},
    []string{"Stories are missing or lack testable acceptance criteria"},
)

// ❌ Vague criteria
WithPassPartialFail([]string{"User stories are good"}, nil, []string{"User stories are bad"})
```

### Providing Examples

- Include both passing and failing examples
- Explain why each example scores as it does
- Use realistic content from your domain

### Versioning

Track rubric versions for reproducibility:

```go
rubricSet := rubric.NewRubricSet("prd-review-v2", "PRD Review v2", "2.0.0")
report.RubricID = "prd-review-v2"
```

## Report Validation

Validate rubric reports for correctness before processing (v0.7.0):

```go
result := rubric.ValidateReport(&report)

if !result.Valid {
    for _, issue := range result.Issues {
        fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Path, issue.Message)
    }
}
```

Validation checks include:

- **Enum values** - Score, severity, and decision status must be valid
- **Required fields** - `metadata.document` and `reviewType` are required
- **Finding titles** - Each finding must have a title
- **Count accuracy** - Reported counts must match actual data
- **Decision consistency** - Decision should align with blocking findings

Use the CLI for quick validation:

```bash
sevaluation lint report.json --strict
```

## Extensions (v0.9.0)

Store domain-specific metadata without modifying the core schema:

```go
report := rubric.NewRubric("dss-spec", "material-v3")

// Set custom extension data
report.SetExtension("coverage", coverageReport)
report.SetExtension("metrics", metricsData)
report.SetExtension("customField", "value")

// Check and retrieve
if report.HasExtension("coverage") {
    data := report.GetExtension("coverage")
}
```

### Coverage Report

A built-in extension type for tracking spec coverage:

```go
// Create coverage report
cr := rubric.NewCoverageReport()
cr.SetSection("components", 10, 8, []string{"card", "dialog"})  // 80%
cr.SetSection("foundations", 4, 4, nil)                          // 100%
cr.SetSection("patterns", 5, 3, []string{"form", "wizard"})     // 60%

// Compute overall (simple average)
cr.ComputeOverall()  // 80%

// Or weighted average
weights := map[string]float64{
    "components":  2.0,  // More important
    "foundations": 1.0,
    "patterns":    1.0,
}
cr.ComputeOverallWeighted(weights)

// Store in rubric
report.SetCoverage(cr)

// Retrieve later (type-safe)
coverage := report.GetCoverage()
coverage.Overall                        // 80
coverage.GetSection("components").Total // 10
```

### Coverage Methods

```go
cr := rubric.NewCoverageReport()
cr.SetSection("a", 10, 10, nil)  // 100%
cr.SetSection("b", 10, 5, nil)   // 50%

// Check thresholds
cr.MeetsThreshold(80)           // false (overall < 80)
cr.AllComplete()                // false (not all 100%)

// Filter sections
above := cr.SectionsAboveThreshold(80)  // ["a"]
below := cr.SectionsBelowThreshold(80)  // ["b"]
```

## Next Steps

- [Multi-Judge Aggregation](multi-judge.md) - Combine evaluations
- [Pairwise Comparison](pairwise.md) - Compare outputs
