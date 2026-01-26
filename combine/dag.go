// Package combine provides functionality for combining multiple evaluations
// into a single report, with DAG-based ordering.
package combine

import (
	"sort"

	"github.com/agentplexus/structured-evaluation/summary"
)

// SortByDAG sorts teams in topological order based on their dependencies.
// Uses Kahn's algorithm. Teams at the same level are sorted alphabetically
// for deterministic output.
func SortByDAG(teams []summary.TeamSection) []summary.TeamSection {
	if len(teams) == 0 {
		return teams
	}

	// Build maps for lookup
	idToTeam := make(map[string]*summary.TeamSection)
	for i := range teams {
		idToTeam[teams[i].ID] = &teams[i]
	}

	// Calculate in-degree (number of dependencies) for each team
	inDegree := make(map[string]int)
	for _, t := range teams {
		if _, exists := inDegree[t.ID]; !exists {
			inDegree[t.ID] = 0
		}
		for _, dep := range t.DependsOn {
			// Only count dependencies that exist in our team set
			if _, exists := idToTeam[dep]; exists {
				inDegree[t.ID]++
			}
		}
	}

	// Build downstream adjacency list (who depends on whom)
	downstream := make(map[string][]string)
	for _, t := range teams {
		for _, dep := range t.DependsOn {
			if _, exists := idToTeam[dep]; exists {
				downstream[dep] = append(downstream[dep], t.ID)
			}
		}
	}

	// Find all teams with no dependencies (in-degree 0)
	var queue []string
	for _, t := range teams {
		if inDegree[t.ID] == 0 {
			queue = append(queue, t.ID)
		}
	}

	// Process in topological order
	var sorted []summary.TeamSection
	for len(queue) > 0 {
		// Sort current level alphabetically for determinism
		sort.Strings(queue)

		// Process first item
		current := queue[0]
		queue = queue[1:]

		if team, exists := idToTeam[current]; exists {
			sorted = append(sorted, *team)
		}

		// Reduce in-degree of downstream teams
		for _, next := range downstream[current] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	// Handle any remaining teams (cycles or orphans)
	seen := make(map[string]bool)
	for _, t := range sorted {
		seen[t.ID] = true
	}
	for _, t := range teams {
		if !seen[t.ID] {
			sorted = append(sorted, t)
		}
	}

	return sorted
}

// SortReportByDAG sorts the teams in a summary report by DAG order.
func SortReportByDAG(report *summary.SummaryReport) {
	report.Teams = SortByDAG(report.Teams)
}
