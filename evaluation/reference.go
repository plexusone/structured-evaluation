package evaluation

// ReferenceData contains ground truth or expected data for evaluation.
// This enables reference-based evaluation where outputs are compared
// against known-good examples.
type ReferenceData struct {
	// ID is the unique identifier for this reference.
	ID string `json:"id,omitempty"`

	// Input is the input/prompt that produced the reference output.
	Input string `json:"input,omitempty"`

	// ExpectedOutput is the gold/reference output.
	ExpectedOutput string `json:"expected_output,omitempty"`

	// ExpectedOutputs allows multiple acceptable outputs.
	ExpectedOutputs []string `json:"expected_outputs,omitempty"`

	// Context provides additional context (e.g., retrieved documents for RAG).
	Context []string `json:"context,omitempty"`

	// Annotations are human-provided labels or scores.
	Annotations []Annotation `json:"annotations,omitempty"`

	// Source indicates where this reference came from.
	Source string `json:"source,omitempty"`

	// Tags categorize or filter references.
	Tags []string `json:"tags,omitempty"`

	// Metadata contains additional reference data.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Annotation represents a human-provided label or score.
type Annotation struct {
	// Name is the annotation type (e.g., "quality", "relevance").
	Name string `json:"name"`

	// Score is a numeric score (if applicable).
	Score float64 `json:"score,omitempty"`

	// Label is a categorical label (if applicable).
	Label string `json:"label,omitempty"`

	// Explanation provides reasoning for the annotation.
	Explanation string `json:"explanation,omitempty"`

	// AnnotatorID identifies who provided this annotation.
	AnnotatorID string `json:"annotator_id,omitempty"`

	// AnnotatorType indicates human vs automated (e.g., "human", "llm", "rule").
	AnnotatorType string `json:"annotator_type,omitempty"`
}

// ReferenceDataset is a collection of reference data items.
type ReferenceDataset struct {
	// ID is the unique identifier for this dataset.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Description explains what this dataset contains.
	Description string `json:"description,omitempty"`

	// Version tracks dataset iterations.
	Version string `json:"version,omitempty"`

	// Items are the reference data items.
	Items []ReferenceData `json:"items"`

	// Tags categorize the dataset.
	Tags []string `json:"tags,omitempty"`

	// Metadata contains additional dataset info.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// NewReferenceData creates a new reference data item.
func NewReferenceData(input, expectedOutput string) *ReferenceData {
	return &ReferenceData{
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

// WithContext adds context documents.
func (r *ReferenceData) WithContext(ctx ...string) *ReferenceData {
	r.Context = append(r.Context, ctx...)
	return r
}

// WithAnnotation adds a human annotation.
func (r *ReferenceData) WithAnnotation(name string, score float64, annotatorID string) *ReferenceData {
	r.Annotations = append(r.Annotations, Annotation{
		Name:          name,
		Score:         score,
		AnnotatorID:   annotatorID,
		AnnotatorType: "human",
	})
	return r
}

// NewReferenceDataset creates a new reference dataset.
func NewReferenceDataset(id, name string) *ReferenceDataset {
	return &ReferenceDataset{
		ID:    id,
		Name:  name,
		Items: []ReferenceData{},
	}
}

// AddItem adds a reference data item to the dataset.
func (d *ReferenceDataset) AddItem(item ReferenceData) {
	d.Items = append(d.Items, item)
}

// GetByID retrieves a reference item by ID.
func (d *ReferenceDataset) GetByID(id string) *ReferenceData {
	for i := range d.Items {
		if d.Items[i].ID == id {
			return &d.Items[i]
		}
	}
	return nil
}
