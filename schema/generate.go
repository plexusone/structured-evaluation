// Package schema provides JSON Schema generation for structured evaluation types.
package schema

import (
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"

	"github.com/plexusone/structured-evaluation/evaluation"
	"github.com/plexusone/structured-evaluation/summary"
)

// GenerateEvaluationSchema generates JSON Schema for EvaluationReport.
func GenerateEvaluationSchema() ([]byte, error) {
	reflector := &jsonschema.Reflector{
		DoNotReference:             true,
		ExpandedStruct:             true,
		RequiredFromJSONSchemaTags: true,
	}

	schema := reflector.Reflect(&evaluation.EvaluationReport{})
	schema.ID = "https://github.com/plexusone/structured-evaluation/schema/evaluation.schema.json"
	schema.Title = "Evaluation Report"
	schema.Description = "Schema for detailed LLM-as-Judge evaluation reports"

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

// WriteSchemaFile writes schema bytes to a file.
func WriteSchemaFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0600)
}
