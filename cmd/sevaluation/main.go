// Command sevaluation provides CLI tools for working with evaluation reports.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/plexusone/structured-evaluation/claims"
	"github.com/plexusone/structured-evaluation/render/box"
	"github.com/plexusone/structured-evaluation/render/detailed"
	htmlrender "github.com/plexusone/structured-evaluation/render/html"
	"github.com/plexusone/structured-evaluation/render/markdown"
	"github.com/plexusone/structured-evaluation/render/terminal"
	"github.com/plexusone/structured-evaluation/rubric"
	"github.com/plexusone/structured-evaluation/schema"
	"github.com/plexusone/structured-evaluation/summary"
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
	Long: `Render an evaluation, summary, or claims report.

Formats for evaluation reports:
  box      - ASCII box format for TUI (deterministic, no colors/emojis)
  detailed - Detailed format with findings
  terminal - ANSI-colored terminal output with UTF8 icons
  markdown - Markdown report format
  json     - Pretty-printed JSON

Formats for summary reports:
  box      - ASCII box format
  json     - Pretty-printed JSON

Formats for claims reports:
  html     - Self-contained HTML page, claims grouped by verdict
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

// Lint command flags
var lintFlags struct {
	strict bool
	format string
}

// Lint command
var lintCmd = &cobra.Command{
	Use:   "lint <file.json>",
	Short: "Lint an evaluation or claims report for correctness",
	Long: `Lint an evaluation or claims report.

Evaluation reports ('categories'):
  - All enum values are valid (score, severity, decision status)
  - Required fields are present
  - Reported counts match actual data
  - Decision is consistent with findings and category results

Claims reports ('claims') — evidence integrity of verified claims:
  - A verified external claim has a resolving URL and a verbatim quote (error)
  - A verified derived claim lists source claims; internal claim has evidence (error)
  - The statistical value appears in the quoted text (warning — targets/ranges
    and unit-scaled values legitimately quote a rule rather than the number)

Exit codes:
  0 - Valid (no errors, warnings allowed unless --strict)
  1 - Invalid (has errors or has warnings with --strict)`,
	Args: cobra.ExactArgs(1),
	RunE: runLint,
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
  rubric.schema.json - Schema for detailed LLM-as-Judge reports
  summary.schema.json    - Schema for GO/NO-GO summary reports`,
	RunE: runSchemaGenerate,
}

func runSchemaGenerate(cmd *cobra.Command, args []string) error {
	outputDir := schemaFlags.outputDir
	if outputDir == "" {
		outputDir = "."
	}

	// Generate rubric schema
	evalSchema, err := schema.GenerateRubricSchema()
	if err != nil {
		return fmt.Errorf("generating rubric schema: %w", err)
	}
	evalPath := outputDir + "/rubric.schema.json"
	if err := schema.WriteSchemaFile(evalPath, evalSchema); err != nil {
		return fmt.Errorf("writing rubric schema: %w", err)
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

	// Generate enums JSON for cross-language validation
	enumsJSON, err := schema.GenerateEnumsJSON()
	if err != nil {
		return fmt.Errorf("generating enums JSON: %w", err)
	}
	enumsPath := outputDir + "/enums.json"
	if err := schema.WriteSchemaFile(enumsPath, enumsJSON); err != nil {
		return fmt.Errorf("writing enums JSON: %w", err)
	}
	fmt.Printf("Generated: %s\n", enumsPath)

	return nil
}

func init() {
	renderCmd.Flags().StringVarP(&renderFlags.format, "format", "f", "detailed", "Output format (box, detailed, json)")
	schemaCmd.Flags().StringVarP(&schemaFlags.outputDir, "output", "o", ".", "Output directory for schema files")
	lintCmd.Flags().BoolVar(&lintFlags.strict, "strict", false, "Treat warnings as errors")
	lintCmd.Flags().StringVarP(&lintFlags.format, "format", "f", "text", "Output format (text, json)")

	rootCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(lintCmd)
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

	// Check for claims report markers
	if _, hasClaims := raw["claims"]; hasClaims {
		return renderClaims(data, renderFlags.format)
	}

	return fmt.Errorf("unknown report type: expected 'categories' (evaluation), 'teams' (summary), or 'claims' (claims)")
}

func lintClaims(data []byte, format string, strict bool) error {
	var report claims.ClaimsReport
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing claims report: %w", err)
	}

	findings := claims.Lint(&report)

	switch format {
	case "json":
		output, err := json.MarshalIndent(findings, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	default:
		if len(findings) == 0 {
			fmt.Println("✓ no lint findings")
		}
		for _, f := range findings {
			icon := "⚠"
			if f.Severity == claims.LintError {
				icon = "✗"
			}
			fmt.Printf("%s [%s] %s: %s\n", icon, f.Severity, f.ClaimID, f.Message)
		}
	}

	if claims.HasErrors(findings) {
		os.Exit(1)
	}
	if strict && claims.HasWarnings(findings) {
		fmt.Println("\n(strict mode: treating warnings as errors)")
		os.Exit(1)
	}
	return nil
}

func renderClaims(data []byte, format string) error {
	var report claims.ClaimsReport
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing claims report: %w", err)
	}

	switch format {
	case "html":
		renderer := htmlrender.New(os.Stdout)
		return renderer.RenderClaims(&report)
	case "json":
		output, err := json.MarshalIndent(&report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	default:
		return fmt.Errorf("format %q not supported for claims reports (use html or json)", format)
	}
}

func renderEvaluation(data []byte, format string) error {
	var report rubric.Rubric
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing evaluation report: %w", err)
	}

	switch format {
	case "box", "ascii":
		renderer := box.NewEvaluationRenderer(os.Stdout)
		return renderer.Render(&report)
	case "detailed":
		renderer := detailed.NewTerminal(os.Stdout)
		return renderer.Render(&report)
	case "terminal":
		renderer := terminal.New(os.Stdout)
		return renderer.Render(&report)
	case "markdown", "md":
		renderer := markdown.New(os.Stdout)
		return renderer.Render(&report)
	case "json":
		output, err := json.MarshalIndent(&report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	default:
		return fmt.Errorf("format %q not supported for evaluation reports (use box, detailed, terminal, markdown, or json)", format)
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
		var report rubric.Rubric
		if err := json.Unmarshal(data, &report); err != nil {
			return fmt.Errorf("parsing evaluation report: %w", err)
		}

		if report.Decision.Passed {
			counts := report.Decision.CategoryCounts
			fmt.Printf("✅ PASSED: %s (%d/%d categories)\n", report.ReviewType, counts.Pass, counts.Total)
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

	// Check claims report
	if _, hasClaims := raw["claims"]; hasClaims {
		var report claims.ClaimsReport
		if err := json.Unmarshal(data, &report); err != nil {
			return fmt.Errorf("parsing claims report: %w", err)
		}

		c := report.Summary.Counts
		if report.Decision.Passed {
			fmt.Printf("✅ %s: %s (%d verified, %d needs-review of %d)\n",
				report.Decision.Status, report.Metadata.DocumentTitle, c.Verified, c.NeedsReview, c.Total)
			return nil
		}
		fmt.Printf("❌ %s: %s — %s\n", report.Decision.Status, report.Metadata.DocumentTitle, report.Decision.Rationale)
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
		var report rubric.Rubric
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

func runLint(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Claims reports have their own lint (evidence-integrity of verified claims).
	if _, hasClaims := raw["claims"]; hasClaims {
		return lintClaims(data, lintFlags.format, lintFlags.strict)
	}

	// Only lint evaluation reports otherwise.
	if _, hasCategories := raw["categories"]; !hasCategories {
		return fmt.Errorf("lint supports evaluation reports ('categories') or claims reports ('claims')")
	}

	var report rubric.Rubric
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parsing evaluation report: %w", err)
	}

	// Run validation
	result := rubric.ValidateReport(&report)

	// Output based on format
	switch lintFlags.format {
	case "json":
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling result: %w", err)
		}
		fmt.Println(string(output))
	default:
		fmt.Print(result.String())
	}

	// Exit with error if invalid or strict mode with warnings
	if !result.Valid {
		os.Exit(1)
	}
	if lintFlags.strict && result.HasWarnings() {
		fmt.Println("\n(strict mode: treating warnings as errors)")
		os.Exit(1)
	}

	return nil
}
