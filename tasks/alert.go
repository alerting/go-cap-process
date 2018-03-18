package tasks

import (
	"bytes"
	"encoding/json"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends"
	mtasks "github.com/RichardKnop/machinery/v1/tasks"

	"github.com/alerting/go-cap"
)

func AddAlert(server *machinery.Server, alert *cap.Alert) (*backends.AsyncResult, error) {
	// Convert the alert to JSON
	var b bytes.Buffer

	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(&alert); err != nil {
		return nil, err
	}

	task := mtasks.Signature{
		Name: "add_alert",
		Args: []mtasks.Arg{
			{
				Type:  "string",
				Value: b.String(),
			},
		},
	}

	return server.SendTask(&task)
}
