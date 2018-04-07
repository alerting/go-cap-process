package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/fs"
	"github.com/alerting/go-cap-process/tasks"
)

func setup(c *cli.Context) error {
	// Connect to the database
	db, err := tasks.CreateDatabase(c)
	if err != nil {
		return err
	}

	return db.Setup()
}

func load(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("Provide at least one alert file to load")
	}

	// Connect to the database
	db, err := tasks.CreateDatabase(c)
	if err != nil {
		return err
	}

	alerts := make([]*cap.Alert, 0)
	for indx, f := range c.Args() {
		alert, err := fs.LoadAlertFile(f)
		if err != nil {
			return err
		}

		err = tasks.ProcessAlert(alert)
		if err != nil {
			return err
		}

		alerts = append(alerts, alert)

		if indx > 0 && indx%200 == 0 {
			err := db.AddAlert(alerts...)
			if err != nil {
				return err
			}

			alerts = make([]*cap.Alert, 0)
		}
	}

	if len(alerts) > 0 {
		err := db.AddAlert(alerts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "cap-load"
	app.Usage = "Load alerts from file into the database"
	app.Authors = []cli.Author{
		{
			Name: "Zachary Seguin",
		},
	}
	app.Copyright = "Copyright (c) 2018 Zachary Seguin"

	app.Flags = tasks.DatabaseFlags

	app.Commands = []cli.Command{
		{
			Name:      "import",
			Aliases:   []string{"i"},
			ArgsUsage: "file [file...]",
			Action:    load,
		},
		{
			Name:      "setup",
			ArgsUsage: "",
			Action:    setup,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
