package tasks

import (
	"errors"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap-process/system"
	"github.com/alerting/go-cap-process/system/canada-naad"
)

var (
	SystemFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "system, s",
			Usage:  "System",
			EnvVar: "CAP_SYSTEM",
		},

		// Canada NAAD
		cli.StringFlag{
			Name:   "canada-naad-fetch",
			Usage:  "Base URL for fetching alerts",
			EnvVar: "CAP_CANADA_NAAD_FETCH",
		},
	}
)

func CreateSystem(c *cli.Context) (system.System, error) {
	systemName := getStringValue(c, "system")

	if systemName == "canada-naad" {
		return canadanaad.CreateSystem(getStringValue(c, "canada-naad-fetch"))
	}

	return nil, errors.New("Unknown system: " + systemName)
}
