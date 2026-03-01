package combine

import (
	"time"

	"github.com/plexusone/structured-evaluation/summary"
)

// AgentResult represents the output from a single agent execution.
// This is the intermediate format agents produce before aggregation.
type AgentResult struct {
	// Schema is the JSON Schema URL.
	Schema string `json:"$schema,omitempty"`

	// AgentID identifies the agent.
	AgentID string `json:"agent_id"`

	// StepID is the workflow step this agent executed.
	StepID string `json:"step_id"`

	// Inputs provided to the agent.
	Inputs map[string]any `json:"inputs,omitempty"`

	// Outputs produced by the agent.
	Outputs map[string]any `json:"outputs,omitempty"`

	// Tasks are the individual task results.
	Tasks []summary.TaskResult `json:"tasks"`

	// Status is the overall status.
	Status summary.Status `json:"status"`

	// ExecutedAt is when the agent ran.
	ExecutedAt time.Time `json:"executed_at"`

	// AgentModel is the LLM model used.
	AgentModel string `json:"agent_model,omitempty"`

	// Duration is the execution time.
	Duration string `json:"duration,omitempty"`

	// Error contains any error message.
	Error string `json:"error,omitempty"`
}

// ToTeamSection converts an AgentResult to a TeamSection.
func (a *AgentResult) ToTeamSection() summary.TeamSection {
	return summary.TeamSection{
		ID:      a.StepID,
		Name:    a.AgentID,
		AgentID: a.AgentID,
		Model:   a.AgentModel,
		Tasks:   a.Tasks,
		Status:  a.Status,
	}
}

// ComputeStatus calculates the status from tasks.
func (a *AgentResult) ComputeStatus() summary.Status {
	a.Status = summary.ComputeStatusFromTasks(a.Tasks)
	return a.Status
}

// AggregateResults combines multiple agent results into a summary report.
func AggregateResults(results []AgentResult, project, version, phase string) *summary.SummaryReport {
	report := summary.NewSummaryReport(project, version, phase)
	report.GeneratedBy = "structured-evaluation"

	for i := range results {
		results[i].ComputeStatus()
		section := results[i].ToTeamSection()
		report.Teams = append(report.Teams, section)
	}

	report.ComputeOverallStatus()
	return report
}

// AggregateWithDAG combines results and sorts by DAG order.
func AggregateWithDAG(results []AgentResult, dag []TeamDependency, project, version, phase string) *summary.SummaryReport {
	report := AggregateResults(results, project, version, phase)

	// Apply DAG dependencies
	depMap := make(map[string][]string)
	for _, d := range dag {
		depMap[d.ID] = d.DependsOn
	}

	for i := range report.Teams {
		if deps, ok := depMap[report.Teams[i].ID]; ok {
			report.Teams[i].DependsOn = deps
		}
	}

	SortReportByDAG(report)
	return report
}

// TeamDependency defines dependencies for a team in the DAG.
type TeamDependency struct {
	ID        string   `json:"id"`
	DependsOn []string `json:"depends_on"`
}
