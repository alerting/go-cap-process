package tasks

import (
	"errors"

	"github.com/urfave/cli"

	"github.com/alerting/go-cap-process/db"
	"github.com/alerting/go-cap-process/db/elastic"
)

var (
	DatabaseFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "database, d",
			Usage:  "Database URL",
			EnvVar: "CAP_DATABASE_TYPE",
			Value:  "elasticsearch",
		},
		cli.StringFlag{
			Name:   "elastic-url",
			Usage:  "Elasticsearch URL",
			EnvVar: "CAP_ELASTIC_URL",
			Value:  "http://localhost:9200",
		},
		cli.StringFlag{
			Name:   "elastic-index",
			Usage:  "Elasticsearch index",
			EnvVar: "CAP_ELASTIC_INDEX",
			Value:  "alerts",
		},
	}
)

func CreateDatabase(c *cli.Context) (db.Database, error) {
	databaseType := getStringValue(c, "database")

	if databaseType == "elastic" || databaseType == "elasticsearch" || databaseType == "es" {
		return elastic.CreateDatabase(getStringValue(c, "elastic-url"), getStringValue(c, "elastic-index"))
	}

	return nil, errors.New("Unknown database type: " + databaseType)
}
