package fs

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alerting/go-cap"
)

const (
	FileTypeUnknown = iota
	FileTypeJSON
	FileTypeXML
)

type FileType int

// LoadAlert loads an alert from the reader.
func LoadAlert(r io.Reader, at FileType) (*cap.Alert, error) {
	var alert cap.Alert

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if at == FileTypeXML {
		err = xml.Unmarshal(b, &alert)
	} else if at == FileTypeJSON {
		err = json.Unmarshal(b, &alert)
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown type: %d", at))
	}

	return &alert, err
}

// LoadAlertFile loads an alert from file.
// Support file types are: .xml and .json
func LoadAlertFile(filename string) (*cap.Alert, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Determine the file type
	var t FileType

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(fi.Name()))
	if ext == ".xml" {
		t = FileTypeXML
	} else if ext == ".json" {
		t = FileTypeJSON
	}

	return LoadAlert(file, t)
}
