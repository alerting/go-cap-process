package main

import (
	"github.com/urfave/cli"

	"github.com/alerting/go-cap-process/tasks"
)

func setup(c *cli.Context) error {
	var err error

	// Create the database
	database, err = tasks.CreateDatabase(c)
	if err != nil {
		return err
	}

	return database.Setup()
}
