# Rubrics

Rubrics define explicit evaluation criteria for consistent assessments across evaluators (human or LLM).

## Overview

A rubric category provides structured guidance for evaluating each category:

```go
type Category struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Criteria    struct {
        Pass    string `json:"pass"`
        Partial string `json:"partial"`
        Fail    string `json:"fail"`
    } `json:"criteria"`
    Examples    []Example `json:"examples,omitempty"`
}
```

## Creating Rubric Categories

```go
cat := rubric.NewCategory("problem_definition", "Problem Definition", "").
    WithDescription("Evaluates clarity and completeness of the problem statement").
    WithPassCriteria("Problem is clearly stated with measurable business impact and affected users identified").
    WithPartialCriteria("Problem is stated but lacks specificity or measurable impact").
    WithFailCriteria("Problem is vague, missing, or not actionable")
```

## Adding Examples

Examples help evaluators understand the criteria:

```go
cat.AddExample(rubric.Example{
    Score:   rubric.ScorePass,
    Text:    "Users spend 3+ hours/week manually reconciling invoices, costing $50k/year in labor",
    Reason:  "Quantifies impact, identifies users, and is actionable",
})

cat.AddExample(rubric.Example{
    Score:   rubric.ScoreFail,
    Text:    "We need to improve the system",
    Reason:  "Vague, no measurable impact, not actionable",
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

## Default PRD RubricSet

```go
rubricSet := rubric.DefaultPRDRubricSet()
```

Includes rubrics for:

- **problem_definition** - Clarity of the problem statement
- **user_stories** - Completeness of user stories
- **success_metrics** - Quantitative success criteria
- **acceptance_criteria** - Testable acceptance criteria
- **scope_definition** - Clear scope boundaries

## Using Rubrics with Reports

```go
// Create report with rubric reference
report := rubric.NewRubric("prd-review", "requirements.md")
report.RubricID = "prd-review-v1"

// Load rubric for evaluation guidance
rubricSet := rubric.DefaultPRDRubricSet()

// Evaluate each category using rubric criteria
for _, cat := range rubricSet.Categories {
    result := evaluateCategory(document, cat)
    report.AddCategoryResult(result)
}
```

## Rubric-Guided LLM Evaluation

When using LLM-as-Judge, include rubric criteria in the prompt:

```go
func buildPrompt(document string, cat Category) string {
    return fmt.Sprintf(`Evaluate the following document for %s.

Criteria:
- PASS: %s
- PARTIAL: %s
- FAIL: %s

Document:
%s

Respond with: score (pass/partial/fail) and reasoning.`,
        cat.Name,
        cat.Criteria.Pass,
        cat.Criteria.Partial,
        cat.Criteria.Fail,
        document,
    )
}
```

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
WithPassCriteria("All user stories follow Given/When/Then format with acceptance criteria")

// ❌ Vague criteria
WithPassCriteria("User stories are good")
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

## Next Steps

- [Multi-Judge Aggregation](multi-judge.md) - Combine evaluations
- [Pairwise Comparison](pairwise.md) - Compare outputs
