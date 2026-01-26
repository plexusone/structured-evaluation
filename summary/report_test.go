package summary

import "testing"

func TestSummaryReport_IsGo(t *testing.T) {
	tests := []struct {
		name     string
		statuses []Status
		want     bool
	}{
		{
			name:     "all GO",
			statuses: []Status{StatusGo, StatusGo, StatusGo},
			want:     true,
		},
		{
			name:     "one WARN",
			statuses: []Status{StatusGo, StatusWarn, StatusGo},
			want:     true,
		},
		{
			name:     "one NO-GO",
			statuses: []Status{StatusGo, StatusNoGo, StatusGo},
			want:     false,
		},
		{
			name:     "empty",
			statuses: []Status{},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &SummaryReport{
				Project: "test",
				Teams:   []TeamSection{},
			}

			// Build tasks and add team using AddTeam which computes status
			tasks := make([]TaskResult, len(tt.statuses))
			for i, s := range tt.statuses {
				tasks[i] = TaskResult{
					ID:     string(rune('a' + i)),
					Status: s,
				}
			}
			report.AddTeam(TeamSection{
				ID:    "team1",
				Name:  "Team 1",
				Tasks: tasks,
			})

			if got := report.IsGo(); got != tt.want {
				t.Errorf("IsGo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_Icon(t *testing.T) {
	tests := []struct {
		status Status
		icon   string
	}{
		{StatusGo, "🟢"},
		{StatusWarn, "🟡"},
		{StatusNoGo, "🔴"},
		{StatusSkip, "⚪"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.Icon(); got != tt.icon {
				t.Errorf("Icon() = %q, want %q", got, tt.icon)
			}
		})
	}
}
