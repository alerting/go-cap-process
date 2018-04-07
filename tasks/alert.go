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
		RetryCount: 20,
	}

	return server.SendTask(&task)
}

func ProcessAlert(alert *cap.Alert) error {
	for indx, info := range alert.Infos {
		// If the info contains no effective time,
		// the specification says to use the sent time.
		if info.Effective == nil {
			alert.Infos[indx].Effective = &alert.Sent
		}
	}

	return nil
}
