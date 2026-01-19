package loop

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/calumpwebb/loopr/internal/config"
	"github.com/calumpwebb/loopr/internal/git"
	"github.com/calumpwebb/loopr/internal/sandbox"
)

type Controller struct {
	config  *config.Config
	sandbox sandbox.Sandbox
	mode    string // "plan" or "build"
	maxIter int
}

func NewController(cfg *config.Config, sb sandbox.Sandbox, mode string, maxIter int) *Controller {
	return &Controller{
		config:  cfg,
		sandbox: sb,
		mode:    mode,
		maxIter: maxIter,
	}
}

func (c *Controller) Run() error {
	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nInterrupted! Exiting gracefully...")
		os.Exit(0)
	}()

	// Load prompt file
	promptFile := c.getPromptFile() // .loopr/PROMPT_plan.md or PROMPT_build.md
	prompt, err := os.ReadFile(promptFile)
	if err != nil {
		return err
	}

	// Get git branch
	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	// Run loop
	for i := 1; i <= c.maxIter; i++ {
		fmt.Printf("\n====== ITERATION %d/%d ======\n\n", i, c.maxIter)

		// Execute Claude
		model := c.getModel()
		if err := c.sandbox.ExecuteClaude(string(prompt), model); err != nil {
			return err
		}

		// Push to git
		fmt.Println("\nPushing to git...")
		if err := git.Push(branch); err != nil {
			return err
		}
		fmt.Println("✓ Pushed to origin/" + branch)
	}

	fmt.Printf("\n✓ Completed %d/%d iterations\n", c.maxIter, c.maxIter)

	// Show next step suggestion
	if c.mode == "plan" {
		fmt.Println("\nNext: loopr build")
	}

	return nil
}

func (c *Controller) getPromptFile() string {
	looprDir := ".loopr"
	if c.config.LooprDir != "" {
		looprDir = c.config.LooprDir
	}
	return filepath.Join(looprDir, fmt.Sprintf("PROMPT_%s.md", c.mode))
}

func (c *Controller) getModel() string {
	// Use config model if specified, otherwise default to sonnet-4
	if c.config.Model != nil {
		if c.mode == "plan" && c.config.Model.Plan != "" {
			return c.config.Model.Plan
		}
		if c.mode == "build" && c.config.Model.Build != "" {
			return c.config.Model.Build
		}
	}
	return "claude-sonnet-4"
}
