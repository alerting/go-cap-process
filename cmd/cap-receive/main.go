package main

import (
	"encoding/xml"
	"log"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/tasks"
)

// We should be receiving 1 alert per minute minimum
var (
	timeout = 80 * time.Second
)

func connect(c *cli.Context) error {
	// Ensure that the hostname:port has been specified
	if c.NArg() != 1 {
		cli.ShowCommandHelpAndExit(c, c.Command.FullName(), 1)
	}

	server, err := tasks.CreateServer(c)
	if err != nil {
		return err
	}

	// Connect to the stream service
	conn, err := net.Dial("tcp", c.Args().First())
	if err != nil {
		return err
	}

	log.Printf("Connected to %s\n", c.Args().First())

	decoder := xml.NewDecoder(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(timeout))

		var alert cap.Alert

		if err = decoder.Decode(&alert); err != nil {
			return err
		}

		log.Printf("Got alert: %s,%s,%s (%s)", alert.Sender, alert.Sent.FormatCAP(), alert.Identifier, alert.Id())

		_, err = tasks.AddAlert(server, &alert)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "cap-receive"
	app.Usage = "Listen to a TCP stream of alerts"
	app.Authors = []cli.Author{
		{
			Name: "Zachary Seguin",
		},
	}
	app.Copyright = "Copyright (c) 2018 Zachary Seguin"

	app.Flags = tasks.ServerFlags

	app.Commands = []cli.Command{
		{
			Name:      "connect",
			Aliases:   []string{"c"},
			Usage:     "Connect to a XML TCP stream.",
			ArgsUsage: "host:port",
			Action:    connect,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
