package summary

import "time"

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
// It aggregates results from multiple teams/agents.
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
