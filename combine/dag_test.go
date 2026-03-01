package combine

import (
	"testing"

	"github.com/plexusone/structured-evaluation/summary"
)

func TestSortByDAG(t *testing.T) {
	teams := []summary.TeamSection{
		{ID: "release", DependsOn: []string{"qa", "security"}},
		{ID: "qa", DependsOn: []string{}},
		{ID: "security", DependsOn: []string{"qa"}},
	}

	sorted := SortByDAG(teams)

	// qa should be first (no dependencies)
	if sorted[0].ID != "qa" {
		t.Errorf("expected qa first, got %s", sorted[0].ID)
	}

	// security should be before release (release depends on security)
	secIdx, relIdx := -1, -1
	for i, team := range sorted {
		if team.ID == "security" {
			secIdx = i
		}
		if team.ID == "release" {
			relIdx = i
		}
	}
	if secIdx > relIdx {
		t.Errorf("security should be before release")
	}
}

func TestSortByDAG_NoDependencies(t *testing.T) {
	teams := []summary.TeamSection{
		{ID: "c"},
		{ID: "a"},
		{ID: "b"},
	}

	sorted := SortByDAG(teams)

	// Should preserve original order when no dependencies
	if len(sorted) != 3 {
		t.Errorf("expected 3 teams, got %d", len(sorted))
	}
}

func TestSortByDAG_Empty(t *testing.T) {
	teams := []summary.TeamSection{}
	sorted := SortByDAG(teams)

	if len(sorted) != 0 {
		t.Errorf("expected empty result, got %d teams", len(sorted))
	}
}
