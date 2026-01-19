package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/calumpwebb/loopr/cmd"
	"github.com/urfave/cli/v2"
)

// version is set via ldflags during release builds
var version = "dev"

func getVersion() string {
	// If version is set via ldflags (release builds), use it
	if version != "dev" {
		return version
	}

	// Otherwise, try to get version from build info (development builds)
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	// Get version from module info
	moduleVersion := info.Main.Version
	if moduleVersion == "(devel)" || moduleVersion == "" {
		moduleVersion = "dev"
	}

	// Try to get git commit from build settings
	var commit string
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			commit = setting.Value
			if len(commit) > 7 {
				commit = commit[:7]
			}
			break
		}
	}

	if commit != "" {
		return fmt.Sprintf("%s (%s)", moduleVersion, commit)
	}
	return moduleVersion
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
