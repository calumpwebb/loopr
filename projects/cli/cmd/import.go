package cmd

import (
	"fmt"
	"os"

	"github.com/calumpwebb/loopr/internal/git"
	"github.com/calumpwebb/loopr/internal/prompts"
	"github.com/calumpwebb/loopr/internal/sandbox"
	"github.com/calumpwebb/loopr/internal/ui"
)

// Import imports tasks from a PRD/spec file
func Import(sourceFile string) {
	// Check if git repo
	if !git.IsGitRepo() {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: Not a git repository"))
		fmt.Println("\nPlease initialize a git repository first:")
		fmt.Println("  git init")
		os.Exit(1)
	}

	// Check if .loopr directory exists
	if _, err := os.Stat(".loopr"); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: .loopr directory not found"))
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Check if source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ ERROR: Source file not found: %s", sourceFile)))
		os.Exit(1)
	}

	// Check if tasks.md exists
	if _, err := os.Stat(".loopr/tasks.md"); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: .loopr/tasks.md not found"))
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Create sandbox (hardcoded to docker for v1)
	sb, err := sandbox.New("docker")
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ ERROR: Failed to create sandbox: %v", err)))
		os.Exit(1)
	}

	// Check Docker available
	if !sb.IsAvailable() {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: Docker is not available"))
		fmt.Println("\nInstall Docker Desktop:")
		fmt.Println("  https://www.docker.com/products/docker-desktop")
		os.Exit(1)
	}

	// Check auth (quick haiku test)
	if !sb.IsAuthenticated() {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: Docker sandbox not authenticated"))
		fmt.Println()

		// Show auth prompt
		if ui.PromptAuthenticate() {
			fmt.Println("\nAuthenticating...")
			if err := sb.Authenticate(); err != nil {
				fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ ERROR: Authentication failed: %v", err)))
				os.Exit(1)
			}
			fmt.Println(ui.SuccessStyle.Render("✓ Authenticated!"))
		} else {
			fmt.Println("Authentication required. Run again when ready.")
			os.Exit(0)
		}
	}

	fmt.Println()
	fmt.Printf("Importing tasks from: %s\n\n", sourceFile)

	// Get import prompt with source file
	prompt := prompts.GetImportPrompt(sourceFile)

	// Execute Claude with sonnet model
	if err := sb.ExecuteClaude(prompt, "sonnet"); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("\n✗ ERROR: Import failed: %v", err)))
		os.Exit(1)
	}

	// Push to git
	branch, err := git.CurrentBranch()
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("\n✗ ERROR: Failed to get git branch: %v", err)))
		os.Exit(1)
	}

	fmt.Println("\nPushing to git...")
	if err := git.Push(branch); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ ERROR: Failed to push: %v", err)))
		os.Exit(1)
	}
	fmt.Println(ui.SuccessStyle.Render("✓ Pushed to origin/" + branch))

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("✓ Tasks imported successfully!"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  - Review .loopr/tasks.md to see imported tasks")
	fmt.Println("  - Run 'loopr plan' to refine the task list")
	fmt.Println("  - Run 'loopr build' to start implementing")
}
