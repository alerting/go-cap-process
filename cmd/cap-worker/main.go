package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap-process/tasks"
)

func main() {
	app := cli.NewApp()

	app.Name = "cap-worker"
	app.Usage = "Tasks worker."
	app.Authors = []cli.Author{
		{
			Name: "Zachary Seguin",
		},
	}
	app.Copyright = "Copyright (c) 2018 Zachary Seguin"

	app.Flags = make([]cli.Flag, 0)
	app.Flags = append(app.Flags, tasks.ServerFlags...)
	app.Flags = append(app.Flags, tasks.SystemFlags...)
	app.Flags = append(app.Flags, tasks.DatabaseFlags...)

	app.Commands = []cli.Command{
		{
			Name:      "work",
			Aliases:   []string{"w"},
			Usage:     "Start a worker.",
			ArgsUsage: "",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "tag",
					Usage:  "Worker tag (ideally each worker should be unique)",
					EnvVar: "CAP_TAG",
					Value:  0,
				},
			},
			Action: work,
		},
		{
			Name:      "setup",
			Usage:     "Setup the database",
			ArgsUsage: "",
			Action:    setup,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
