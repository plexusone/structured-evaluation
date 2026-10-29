package rubric

import "time"

// JudgeMetadata tracks information about the judge that produced an
// evaluation — a deterministic check (Code-as-a-Judge), an LLM
// (LLM-as-a-Judge), or a human reviewer (Human-as-a-Judge) — so all three
// are recorded with the same structured identity rather than only the LLM
// case being first-class. This enables reproducibility, debugging, and
// comparison of different judge configurations.
type JudgeMetadata struct {
	// JudgeID is a unique identifier for this judge configuration.
	JudgeID string `json:"judge_id,omitempty"`

	// Method classifies how this judge evaluates: EvalMethodDeterministic
	// (mechanical check), EvalMethodSemantic (LLM judgment), or
	// EvalMethodHuman (human decision). Empty is treated as
	// EvalMethodSemantic — JudgeMetadata's original, sole meaning before
	// this field existed — so a value created before this field existed
	// keeps meaning exactly what it always meant.
	Method EvaluationMethod `json:"method,omitempty"`

	// Model is the LLM model used (e.g., "claude-3-opus-20240229",
	// "gpt-4-turbo"). Applies when Method is EvalMethodSemantic (or empty).
	Model string `json:"model,omitempty"`

	// ModelProvider is the API provider (e.g., "anthropic", "openai", "bedrock").
	ModelProvider string `json:"model_provider,omitempty"`

	// ModelVersion is the specific model version if applicable.
	ModelVersion string `json:"model_version,omitempty"`

	// PromptTemplate is the name/ID of the prompt template used.
	PromptTemplate string `json:"prompt_template,omitempty"`

	// PromptVersion is the version of the prompt template.
	PromptVersion string `json:"prompt_version,omitempty"`

	// SystemPrompt is the system prompt used (or hash/reference if too long).
	SystemPrompt string `json:"system_prompt,omitempty"`

	// Temperature is the sampling temperature used.
	Temperature float64 `json:"temperature,omitempty"`

	// MaxTokens is the max tokens setting.
	MaxTokens int `json:"max_tokens,omitempty"`

	// CheckID identifies the specific deterministic check or tool that
	// produced this evaluation (e.g. "golangci-lint", "schema-validate",
	// "sas-validate-security-profile"). Applies when Method is
	// EvalMethodDeterministic.
	CheckID string `json:"check_id,omitempty"`

	// CheckVersion is the version of the deterministic check or tool.
	CheckVersion string `json:"check_version,omitempty"`

	// ReviewerID identifies the human reviewer who produced this
	// evaluation (e.g. a username or employee ID). Applies when Method is
	// EvalMethodHuman.
	ReviewerID string `json:"reviewer_id,omitempty"`

	// ReviewerName is the human reviewer's display name.
	ReviewerName string `json:"reviewer_name,omitempty"`

	// RubricID references the rubric set used for scoring.
	RubricID string `json:"rubric_id,omitempty"`

	// RubricVersion is the version of the rubric used.
	RubricVersion string `json:"rubric_version,omitempty"`

	// EvaluatedAt is when this evaluation was performed.
	EvaluatedAt time.Time `json:"evaluated_at,omitempty"`

	// Latency is the evaluation duration.
	Latency time.Duration `json:"latency,omitempty"`

	// TokensUsed tracks token consumption.
	TokensUsed *TokenUsage `json:"tokens_used,omitempty"`

	// TraceID links to observability trace (e.g., for Opik/Phoenix/Langfuse).
	TraceID string `json:"trace_id,omitempty"`

	// SpanID links to observability span.
	SpanID string `json:"span_id,omitempty"`
}

// TokenUsage tracks token consumption for an evaluation.
type TokenUsage struct {
	// InputTokens is the number of input/prompt tokens.
	InputTokens int `json:"input_tokens"`

	// OutputTokens is the number of output/completion tokens.
	OutputTokens int `json:"output_tokens"`

	// TotalTokens is the total tokens used.
	TotalTokens int `json:"total_tokens"`

	// CacheReadTokens is tokens read from cache (if applicable).
	CacheReadTokens int `json:"cache_read_tokens,omitempty"`

	// CacheWriteTokens is tokens written to cache (if applicable).
	CacheWriteTokens int `json:"cache_write_tokens,omitempty"`
}

// NewJudgeMetadata creates judge metadata for an LLM judge
// (LLM-as-a-Judge) with required fields.
func NewJudgeMetadata(model string) *JudgeMetadata {
	return &JudgeMetadata{
		Method:      EvalMethodSemantic,
		Model:       model,
		EvaluatedAt: time.Now().UTC(),
	}
}

// NewDeterministicJudgeMetadata creates judge metadata for a mechanical,
// code-based check (Code-as-a-Judge) — no LLM involved. checkID identifies
// the specific check or tool, e.g. "golangci-lint" or "schema-validate".
func NewDeterministicJudgeMetadata(checkID string) *JudgeMetadata {
	return &JudgeMetadata{
		Method:      EvalMethodDeterministic,
		CheckID:     checkID,
		EvaluatedAt: time.Now().UTC(),
	}
}

// NewHumanJudgeMetadata creates judge metadata for a human reviewer
// (Human-as-a-Judge). reviewerID identifies the reviewer, e.g. a username
// or employee ID.
func NewHumanJudgeMetadata(reviewerID string) *JudgeMetadata {
	return &JudgeMetadata{
		Method:      EvalMethodHuman,
		ReviewerID:  reviewerID,
		EvaluatedAt: time.Now().UTC(),
	}
}

// WithProvider sets the model provider.
func (j *JudgeMetadata) WithProvider(provider string) *JudgeMetadata {
	j.ModelProvider = provider
	return j
}

// WithPrompt sets the prompt template info.
func (j *JudgeMetadata) WithPrompt(template, version string) *JudgeMetadata {
	j.PromptTemplate = template
	j.PromptVersion = version
	return j
}

// WithRubric sets the rubric reference.
func (j *JudgeMetadata) WithRubric(id, version string) *JudgeMetadata {
	j.RubricID = id
	j.RubricVersion = version
	return j
}

// WithTemperature sets the sampling temperature.
func (j *JudgeMetadata) WithTemperature(temp float64) *JudgeMetadata {
	j.Temperature = temp
	return j
}

// WithTokenUsage sets the token usage.
func (j *JudgeMetadata) WithTokenUsage(input, output int) *JudgeMetadata {
	j.TokensUsed = &TokenUsage{
		InputTokens:  input,
		OutputTokens: output,
		TotalTokens:  input + output,
	}
	return j
}

// WithTrace links to observability.
func (j *JudgeMetadata) WithTrace(traceID, spanID string) *JudgeMetadata {
	j.TraceID = traceID
	j.SpanID = spanID
	return j
}

// SetLatency records the evaluation duration.
func (j *JudgeMetadata) SetLatency(d time.Duration) {
	j.Latency = d
}

// WithCheckVersion sets the deterministic check or tool's version.
func (j *JudgeMetadata) WithCheckVersion(version string) *JudgeMetadata {
	j.CheckVersion = version
	return j
}

// WithReviewerName sets the human reviewer's display name.
func (j *JudgeMetadata) WithReviewerName(name string) *JudgeMetadata {
	j.ReviewerName = name
	return j
}

// Validate checks that JudgeMetadata carries the fields its Method
// requires, returning a human-readable issue for each problem found. An
// empty Method is treated as EvalMethodSemantic, matching JudgeMetadata's
// original, sole meaning.
func (j *JudgeMetadata) Validate() []string {
	var issues []string

	method := j.Method
	if method == "" {
		method = EvalMethodSemantic
	}

	switch method {
	case EvalMethodSemantic:
		if j.Model == "" {
			issues = append(issues, "semantic judge requires Model")
		}
	case EvalMethodDeterministic:
		if j.CheckID == "" {
			issues = append(issues, "deterministic judge requires CheckID")
		}
	case EvalMethodHuman:
		if j.ReviewerID == "" {
			issues = append(issues, "human judge requires ReviewerID")
		}
	default:
		issues = append(issues, "unknown Method: "+string(method))
	}

	return issues
}
