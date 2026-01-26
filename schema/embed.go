package schema

import _ "embed"

// EvaluationSchemaJSON is the embedded JSON Schema for EvaluationReport.
//
//go:embed evaluation.schema.json
var EvaluationSchemaJSON []byte

// SummarySchemaJSON is the embedded JSON Schema for SummaryReport.
//
//go:embed summary.schema.json
var SummarySchemaJSON []byte
