package evaluation

import "time"

// JudgeMetadata tracks information about the LLM judge that produced an evaluation.
// This enables reproducibility, debugging, and comparison of different judge configurations.
type JudgeMetadata struct {
	// JudgeID is a unique identifier for this judge configuration.
	JudgeID string `json:"judge_id,omitempty"`

	// Model is the LLM model used (e.g., "claude-3-opus-20240229", "gpt-4-turbo").
	Model string `json:"model"`

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

// NewJudgeMetadata creates judge metadata with required fields.
func NewJudgeMetadata(model string) *JudgeMetadata {
	return &JudgeMetadata{
		Model:       model,
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
