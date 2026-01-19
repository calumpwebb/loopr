package cmd

import (
	"fmt"
	"os"

	"github.com/calumpwebb/loopr/internal/config"
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
		fmt.Println("✗ Not a git repository")
		fmt.Println("\nPlease initialize a git repository first:")
		fmt.Println("  git init")
		os.Exit(1)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("✗ Failed to load config: %v\n", err)
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Create sandbox
	sb, err := sandbox.New(cfg.Sandbox)
	if err != nil {
		fmt.Printf("✗ Failed to create sandbox: %v\n", err)
		os.Exit(1)
	}

	// Check Docker available
	if !sb.IsAvailable() {
		fmt.Println("✗ Docker is not available")
		fmt.Println("\nInstall Docker Desktop:")
		fmt.Println("  https://www.docker.com/products/docker-desktop")
		os.Exit(1)
	}

	// Check auth (quick haiku test)
	if !sb.IsAuthenticated() {
		fmt.Println("✗ Docker sandbox not authenticated")
		fmt.Println()

		// Show auth prompt
		if ui.PromptAuthenticate() {
			fmt.Println("\nAuthenticating...")
			if err := sb.Authenticate(); err != nil {
				fmt.Printf("✗ Authentication failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✓ Authenticated!")
		} else {
			fmt.Println("Authentication required. Run again when ready.")
			os.Exit(0)
		}
	}

	// Prompt for iterations
	maxIterations := ui.PromptIterations()

	// Run loop with live dashboard
	controller := loop.NewController(cfg, sb, mode, maxIterations)
	if err := controller.Run(); err != nil {
		fmt.Printf("\n✗ Loop failed: %v\n", err)
		os.Exit(1)
	}
}
