package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/calumpwebb/loopr/internal/config"
	"github.com/calumpwebb/loopr/internal/sandbox"
	"github.com/calumpwebb/loopr/internal/ui"
	"github.com/calumpwebb/loopr/templates"
)

// Init initializes a new loopr workflow in the current directory
func Init() error {
	// Check if .loopr exists
	if config.LooprDirExists() {
		// Show overwrite/cancel prompt
		if !ui.PromptOverwrite() {
			fmt.Println("Cancelled.")
			return nil
		}
		fmt.Println(ui.WarningStyle.Render("⚠  This will overwrite template files!"))
	}

	// Prompt for sandbox type
	fmt.Println()
	fmt.Println("Welcome to Loopr!")
	fmt.Println("Ralph Loop orchestration for Claude")
	fmt.Println()

	sandboxType := ui.PromptSandbox()

	// Validate Docker available (fail fast)
	fmt.Println()
	fmt.Println("Checking Docker...")
	sb := sandbox.NewDocker()
	if !sb.IsAvailable() {
		fmt.Println(ui.ErrorStyle.Render("✗ Docker is not available"))
		fmt.Println()
		fmt.Println("Install Docker Desktop:")
		fmt.Println("  https://www.docker.com/products/docker-desktop")
		fmt.Println()
		fmt.Println("Then run: loopr init")
		return fmt.Errorf("docker is not available")
	}
	fmt.Println(ui.SuccessStyle.Render("✓ Docker is available"))

	// Create .loopr/ directory
	looprDir := ".loopr"
	if err := os.MkdirAll(looprDir, 0755); err != nil {
		return fmt.Errorf("failed to create .loopr directory: %w", err)
	}

	// Extract templates using templates.ExtractTo()
	fmt.Println()
	fmt.Println("Creating files...")
	if err := templates.ExtractTo(looprDir); err != nil {
		return fmt.Errorf("failed to extract templates: %w", err)
	}

	// List created files with checkmarks
	files := []string{
		"PROMPT_build.md",
		"PROMPT_plan.md",
		"AGENTS.md",
		"config.json",
		"specs/README.md",
		"specs/example-spec.md",
	}

	for _, file := range files {
		fmt.Println(ui.SuccessStyle.Render("✓") + " .loopr/" + file)
	}

	// Create CLAUDE.md at root by reading template and writing to root
	claudeTemplatePath := filepath.Join(looprDir, "CLAUDE.md.template")
	claudeContent, err := os.ReadFile(claudeTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read CLAUDE.md.template: %w", err)
	}

	claudePath := "CLAUDE.md"
	if err := os.WriteFile(claudePath, claudeContent, 0644); err != nil {
		return fmt.Errorf("failed to write CLAUDE.md: %w", err)
	}
	fmt.Println(ui.SuccessStyle.Render("✓") + " CLAUDE.md")

	// Display success message
	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("✓ All set! Next: loopr plan"))
	fmt.Println()

	// Note: sandboxType is prompted but currently only docker is supported
	// In the future, we could write this to config.json if needed
	_ = sandboxType

	return nil
}
