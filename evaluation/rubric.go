package evaluation

import "encoding/json"

// EvaluationType defines how evaluation is performed.
type EvaluationType string

const (
	// EvaluationTypeAnalytic scores each category independently (recommended for LLM-as-Judge).
	EvaluationTypeAnalytic EvaluationType = "analytic"

	// EvaluationTypeHolistic provides a single overall score.
	EvaluationTypeHolistic EvaluationType = "holistic"
)

// ScaleType defines the type of scoring scale.
type ScaleType string

const (
	// ScaleTypeCategorical uses discrete categories (pass/partial/fail).
	// Recommended for LLM-as-Judge - better calibrated than numeric scales.
	ScaleTypeCategorical ScaleType = "categorical"

	// ScaleTypeChecklist uses a list of required/optional items.
	ScaleTypeChecklist ScaleType = "checklist"

	// ScaleTypeBinary is simple pass/fail.
	ScaleTypeBinary ScaleType = "binary"
)

// RubricSet is a collection of rubrics for a complete evaluation.
// Follows Go-first principles: Go types are source of truth, JSON Schema generated from them.
type RubricSet struct {
	// ID uniquely identifies this rubric set.
	ID string `json:"id"`

	// Name is the human-readable name.
	Name string `json:"name"`

	// Version is the semantic version of this rubric.
	Version string `json:"version"`

	// Description explains what this rubric set evaluates.
	Description string `json:"description,omitempty"`

	// EvaluationType is "analytic" (per-category) or "holistic" (single score).
	// Analytic is recommended for LLM-as-Judge.
	EvaluationType EvaluationType `json:"evaluationType,omitempty"`

	// PassCriteria defines requirements for overall pass/fail.
	PassCriteria RubricPassCriteria `json:"passCriteria"`

	// Categories are the evaluation dimensions.
	Categories []Category `json:"categories"`

	// JudgePromptTemplate is the prompt template for LLM evaluation.
	// Supports placeholders: {content}, {categories}, etc.
	JudgePromptTemplate string `json:"judgePromptTemplate,omitempty"`

	// Metadata contains additional information about the rubric.
	Metadata *RubricMetadata `json:"metadata,omitempty"`
}

// RubricPassCriteria defines requirements for overall pass/fail determination.
type RubricPassCriteria struct {
	// MinCategoriesPassing is "all", "all_required", or a number.
	MinCategoriesPassing string `json:"minCategoriesPassing,omitempty"`

	// MaxFindings limits findings by severity.
	MaxFindings *FindingLimits `json:"maxFindingsSeverity,omitempty"`
}

// FindingLimits sets maximum allowed findings per severity.
// Use -1 for unlimited.
type FindingLimits struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low,omitempty"`
}

// RubricMetadata contains additional rubric information.
type RubricMetadata struct {
	CreatedAt string   `json:"createdAt,omitempty"`
	Author    string   `json:"author,omitempty"`
	BasedOn   []string `json:"basedOn,omitempty"`
}

// Category is a single evaluation dimension.
type Category struct {
	// ID uniquely identifies this category within the rubric.
	ID string `json:"id"`

	// Name is the human-readable category name.
	Name string `json:"name"`

	// Description explains what this category measures.
	Description string `json:"description"`

	// Weight is the relative importance (default 1.0).
	Weight float64 `json:"weight,omitempty"`

	// Required indicates if this category must pass for overall pass.
	Required bool `json:"required,omitempty"`

	// Scale defines how this category is scored.
	Scale Scale `json:"scale"`

	// EvaluationPrompt is a specific prompt for evaluating this category.
	EvaluationPrompt string `json:"evaluationPrompt,omitempty"`

	// Examples provides few-shot examples for LLM evaluation.
	// Research shows 1 example per level improves LLM alignment.
	Examples *CategoryExamples `json:"examples,omitempty"`
}

// Scale defines the scoring mechanism for a category.
type Scale struct {
	// Type is "categorical", "checklist", or "binary".
	// Categorical with 2-3 options is recommended for LLM-as-Judge.
	Type ScaleType `json:"type"`

	// Options are the scoring options (for categorical scales).
	Options []ScaleOption `json:"options,omitempty"`

	// RequiredItems are items that must be present (for checklist scales).
	RequiredItems []string `json:"requiredItems,omitempty"`

	// OptionalItems are items that add value (for checklist scales).
	OptionalItems []string `json:"optionalItems,omitempty"`

	// PassingThreshold defines pass criteria (for checklist scales).
	PassingThreshold *ChecklistThreshold `json:"passingThreshold,omitempty"`
}

// ScaleOption is a single option in a categorical scale.
type ScaleOption struct {
	// Value is the machine-readable value (e.g., "pass", "partial", "fail").
	Value string `json:"value"`

	// Label is the human-readable label.
	Label string `json:"label"`

	// Criteria are specific requirements for this score level.
	Criteria []string `json:"criteria"`
}

// ChecklistThreshold defines pass criteria for checklist scales.
type ChecklistThreshold struct {
	// Required is "all" or a number of required items that must be present.
	Required string `json:"required,omitempty"`

	// Optional is the minimum number of optional items needed.
	Optional int `json:"optional,omitempty"`
}

// CategoryExamples provides few-shot examples for a category.
// Research shows 1 example per level improves LLM alignment.
type CategoryExamples struct {
	Pass    *Example `json:"pass,omitempty"`
	Partial *Example `json:"partial,omitempty"`
	Fail    *Example `json:"fail,omitempty"`
}

// Example is a few-shot example for LLM evaluation.
type Example struct {
	// Excerpt is example content from a document.
	Excerpt string `json:"excerpt"`

	// Reasoning explains why this gets this score.
	// Including reasoning improves LLM alignment (chain-of-thought).
	Reasoning string `json:"reasoning"`
}

// NewRubricSet creates a new rubric set with required fields.
func NewRubricSet(id, name, version string) *RubricSet {
	return &RubricSet{
		ID:             id,
		Name:           name,
		Version:        version,
		EvaluationType: EvaluationTypeAnalytic,
		Categories:     []Category{},
		PassCriteria: RubricPassCriteria{
			MinCategoriesPassing: "all_required",
			MaxFindings: &FindingLimits{
				Critical: 0,
				High:     0,
				Medium:   -1, // Unlimited by default
			},
		},
	}
}

// AddCategory adds a category to the rubric set.
func (rs *RubricSet) AddCategory(cat Category) *RubricSet {
	rs.Categories = append(rs.Categories, cat)
	return rs
}

// SetPassCriteria sets the pass criteria.
func (rs *RubricSet) SetPassCriteria(criteria RubricPassCriteria) *RubricSet {
	rs.PassCriteria = criteria
	return rs
}

// SetJudgePrompt sets the judge prompt template.
func (rs *RubricSet) SetJudgePrompt(template string) *RubricSet {
	rs.JudgePromptTemplate = template
	return rs
}

// SetMetadata sets the rubric metadata.
func (rs *RubricSet) SetMetadata(meta *RubricMetadata) *RubricSet {
	rs.Metadata = meta
	return rs
}

// ToJSON serializes a rubric set to JSON.
func (rs *RubricSet) ToJSON() ([]byte, error) {
	return json.MarshalIndent(rs, "", "  ")
}

// Validate checks the rubric for common issues.
func (rs *RubricSet) Validate() []string {
	var issues []string

	if rs.ID == "" {
		issues = append(issues, "rubric ID is required")
	}
	if rs.Name == "" {
		issues = append(issues, "rubric name is required")
	}
	if rs.Version == "" {
		issues = append(issues, "rubric version is required")
	}
	if len(rs.Categories) == 0 {
		issues = append(issues, "at least one category is required")
	}

	for i, cat := range rs.Categories {
		if cat.ID == "" {
			issues = append(issues, "category "+itoa(i)+": ID is required")
		}
		if cat.Name == "" {
			issues = append(issues, "category "+cat.ID+": name is required")
		}
		if len(cat.Scale.Options) == 0 && cat.Scale.Type == ScaleTypeCategorical {
			issues = append(issues, "category "+cat.ID+": categorical scale requires options")
		}
	}

	return issues
}

// GetCategory returns a category by ID, or nil if not found.
func (rs *RubricSet) GetCategory(id string) *Category {
	for i := range rs.Categories {
		if rs.Categories[i].ID == id {
			return &rs.Categories[i]
		}
	}
	return nil
}

// GetRequiredCategories returns all required categories.
func (rs *RubricSet) GetRequiredCategories() []Category {
	var required []Category
	for _, cat := range rs.Categories {
		if cat.Required {
			required = append(required, cat)
		}
	}
	return required
}

// NewCategory creates a new category with a categorical scale.
func NewCategory(id, name, description string) *Category {
	return &Category{
		ID:          id,
		Name:        name,
		Description: description,
		Weight:      1.0,
		Scale: Scale{
			Type:    ScaleTypeCategorical,
			Options: []ScaleOption{},
		},
	}
}

// SetRequired marks this category as required for pass.
func (c *Category) SetRequired(required bool) *Category {
	c.Required = required
	return c
}

// SetWeight sets the category weight.
func (c *Category) SetWeight(weight float64) *Category {
	c.Weight = weight
	return c
}

// SetEvaluationPrompt sets the evaluation prompt for this category.
func (c *Category) SetEvaluationPrompt(prompt string) *Category {
	c.EvaluationPrompt = prompt
	return c
}

// AddOption adds a scale option to a categorical category.
func (c *Category) AddOption(value, label string, criteria ...string) *Category {
	c.Scale.Options = append(c.Scale.Options, ScaleOption{
		Value:    value,
		Label:    label,
		Criteria: criteria,
	})
	return c
}

// WithPassPartialFail sets up a standard pass/partial/fail scale.
func (c *Category) WithPassPartialFail(passCriteria, partialCriteria, failCriteria []string) *Category {
	c.Scale.Type = ScaleTypeCategorical
	c.Scale.Options = []ScaleOption{
		{Value: "pass", Label: "Pass", Criteria: passCriteria},
		{Value: "partial", Label: "Partial", Criteria: partialCriteria},
		{Value: "fail", Label: "Fail", Criteria: failCriteria},
	}
	return c
}

// WithBinary sets up a binary pass/fail scale.
func (c *Category) WithBinary(passCriteria, failCriteria []string) *Category {
	c.Scale.Type = ScaleTypeBinary
	c.Scale.Options = []ScaleOption{
		{Value: "pass", Label: "Pass", Criteria: passCriteria},
		{Value: "fail", Label: "Fail", Criteria: failCriteria},
	}
	return c
}

// WithChecklist sets up a checklist scale.
func (c *Category) WithChecklist(required, optional []string, threshold *ChecklistThreshold) *Category {
	c.Scale.Type = ScaleTypeChecklist
	c.Scale.RequiredItems = required
	c.Scale.OptionalItems = optional
	c.Scale.PassingThreshold = threshold
	return c
}

// SetExamples sets few-shot examples for the category.
func (c *Category) SetExamples(examples *CategoryExamples) *Category {
	c.Examples = examples
	return c
}

// GetOptionForValue returns the scale option for a given value.
func (c *Category) GetOptionForValue(value string) *ScaleOption {
	for i := range c.Scale.Options {
		if c.Scale.Options[i].Value == value {
			return &c.Scale.Options[i]
		}
	}
	return nil
}
