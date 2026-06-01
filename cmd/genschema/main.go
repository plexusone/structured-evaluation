// Command genschema generates JSON schemas for all report types.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/plexusone/structured-evaluation/schema"
)

func main() {
	schemaDir := "schema"
	if len(os.Args) > 1 {
		schemaDir = os.Args[1]
	}
	// Clean the path to prevent path traversal (gosec G703)
	schemaDir = filepath.Clean(schemaDir)

	// Generate rubric schema
	evalData, err := schema.GenerateRubricSchema()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating rubric schema: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(schemaDir, "rubric.schema.json"), evalData, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing rubric schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated rubric.schema.json")

	// Generate summary schema
	summaryData, err := schema.GenerateSummarySchema()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating summary schema: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(schemaDir, "summary.schema.json"), summaryData, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing summary schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated summary.schema.json")

	// Generate claims schema
	claimsData, err := schema.GenerateClaimsSchema()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating claims schema: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(schemaDir, "claims.schema.json"), claimsData, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing claims schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated claims.schema.json")

	fmt.Println("All schemas generated successfully")
}
