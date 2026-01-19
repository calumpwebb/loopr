package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/calumpwebb/loopr/cmd"
	"github.com/urfave/cli/v2"
)

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	// Get version from module info
	version := info.Main.Version
	if version == "(devel)" || version == "" {
		version = "dev"
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
		return fmt.Sprintf("%s (%s)", version, commit)
	}
	return version
}

func main() {
	version := getVersion()

	app := &cli.App{
		Name:    "loopr",
		Usage:   "Autonomous development orchestration with Claude and Ralph Loop",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version information",
				Action: func(c *cli.Context) error {
					fmt.Printf("loopr version %s\n", version)
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
