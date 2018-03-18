package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"

	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/fs"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "cap-convert"
	app.Usage = "Convert CAP alerts from XML to JSON and vice-versa"
	app.Commands = []cli.Command{
		{
			Name:    "to-json",
			Aliases: []string{"json", "j"},
			Usage:   "Convert a CAP alert from XML to JSON",
			Action: func(c *cli.Context) error {
				for _, f := range c.Args() {
					var alert *cap.Alert
					var err error

					// Process file, - means stdin
					if f == "-" {
						alert, err = fs.LoadAlert(os.Stdin, fs.FileTypeXML)
					} else {
						alert, err = fs.LoadAlertFile(f)
					}

					if err != nil {
						return err
					}

					enc := json.NewEncoder(os.Stdout)
					return enc.Encode(alert)
				}
				return nil
			},
		},
		{
			Name:    "to-xml",
			Aliases: []string{"xml", "x"},
			Usage:   "Convert a CAP alert from JSON to XML",
			Action: func(c *cli.Context) error {
				for _, f := range c.Args() {
					var alert *cap.Alert
					var err error

					// Process file, - means stdin
					if f == "-" {
						alert, err = fs.LoadAlert(os.Stdin, fs.FileTypeJSON)
					} else {
						alert, err = fs.LoadAlertFile(f)
					}

					if err != nil {
						return err
					}

					enc := xml.NewEncoder(os.Stdout)
					return enc.Encode(alert)
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
