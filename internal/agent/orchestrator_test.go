package agent

import (
	"strings"
	"testing"
)

func TestBuildRecursiveQuery(t *testing.T) {
	tests := []struct {
		name         string
		report       string
		wantEmpty    bool
		wantContains []string
	}{
		{
			name:      "both none stops the loop",
			report:    "SUMMARY: all good.\nCONTRADICTIONS: none\nMISSING: none",
			wantEmpty: true,
		},
		{
			name:      "none is case-insensitive",
			report:    "CONTRADICTIONS: None\nMISSING: NONE",
			wantEmpty: true,
		},
		{
			name:      "too-short findings are ignored",
			report:    "CONTRADICTIONS: no\nMISSING: x",
			wantEmpty: true,
		},
		{
			name:      "no structured markers at all",
			report:    "Just a free-form summary with no tail markers.",
			wantEmpty: true,
		},
		{
			name:         "contradiction only",
			report:       "CONTRADICTIONS: sources disagree on the release date\nMISSING: none",
			wantContains: []string{"Resolve this contradiction: sources disagree on the release date"},
		},
		{
			name:         "missing only",
			report:       "CONTRADICTIONS: none\nMISSING: no independent benchmark numbers",
			wantContains: []string{"Investigate these missing points: no independent benchmark numbers"},
		},
		{
			name:   "both present are combined",
			report: "CONTRADICTIONS: claim A vs claim B\nMISSING: third-party data",
			wantContains: []string{
				"Resolve this contradiction: claim A vs claim B",
				"Investigate these missing points: third-party data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRecursiveQuery(tt.report)

			if tt.wantEmpty {
				if got != "" {
					t.Fatalf("expected empty query, got %q", got)
				}
				return
			}
			if got == "" {
				t.Fatalf("expected a non-empty query, got empty")
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("query %q does not contain %q", got, want)
				}
			}
		})
	}
}

// Regression: the greedy original regex let the CONTRADICTIONS capture swallow
// the MISSING section. The contradiction clause must stop before "MISSING:".
func TestBuildRecursiveQuery_NoBleedIntoMissing(t *testing.T) {
	report := "CONTRADICTIONS: claim A vs claim B\nMISSING: third-party data"
	got := buildRecursiveQuery(report)
	if strings.Contains(got, "Resolve this contradiction: claim A vs claim B\nMISSING") {
		t.Fatalf("contradiction capture bled into the MISSING section: %q", got)
	}
}
