package main

import (
	"fmt"
	"log"
	"os"

	"github.com/calumpwebb/loopr/cmd"
	"github.com/urfave/cli/v2"
)

// version is set via ldflags during release builds
var (
	version = "dev" // Set via ldflags: -X main.version=v0.1.0
)

func getVersion() string {
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
					fmt.Printf("loopr version: %s\n", currentVersion)
					return nil
				},
			},
			{
				Name:  "guide",
				Usage: "<- AI (and humans) START HERE!",
				Action: func(c *cli.Context) error {
					cmd.Guide()
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
