// Package cli provides command-line interface functionality for the vulnerable target application.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/happyhackingspace/vulnerable-target/pkg/templates"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all templates in the templates directory",
	Long:  "Validate all templates by checking their index.yaml files for proper structure and content using the built-in validator",
	Run: func(_ *cobra.Command, _ []string) {
		validateAllTemplates()
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

// validateAllTemplates validates all templates in the templates directory
func validateAllTemplates() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	templatesDir := filepath.Join(wd, "templates")

	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		fmt.Printf("Templates directory not found: %s\n", templatesDir)
		os.Exit(1)
	}

	// Read all template directories
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		fmt.Printf("Error reading templates directory: %v\n", err)
		os.Exit(1)
	}

	var validationErrors []string
	var validatedCount int

	fmt.Println("Validating templates...")
	fmt.Println(strings.Repeat("=", 50))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templateName := entry.Name()
		templatePath := filepath.Join(templatesDir, templateName)
		indexPath := filepath.Join(templatePath, "index.yaml")

		fmt.Printf("Validating template: %s\n", templateName)

		// Check if index.yaml exists
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			validationErrors = append(validationErrors, fmt.Sprintf("Template '%s': index.yaml file not found", templateName))
			fmt.Printf("  ❌ index.yaml file not found\n")
			continue
		}

		// Load and validate the template
		template, err := templates.LoadTemplate(templatePath)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Template '%s': %v", templateName, err))
			fmt.Printf("  ❌ Validation failed: %v\n", err)
			continue
		}

		// Check if template ID matches directory name
		if template.ID != templateName {
			validationErrors = append(validationErrors, fmt.Sprintf("Template '%s': ID '%s' does not match directory name", templateName, template.ID))
			fmt.Printf("  ❌ ID mismatch: expected '%s', got '%s'\n", templateName, template.ID)
			continue
		}

		validatedCount++
		fmt.Printf("  ✅ Validation passed\n")
	}

	fmt.Println(strings.Repeat("=", 50))

	if len(validationErrors) > 0 {
		fmt.Printf("\n❌ Validation failed with %d error(s):\n\n", len(validationErrors))
		for i, err := range validationErrors {
			fmt.Printf("%d. %s\n", i+1, err)
		}
		fmt.Printf("\nTotal: %d templates validated, %d failed\n", len(entries), len(validationErrors))
		os.Exit(1)
	}

	fmt.Printf("\n✅ All templates validated successfully!\n")
	fmt.Printf("Total: %d templates validated\n", validatedCount)
}
