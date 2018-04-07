package db

import (
	"github.com/alerting/go-cap"
)

type Database interface {
	Setup() error

	AddAlert(alert ...*cap.Alert) error
	AlertExists(reference *cap.Reference) (bool, error)
	GetAlert(reference *cap.Reference) (*cap.Alert, error)

	NewInfoFinder() InfoFinder
}
