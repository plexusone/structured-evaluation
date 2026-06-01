// Package schema provides JSON Schema generation for structured evaluation types.
package schema

import (
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"

	"github.com/plexusone/structured-evaluation/claims"
	"github.com/plexusone/structured-evaluation/rubric"
	"github.com/plexusone/structured-evaluation/summary"
)

// GenerateRubricSchema generates JSON Schema for Rubric.
func GenerateRubricSchema() ([]byte, error) {
	reflector := &jsonschema.Reflector{
		DoNotReference:             true,
		ExpandedStruct:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := reflector.Reflect(&rubric.Rubric{})
	schema.ID = "https://github.com/plexusone/structured-evaluation/schema/rubric.schema.json"
	schema.Title = "Rubric Report"
	schema.Description = "Schema for rubric-based LLM-as-Judge evaluation reports"

	return json.MarshalIndent(schema, "", "  ")
}

// GenerateSummarySchema generates JSON Schema for SummaryReport.
func GenerateSummarySchema() ([]byte, error) {
	reflector := &jsonschema.Reflector{
		DoNotReference:             true,
		ExpandedStruct:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := reflector.Reflect(&summary.SummaryReport{})
	schema.ID = "https://github.com/plexusone/structured-evaluation/schema/summary.schema.json"
	schema.Title = "Summary Report"
	schema.Description = "Schema for GO/NO-GO summary reports from deterministic checks"

	return json.MarshalIndent(schema, "", "  ")
}

// GenerateClaimsSchema generates JSON Schema for ClaimsReport.
func GenerateClaimsSchema() ([]byte, error) {
	reflector := &jsonschema.Reflector{
		DoNotReference:             true,
		ExpandedStruct:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := reflector.Reflect(&claims.ClaimsReport{})
	schema.ID = "https://github.com/plexusone/structured-evaluation/schema/claims.schema.json"
	schema.Title = "Claims Report"
	schema.Description = "Schema for claim extraction and source validation reports"

	return json.MarshalIndent(schema, "", "  ")
}

// WriteSchemaFile writes schema bytes to a file.
func WriteSchemaFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0600)
}
