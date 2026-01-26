// Command sevaluation provides CLI tools for working with evaluation reports.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/agentplexus/structured-evaluation/evaluation"
	"github.com/agentplexus/structured-evaluation/render/box"
	"github.com/agentplexus/structured-evaluation/render/detailed"
	"github.com/agentplexus/structured-evaluation/schema"
	"github.com/agentplexus/structured-evaluation/summary"
)

var version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "sevaluation",
	Short: "Structured evaluation report tools",
	Long: `sevaluation is a CLI tool for working with structured evaluation reports.

It supports two report types:
  - Summary reports: GO/NO-GO status per task (for deterministic checks)
  - Evaluation reports: Detailed LLM-as-Judge reviews with findings

Commands:
  validate   - Validate report structure
  render     - Render report to terminal or markdown
  check      - Check pass/fail status (exit code 0/1)
  combine    - Combine multiple reports with DAG ordering

Examples:
  sevaluation render report.json --format=box
  sevaluation render report.json --format=detailed
  sevaluation check report.json`,
	Version: version,
}

// Render command
var renderFlags struct {
	format string
}

var renderCmd = &cobra.Command{
	Use:   "render <file.json>",
	Short: "Render a report to terminal",
	Long: `Render an evaluation or summary report.

Formats:
  box      - Summary box format (for summary reports)
  detailed - Detailed format with findings (for evaluation reports)
  json     - Pretty-printed JSON`,
	Args: cobra.ExactArgs(1),
	RunE: runRender,
}

// Check command
var checkCmd = &cobra.Command{
	Use:   "check <file.json>",
	Short: "Check if evaluation passes",
	Long: `Check if an evaluation report passes its criteria.

Exit codes:
  0 - Passed
  1 - Failed or blocked`,
	Args: cobra.ExactArgs(1),
	RunE: runCheck,
}

// Validate command
var validateCmd = &cobra.Command{
	Use:   "validate <file.json>",
	Short: "Validate report structure",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

// Schema command
var schemaFlags struct {
	outputDir string
}

var schemaCmd = &cobra.Command{
	Use:   "schema generate",
	Short: "Generate JSON Schema files",
	Long: `Generate JSON Schema files for evaluation and summary reports.

Outputs:
  evaluation.schema.json - Schema for detailed LLM-as-Judge reports
  summary.schema.json    - Schema for GO/NO-GO summary reports`,
	RunE: runSchemaGenerate,
}

func runSchemaGenerate(cmd *cobra.Command, args []string) error {
	outputDir := schemaFlags.outputDir
	if outputDir == "" {
		outputDir = "."
	}

	// Generate evaluation schema
	evalSchema, err := schema.GenerateEvaluationSchema()
	if err != nil {
		return fmt.Errorf("generating evaluation schema: %w", err)
	}
	evalPath := outputDir + "/evaluation.schema.json"
	if err := schema.WriteSchemaFile(evalPath, evalSchema); err != nil {
		return fmt.Errorf("writing evaluation schema: %w", err)
	}
	fmt.Printf("Generated: %s\n", evalPath)

	// Generate summary schema
	summarySchema, err := schema.GenerateSummarySchema()
	if err != nil {
		return fmt.Errorf("generating summary schema: %w", err)
	}
	summaryPath := outputDir + "/summary.schema.json"
	if err := schema.WriteSchemaFile(summaryPath, summarySchema); err != nil {
		return fmt.Errorf("writing summary schema: %w", err)
	}
	fmt.Printf("Generated: %s\n", summaryPath)

	return nil
}

func init() {
	renderCmd.Flags().StringVarP(&renderFlags.format, "format", "f", "detailed", "Output format (box, detailed, json)")
	schemaCmd.Flags().StringVarP(&schemaFlags.outputDir, "output", "o", ".", "Output directory for schema files")

	rootCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(schemaCmd)
}

func runRender(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	// Try to detect report type
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Check for evaluation report markers
	if _, hasCategories := raw["categories"]; hasCategories {
		return renderEvaluation(data, renderFlags.format)
	}

	// Check for summary report markers
	if _, hasTeams := raw["teams"]; hasTeams {
		return renderSummary(data, renderFlags.format)
	}

	return fmt.Errorf("unknown report type: expected 'categories' (evaluation) or 'teams' (summary)")
}

func renderEvaluation(data []byte, format string) error {
	var report evaluation.EvaluationReport
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing evaluation report: %w", err)
	}

	switch format {
	case "detailed", "terminal":
		renderer := detailed.NewTerminal(os.Stdout)
		return renderer.Render(&report)
	case "json":
		output, err := json.MarshalIndent(&report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	default:
		return fmt.Errorf("format %q not supported for evaluation reports (use detailed or json)", format)
	}
}

func renderSummary(data []byte, format string) error {
	var report summary.SummaryReport
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing summary report: %w", err)
	}

	switch format {
	case "box", "summary":
		renderer := box.New(os.Stdout)
		return renderer.Render(&report)
	case "json":
		output, err := json.MarshalIndent(&report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	default:
		return fmt.Errorf("format %q not supported for summary reports (use box or json)", format)
	}
}

func runCheck(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Check evaluation report
	if _, hasCategories := raw["categories"]; hasCategories {
		var report evaluation.EvaluationReport
		if err := json.Unmarshal(data, &report); err != nil {
			return fmt.Errorf("parsing evaluation report: %w", err)
		}

		if report.Decision.Passed {
			fmt.Printf("✅ PASSED: %s (%.1f/10)\n", report.ReviewType, report.WeightedScore)
			return nil
		}
		fmt.Printf("❌ FAILED: %s - %s\n", report.ReviewType, report.Decision.Rationale)
		os.Exit(1)
	}

	// Check summary report
	if _, hasTeams := raw["teams"]; hasTeams {
		var report summary.SummaryReport
		if err := json.Unmarshal(data, &report); err != nil {
			return fmt.Errorf("parsing summary report: %w", err)
		}

		if report.IsGo() {
			fmt.Printf("🟢 GO: %s %s\n", report.Project, report.Version)
			return nil
		}
		fmt.Printf("🔴 NO-GO: %s %s\n", report.Project, report.Version)
		os.Exit(1)
	}

	return fmt.Errorf("unknown report type")
}

func runValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Try parsing as evaluation
	if _, hasCategories := raw["categories"]; hasCategories {
		var report evaluation.EvaluationReport
		if err := json.Unmarshal(data, &report); err != nil {
			return fmt.Errorf("invalid evaluation report: %w", err)
		}
		fmt.Printf("Valid evaluation report: %s (%s)\n", report.Metadata.Document, report.ReviewType)
		return nil
	}

	// Try parsing as summary
	if _, hasTeams := raw["teams"]; hasTeams {
		var report summary.SummaryReport
		if err := json.Unmarshal(data, &report); err != nil {
			return fmt.Errorf("invalid summary report: %w", err)
		}
		fmt.Printf("Valid summary report: %s %s\n", report.Project, report.Version)
		return nil
	}

	return fmt.Errorf("unknown report type: expected 'categories' or 'teams'")
}
