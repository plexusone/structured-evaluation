package rubric

// CoverageSection represents coverage metrics for a single section/area.
// This is a generic type that can be used for:
// - Spec coverage (components, foundations, patterns)
// - Code coverage (functions, lines, branches)
// - Test coverage (scenarios, edge cases)
// - Documentation coverage (API docs, guides)
type CoverageSection struct {
	// Total is the total number of items in this section.
	Total int `json:"total"`

	// Complete is the number of items that are complete/covered.
	Complete int `json:"complete"`

	// Percentage is the coverage percentage (0-100).
	Percentage int `json:"percentage"`

	// Missing lists the IDs or names of missing/incomplete items.
	Missing []string `json:"missing,omitempty"`
}

// CoverageReport aggregates coverage across multiple sections.
type CoverageReport struct {
	// Sections contains coverage for each named section.
	// Example keys: "components", "foundations", "functions", "lines"
	Sections map[string]CoverageSection `json:"sections"`

	// Overall is the aggregate coverage percentage (0-100).
	Overall int `json:"overall"`
}

// NewCoverageReport creates an empty coverage report.
func NewCoverageReport() *CoverageReport {
	return &CoverageReport{
		Sections: make(map[string]CoverageSection),
	}
}

// AddSection adds a coverage section to the report.
func (cr *CoverageReport) AddSection(name string, section CoverageSection) {
	cr.Sections[name] = section
}

// SetSection is an alias for AddSection for fluent API.
func (cr *CoverageReport) SetSection(name string, total, complete int, missing []string) *CoverageReport {
	percentage := 0
	if total > 0 {
		percentage = (complete * 100) / total
	}
	cr.Sections[name] = CoverageSection{
		Total:      total,
		Complete:   complete,
		Percentage: percentage,
		Missing:    missing,
	}
	return cr
}

// ComputeOverall calculates the overall coverage percentage.
// Uses a simple average of all section percentages.
func (cr *CoverageReport) ComputeOverall() int {
	if len(cr.Sections) == 0 {
		return 0
	}

	total := 0
	for _, section := range cr.Sections {
		total += section.Percentage
	}
	cr.Overall = total / len(cr.Sections)
	return cr.Overall
}

// ComputeOverallWeighted calculates overall coverage with weights.
// The weights map should have keys matching section names.
// Sections not in the weights map get weight 1.0.
func (cr *CoverageReport) ComputeOverallWeighted(weights map[string]float64) int {
	if len(cr.Sections) == 0 {
		return 0
	}

	totalWeight := 0.0
	weightedSum := 0.0

	for name, section := range cr.Sections {
		weight := 1.0
		if w, ok := weights[name]; ok {
			weight = w
		}
		totalWeight += weight
		weightedSum += float64(section.Percentage) * weight
	}

	if totalWeight == 0 {
		return 0
	}
	cr.Overall = int(weightedSum / totalWeight)
	return cr.Overall
}

// GetSection retrieves a section by name.
// Returns an empty section if not found.
func (cr *CoverageReport) GetSection(name string) CoverageSection {
	if section, ok := cr.Sections[name]; ok {
		return section
	}
	return CoverageSection{}
}

// HasSection checks if a section exists.
func (cr *CoverageReport) HasSection(name string) bool {
	_, ok := cr.Sections[name]
	return ok
}

// AllComplete returns true if all sections have 100% coverage.
func (cr *CoverageReport) AllComplete() bool {
	for _, section := range cr.Sections {
		if section.Percentage < 100 {
			return false
		}
	}
	return true
}

// MeetsThreshold returns true if overall coverage meets the threshold.
func (cr *CoverageReport) MeetsThreshold(threshold int) bool {
	return cr.Overall >= threshold
}

// SectionsAboveThreshold returns section names with coverage >= threshold.
func (cr *CoverageReport) SectionsAboveThreshold(threshold int) []string {
	var result []string
	for name, section := range cr.Sections {
		if section.Percentage >= threshold {
			result = append(result, name)
		}
	}
	return result
}

// SectionsBelowThreshold returns section names with coverage < threshold.
func (cr *CoverageReport) SectionsBelowThreshold(threshold int) []string {
	var result []string
	for name, section := range cr.Sections {
		if section.Percentage < threshold {
			result = append(result, name)
		}
	}
	return result
}

// ExtensionKeyCoverage is the standard extension key for coverage data.
const ExtensionKeyCoverage = "coverage"

// SetCoverage is a convenience method to set coverage on a Rubric.
func (r *Rubric) SetCoverage(coverage *CoverageReport) {
	r.SetExtension(ExtensionKeyCoverage, coverage)
}

// GetCoverage retrieves coverage data from a Rubric's extensions.
// Returns nil if not set or if type assertion fails.
func (r *Rubric) GetCoverage() *CoverageReport {
	ext := r.GetExtension(ExtensionKeyCoverage)
	if ext == nil {
		return nil
	}

	// Handle both direct type and map[string]any (from JSON unmarshal)
	if cr, ok := ext.(*CoverageReport); ok {
		return cr
	}

	// Try to convert from map[string]any (common after JSON round-trip)
	if m, ok := ext.(map[string]any); ok {
		return coverageReportFromMap(m)
	}

	return nil
}

// coverageReportFromMap converts a map to CoverageReport.
// Used when coverage data has been JSON marshaled/unmarshaled.
func coverageReportFromMap(m map[string]any) *CoverageReport {
	cr := NewCoverageReport()

	if overall, ok := m["overall"].(float64); ok {
		cr.Overall = int(overall)
	}

	if sections, ok := m["sections"].(map[string]any); ok {
		for name, sectionData := range sections {
			if sectionMap, ok := sectionData.(map[string]any); ok {
				section := CoverageSection{}
				if total, ok := sectionMap["total"].(float64); ok {
					section.Total = int(total)
				}
				if complete, ok := sectionMap["complete"].(float64); ok {
					section.Complete = int(complete)
				}
				if percentage, ok := sectionMap["percentage"].(float64); ok {
					section.Percentage = int(percentage)
				}
				if missing, ok := sectionMap["missing"].([]any); ok {
					for _, item := range missing {
						if s, ok := item.(string); ok {
							section.Missing = append(section.Missing, s)
						}
					}
				}
				cr.Sections[name] = section
			}
		}
	}

	return cr
}
