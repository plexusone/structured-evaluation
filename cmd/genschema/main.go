// Command genschema generates JSON schemas for all report types.
package main

import (
	"fmt"
	"os"

	"github.com/plexusone/structured-evaluation/schema"
)

func main() {
	schemaDir := "schema"
	if len(os.Args) > 1 {
		schemaDir = os.Args[1]
	}

	if err := run(schemaDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All schemas generated successfully")
}

func run(schemaDir string) error {
	// os.Root confines every write to schemaDir, so a schema file name can
	// never escape it (gosec G703).
	root, err := os.OpenRoot(schemaDir)
	if err != nil {
		return fmt.Errorf("open schema dir %q: %w", schemaDir, err)
	}
	// Closing a read-only directory handle loses no data, so its error is
	// not actionable.
	defer func() { _ = root.Close() }()

	generators := []struct {
		file     string
		generate func() ([]byte, error)
	}{
		{"rubric.schema.json", schema.GenerateRubricSchema},
		{"rubricset.schema.json", schema.GenerateRubricSetSchema},
		{"summary.schema.json", schema.GenerateSummarySchema},
		{"claims.schema.json", schema.GenerateClaimsSchema},
	}
	for _, g := range generators {
		data, err := g.generate()
		if err != nil {
			return fmt.Errorf("generate %s: %w", g.file, err)
		}
		if err := writeFile(root, g.file, data); err != nil {
			return fmt.Errorf("write %s: %w", g.file, err)
		}
		fmt.Printf("Generated %s\n", g.file)
	}
	return nil
}

func writeFile(root *os.Root, name string, data []byte) error {
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close() // the write error is the one worth reporting
		return err
	}
	return f.Close()
}
