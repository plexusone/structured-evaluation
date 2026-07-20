# Fixtures

`rubric-report.json` is real output from `rubric.Rubric`, not hand-written —
generated via:

```go
r := rubric.NewRubric("prd", "prd.md")
r.AddCategoryResult(rubric.CategoryResult{
	Category:  "clarity",
	Score:     rubric.ScorePass,
	IntScore:  rubric.ScoreExcellent,
	Reasoning: "Clear and specific.",
	Findings: []rubric.Finding{
		{ID: "f1", Category: "clarity", Severity: rubric.SeverityLow,
			Title: "Minor wording", Description: "Could be tighter.",
			Recommendation: "Tighten phrasing."},
	},
})
r.Finalize(nil, "pdlc check")
```

Regenerate after any change to `rubric.Rubric`'s shape so this test keeps
validating against real output, not a stale hand-maintained fixture.
