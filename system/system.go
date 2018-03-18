package system

import (
	"github.com/alerting/go-cap"
)

type System interface {
	GetAlert(reference *cap.Reference) (*cap.Alert, error)
}
