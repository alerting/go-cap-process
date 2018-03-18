package main

import (
	"encoding/json"

	"github.com/RichardKnop/machinery/v1/log"
	"github.com/urfave/cli"

	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/tasks"
)

func addAlert(alertJSON string) error {
	// Convert alert
	var alert cap.Alert

	if err := json.Unmarshal([]byte(alertJSON), &alert); err != nil {
		return err
	}

	log.INFO.Printf("Got alert: %s,%s,%s (%s)", alert.Sender, alert.Sent.FormatCAP(), alert.Identifier, alert.Id())

	return nil
}

func work(c *cli.Context) error {
	// Create the server
	server, err := tasks.CreateServer(c)
	if err != nil {
		return err
	}

	// Register tasks
	tasks := map[string]interface{}{
		"add_alert": addAlert,
	}

	if err = server.RegisterTasks(tasks); err != nil {
		return err
	}

	// Start the worker
	worker := server.NewWorker("worker", c.Int("tag"))
	return worker.Launch()
}
