# Rubrics

Rubrics define explicit evaluation criteria for consistent assessments across evaluators (human or LLM).

## Overview

A rubric provides structured guidance for evaluating each category:

```go
type Rubric struct {
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

## Creating Rubrics

```go
rubric := evaluation.NewRubric("problem_definition", "Problem Definition").
    WithDescription("Evaluates clarity and completeness of the problem statement").
    WithPassCriteria("Problem is clearly stated with measurable business impact and affected users identified").
    WithPartialCriteria("Problem is stated but lacks specificity or measurable impact").
    WithFailCriteria("Problem is vague, missing, or not actionable")
```

## Adding Examples

Examples help evaluators understand the criteria:

```go
rubric.AddExample(evaluation.Example{
    Score:   evaluation.ScorePass,
    Text:    "Users spend 3+ hours/week manually reconciling invoices, costing $50k/year in labor",
    Reason:  "Quantifies impact, identifies users, and is actionable",
})

rubric.AddExample(evaluation.Example{
    Score:   evaluation.ScoreFail,
    Text:    "We need to improve the system",
    Reason:  "Vague, no measurable impact, not actionable",
})
```

## RubricSet

Group rubrics for a specific review type:

```go
type RubricSet struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Categories  []Rubric `json:"categories"`
}
```

### Creating a RubricSet

```go
rubricSet := evaluation.NewRubricSet("prd-review", "PRD Review").
    WithDescription("Evaluates Product Requirements Documents").
    AddRubric(problemDefinitionRubric).
    AddRubric(userStoriesRubric).
    AddRubric(successMetricsRubric).
    AddRubric(acceptanceCriteriaRubric)
```

## Default PRD RubricSet

```go
rubricSet := evaluation.DefaultPRDRubricSet()
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
report := evaluation.NewEvaluationReport("prd-review", "requirements.md")
report.RubricID = "prd-review-v1"

// Load rubric for evaluation guidance
rubricSet := evaluation.DefaultPRDRubricSet()

// Evaluate each category using rubric criteria
for _, rubric := range rubricSet.Categories {
    result := evaluateCategory(document, rubric)
    report.AddCategory(result)
}
```

## Rubric-Guided LLM Evaluation

When using LLM-as-Judge, include rubric criteria in the prompt:

```go
func buildPrompt(document string, rubric Rubric) string {
    return fmt.Sprintf(`Evaluate the following document for %s.

Criteria:
- PASS: %s
- PARTIAL: %s
- FAIL: %s

Document:
%s

Respond with: score (pass/partial/fail) and reasoning.`,
        rubric.Name,
        rubric.Criteria.Pass,
        rubric.Criteria.Partial,
        rubric.Criteria.Fail,
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
rubricSet := evaluation.NewRubricSet("prd-review-v2", "PRD Review v2")
report.RubricID = "prd-review-v2"
```

## Next Steps

- [Multi-Judge Aggregation](multi-judge.md) - Combine evaluations
- [Pairwise Comparison](pairwise.md) - Compare outputs
