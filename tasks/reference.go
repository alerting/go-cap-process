package tasks

import (
	"bytes"
	"encoding/json"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends"
	mtasks "github.com/RichardKnop/machinery/v1/tasks"

	"github.com/alerting/go-cap"
)

func EnsureReference(server *machinery.Server, reference *cap.Reference) (*backends.AsyncResult, error) {
	// Convert the reference to JSON
	var b bytes.Buffer

	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(&reference); err != nil {
		return nil, err
	}

	task := mtasks.Signature{
		UUID: reference.Id(),
		Name: "ensure_reference",
		Args: []mtasks.Arg{
			{
				Type:  "string",
				Value: b.String(),
			},
		},
		RetryCount: 0,
	}

	return server.SendTask(&task)
}
