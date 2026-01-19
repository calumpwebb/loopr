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
	// Quick haiku test
	cmd := exec.Command(
		"docker", "sandbox", "run",
		"claude", "-p",
		"--model=claude-haiku-4",
		"--system-prompt=You must reply with only 'OK'. Nothing else.",
		"Say OK",
	)

	// Suppress output, just check exit code
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	return err == nil
}

func (d *DockerSandbox) Authenticate() error {
	// Run interactive auth
	cmd := exec.Command(
		"docker", "sandbox", "run",
		"claude", "--version",
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
