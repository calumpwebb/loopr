package cmd

import (
	"fmt"
	"os"

	"github.com/calumpwebb/loopr/internal/git"
	"github.com/calumpwebb/loopr/internal/loop"
	"github.com/calumpwebb/loopr/internal/sandbox"
	"github.com/calumpwebb/loopr/internal/ui"
)

func Plan() {
	runLoop("plan")
}

func Build() {
	runLoop("build")
}

func runLoop(mode string) {
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

	// Prompt for iterations
	maxIterations := ui.PromptIterations()

	// Run loop
	controller := loop.NewController(sb, mode, maxIterations)
	if err := controller.Run(); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("\n✗ ERROR: Loop failed: %v", err)))
		os.Exit(1)
	}
}
