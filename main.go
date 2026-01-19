package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/calumpwebb/loopr/cmd"
	"github.com/urfave/cli/v2"
)

const (
	githubURL = "https://github.com/calumpwebb/loopr"
)

// version and commit are set via ldflags during release builds
var (
	version = "dev" // Set via ldflags: -X main.version=v0.1.0
	commit  = ""    // Set via ldflags: -X main.commit=abc1234
)

func getVersion() string {
	// For release builds, both version and commit are set via ldflags
	if version != "dev" && commit != "" {
		return fmt.Sprintf("%s (%s)", version, commit)
	}

	// For dev builds, try to get commit from build info
	if version == "dev" {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return "unknown"
		}

		// Try to get git commit from build settings
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				c := setting.Value
				if len(c) > 7 {
					c = c[:7]
				}
				return fmt.Sprintf("%s (%s)", version, c)
			}
		}
	}

	// Fallback: just return version
	return version
}

func main() {
	currentVersion := getVersion()

	app := &cli.App{
		Name:    "loopr",
		Usage:   "Autonomous development orchestration with Claude and Ralph Loop",
		Version: currentVersion,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version information",
				Action: func(c *cli.Context) error {
					fmt.Printf("loopr version %s\n", currentVersion)
					fmt.Printf("GitHub: %s\n", githubURL)
					return nil
				},
			},
			{
				Name:  "init",
				Usage: "Initialize loopr in current directory",
				Action: func(c *cli.Context) error {
					return cmd.Init()
				},
			},
			{
				Name:  "plan",
				Usage: "Generate implementation plan",
				Action: func(c *cli.Context) error {
					cmd.Plan()
					return nil
				},
			},
			{
				Name:  "build",
				Usage: "Build from plan",
				Action: func(c *cli.Context) error {
					cmd.Build()
					return nil
				},
			},
			{
				Name:      "import",
				Usage:     "Import tasks from a PRD/spec file",
				ArgsUsage: "<file>",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						fmt.Println("Usage: loopr import <file>")
						fmt.Println("\nExample:")
						fmt.Println("  loopr import docs/feature-spec.md")
						fmt.Println("  loopr import .loopr/prd/new-feature.md")
						os.Exit(1)
					}
					sourceFile := c.Args().First()
					cmd.Import(sourceFile)
					return nil
				},
			},
			{
				Name:  "archive",
				Usage: "Archive completed tasks from tasks.md",
				Action: func(c *cli.Context) error {
					cmd.Archive()
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "Show current task status",
				Action: func(c *cli.Context) error {
					cmd.Status()
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "Update loopr to the latest version",
				Action: func(c *cli.Context) error {
					return cmd.Update(currentVersion)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
