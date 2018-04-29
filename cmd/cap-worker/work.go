package main

import (
	"encoding/json"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/urfave/cli"

	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/db"
	capsys "github.com/alerting/go-cap-process/system"
	"github.com/alerting/go-cap-process/tasks"
)

var (
	server   *machinery.Server
	database db.Database
	system   capsys.System
)

func ensureReference(referenceJSON string) error {
	// Convert alert
	var reference cap.Reference

	if err := json.Unmarshal([]byte(referenceJSON), &reference); err != nil {
		return err
	}

	log.INFO.Printf("Ensuring referenced alert exists: %s,%s,%s (%s)",
		reference.Sender,
		reference.Sent.FormatCAP(),
		reference.Identifier,
		reference.Id())

	exists, err := database.AlertExists(&reference)
	if err != nil {
		return err
	}

	// Don't need to continue if we already have it.
	if exists {
		log.INFO.Println("Alert already exists")
		return nil
	}

	log.INFO.Println("Alert does not exist, fetching")
	alert, err := system.GetAlert(&reference)
	if err != nil {
		return err
	}

	_, err = tasks.AddAlert(server, alert)
	return err
}

func addAlert(alertJSON string) error {
	// Convert alert
	var alert cap.Alert

	if err := json.Unmarshal([]byte(alertJSON), &alert); err != nil {
		return err
	}

	log.INFO.Printf("Got alert: %s,%s,%s (%s)",
		alert.Sender,
		alert.Sent.FormatCAP(),
		alert.Identifier,
		alert.Id())

	// Perform any processing/cleanup
	if err := tasks.ProcessAlert(&alert); err != nil {
		return err
	}

	// Queue tasks to ensure that references have been loaded
	for _, reference := range alert.References {
		_, err := tasks.EnsureReference(server, reference)
		if err != nil {
			return err
		}
	}

	// Stop processing system messages
	if alert.Status == cap.StatusSystem {
		log.INFO.Printf("Got system message, will stop processing")
		return nil
	}

	// TODO: Fetch resources

	return database.AddAlert(&alert)
}

func work(c *cli.Context) error {
	var err error

	// Create the system
	system, err = tasks.CreateSystem(c)
	if err != nil {
		return err
	}

	// Create the database
	database, err = tasks.CreateDatabase(c)
	if err != nil {
		return err
	}

	// Create the server
	server, err = tasks.CreateServer(c)
	if err != nil {
		return err
	}

	// Register tasks
	tasks := map[string]interface{}{
		"add_alert":        addAlert,
		"ensure_reference": ensureReference,
	}

	if err = server.RegisterTasks(tasks); err != nil {
		return err
	}

	// Start the worker
	worker := server.NewWorker("worker", c.Int("tag"))
	return worker.Launch()
}
