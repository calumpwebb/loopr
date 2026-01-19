package main

import (
	"log"
	"os"

	"github.com/calumpwebb/loopr/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "loopr",
		Usage: "Autonomous development orchestration with Claude and Ralph Loop",
		Commands: []*cli.Command{
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
