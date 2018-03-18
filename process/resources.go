package process

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/alerting/go-cap"
)

// fetchRemoteResource fetches a remote resource and
// saves it to the file at filePath.
func fetchRemoteResource(uri *url.URL, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	res, err := http.Get(uri.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	return f.Sync()
}

// decodeEncodedResource decodes a base64 encoded file
// and saves it a file at filePath.
func decodeEncodedResource(contents string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	res, err := base64.StdEncoding.DecodeString(contents)
	if err != nil {
		return err
	}

	_, err = f.Write(res)
	if err != nil {
		return err
	}

	return f.Sync()
}

// FetchResources fetches resources contained within the alert
// and outputs them in the the outputDirectory. Files are
// stored by their checksum + extension.
//
// Note: the first return value counts only new resources.
func FetchResources(alert *cap.Alert, outputDir string) (int, error) {
	fetched := 0

	for _, info := range alert.Infos {
		for _, resource := range info.Resources {
			// We can ignore application/x-url resources
			if resource.MimeType == "application/x-url" {
				continue
			}

			// Identify the path where to save the file
			uri, err := url.Parse(resource.Uri)
			if err != nil {
				return fetched, err
			}

			var filePath string

			if resource.Checksum() == "" {
				filePath = strings.ToLower(filepath.Join(outputDir,
					filepath.Base(uri.Path)))
			} else {
				filePath = strings.ToLower(filepath.Join(outputDir,
					resource.Checksum()+filepath.Ext(uri.Path)))
			}

			// If we already have the resource, we can ignore it.
			if _, err = os.Stat(filePath); err != nil {
				// Ensure that we have a valid resource
				if resource.DerefUri == nil && uri.Hostname() == "" {
					return fetched, errors.New("Invalid resource: no hostname or DerefUri")
				}

				// Handle the resource based on type
				if resource.DerefUri != nil {
					err := decodeEncodedResource(*resource.DerefUri, filePath)
					if err != nil {
						return fetched, err
					}
				} else {
					err := fetchRemoteResource(uri, filePath)
					if err != nil {
						return fetched, err
					}
				}

				fetched++
			}

		}
	}

	return fetched, nil
}
