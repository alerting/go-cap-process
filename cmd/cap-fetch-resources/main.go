package main

import (
	"flag"
	"log"

	"github.com/alerting/go-cap-process/fs"
	"github.com/alerting/go-cap-process/process"
)

func main() {
	// Parse command line
	outDir := flag.String("out-dir", "resources", "Directory to place resources")
	flag.Parse()

	for _, f := range flag.Args() {
		log.Printf("Processing %s\n", f)
		alert, err := fs.LoadAlertFile(f)
		if err != nil {
			log.Fatal(err)
		}

		fetched, err := process.FetchResources(alert, *outDir)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("..Got %d new resource(s)\n", fetched)
	}
}
