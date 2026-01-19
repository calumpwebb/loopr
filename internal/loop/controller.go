package loop

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/calumpwebb/loopr/internal/git"
	"github.com/calumpwebb/loopr/internal/prompts"
	"github.com/calumpwebb/loopr/internal/sandbox"
)

type Controller struct {
	sandbox sandbox.Sandbox
	mode    string // "plan" or "build"
	maxIter int
}

func NewController(sb sandbox.Sandbox, mode string, maxIter int) *Controller {
	return &Controller{
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

	// Get embedded prompt based on mode
	var prompt string
	switch c.mode {
	case "plan":
		prompt = prompts.PlanPrompt
	case "build":
		prompt = prompts.BuildPrompt
	default:
		return fmt.Errorf("unknown mode: %s", c.mode)
	}

	// Get git branch
	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	// Run loop
	for i := 1; i <= c.maxIter; i++ {
		fmt.Printf("\n====== ITERATION %d/%d ======\n\n", i, c.maxIter)

		// Execute Claude with sonnet model
		if err := c.sandbox.ExecuteClaude(prompt, "sonnet"); err != nil {
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
