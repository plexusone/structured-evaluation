package schema

import _ "embed"

// RubricSchemaJSON is the embedded JSON Schema for Rubric.
//
//go:embed rubric.schema.json
var RubricSchemaJSON []byte

// SummarySchemaJSON is the embedded JSON Schema for SummaryReport.
//
//go:embed summary.schema.json
var SummarySchemaJSON []byte

// ClaimsSchemaJSON is the embedded JSON Schema for ClaimsReport.
//
//go:embed claims.schema.json
var ClaimsSchemaJSON []byte

// RubricSetSchemaJSON is the embedded JSON Schema for RubricSet — the rubric
// definition, including the rich weighted-criteria form.
//
//go:embed rubricset.schema.json
var RubricSetSchemaJSON []byte
