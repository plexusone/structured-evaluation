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

// EnumDefinitions contains all enum values for the structured evaluation types.
// These are exported for cross-language tooling that needs to know valid values.
var EnumDefinitions = map[string][]string{
	"ScoreValue":     rubric.ValidScoreValues(),
	"Severity":       rubric.ValidSeverityValues(),
	"DecisionStatus": rubric.ValidDecisionStatusValues(),
	"EvaluationType": rubric.ValidEvaluationTypeValues(),
	"ScaleType":      rubric.ValidScaleTypeValues(),
}

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

	// Add enum constraints to the generated schema
	addEnumConstraints(schema)

	return json.MarshalIndent(schema, "", "  ")
}

// addEnumConstraints walks the schema and adds enum constraints for known types.
func addEnumConstraints(schema *jsonschema.Schema) {
	if schema == nil {
		return
	}

	// Process definitions if present
	if schema.Definitions != nil {
		for _, def := range schema.Definitions {
			addEnumConstraints(def)
		}
	}

	// Process properties using ordered map iteration
	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			key := pair.Key
			prop := pair.Value
			addEnumConstraintByFieldName(key, prop)
			addEnumConstraints(prop)
		}
	}

	// Process items (for arrays)
	if schema.Items != nil {
		addEnumConstraints(schema.Items)
	}

	// Process additional properties
	if schema.AdditionalProperties != nil {
		addEnumConstraints(schema.AdditionalProperties)
	}
}

// addEnumConstraintByFieldName adds enum constraint based on field name.
func addEnumConstraintByFieldName(fieldName string, schema *jsonschema.Schema) {
	if schema == nil || schema.Type != "string" {
		return
	}

	switch fieldName {
	case "score":
		schema.Enum = toInterfaceSlice(rubric.ValidScoreValues())
		schema.Description = "Category score: pass, partial, or fail"
	case "severity":
		schema.Enum = toInterfaceSlice(rubric.ValidSeverityValues())
		schema.Description = "Finding severity: critical, high, medium, low, or info"
	case "status":
		// Only apply to decision status fields
		schema.Enum = toInterfaceSlice(rubric.ValidDecisionStatusValues())
		schema.Description = "Decision status: pass, conditional, fail, or human_review"
	case "overallDecision":
		schema.Enum = toInterfaceSlice(rubric.ValidDecisionStatusValues())
		schema.Description = "Overall decision status"
	case "evaluationType":
		schema.Enum = toInterfaceSlice(rubric.ValidEvaluationTypeValues())
		schema.Description = "Evaluation type: analytic or holistic"
	case "type":
		// Check if this is a scale type field by looking at context
		// This is a heuristic - may need refinement
		if len(schema.Enum) == 0 {
			schema.Enum = toInterfaceSlice(rubric.ValidScaleTypeValues())
			schema.Description = "Scale type: categorical, checklist, binary, or likert"
		}
	}
}

// toInterfaceSlice converts string slice to interface slice for JSON Schema enum.
func toInterfaceSlice(s []string) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
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

// GenerateEnumsJSON generates a JSON file containing all enum definitions.
// This is useful for cross-language tooling that needs to validate enum values.
func GenerateEnumsJSON() ([]byte, error) {
	enums := map[string]any{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"title":       "Structured Evaluation Enum Definitions",
		"description": "Valid enum values for structured evaluation types. Use these for validation in any language.",
		"enums":       EnumDefinitions,
		"validation": map[string]any{
			"score": map[string]any{
				"values":      rubric.ValidScoreValues(),
				"description": "Category score values",
				"usage":       "CategoryResult.score field",
			},
			"severity": map[string]any{
				"values":      rubric.ValidSeverityValues(),
				"description": "Finding severity levels",
				"usage":       "Finding.severity field",
				"blocking":    []string{"critical", "high"},
			},
			"decisionStatus": map[string]any{
				"values":      rubric.ValidDecisionStatusValues(),
				"description": "Decision status values",
				"usage":       "Decision.status and overallDecision fields",
				"passing":     []string{"pass", "conditional"},
			},
			"evaluationType": map[string]any{
				"values":      rubric.ValidEvaluationTypeValues(),
				"description": "Evaluation methodology types",
				"usage":       "RubricSet.evaluationType field",
			},
			"scaleType": map[string]any{
				"values":      rubric.ValidScaleTypeValues(),
				"description": "Scale types for rubric categories",
				"usage":       "Category.scale.type field",
			},
		},
	}

	return json.MarshalIndent(enums, "", "  ")
}
