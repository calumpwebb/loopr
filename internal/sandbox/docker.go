package sandbox

import (
	"os"
	"os/exec"
)

type DockerSandbox struct{}

func NewDocker() *DockerSandbox {
	return &DockerSandbox{}
}

func (d *DockerSandbox) IsAvailable() bool {
	cmd := exec.Command("docker", "ps")
	err := cmd.Run()
	return err == nil
}

func (d *DockerSandbox) IsAuthenticated() bool {
	// Check if Claude Code is fully set up by running a quick command
	// This detects both API auth AND first-time setup completion
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	cmd := exec.Command(
		"docker", "sandbox", "run",
		"-w", cwd,
		"claude",
		"--model=haiku",
		"--dangerously-skip-permissions",
		"Say OK",
	)

	// Suppress output, just check exit code
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	err = cmd.Run()
	return err == nil
}

func (d *DockerSandbox) Authenticate() error {
	// Run interactive auth with full Claude Code setup
	// This triggers first-time setup (theme selection, etc.) during auth phase
	cmd := exec.Command(
		"docker", "sandbox", "run",
		"claude",
		"-p",
		"--model=haiku",
		"--system-prompt=Reply with ONLY: 'Authentication successful!'",
		"Authenticate",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *DockerSandbox) ExecuteClaude(prompt string, model string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"docker", "sandbox", "run",
		"-w", cwd,
		"claude",
		"--dangerously-skip-permissions",
		"--model="+model,
		prompt,
	)

	// CRITICAL: Stream to terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
