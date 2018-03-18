package tasks

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/urfave/cli"
)

var (
	Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "broker, b",
			Usage:  "Message broker URL",
			EnvVar: "CAP_BROKER_URL",
			Value:  "redis://127.0.0.1:6379",
		},
		cli.StringFlag{
			Name:   "queue, q",
			Usage:  "Queue name",
			EnvVar: "CAP_QUEUE",
			Value:  "alerts",
		},
		cli.StringFlag{
			Name:   "result-backend, r",
			Usage:  "Result backend URL",
			EnvVar: "CAP_RESULTS_BACKEND",
			Value:  "redis://127.0.0.1:6379",
		},
		cli.IntFlag{
			Name:   "results-expiry, e",
			Usage:  "Time when results expire (in seconds)",
			EnvVar: "CAP_RESULTS_EXPIRY",
			Value:  120,
		},
	}
)

func getStringValue(c *cli.Context, arg string) string {
	if c.IsSet(arg) {
		return c.String(arg)
	} else if c.GlobalIsSet(arg) {
		return c.GlobalString(arg)
	}

	return ""
}

func getIntValue(c *cli.Context, arg string) int {
	if c.IsSet(arg) {
		return c.Int(arg)
	} else if c.GlobalIsSet(arg) {
		return c.GlobalInt(arg)
	}

	return 0
}

func CreateServer(c *cli.Context) (*machinery.Server, error) {
	// Create config from CLI flags
	var conf config.Config

	conf.Broker = getStringValue(c, "broker")
	conf.DefaultQueue = getStringValue(c, "queue")
	conf.ResultBackend = getStringValue(c, "result-backend")
	conf.ResultsExpireIn = getIntValue(c, "result-backend")

	return machinery.NewServer(&conf)
}
