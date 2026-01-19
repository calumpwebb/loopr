package ui

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
)

// PromptSandbox prompts the user to select a sandbox type
func PromptSandbox() string {
	var sandboxType string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which sandbox would you like to use?").
				Options(
					huh.NewOption("docker", "docker"),
				).
				Value(&sandboxType),
		),
	)

	form.Run()
	return sandboxType
}

// PromptOverwrite prompts the user to confirm overwriting the .loopr/ directory
func PromptOverwrite() bool {
	var overwrite bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("⚠  .loopr/ directory already exists. Overwrite all files?").
				Value(&overwrite),
		),
	)

	form.Run()
	return overwrite
}

// PromptAuthenticate prompts the user to confirm authentication
func PromptAuthenticate() bool {
	var authenticate bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Authenticate now?").
				Value(&authenticate),
		),
	)

	form.Run()
	return authenticate
}

// PromptIterations prompts the user to input an iteration count with validation
func PromptIterations() int {
	var iterationsStr string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("How many iterations?").
				Value(&iterationsStr).
				Validate(func(s string) error {
					i, err := strconv.Atoi(s)
					if err != nil || i < 1 {
						return errors.New("must be >= 1")
					}
					return nil
				}),
		),
	)

	form.Run()
	iterations, _ := strconv.Atoi(iterationsStr)
	return iterations
}

// PromptUpdate prompts the user to confirm updating loopr
func PromptUpdate() bool {
	var update bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Update now?").
				Value(&update),
		),
	)

	form.Run()
	return update
}
