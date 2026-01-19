package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "loopr",
		Usage: "A tool for creating reproducible LLM workflows",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize a new loopr workflow",
				Action: func(c *cli.Context) error {
					fmt.Println("Initializing loopr workflow...")
					return nil
				},
			},
			{
				Name:  "plan",
				Usage: "Create or update workflow plan",
				Action: func(c *cli.Context) error {
					fmt.Println("Creating workflow plan...")
					return nil
				},
			},
			{
				Name:  "build",
				Usage: "Execute workflow based on plan",
				Action: func(c *cli.Context) error {
					fmt.Println("Building workflow...")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
