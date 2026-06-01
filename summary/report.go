package summary

import (
	"encoding/json"
	"time"
)

// TeamSection represents results from a single agent or validation area.
type TeamSection struct {
	// ID is the unique identifier (workflow step ID).
	ID string `json:"id"`

	// Name is the human-readable name.
	Name string `json:"name"`

	// AgentID is the agent that produced this section.
	AgentID string `json:"agent_id,omitempty"`

	// Model is the LLM model used (if applicable).
	Model string `json:"model,omitempty"`

	// DependsOn lists upstream team IDs (for DAG ordering).
	DependsOn []string `json:"depends_on,omitempty"`

	// Tasks are the individual check results.
	Tasks []TaskResult `json:"tasks"`

	// Status is the computed overall status for this section.
	Status Status `json:"status"`
}

// ComputeStatus calculates the status from tasks.
func (t *TeamSection) ComputeStatus() Status {
	t.Status = ComputeStatusFromTasks(t.Tasks)
	return t.Status
}

// SummaryReport is the top-level report for summary-style evaluations.
// It aggregates results from multiple teams/agents and can embed full
// EvaluationReport and ClaimsReport for complete fidelity.
type SummaryReport struct {
	// Schema is the JSON Schema URL for validation.
	Schema string `json:"$schema,omitempty"`

	// Project identifies the project being evaluated.
	Project string `json:"project"`

	// Version is the version being evaluated.
	Version string `json:"version,omitempty"`

	// Target is a human-readable target description.
	Target string `json:"target,omitempty"`

	// Phase describes the evaluation phase (e.g., "RELEASE VALIDATION").
	Phase string `json:"phase,omitempty"`

	// Teams are the individual team/agent sections.
	Teams []TeamSection `json:"teams"`

	// Status is the computed overall status.
	Status Status `json:"status"`

	// GeneratedAt is when the report was created.
	GeneratedAt time.Time `json:"generated_at"`

	// GeneratedBy identifies what created this report.
	GeneratedBy string `json:"generated_by,omitempty"`

	// EmbeddedReports contains full-fidelity embedded reports.
	// This allows the SummaryReport to serve as a container for
	// detailed reports while providing a summary view.
	EmbeddedReports *EmbeddedReports `json:"embeddedReports,omitempty"`
}

// EmbeddedReports contains full-fidelity embedded reports.
// Reports are stored as json.RawMessage to avoid circular imports
// and allow flexible report types.
type EmbeddedReports struct {
	// Evaluations contains embedded EvaluationReport(s).
	// Key is a report identifier (e.g., "prd-review", "article-quality").
	Evaluations map[string]json.RawMessage `json:"evaluations,omitempty"`

	// Claims contains embedded ClaimsReport(s).
	// Key is a report identifier (e.g., "source-validation").
	Claims map[string]json.RawMessage `json:"claims,omitempty"`

	// Custom contains any other embedded reports.
	// Key is a report identifier, value is the full report JSON.
	Custom map[string]json.RawMessage `json:"custom,omitempty"`
}

// ComputeOverallStatus calculates the overall status from all teams.
func (r *SummaryReport) ComputeOverallStatus() Status {
	statuses := make([]Status, len(r.Teams))
	for i, t := range r.Teams {
		statuses[i] = t.Status
	}
	r.Status = ComputeStatus(statuses)
	return r.Status
}

// IsGo returns true if the overall status allows proceeding.
func (r *SummaryReport) IsGo() bool {
	for _, t := range r.Teams {
		if t.Status == StatusNoGo {
			return false
		}
	}
	return true
}

// FinalMessage returns a formatted final status message.
func (r *SummaryReport) FinalMessage() string {
	if r.IsGo() {
		if r.Version != "" {
			return "🚀 GO for " + r.Version + " 🚀"
		}
		return "🚀 GO 🚀"
	}
	return "❌ NO-GO ❌"
}

// NewSummaryReport creates a new summary report with defaults.
func NewSummaryReport(project, version, phase string) *SummaryReport {
	return &SummaryReport{
		Project:     project,
		Version:     version,
		Phase:       phase,
		Teams:       []TeamSection{},
		GeneratedAt: time.Now().UTC(),
	}
}

// AddTeam adds a team section to the report.
func (r *SummaryReport) AddTeam(team TeamSection) {
	team.ComputeStatus()
	r.Teams = append(r.Teams, team)
}

// EnsureEmbeddedReports initializes the EmbeddedReports field if nil.
func (r *SummaryReport) EnsureEmbeddedReports() {
	if r.EmbeddedReports == nil {
		r.EmbeddedReports = &EmbeddedReports{
			Evaluations: make(map[string]json.RawMessage),
			Claims:      make(map[string]json.RawMessage),
			Custom:      make(map[string]json.RawMessage),
		}
	}
	if r.EmbeddedReports.Evaluations == nil {
		r.EmbeddedReports.Evaluations = make(map[string]json.RawMessage)
	}
	if r.EmbeddedReports.Claims == nil {
		r.EmbeddedReports.Claims = make(map[string]json.RawMessage)
	}
	if r.EmbeddedReports.Custom == nil {
		r.EmbeddedReports.Custom = make(map[string]json.RawMessage)
	}
}

// EmbedEvaluationReport embeds an EvaluationReport with the given key.
// The report is marshaled to JSON for storage.
func (r *SummaryReport) EmbedEvaluationReport(key string, report any) error {
	r.EnsureEmbeddedReports()
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	r.EmbeddedReports.Evaluations[key] = data
	return nil
}

// EmbedClaimsReport embeds a ClaimsReport with the given key.
// The report is marshaled to JSON for storage.
func (r *SummaryReport) EmbedClaimsReport(key string, report any) error {
	r.EnsureEmbeddedReports()
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	r.EmbeddedReports.Claims[key] = data
	return nil
}

// EmbedCustomReport embeds a custom report with the given key.
// The report is marshaled to JSON for storage.
func (r *SummaryReport) EmbedCustomReport(key string, report any) error {
	r.EnsureEmbeddedReports()
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	r.EmbeddedReports.Custom[key] = data
	return nil
}

// GetEmbeddedEvaluation retrieves and unmarshals an embedded EvaluationReport.
// Returns nil if not found. The target should be a pointer to the report struct.
func (r *SummaryReport) GetEmbeddedEvaluation(key string, target any) error {
	if r.EmbeddedReports == nil || r.EmbeddedReports.Evaluations == nil {
		return nil
	}
	data, ok := r.EmbeddedReports.Evaluations[key]
	if !ok {
		return nil
	}
	return json.Unmarshal(data, target)
}

// GetEmbeddedClaims retrieves and unmarshals an embedded ClaimsReport.
// Returns nil if not found. The target should be a pointer to the report struct.
func (r *SummaryReport) GetEmbeddedClaims(key string, target any) error {
	if r.EmbeddedReports == nil || r.EmbeddedReports.Claims == nil {
		return nil
	}
	data, ok := r.EmbeddedReports.Claims[key]
	if !ok {
		return nil
	}
	return json.Unmarshal(data, target)
}

// GetEmbeddedCustom retrieves and unmarshals an embedded custom report.
// Returns nil if not found. The target should be a pointer to the report struct.
func (r *SummaryReport) GetEmbeddedCustom(key string, target any) error {
	if r.EmbeddedReports == nil || r.EmbeddedReports.Custom == nil {
		return nil
	}
	data, ok := r.EmbeddedReports.Custom[key]
	if !ok {
		return nil
	}
	return json.Unmarshal(data, target)
}

// HasEmbeddedReports returns true if any embedded reports exist.
func (r *SummaryReport) HasEmbeddedReports() bool {
	if r.EmbeddedReports == nil {
		return false
	}
	return len(r.EmbeddedReports.Evaluations) > 0 ||
		len(r.EmbeddedReports.Claims) > 0 ||
		len(r.EmbeddedReports.Custom) > 0
}
